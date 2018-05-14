package db

import (
	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func findCalendarsFromAccount(account api.AccountManager) (calendars []api.CalendarManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("could not connect to db: %s", err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("select calendars.* from calendars join accounts a on calendars.account_email = a.email where a.id=$1", account.GetInternalID())

	for rows.Next() {
		var email string
		var id int
		var calendar api.CalendarManager
		err = rows.Scan(&email, &id)
		switch account.GetKind() {
		case api.GOOGLE:
			calendar = &api.GoogleCalendar{}
		case api.OUTLOOK:
			calendar = &api.OutlookCalendar{}
		default:
			return nil, &WrongKindError{email}
		}
		calendars = append(calendars, calendar)
	}
	return
}
