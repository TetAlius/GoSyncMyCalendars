package db

import (
	"database/sql"
	"fmt"

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

func connect() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	return sql.Open("postgres", dbinfo)
}
