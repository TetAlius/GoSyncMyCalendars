package db

import (
	"errors"
	"fmt"

	"net/http"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

func saveAccount(db *sql.DB, account api.AccountManager, user User) (err error) {
	principal := len(user.Accounts) == 0
	stmt, err := db.Prepare("insert into accounts(user_uuid,token_type,refresh_token,email,kind,access_token, principal) values ($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.UUID, account.GetTokenType(), account.GetRefreshToken(), account.Mail(), account.GetKind(), account.GetAccessToken(), principal)
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

func updateAccount(db *sql.DB, account api.AccountManager) (err error) {
	stmt, err := db.Prepare("update accounts set (token_type,refresh_token,access_token) = ($1,$2,$3) where accounts.email = $4;")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(account.GetTokenType(), account.GetRefreshToken(), account.GetAccessToken(), account.Mail())
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

func getAccountsByUser(db *sql.DB, userUUID uuid.UUID) (accounts []api.AccountManager, err error) {
	rows, err := db.Query("SELECT accounts.token_type, accounts.refresh_token,accounts.email,accounts.kind,accounts.access_token,accounts.id, accounts.principal FROM accounts where user_uuid = $1 order by accounts.principal DESC, accounts.email ASC", userUUID)
	if err != nil {
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
		var account api.AccountManager
		var principal bool
		err = rows.Scan(&tokenType, &refreshToken, &email, &kind, &accessToken, &id, &principal)
		switch kind {
		case api.GOOGLE:
			account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		case api.OUTLOOK:
			account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		default:
			return nil, &WrongKindError{email}
		}
		accounts = append(accounts, account)
		calendars, _ := findCalendarsFromAccount(db, account)
		if len(calendars) != 0 {
			account.SetCalendars(calendars)
		}
	}
	return
}

func getPrincipalAccountByUser(db *sql.DB, userUUID uuid.UUID) (principalAccount api.AccountManager, err error) {
	rows, err := db.Query("SELECT accounts.token_type, accounts.refresh_token,accounts.email,accounts.kind,accounts.access_token,accounts.id, accounts.principal FROM accounts where user_uuid = $1 and accounts.principal =true order by accounts.email ASC", userUUID)
	if err != nil {
		return
	}

	defer rows.Close()
	if rows.Next() {
		var email string
		var kind int
		var id int
		var tokenType string
		var refreshToken string
		var accessToken string
		var principalAccount api.AccountManager
		var principal bool
		err = rows.Scan(&tokenType, &refreshToken, &email, &kind, &accessToken, &id, &principal)
		switch kind {
		case api.GOOGLE:
			principalAccount = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		case api.OUTLOOK:
			principalAccount = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		default:
			return nil, &WrongKindError{email}
		}
		calendars, _ := findCalendarsFromAccount(db, principalAccount)
		if len(calendars) != 0 {
			principalAccount.SetCalendars(calendars)
		}
	}
	return

}
func findAccountFromUser(db *sql.DB, user *User, internalID int) (account api.AccountManager, err error) {
	rows, err := db.Query("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token, accounts.principal FROM accounts where user_uuid = $1 and id = $2", user.UUID, internalID)

	defer rows.Close()
	if rows.Next() {
		var email string
		var kind int
		var id int
		var tokenType string
		var refreshToken string
		var accessToken string
		var principal bool
		err = rows.Scan(&email, &kind, &id, &tokenType, &refreshToken, &accessToken, &principal)
		switch kind {
		case api.GOOGLE:
			account = api.RetrieveGoogleAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		case api.OUTLOOK:
			account = api.RetrieveOutlookAccount(tokenType, refreshToken, email, kind, accessToken, id, principal)
		default:
			return nil, &WrongKindError{email}
		}
	} else {
		return nil, &NotFoundError{http.StatusNotFound}
	}
	calendars, err := findCalendarsFromAccount(db, account)
	if len(calendars) != 0 {
		account.SetCalendars(calendars)
	}
	return

}

func addCalendarToAccount(db *sql.DB, account api.AccountManager, calendar api.CalendarManager) (err error) {
	stmt, err := db.Prepare("insert into calendars(uuid, account_email, name, id) values ($1,$2,$3,$4);")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	uid := uuid.New().String()
	res, err := stmt.Exec(uid, account.Mail(), calendar.GetName(), calendar.GetID())
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
		return errors.New(fmt.Sprintf("could not create new calendar with id: %s and name: %s", calendar.GetID(), calendar.GetName()))
	}
	calendar.SetUUID(uid)

	return
}

func getAccountFromCalendarID(db *sql.DB, user *User, calendarUUID string) (account api.AccountManager, err error) {
	rows, err := db.Query("SELECT accounts.email,accounts.kind,accounts.id, accounts.token_type,accounts.refresh_token,accounts.access_token FROM accounts join calendars c2 on accounts.email = c2.account_email where c2.uuid = $1 and user_uuid = $2", calendarUUID, user.UUID)

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
			account = &api.GoogleAccount{Email: email, Kind: kind, InternID: id, TokenType: tokenType, RefreshToken: refreshToken, AccessToken: accessToken}
		case api.OUTLOOK:
			account = &api.OutlookAccount{AnchorMailbox: email, Kind: kind, InternID: id, TokenType: tokenType, RefreshToken: refreshToken, AccessToken: accessToken}
		default:
			return nil, &WrongKindError{email}
		}
	} else {
		return nil, &NotFoundError{http.StatusNotFound}
	}
	return

}
