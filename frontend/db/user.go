package db

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
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

func GetUserFromToken(token string) (user *User, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	email := token
	user, err = findUserByMail(db, email)
	if _, ok := err.(*customErrors.NotFoundError); ok {
		log.Debugf("no user found with email %s", email)
		user = &User{UUID: uuid.New(), Name: "TESTING", Email: email, Surname: "asasd"}
		err = user.creteUser(db)
	}
	if err != nil {
		log.Errorf("error retrieving email: %s", email)
		return nil, err

	}
	log.Infof("user with email %s successfully retrieve from DB", user.Email)
	err = user.setAccounts(db)
	if err != nil {
		log.Errorf("error retrieving accounts: %s", err.Error())
	}

	return
}

func (user *User) creteUser(db *sql.DB) (err error) {
	stmt, err := db.Prepare("insert into users(uuid,email,name,surname) values ($1,$2,$3,$4);")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.UUID, user.Email, user.Name, user.Surname)
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

func findUserByMail(db *sql.DB, email string) (user *User, err error) {
	rows, err := db.Query("SELECT users.uuid,users.name,users.surname,users.email from users where users.email = $1;", email)
	if err != nil {
		log.Errorf("error querying: %s", err.Error())
		return
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

func (user *User) AddCalendarsRelation(parentCalendarUUID string, calendarIDs []string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("update calendars set (parent_calendar_uuid) = ($1) from accounts as a where calendars.uuid = $2 and calendars.account_email = a.email and a.user_uuid = $3")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return err
	}
	defer stmt.Close()
	for _, childID := range calendarIDs {
		_, err := uuid.Parse(childID)
		if err != nil {
			continue
		}
		res, err := stmt.Exec(parentCalendarUUID, childID, user.UUID)
		if err != nil {
			log.Errorf("error executing query: %s", err.Error())
			return err
		}

		affect, err := res.RowsAffected()
		if err != nil {
			log.Errorf("error retrieving rows affected: %s", err.Error())
			return err
		}
		if affect != 1 {
			return errors.New(fmt.Sprintf("could not create relations with parent: %s and child: %s", parentCalendarUUID, childID))
		}

	}
	return

}
func (user *User) UpdateCalendar(calendarID string, parentID string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	return updateCalendarFromUser(db, user, calendarID, parentID)

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
