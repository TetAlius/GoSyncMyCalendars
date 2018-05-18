package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func RetrieveAccount(userUUID string, ID string) (account api.AccountManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	account, err = findAccountFromUser(db, userUUID, ID)
	if err != nil {
		log.Errorf("error retrieving account: %s", ID)
		return nil, err
	}
	return
}

func findAccountFromUser(db *sql.DB, userUUID string, internalID string) (account api.AccountManager, err error) {
	rows, err := db.Query("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token FROM accounts where user_uuid = $1 and id = $2", userUUID, internalID)
	if err != nil {
		log.Errorf("error executing query: %s", err.Error())
		return nil, err
	}

	defer rows.Close()
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
			account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken, id, false)
		case api.OUTLOOK:
			account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken, id, false)
		default:
			return nil, &WrongKindError{email}
		}
	} else {
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	}
	return

}

func UpdateAccountFromUser(account api.AccountManager, userUUID string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("update accounts set (token_type,refresh_token,access_token) = ($1,$2,$3) where accounts.email = $4 and accounts.user_uuid =$5;")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), account.Mail(), userUUID)
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
		return errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail()))
	}
	return

}
