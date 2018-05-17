package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

type Account struct {
	User         *User
	TokenType    string
	RefreshToken string
	Email        string
	Kind         int
	AccessToken  string
	ID           int
	Principal    bool
	Calendars    []Calendar
}

func (account Account) save(db *sql.DB) (err error) {
	stmt, err := db.Prepare("insert into accounts(user_uuid,token_type,refresh_token,email,kind,access_token, principal) values ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.User.UUID, account.TokenType, account.RefreshToken, account.Email, account.Kind, account.AccessToken, account.Principal)
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
		return errors.New(fmt.Sprintf("could not create new account for user %s with name: %s", account.User.Email, account.Email))
	}
	return

}
