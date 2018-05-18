package db

import (
	"database/sql"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

type Calendar struct {
	UUID         uuid.UUID
	Account      Account
	Name         string
	ID           string
	ParentUUID   uuid.UUID
	Events       []Event
	Subscription Subscription
}

func (calendar Calendar) deleteFromUser(db *sql.DB, user *User) (err error) {
	stmt, err := db.Prepare("delete from calendars using accounts where calendars.account_email = accounts.email and accounts.user_uuid = $1 and calendars.uuid = $2")
	defer stmt.Close()
	if err != nil {
		log.Errorf("error preparing sql: %s", err.Error())
		return
	}
	res, err := stmt.Exec(user.UUID, calendar.UUID)
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
