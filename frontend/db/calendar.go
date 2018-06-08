package db

import (
	"errors"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

type Calendar struct {
	UUID             uuid.UUID
	AccountEmail     string
	Account          Account
	Name             string
	ID               string
	ParentUUID       uuid.UUID
	Events           []Event
	SubscriptionUUID uuid.UUID
	Calendars        []Calendar
}

func newCalendar(id string, name string, uid uuid.UUID, accountEmail string, account Account, subscriptionUUID uuid.UUID) Calendar {
	return Calendar{
		ID:               id,
		Name:             name,
		UUID:             uid,
		AccountEmail:     accountEmail,
		Account:          account,
		SubscriptionUUID: subscriptionUUID,
	}

}

func (data Database) findCalendars(account *Account) (err error) {
	rows, err := data.client.Query("select calendars.id, calendars.name, calendars.uuid, s2.uuid from calendars join accounts a on calendars.account_email = a.email left outer join subscriptions s2 on calendars.uuid = s2.calendar_uuid where a.id=$1 order by calendars.name ASC", account.ID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorln("error selecting findCalendarsFromAccount")
		return
	}

	defer rows.Close()
	var calendars []Calendar
	for rows.Next() {
		var id string
		var name string
		var uid uuid.UUID
		var subscription uuid.UUID
		err = rows.Scan(&id, &name, &uid, &subscription)
		if err != nil {
			//TODO
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			continue
		}
		calendar := newCalendar(id, name, uid, account.Email, *account, subscription)

		data.setSynchronizedCalendars(&calendar, account.Principal)
		calendars = append(calendars, calendar)
	}
	account.Calendars = calendars
	return

}

func (data Database) deleteFromUser(calendar Calendar, user *User) (err error) {
	stmt, err := data.client.Prepare("delete from calendars using accounts where calendars.account_email = accounts.email and accounts.user_uuid = $1 and calendars.uuid = $2")
	defer stmt.Close()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error preparing sql: %s", err.Error())
		return
	}
	res, err := stmt.Exec(user.UUID, calendar.UUID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error executing delete: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error getting affected rows: %s", err.Error())
		return
	}
	log.Debugf("affected %d rows", affect)

	return

}

func (data Database) updateCalendarFromUser(user *User, calendarUUID string, parentUUID string) (err error) {
	stmt, err := data.client.Prepare("update calendars set parent_calendar_uuid = $1 from accounts where calendars.account_email = accounts.email and accounts.user_uuid = $2 and calendars.uuid = $3;")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	var parent interface{}
	parent = parentUUID
	if len(parentUUID) == 0 {
		parent = nil
	}

	res, err := stmt.Exec(parent, user.UUID, calendarUUID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not update calendar with UUID: %s", calendarUUID)), map[string]string{"database": "frontend"})
		return errors.New(fmt.Sprintf("could not update calendar with UUID: %s", calendarUUID))
	}
	return
}

func (data Database) setSynchronizedCalendars(calendar *Calendar, principal bool) (err error) {
	var query string
	if principal {
		query = "select calendars.id, calendars.name, calendars.uuid, a.kind, a.email, s2.uuid from calendars join accounts a on calendars.account_email = a.email left outer join subscriptions s2 on calendars.uuid = s2.calendar_uuid where calendars.parent_calendar_uuid = $1"
	} else {
		query = "select calendars.id, calendars.name, calendars.uuid, a.kind, a.email, s2.uuid from calendars join accounts a on calendars.account_email = a.email left outer join subscriptions s2 on calendars.uuid = s2.calendar_uuid where calendars.parent_calendar_uuid = (Select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1) OR calendars.uuid = (select calendars.parent_calendar_uuid from calendars where calendars.uuid = $1)"
	}
	rows, err := data.client.Query(query, calendar.UUID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorln("error selecting setSynchronizedCalendars")
		return
	}
	var calendars []Calendar
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var uid uuid.UUID
		var cal Calendar
		var kind int
		var accountEmail string
		var subscriptionUUID uuid.UUID
		err = rows.Scan(&id, &name, &uid, &kind, &accountEmail, &subscriptionUUID)
		if err != nil {
			//TODO
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			continue
		}

		cal = newCalendar(id, name, uid, accountEmail, Account{Email: accountEmail}, subscriptionUUID)
		calendars = append(calendars, cal)
	}
	calendar.Calendars = calendars
	return
}
