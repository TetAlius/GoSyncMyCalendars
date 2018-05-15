package db

import (
	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func findCalendarsFromAccount(db *sql.DB, account api.AccountManager) (calendars []api.CalendarManager, err error) {
	rows, err := db.Query("select calendars.id, calendars.name, calendars.uuid from calendars join accounts a on calendars.account_email = a.email where a.id=$1 order by calendars.name ASC", account.GetInternalID())
	if err != nil {
		log.Errorln("error selecting findCalendarsFromAccount")
		return
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var uuid string
		var calendar api.CalendarManager
		err = rows.Scan(&id, &name, &uuid)
		switch account.GetKind() {
		case api.GOOGLE:
			calendar = &api.GoogleCalendar{ID: id, Name: name}
		case api.OUTLOOK:
			calendar = &api.OutlookCalendar{ID: id, Name: name}
		default:
			return nil, &WrongKindError{name}
		}
		calendar.SetUUID(uuid)
		calendar.SetAccount(account)
		synced, err := getSynchronizedCalendars(db, calendar, account.Principal())
		if err != nil {
			log.Errorf("error getting relation: %s", err.Error())
			return nil, err
		}
		calendar.SetCalendars(synced)
		calendars = append(calendars, calendar)
	}
	return
}

func getSynchronizedCalendars(db *sql.DB, calendar api.CalendarManager, principal bool) (calendars []api.CalendarManager, err error) {
	var query string
	if principal {
		query = "select calendars.id, calendars.name, calendars.uuid, a.kind from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = $1"
	} else {
		query = "select calendars.id, calendars.name, calendars.uuid, a.kind from calendars join accounts a on calendars.account_email = a.email where calendars.parent_calendar_uuid = (Select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.uuid = (select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1)"
	}
	rows, err := db.Query(query, calendar.GetUUID())
	if err != nil {
		log.Errorln("error selecting setSynchronizedCalendars")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var uuid string
		var calendar api.CalendarManager
		var kind int
		err = rows.Scan(&id, &name, &uuid, &kind)
		switch kind {
		case api.GOOGLE:
			calendar = &api.GoogleCalendar{ID: id, Name: name}
		case api.OUTLOOK:
			calendar = &api.OutlookCalendar{ID: id, Name: name}
		default:
			return nil, &WrongKindError{name}
		}
		calendar.SetUUID(uuid)
		calendars = append(calendars, calendar)

	}
	return
}

func findCalendarFromUser(db *sql.DB, user *User, calendarUUID string) (calendar api.CalendarManager, err error) {
	//TODO:
	//rows, err := db.Query("select calendars.id, calendars.name, calendars.uuid from calendars join accounts a on calendars.account_email = a.email where a.id=$1", account.GetInternalID())
	//defer rows.Close()
	//for rows.Next() {
	//	var id string
	//	var name string
	//	var uuid string
	//	var calendar api.CalendarManager
	//	err = rows.Scan(&id, &name, &uuid)
	//	switch account.GetKind() {
	//	case api.GOOGLE:
	//		calendar = &api.GoogleCalendar{ID: id, Name: name, UUID: uuid}
	//	case api.OUTLOOK:
	//		calendar = &api.OutlookCalendar{ID: id, Name: name, UUID: uuid}
	//	default:
	//		return nil, &WrongKindError{name}
	//	}
	//}
	return
}

func deleteCalendarFromUser(db *sql.DB, user *User, calendarUUID string) (err error) {
	stmt, err := db.Prepare("delete from calendars using accounts where calendars.account_email = accounts.email and accounts.user_uuid = $1 and calendars.uuid = $2")
	defer stmt.Close()
	if err != nil {
		log.Errorf("error preparing sql: %s", err.Error())
		return
	}
	res, err := stmt.Exec(user.UUID, calendarUUID)
	if err != nil {
		log.Errorf("error executing delete: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Errorf("error getting affected rows: %s", err.Error())
		return
	}
	log.Debugf("affected %d rows", affect)

	return

}
