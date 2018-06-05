package db

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (data Database) RetrieveSyncedEventsWithSubscription(eventID string, subscriptionID string, calendar api.CalendarManager) (events []api.EventManager, found bool, err error) {
	var principalEventID int
	found = true
	err = data.client.QueryRow("SELECT COALESCE(events.parent_event_internal_id, events.internal_id) from events join calendars c2 on events.calendar_uuid = c2.uuid join subscriptions s2 on c2.uuid = s2.calendar_uuid where events.id = $1 and s2.id=$2", eventID, subscriptionID).Scan(&principalEventID)
	switch {
	case err == sql.ErrNoRows:
		log.Warningf("principal event from event ID: %s and subscription ID: %s not found", eventID, subscriptionID)
		found = false
		err = nil
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting principal event from event id: %s and subscription ID: %s", eventID, subscriptionID)
		return nil, false, err
	}
	if found {
		events, err = data.getSynchronizedEventsFromEvent(principalEventID, eventID)
	} else {
		calendars, err := data.getSynchronizedCalendars(calendar)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error retrieving synced calendars: %s", calendar.GetID())
			return nil, false, err
		}
		for _, calendar := range calendars {
			event := calendar.CreateEmptyEvent("")
			events = append(events, event)
		}

	}
	return
}

func (data Database) getSynchronizedEventsFromEvent(principalEventID int, eventID string) (events []api.EventManager, err error) {
	stmt, err := data.client.Prepare("select events.id, a.kind, a.token_type, a.refresh_token, a.email, a.access_token, c2.id, c2.uuid from events join calendars c2 on events.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email where events.internal_id = $1 or events.parent_event_internal_id=$1 and events.id!=$2")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting synced events from principalID: %d", principalEventID)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(principalEventID, eventID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting synced events from principalID: %d", principalEventID)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var kind int
		var tokenType string
		var refreshToken string
		var email string
		var accessToken string
		var calendarID string
		var calendarUUID string
		err = rows.Scan(&id, &kind, &tokenType, &refreshToken, &email, &accessToken, &calendarID, &calendarUUID)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error scanning synced events from principalID: %d", principalEventID)
			return nil, err
		}
		var eventSync api.EventManager
		var calendar api.CalendarManager
		switch kind {
		case api.GOOGLE:
			account := api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken)
			calendar = api.RetrieveGoogleCalendar(calendarID, calendarUUID, "", account)
			eventSync = &api.GoogleEvent{ID: id}
		case api.OUTLOOK:
			account := api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken)
			calendar = api.RetrieveOutlookCalendar(calendarID, calendarUUID, account)
			eventSync = &api.OutlookEvent{ID: id}
		default:
			err = &customErrors.WrongKindError{Mail: fmt.Sprintf("wrong kind of account for events with parent ID: %d", principalEventID)}
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			return nil, err
		}
		err = eventSync.SetCalendar(calendar)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error setting calendar for event ID: %s", eventID)
		}
		events = append(events, eventSync)
	}
	return
}

func (data Database) prepareEventsToSync(subscriptionID string) (events []api.EventManager, err error) {

	stmt, err := data.client.Prepare("select a.kind, a.token_type, a.refresh_token, a.email, a.access_token, c2.id from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email where subs")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing events to sync from subscriptionID: %s", subscriptionID)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(subscriptionID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing events to sync from subscriptionID: %s", subscriptionID)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {

	}
	return
}

func (data Database) SavePrincipalEvent(event api.EventManager) (err error) {
	transaction, err := data.client.Begin()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error starting transaction: %s", err.Error())
		return
	}
	err = data.savePrincipalEvent(transaction, event)
	if err != nil {
		transaction.Rollback()
		return err
	} else {
		transaction.Commit()
	}

	return

}

func (data Database) savePrincipalEvents(transaction *sql.Tx, events []api.EventManager) (err error) {
	for _, event := range events {
		err = data.savePrincipalEvent(transaction, event)
		if err != nil {
			break
		}
	}
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
	}
	return

}
func (data Database) savePrincipalEvent(transaction *sql.Tx, event api.EventManager) (err error) {
	lastInsertId := 0
	updatedAt, err := event.GetUpdatedAt()
	if err != nil {
		log.Errorf("error getting updated at for event: %s", event.GetID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return err
	}
	err = transaction.QueryRow("INSERT INTO events (calendar_uuid, id, updated_at) VALUES($1, $2, $3) RETURNING internal_id", event.GetCalendar().GetUUID(), event.GetID(), updatedAt).Scan(&lastInsertId)
	switch {
	case err == sql.ErrNoRows:
		err = fmt.Errorf("could not insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("could not insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
		return err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
		return err
	}
	event.SetInternalID(lastInsertId)
	return
}

func (data Database) SaveEventsRelation(from api.EventManager, to api.EventManager) (err error) {
	transaction, err := data.client.Begin()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error starting transaction: %s", err.Error())
		return
	}
	err = data.saveEventsRelation(transaction, from, to)
	if err != nil {
		transaction.Rollback()
		return err
	} else {
		transaction.Commit()
	}
	return

}
func (data Database) saveEventsRelation(transaction *sql.Tx, from api.EventManager, to api.EventManager) (err error) {
	updatedAt, err := to.GetUpdatedAt()
	if err != nil {
		log.Errorf("error getting updated at for event: %s", to.GetID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return
	}
	stmt, err := transaction.Prepare("insert into events(calendar_uuid, id, parent_event_internal_id, updated_at) values ($1,$2,$3,$4)")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(to.GetCalendar().GetUUID(), to.GetID(), from.GetInternalID(), updatedAt)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	//TODO: Change this...
	if affect != 1 {
		err = errors.New(fmt.Sprintf("could not update account with id: %s", to.GetID()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return
	}
	return
}

func (data Database) deleteEventsFromSubscription(transaction *sql.Tx, subscription api.SubscriptionManager) (err error) {
	stmt, err := transaction.Prepare("delete from events using subscriptions, calendars where subscriptions.uuid = $1 and subscriptions.calendar_uuid = calendars.uuid and events.calendar_uuid = calendars.uuid")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing statement: %s", err.Error())
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(subscription.GetUUID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting rows affected: %s", err.Error())
	}
	log.Infof("deleted %d events from subscription: %s", rows, subscription.GetUUID())

	return

}

func (data Database) EventAlreadyUpdated(event api.EventManager) bool {
	var exists bool
	updatedAt, err := event.GetUpdatedAt()
	if err != nil {
		log.Errorf("error getting updated at for event: %s", event.GetID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return false
	}
	err = data.client.QueryRow("select true from events where events.id=$1 and events.updated_at != $2", event.GetID(), updatedAt).Scan(&exists)
	switch {
	case err == sql.ErrNoRows:
		return true
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting synchronized date for event: %s", event.GetID())
		return false
	}
	return false
}

func (data Database) UpdateModificationDate(event api.EventManager) error {
	updatedAt, err := event.GetUpdatedAt()
	if err != nil {
		log.Errorf("error getting updated at for event: %s", event.GetID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return err
	}
	stmt, err := data.client.Prepare("update events set updated_at= $1 where events.id=$2")

	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(updatedAt, event.GetID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return err
	}
	//TODO: Change this...
	if affect != 1 {
		err = errors.New(fmt.Sprintf("could not update account with id: %s", event.GetID()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return err
	}
	return nil
}

func (data Database) DeleteEvent(event api.EventManager) error {
	stmt, err := data.client.Prepare("delete from events where events.id =$1")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(event.GetID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return err
	}
	//TODO: Change this...
	if affect != 1 {
		err = errors.New(fmt.Sprintf("could not delete event with id: %s", event.GetID()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return err
	}
	return nil

}

func (data Database) ExistsEvent(event api.EventManager) bool {
	var exists bool
	err := data.client.QueryRow("select true from events join calendars c2 on events.calendar_uuid = c2.uuid where events.id = $1 and c2.id=$2", event.GetID(), event.GetCalendar().GetID()).Scan(&exists)

	switch {
	case err == sql.ErrNoRows:
		log.Warningf("event ID: %s not found", event.GetID())
		exists = false
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error querying event id: %s ", event.GetID())
		return false
	}
	return exists
}

func (data Database) GetGoogleEventIDs(subscriptionID string) (eventIDs []string, err error) {
	stmt, err := data.client.Prepare("select events.id from events join calendars c2 on events.calendar_uuid = c2.uuid join subscriptions s2 on c2.uuid = s2.calendar_uuid where c2.id=$1")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting stored events from subscriptionID: %s", subscriptionID)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(subscriptionID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting stored events from subscriptionID: %s", subscriptionID)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error scanning stored events from subscriptionID: %s", subscriptionID)
			return nil, err
		}
		eventIDs = append(eventIDs, id)
	}
	return
}
