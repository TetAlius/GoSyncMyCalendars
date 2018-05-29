package db

import (
	"errors"
	"fmt"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
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

func (data Database) getAccountsByUser(userUUID uuid.UUID) (principalAccount Account, accounts []Account, err error) {
	rows, err := data.client.Query("SELECT accounts.token_type, accounts.refresh_token,accounts.email,accounts.kind,accounts.access_token,accounts.id, accounts.principal FROM accounts where user_uuid = $1 order by accounts.principal DESC, lower(accounts.email) ASC", userUUID)
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("could not query select: %s", err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		var email string
		var kind int
		var id int
		var tokenType string
		var refreshToken string
		var accessToken string
		var principal bool
		var account Account
		err = rows.Scan(&tokenType, &refreshToken, &email, &kind, &accessToken, &id, &principal)
		if err != nil {
			data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
			log.Errorf("error retrieving accounts for user: %s", userUUID)
			return account, nil, err
		}
		account = Account{
			TokenType:    tokenType,
			RefreshToken: refreshToken,
			Email:        email,
			AccessToken:  accessToken,
			Kind:         kind,
			ID:           id,
			Principal:    principal,
		}
		err = data.findCalendars(&account)
		if err != nil {
			log.Errorf("error adding calendars: %s", account.ID)
			return
		}
		if principal {
			principalAccount = account
		} else {
			accounts = append(accounts, account)
		}
	}
	return
}

func (data Database) FindCalendars(account *Account) (err error) {
	return data.findCalendars(account)

}

func (data Database) save(account Account) (err error) {
	stmt, err := data.client.Prepare("insert into accounts(user_uuid,token_type,refresh_token,email,kind,access_token, principal) values ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.User.UUID, account.TokenType, account.RefreshToken, account.Email, account.Kind, account.AccessToken, account.Principal)
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureError(errors.New(fmt.Sprintf("could not create new account for user %s with name: %s", account.User.Email, account.Email)), map[string]string{"database": "frontend"})
		return errors.New(fmt.Sprintf("could not create new account for user %s with name: %s", account.User.Email, account.Email))
	}
	return

}

func (data Database) findAccount(account *Account) (err error) {
	var email string
	var kind int
	var id int
	var principal bool
	err = data.client.QueryRow("SELECT accounts.email,accounts.kind,accounts.id, accounts.principal FROM accounts where user_uuid = $1 and id = $2", account.User.UUID, account.ID).Scan(&email, &kind, &id, &principal)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("No account with that id: %d.", id)
		return &customErrors.NotFoundError{Message: fmt.Sprintf("No account with that id: %d.", id)}
	case err != nil:
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error looking for account with id: %d", id)
		return
	}
	account.Email = email
	account.Kind = kind
	account.ID = id
	account.Principal = principal
	return

}

func (data Database) addCalendar(account Account, calendar Calendar) (err error) {
	stmt, err := data.client.Prepare("insert into calendars(uuid, account_email, name, id) values ($1,$2,$3,$4);")
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	calendar.UUID = uuid.New()
	res, err := stmt.Exec(calendar.UUID, account.Email, calendar.Name, calendar.ID)
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureError(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureError(errors.New(fmt.Sprintf("could not create new calendar with id: %s and name: %s", calendar.ID, calendar.Name)), map[string]string{"database": "frontend"})
		return errors.New(fmt.Sprintf("could not create new calendar with id: %s and name: %s", calendar.ID, calendar.Name))
	}

	return

}
