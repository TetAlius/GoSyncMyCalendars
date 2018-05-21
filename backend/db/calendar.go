package db

import (
	"database/sql"
	"fmt"

	"net/http"

	"errors"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func RetrieveCalendars(userEmail string, userUUID string, calendarUUID string) (calendar api.CalendarManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	calendar, err = findCalendarFromUser(db, userEmail, userUUID, calendarUUID)
	if err != nil {
		log.Errorf("error retrieving calendar %s from user %s", calendarUUID, userUUID)
	}
	defer db.Close()
	return
}

func UpdateAllCalendarsFromUser(userUUID string, userEmail string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}

	rows, err := db.Query("SELECT calendars.id, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email=$2", userUUID, userEmail)
	if err != nil {
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
			log.Errorf("kind of calendar is not valid: %d", kind)
			return &customErrors.WrongKindError{Mail: email}
		}
		account.Refresh()
		UpdateAccountFromUser(account, userUUID)
		calendar, err := account.GetCalendar(id)
		if err != nil {
			log.Errorf("error: %s", err.Error())
		} else {
			updateCalendarFromUser(db, calendar, userUUID)
		}

	}

	defer db.Close()
	return
}

func findCalendarFromUser(db *sql.DB, userEmail string, userUUID string, calendarUUID string) (calendar api.CalendarManager, err error) {
	rows, err := db.Query("SELECT calendars.id, a.kind, a.token_type, a.refresh_token, a.email, a.access_token, a.principal from calendars join accounts a on calendars.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email=$2 and calendars.uuid =$3", userUUID, userEmail, calendarUUID)
	if err != nil {
		log.Errorf("error querying get calendar")
		return
	}
	defer rows.Close()
	var principal bool
	if rows.Next() {
		var id string
		var tokenType string
		var refreshToken string
		var email string
		var kind int
		var accessToken string
		rows.Scan(&id, &kind, &tokenType, &refreshToken, &email, &accessToken, &principal)
		switch kind {
		case api.GOOGLE:
			calendar = &api.GoogleCalendar{ID: id}
			calendar.SetAccount(api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
		case api.OUTLOOK:
			calendar = &api.OutlookCalendar{ID: id}
			calendar.SetAccount(api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
		default:
			log.Errorf("kind of calendar is not valid: %d", kind)
			return nil, &customErrors.WrongKindError{Mail: email}
		}
	} else {
		log.Errorf("calendar %s not found", calendarUUID)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	}
	calendar.SetUUID(calendarUUID)
	calendars, err := getSynchronizedCalendars(db, calendar, principal)
	if err != nil {
		log.Errorf("error retrieving sync calendars from %s", calendarUUID)
	}
	calendar.SetCalendars(calendars)
	return
}

func getSynchronizedCalendars(db *sql.DB, calendar api.CalendarManager, principal bool) (calendars []api.CalendarManager, err error) {
	var query string
	if principal {
		query = "select calendars.id, calendars.uuid, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = $1"
	} else {
		query = "select calendars.id, calendars.uuid, a.kind, a.token_type, a.refresh_token, a.email, a.access_token from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = (Select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.uuid = (select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1)"
	}
	rows, err := db.Query(query, calendar.GetUUID())
	if err != nil {
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
			calendar = &api.GoogleCalendar{ID: id}
			calendar.SetAccount(api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken))
		case api.OUTLOOK:
			calendar = &api.OutlookCalendar{ID: id}
			calendar.SetAccount(api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken))
		default:
			return nil, &WrongKindError{calendar.GetName()}
		}
		calendar.SetUUID(uid)
		calendars = append(calendars, calendar)
	}
	return
}

func UpdateCalendarFromUser(calendar api.CalendarManager, userUUID string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	err = updateCalendarFromUser(db, calendar, userUUID)
	defer db.Close()
	return
}

func updateCalendarFromUser(db *sql.DB, calendar api.CalendarManager, userUUID string) (err error) {
	stmt, err := db.Prepare("update calendars set name = $1 from accounts where calendars.account_email = accounts.email and accounts.user_uuid =$2 and calendars.id =$3;")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(calendar.GetName(), userUUID, calendar.GetID())
	if err != nil {
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect < 1 {
		log.Errorf("could not update calendar with id: %s from user: %s", calendar.GetID(), userUUID)
		return errors.New(fmt.Sprintf("could not update calendar with id: %s from user: %s", calendar.GetID(), userUUID))
	}
	return

}

func SaveSubscription(subscription api.SubscriptionManager, calendar api.CalendarManager) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into subscriptions(uuid,calendar_uuid,id) values ($1,$2,$3)")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(subscription.GetUUID(), calendar.GetUUID(), subscription.GetID())
	if err != nil {
		log.Errorf("error executing query: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		return errors.New(fmt.Sprintf("could not create new subscription for calendar name: %s", calendar.GetName()))
	}
	return
}
