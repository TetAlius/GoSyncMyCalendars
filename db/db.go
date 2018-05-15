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

type WrongKindError struct {
	Mail string
}

func (err *WrongKindError) Error() string {
	return fmt.Sprintf("wrong kind of account %s", err.Mail)
}

func connect() (db *sql.DB, err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return
	}
	log.Infof("DB: opening connection")
	return
}
