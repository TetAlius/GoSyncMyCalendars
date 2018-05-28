package db

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (data Database) RetrieveAccount(userUUID string, ID string) (account api.AccountManager, err error) {
	account, err = data.findAccountFromUser(userUUID, ID)
	if err != nil {
		log.Errorf("error retrieving account: %s", ID)
		return nil, err
	}
	return
}

func (data Database) findAccountFromUser(userUUID string, internalID string) (account api.AccountManager, err error) {
	var email string
	var kind int
	var id int
	var tokenType string
	var refreshToken string
	var accessToken string
	err = data.DB.QueryRow("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token FROM accounts where user_uuid = $1 and id = $2", userUUID, internalID).Scan(&email, &kind, &id, &tokenType, &refreshToken, &accessToken)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("No account from user: %s with that id: %d.", userUUID, id)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	case err != nil:
		log.Debugf("error looking for account from user: %s with id: %d.", userUUID, id)
		return
	}
	switch kind {
	case api.GOOGLE:
		account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken)
	case api.OUTLOOK:
		account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken)
	default:
		return nil, &customErrors.WrongKindError{email}
	}
	return

}

func (data Database) UpdateAccountFromSubscription(account api.AccountManager, subscription api.SubscriptionManager) (err error) {
	stmt, err := data.DB.Prepare("update accounts set token_type = $1, refresh_token = $2, access_token = $3 from subscriptions, calendars where subscriptions.uuid = $4 and subscriptions.calendar_uuid = calendars.uuid and calendars.account_email = accounts.email and accounts.email=$5")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), subscription.GetUUID(), account.Mail())
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

func (data Database) UpdateAccountFromUser(account api.AccountManager, userUUID string) (err error) {
	stmt, err := data.DB.Prepare("update accounts set (token_type,refresh_token,access_token) = ($1,$2,$3) where accounts.email = $4 and accounts.user_uuid =$5;")
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
