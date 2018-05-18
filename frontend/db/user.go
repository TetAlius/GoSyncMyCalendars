package db

import (
	"database/sql"
	"encoding/base64"
	"net/http"

	"strings"

	"net/url"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

type User struct {
	UUID             uuid.UUID
	Email            string
	Name             string
	Surname          string
	PrincipalAccount Account
	Accounts         []Account
}

func RetrieveUser(uuid string) (user *User, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	user, err = findUserByID(db, uuid)
	if err != nil {
		log.Errorf("error retrieving user %s", uuid)
		return
	}

	return
}

func (user *User) AddAccount(account Account) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	var accounts int
	err = db.QueryRow("SELECT count(accounts.id) FROM accounts where accounts.user_uuid = $1", user.UUID).Scan(&accounts)
	if err != nil {
		log.Errorf("error getting number of accounts from user: %s", err.Error())
		return
	}
	principal := accounts < 1
	account.Principal = principal
	err = account.save(db)
	if err != nil {
		log.Errorf("could not save account: %s", err.Error())
		return
	}
	if principal {
		user.PrincipalAccount = account
	} else {
		user.Accounts = append(user.Accounts, account)
	}

	return
}
func (user *User) SetAccounts() (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	err = user.setAccounts(db)
	if err != nil {
		log.Errorf("error retrieving accounts: %s", err.Error())
	}
	return
}

func (user *User) setAccounts(db *sql.DB) (err error) {
	principal, accounts, err := getAccountsByUser(db, user.UUID)
	if err != nil {
		return err
	}
	user.PrincipalAccount = principal
	user.Accounts = accounts
	return
}
func (user *User) FindAccount(ID int) (account Account, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	account = Account{User: user, ID: ID}
	err = account.findAccount(db)

	if err != nil {
		log.Errorf("error retrieving account: %s for user: %s", ID, user.UUID)
		return
	}
	return

}
func (user *User) AddCalendarsToAccount(account Account, values []string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	for _, value := range values {
		decodedToken, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}
		query, err := url.QueryUnescape(string(decodedToken[:]))
		if err != nil {
			log.Errorf("error unscaping url: %s", query)
		}
		info := strings.Split(query, ":::")

		err = account.addCalendar(db, Calendar{ID: info[0], Name: info[1]})
		if err != nil {
			log.Errorf("error adding calendar: %s", err.Error())
		}
	}
	return

}
func (user *User) DeleteCalendar(id string) (err error) {
	//	TODO: subscrition delete
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	uid, err := uuid.Parse(id)
	calendar := Calendar{UUID: uid}
	return calendar.deleteFromUser(db, user)
}

func findUserByID(db *sql.DB, id string) (user *User, err error) {
	rows, err := db.Query("SELECT users.uuid, users.name,users.surname, users.email from users where users.uuid = $1;", id)
	if err != nil {
		log.Errorf("error querying: %s", err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		var uid uuid.UUID
		var name string
		var surname string
		var email string
		err = rows.Scan(&uid, &name, &surname, &email)
		if err != nil {
			log.Errorf("error on scan: %s", err.Error())
			return
		}
		user = &User{UUID: uid, Name: name, Surname: surname, Email: email}
	} else {
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	}

	return

}
