package db

import (
	"database/sql"
	"net/http"

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
		return nil, &NotFoundError{Code: http.StatusNotFound}
	}

	return

}
