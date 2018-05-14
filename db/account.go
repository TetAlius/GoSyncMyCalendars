package db

import (
	"errors"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

func saveAccount(account api.AccountManager, user User) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("could not connect to db: %s", err.Error())
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into accounts(user_uuid,token_type,refresh_token,email,kind,access_token) values ($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}

	res, err := stmt.Exec(user.UUID, account.GetTokenType(), account.GetRefreshToken(), account.Mail(), account.GetKind(), account.GetAccessToken())
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
		return errors.New(fmt.Sprintf("could not create new user with mail: %s", user.Email))
	}
	return
}

func getAccountsByUser(userUUID uuid.UUID) (accounts []api.AccountManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("could not connect to db: %s", err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT accounts.email,accounts.kind,accounts.id FROM accounts where user_uuid = $1", userUUID)

	for rows.Next() {
		var email string
		var kind int
		var id int
		var account api.AccountManager
		err = rows.Scan(&email, &kind, &id)
		switch kind {
		case api.GOOGLE:
			account = &api.GoogleAccount{Email: email, Kind: kind, InternID: id}
		case api.OUTLOOK:
			account = &api.OutlookAccount{AnchorMailbox: email, Kind: kind, InternID: id}
		default:
			return nil, &WrongKindError{email}
		}
		accounts = append(accounts, account)
	}
	return
}

func findAccountFromUser(user *User, internalID int) (account api.AccountManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("could not connect to db: %s", err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token FROM accounts where user_uuid = $1 and id = $2", user.UUID, internalID)

	if rows.Next() {
		var email string
		var kind int
		var id int
		var tokenType string
		var refreshToken string
		var accessToken string

		err = rows.Scan(&email, &kind, &id, &tokenType, &refreshToken, &accessToken)
		switch kind {
		case api.GOOGLE:
			account = &api.GoogleAccount{Email: email, Kind: kind, InternID: id, TokenType: tokenType, RefreshToken: refreshToken, AccessToken: accessToken}
		case api.OUTLOOK:
			account = &api.OutlookAccount{AnchorMailbox: email, Kind: kind, InternID: id, TokenType: tokenType, RefreshToken: refreshToken, AccessToken: accessToken}
		default:
			return nil, &WrongKindError{email}
		}
	}
	calendars, err := findCalendarsFromAccount(account)
	account.SetCalendars(calendars)
	return

}
