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

func (data Database) findCalendarFromUser(userEmail string, userUUID string, calendarUUID string) (calendar api.CalendarManager, err error) {
	var id string
	var tokenType string
	var refreshToken string
	var email string
	var kind int
	var accessToken string
	var principal bool
	err = data.client.QueryRow("SELECT calendars.id, a.kind, a.token_type, a.refresh_token, a.email, a.access_token, a.principal from calendars join accounts a on calendars.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email=$2 and calendars.uuid =$3", userUUID, userEmail, calendarUUID).Scan(&id, &kind, &tokenType, &refreshToken, &email, &accessToken, &principal)
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
		calendar = api.RetrieveGoogleCalendar(id, api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
	case api.OUTLOOK:
		calendar = api.RetrieveOutlookCalendar(id, api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
	default:
		data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: email}, map[string]string{"database": "backend"})
		log.Errorf("kind of calendar is not valid: %d", kind)
		return nil, &customErrors.WrongKindError{Mail: email}
	}
	calendar.SetUUID(calendarUUID)
	calendars, err := data.getSynchronizedCalendars(calendar, principal)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving sync calendars from %s", calendarUUID)
		return
	}
	calendar.SetCalendars(calendars)
	return
}

func (data Database) getSynchronizedCalendars(calendar api.CalendarManager, principal bool) (calendars []api.CalendarManager, err error) {
	var query string
	if principal {
		query = "select calendars.id, calendars.uuid, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = $1"
	} else {
		query = "select calendars.id, calendars.uuid, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = (Select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.uuid = (select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1)"
	}
	rows, err := data.client.Query(query, calendar.GetUUID())
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
			calendar = api.RetrieveGoogleCalendar(id, api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
		case api.OUTLOOK:
			calendar = api.RetrieveOutlookCalendar(id, api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
		default:
			data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: calendar.GetName()}, map[string]string{"database": "backend"})
			return nil, &customErrors.WrongKindError{Mail: calendar.GetName()}
		}
		calendar.SetUUID(uid)
		calendars = append(calendars, calendar)
	}
	return
}

func (data Database) UpdateCalendarFromUser(calendar api.CalendarManager, userUUID string) (err error) {
	err = data.updateCalendarFromUser(calendar, userUUID)
	return
}

func (data Database) updateCalendarFromUser(calendar api.CalendarManager, userUUID string) (err error) {
	stmt, err := data.client.Prepare("update calendars set name = $1 from accounts where calendars.account_email = accounts.email and accounts.user_uuid =$2 and calendars.id =$3;")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(calendar.GetName(), userUUID, calendar.GetID())
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
	stmt, err := transaction.Prepare("insert into subscriptions(uuid,calendar_uuid,id, type, expiration_date) values ($1,$2,$3,$4,$5)")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(subscription.GetUUID(), calendar.GetUUID(), subscription.GetID(), subscription.GetType(), subscription.GetExpirationDate())
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
