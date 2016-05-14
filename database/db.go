package database

import (
	"database/sql"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//SaveUser saves the user data onto the DB
func SaveUser(db *sql.DB) (err error) {
	row, err := db.Query(`INSERT INTO users(mail, diggested_password) VALUES('marta@marta.com', 'password')`)
	log.Debugln(row)
	return
}

//CreateEvent TODO
func CreateEvent(db *sql.DB, idCal string, id string, mail string) {

}

//UpdateEvent TODO
func UpdateEvent(db *sql.DB, id string, info interface{}) {

}

//GetEvent TODO
func GetEvent(db *sql.DB, id string) {

}

//DeleteEvent TODO
func DeleteEvent(db *sql.DB, id string) {

}
