package db

import (
	"database/sql"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	_ "github.com/lib/pq"
)

const (
	USER     = "postgres"
	PASSWORD = "postgres"
	NAME     = "postgres"
)

func connect() (db *sql.DB, err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
