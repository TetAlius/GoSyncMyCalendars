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

func (data Database) RetrieveUser(uuid string) (user *User, err error) {
	user, err = data.findUserByID(uuid)
	if err != nil {
		log.Errorf("error retrieving user %s", uuid)
		return
	}

	return
}

func (data Database) GetUserFromToken(token string) (user *User, err error) {
	email := token
	user, err = data.findUserByMail(email)
	if _, ok := err.(*customErrors.NotFoundError); ok {
		log.Debugf("no user found with email %s", email)
		user = &User{UUID: uuid.New(), Name: "Name", Email: email, Surname: "Surname"}
		err = data.createUser(user)
	}
	if err != nil {
		log.Errorf("error retrieving email: %s", email)
		return nil, err

	}
	log.Infof("user with email %s successfully retrieve from DB", user.Email)
	err = data.SetUserAccounts(user)
	if err != nil {
		log.Errorf("error retrieving accounts: %s", err.Error())
	}

	return

}

func (data Database) createUser(user *User) (err error) {
	stmt, err := data.DB.Prepare("insert into users(uuid,email,name,surname) values ($1,$2,$3,$4);")
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

func (data Database) findUserByMail(email string) (user *User, err error) {
	var uid uuid.UUID
	var name string
	var surname string
	var mail string
	err = data.DB.QueryRow("SELECT users.uuid,users.name,users.surname,users.email from users where users.email = $1;", email).Scan(&uid, &name, &surname, &mail)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("No user with that email: %s.", email)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	case err != nil:
		log.Errorf("error looking for user with email: %s", email)
		return
	}
	user = &User{UUID: uid, Name: name, Surname: surname, Email: mail}
	return
}

func (data Database) AddAccount(user *User, account Account) (err error) {
	var accounts int
	err = data.DB.QueryRow("SELECT count(accounts.id) FROM accounts where accounts.user_uuid = $1", user.UUID).Scan(&accounts)
	if err != nil {
		log.Errorf("error getting number of accounts from user: %s", err.Error())
		return
	}
	principal := accounts < 1
	account.Principal = principal
	err = data.save(account)
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

func (data Database) SetUserAccounts(user *User) (err error) {
	principal, accounts, err := data.getAccountsByUser(user.UUID)
	if err != nil {
		return err
	}
	user.PrincipalAccount = principal
	user.Accounts = accounts
	return

}

func (data Database) FindAccount(user *User, ID int) (account Account, err error) {
	account = Account{User: user, ID: ID}
	err = data.findAccount(&account)

	if err != nil {
		log.Errorf("error retrieving account: %s for user: %s", ID, user.UUID)
		return
	}
	err = data.findCalendars(&account)
	return

}

func (data Database) AddCalendarsToAccount(user *User, account Account, values []string) (err error) {
	for _, value := range values {
		decodedToken, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}
		query, err := url.QueryUnescape(string(decodedToken[:]))
		if err != nil {
			log.Errorf("error unescaped url: %s", query)
		}
		info := strings.Split(query, ":::")

		data.addCalendar(account, Calendar{ID: info[0], Name: info[1]})
		if err != nil {
			log.Errorf("error adding calendar: %s", err.Error())
		}
	}
	return

}

func (data Database) DeleteCalendar(user *User, id string) (err error) {
	uid, err := uuid.Parse(id)
	calendar := Calendar{UUID: uid}
	return data.deleteFromUser(calendar, user)
}

func (data Database) AddCalendarsRelation(user *User, parentCalendarUUID string, calendarIDs []string) (err error) {
	stmt, err := data.DB.Prepare("update calendars set (parent_calendar_uuid) = ($1) from accounts as a where calendars.uuid = $2 and calendars.account_email = a.email and a.user_uuid = $3")
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

func (data Database) UpdateCalendar(user *User, calendarID string, parentID string) (err error) {
	return data.updateCalendarFromUser(user, calendarID, parentID)

}

func (data Database) findUserByID(id string) (user *User, err error) {
	var uid uuid.UUID
	var name string
	var surname string
	var email string
	err = data.DB.QueryRow("SELECT users.uuid, users.name,users.surname, users.email from users where users.uuid = $1;", id).Scan(&uid, &name, &surname, &email)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("No user with that id: %s.", id)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	case err != nil:
		log.Errorf("error looking for user with id: %s", id)
		return
	}
	user = &User{UUID: uid, Name: name, Surname: surname, Email: email}
	return

}
