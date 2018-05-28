package db

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (data Database) savePrincipalEvents(transaction *sql.Tx, events []api.EventManager) (err error) {
	stmt, err := transaction.Prepare("insert into events(calendar_uuid, id) values ($1,$2)")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	for _, event := range events {
		res, err := stmt.Exec(event.GetCalendar().GetUUID(), event.GetID())
		if err != nil {
			log.Errorf("error executing query: %s", err.Error())
			break
		}

		affect, err := res.RowsAffected()
		if err != nil {
			log.Errorf("error retrieving rows affected: %s", err.Error())
			break
		}
		if affect != 1 {
			err = errors.New(fmt.Sprintf("could not update account with id: %s", event.GetID()))
			break
		}
	}
	return

}
func (data Database) saveEventsRelation(transaction *sql.Tx, from api.EventManager, to api.EventManager) (err error) {
	var internalID int
	err = transaction.QueryRow("select events.internal_id from events where events.id = $1 and events.calendar_uuid = $2", from.GetID(), from.GetCalendar().GetUUID()).Scan(&internalID)
	if err != nil {
		log.Errorf("could not retrieve internal id from event: %s", from.GetID())
		return
	}
	stmt, err := transaction.Prepare("insert into events(calendar_uuid, id, parent_event_internal_id) values ($1,$2,$3)")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(to.GetCalendar().GetUUID(), to.GetID(), internalID)
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
		err = errors.New(fmt.Sprintf("could not update account with id: %s", to.GetID()))
		return
	}
	return
}
