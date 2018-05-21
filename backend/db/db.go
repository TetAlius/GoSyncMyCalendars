package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type WrongKindError struct {
	Mail string
}

func (err *WrongKindError) Error() string {
	return fmt.Sprintf("wrong kind of account %s", err.Mail)
}
