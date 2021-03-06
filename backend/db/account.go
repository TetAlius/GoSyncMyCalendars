package db

import (
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// Method that retrieves an AccountManager from the DB
func (data Database) RetrieveAccount(userUUID string, ID string) (account api.AccountManager, err error) {
	account, err = data.findAccountFromUser(userUUID, ID)
	if err != nil {
		log.Errorf("error retrieving account: %s", ID)
		return nil, err
	}
	return
}

// Returns account from user given the internal id
func (data Database) findAccountFromUser(userUUID string, internalID string) (account api.AccountManager, err error) {
	var email string
	var kind int
	var id int
	var tokenType string
	var refreshToken string
	var accessToken string
	err = data.client.QueryRow("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token FROM accounts where user_uuid = $1 and id = $2", userUUID, internalID).Scan(&email, &kind, &id, &tokenType, &refreshToken, &accessToken)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("No account from user: %s with that id: %d.", userUUID, id)}
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("No account from user: %s with that id: %d.", userUUID, id)
		return nil, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("error looking for account from user: %s with id: %d.", userUUID, id)
		return
	}
	switch kind {
	case api.GOOGLE:
		account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken)
	case api.OUTLOOK:
		account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken)
	default:
		data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: email}, map[string]string{"database": "backend"})
		return nil, &customErrors.WrongKindError{Mail: email}
	}
	return

}

// Method that updates the account info of a subscription
func (data Database) UpdateAccountFromSubscription(account api.AccountManager, subscription api.SubscriptionManager) (err error) {
	stmt, err := data.client.Prepare("update accounts set token_type = $1, refresh_token = $2, access_token = $3 from subscriptions, calendars where subscriptions.uuid = $4 and subscriptions.calendar_uuid = calendars.uuid and calendars.account_email = accounts.email and accounts.email=$5")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), subscription.GetUUID(), account.Mail())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail())), map[string]string{"database": "backend"})
		return errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail()))
	}
	return

}

// Method that updates the account info of a user
func (data Database) UpdateAccountFromUser(account api.AccountManager, userUUID string) (err error) {
	stmt, err := data.client.Prepare("update accounts set (token_type,refresh_token,access_token) = ($1,$2,$3) where accounts.email = $4 and accounts.user_uuid =$5;")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), account.Mail(), userUUID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail())), map[string]string{"database": "backend"})
		return errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail()))
	}
	return

}

// Method that updates the account info
func (data Database) UpdateAccount(account api.AccountManager) {
	stmt, err := data.client.Prepare("update accounts set (token_type,refresh_token,access_token) = ($1,$2,$3) where accounts.email = $4;")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), account.Mail())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		err = errors.New(fmt.Sprintf("could not update account with mail: %s", account.Mail()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return
	}
	return

}
