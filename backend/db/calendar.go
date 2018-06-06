package db

import (
	"database/sql"
	"fmt"

	"errors"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (data Database) RetrieveCalendars(userEmail string, userUUID string, calendarUUID string) (calendar api.CalendarManager, err error) {
	calendar, err = data.findCalendarFromUser(userEmail, userUUID, calendarUUID)
	if err != nil {
		log.Errorf("error retrieving calendar %s from user %s", calendarUUID, userUUID)
	}
	return
}

func (data Database) UpdateAllCalendarsFromUser(userUUID string, userEmail string) (err error) {
	rows, err := data.client.Query("SELECT calendars.id, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email=$2", userUUID, userEmail)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error querying get calendar: %s", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var tokenType string
		var refreshToken string
		var email string
		var kind int
		var accessToken string
		var account api.AccountManager
		rows.Scan(&id, &kind, &tokenType, &refreshToken, &email, &accessToken)
		switch kind {
		case api.GOOGLE:
			account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken)
		case api.OUTLOOK:
			account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken)
		default:
			data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: email}, map[string]string{"database": "backend"})
			log.Errorf("kind of calendar is not valid: %d", kind)
			return &customErrors.WrongKindError{Mail: email}
		}
		//TODO: manage errors
		account.Refresh()
		data.UpdateAccountFromUser(account, userUUID)
		calendar, err := account.GetCalendar(id)
		if err != nil {
			log.Errorf("error: %s", err.Error())
		} else {
			data.updateCalendarFromUser(calendar, userUUID)
		}

	}
	return
}

func (data Database) RetrieveCalendarFromSubscription(subscriptionID string) (calendar api.CalendarManager, err error) {
	var tokenType string
	var refreshToken string
	var email string
	var kind int
	var accessToken string
	var calendarID string
	var syncToken string
	var uid string
	err = data.client.QueryRow("SELECT a.token_type, a.refresh_token,a.email,a.kind,a.access_token, calendars.id, calendars.sync_token, calendars.uuid from calendars join subscriptions s2 on calendars.uuid = s2.calendar_uuid join accounts a on calendars.account_email = a.email where s2.id = $1", subscriptionID).
		Scan(&tokenType, &refreshToken, &email, &kind, &accessToken, &calendarID, &syncToken, &uid)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("calendar from subscription with ID: %s not found", subscriptionID)}
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("calendar from subscription with ID: %s not found", subscriptionID)
		return nil, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("error getting calendar from subscription with ID: %s", subscriptionID)
		return nil, err
	}
	switch kind {
	case api.OUTLOOK:
		account := api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken)
		calendar = api.RetrieveOutlookCalendar(calendarID, uid, account)
	case api.GOOGLE:
		account := api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken)
		calendar = api.RetrieveGoogleCalendar(calendarID, uid, syncToken, account)
	default:
		return nil, &customErrors.WrongKindError{Mail: fmt.Sprintf("error getting calendar with subscription ID: %s", subscriptionID)}
	}

	return
}

func (data Database) findCalendarFromUser(userEmail string, userUUID string, calendarUUID string) (calendar api.CalendarManager, err error) {
	var id string
	var tokenType string
	var refreshToken string
	var email string
	var kind int
	var accessToken string
	var uid string
	err = data.client.QueryRow("SELECT calendars.id, calendars.uuid,a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email=$2 and calendars.uuid =$3", userUUID, userEmail, calendarUUID).Scan(&id, &uid, &kind, &tokenType, &refreshToken, &email, &accessToken)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("No account from user: %s with that uuid: %s.", userUUID, calendarUUID)}
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("No account from user: %s with that id: %d.", userUUID, id)
		return nil, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("error looking for account from user: %s with id: %d.", userUUID, id)
		return
	}
	switch kind {
	case api.GOOGLE:
		calendar = api.RetrieveGoogleCalendar(id, uid, "", api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
	case api.OUTLOOK:
		calendar = api.RetrieveOutlookCalendar(id, uid, api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
	default:
		data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: email}, map[string]string{"database": "backend"})
		log.Errorf("kind of calendar is not valid: %d", kind)
		return nil, &customErrors.WrongKindError{Mail: email}
	}
	calendar.SetUUID(calendarUUID)
	calendars, err := data.getSynchronizedCalendars(calendar)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving sync calendars from %s", calendarUUID)
		return
	}
	calendar.SetCalendars(calendars)
	return
}

func (data Database) getSynchronizedCalendars(calendar api.CalendarManager) (calendars []api.CalendarManager, err error) {
	rows, err := data.client.Query("select calendars.id, calendars.uuid, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email where (calendars.parent_calendar_uuid = (Select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.uuid = (select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.parent_calendar_uuid = $1) AND calendars.uuid != $1", calendar.GetUUID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error selecting setSynchronizedCalendars: %s", err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var uid string
		var tokenType string
		var refreshToken string
		var email string
		var accessToken string
		var calendar api.CalendarManager
		var kind int
		err = rows.Scan(&id, &uid, &kind, &tokenType, &refreshToken, &email, &accessToken)
		switch kind {
		case api.GOOGLE:
			calendar = api.RetrieveGoogleCalendar(id, uid, "", api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
		case api.OUTLOOK:
			calendar = api.RetrieveOutlookCalendar(id, uid, api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
		default:
			data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: calendar.GetName()}, map[string]string{"database": "backend"})
			return nil, &customErrors.WrongKindError{Mail: calendar.GetName()}
		}
		calendars = append(calendars, calendar)
	}
	return
}

func (data Database) UpdateCalendarFromUser(calendar api.CalendarManager, userUUID string) (err error) {
	err = data.updateCalendarFromUser(calendar, userUUID)
	return
}

func (data Database) updateCalendarFromUser(calendar api.CalendarManager, userUUID string) (err error) {
	stmt, err := data.client.Prepare("update calendars set name = $1, sync_token=$2 from accounts where calendars.account_email = accounts.email and accounts.user_uuid =$3 and calendars.id =$4;")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(calendar.GetName(), calendar.GetSyncToken(), userUUID, calendar.GetID())
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
	if affect < 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not update calendar with id: %s from user: %s", calendar.GetID(), userUUID)), map[string]string{"database": "backend"})
		log.Errorf("could not update calendar with id: %s from user: %s", calendar.GetID(), userUUID)
		return errors.New(fmt.Sprintf("could not update calendar with id: %s from user: %s", calendar.GetID(), userUUID))
	}
	return

}

func (data Database) saveSubscription(transaction *sql.Tx, subscription api.SubscriptionManager, calendar api.CalendarManager) (err error) {
	stmt, err := transaction.Prepare("insert into subscriptions(uuid,calendar_uuid,id, type, expiration_date, resource_id) values ($1,$2,$3,$4,$5,$6)")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(subscription.GetUUID(), calendar.GetUUID(), subscription.GetID(), subscription.GetType(), subscription.GetExpirationDate(), subscription.GetResourceID())
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
	if affect != 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not create new subscription for calendar name: %s", calendar.GetName())), map[string]string{"database": "backend"})
		return errors.New(fmt.Sprintf("could not create new subscription for calendar name: %s", calendar.GetName()))
	}
	return
}
