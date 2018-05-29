package db

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (data Database) savePrincipalEvents(transaction *sql.Tx, events []api.EventManager) (err error) {
	for _, event := range events {
		lastInsertId := 0
		err = transaction.QueryRow("INSERT INTO events (calendar_uuid, id) VALUES($1, $2) RETURNING internal_id", event.GetCalendar().GetUUID(), event.GetID()).Scan(&lastInsertId)
		switch {
		case err == sql.ErrNoRows:
			err = fmt.Errorf("could not insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Debugf("could not insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
			return err
		case err != nil:
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Debugf("error insert event with id: %s and calendar UUID: %s", event.GetID(), event.GetCalendar().GetUUID())
			return
		}
		event.SetInternalID(lastInsertId)
	}
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
	}
	return

}
func (data Database) saveEventsRelation(transaction *sql.Tx, from api.EventManager, to api.EventManager) (err error) {
	stmt, err := transaction.Prepare("insert into events(calendar_uuid, id, parent_event_internal_id) values ($1,$2,$3)")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(to.GetCalendar().GetUUID(), to.GetID(), from.GetInternalID())
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
		err = errors.New(fmt.Sprintf("could not update account with id: %s", to.GetID()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return
	}
	return
}
