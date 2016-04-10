package calendar

import (
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"github.com/TetAlius/GoSyncMyCalendars/database"
)

func createErrorForAccount(account string) (err error) {
	return customErrors.AccError{Account: account}
}

// GetSentenceToCalendarInsert TODO
func getSentenceToCalendarInsert(account string) (sentence string, err error) {
	switch account {
	case database.OUTLOOK:
		return `INSERT INTO outlook_calendars() VALUES()`, err
	case database.GOOGLE:
		return `INSERT INTO google_calendars(id) VALUES($1)`, err
	}
	return "", createErrorForAccount(account)
}

// GetSentenceToCalendarUpdate TODO
func getSentenceToCalendarUpdate(account string) (sentence string, err error) {
	switch account {
	case database.OUTLOOK:
		return `UPDATE outlook_calendars SET kind = $1 WHERE kind = $2`, err
	case database.GOOGLE:
		return `UPDATE google_calendars SET kind = $1 WHERE kind = $2`, err
	}
	return "", createErrorForAccount(account)
}

// GetSentenceToCalendarRead TODO
func getSentenceToCalendarRead(account string) (sentence string, err error) {
	switch account {
	case database.OUTLOOK:
		return `SELECT * FROM outlook_calendars WHERE id = $1`, err
	case database.GOOGLE:
		return `SELECT * FROM google_calendars WHERE id = $1`, err
	}
	return "", createErrorForAccount(account)
}
