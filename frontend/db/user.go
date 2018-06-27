package db

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"net/url"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

// User mapped from db
type User struct {
	// UUID of the user
	UUID uuid.UUID
	// Email of the user
	Email string
	// Name of the user
	Name string
	// Principal Account of the user
	PrincipalAccount Account
	// Accounts associated with user that are not principal
	Accounts []Account
}

// Method that retrieves the user from the db
func (data Database) RetrieveUser(uuid string) (user *User, err error) {
	user, err = data.findUserByID(uuid)
	if err != nil {
		log.Errorf("error retrieving user %s", uuid)
		return
	}

	return
}

// Method that finds an user or creates one if it is not stored
func (data Database) FindOrCreateUser(user *User) (err error) {
	err = data.findUserByMail(user)
	if _, ok := err.(*customErrors.NotFoundError); ok {
		log.Debugf("no user found with email %s", user.Email)
		err = data.createUser(user)
	}
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving email: %s", user.Email)
		return

	}
	log.Infof("user with email %s successfully retrieve from DB", user.Email)
	err = data.SetUserAccounts(user)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving accounts: %s", err.Error())
	}

	return

}

// Method that creates a user on the db
func (data Database) createUser(user *User) (err error) {
	user.UUID = uuid.New()
	stmt, err := data.client.Prepare("insert into users(uuid,email,name) values ($1,$2,$3);")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.UUID, user.Email, user.Name)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not create new user with mail: %s", user.Email)), map[string]string{"database": "frontend"})
		return errors.New(fmt.Sprintf("could not create new user with mail: %s", user.Email))
	}
	return

}

// Method that retrieves a user by its email
func (data Database) findUserByMail(user *User) (err error) {
	var uid uuid.UUID
	var name string
	var mail string
	err = data.client.QueryRow("SELECT users.uuid,users.name,users.email from users where users.email = $1;", user.Email).Scan(&uid, &name, &mail)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("No user with that email: %s.", user.Email)}
		//data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Debugf("No user with that email: %s.", user.Email)
		return err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error looking for user with email: %s", user.Email)
		return
	}
	*user = User{UUID: uid, Name: name, Email: mail}
	return
}

// Method that saves an account to the user on db
func (data Database) AddAccount(user *User, account Account) (id int, err error) {
	var accounts int
	err = data.client.QueryRow("SELECT count(accounts.id) FROM accounts where accounts.user_uuid = $1", user.UUID).Scan(&accounts)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error getting number of accounts from user: %s", err.Error())
		return
	}
	principal := accounts < 1
	account.Principal = principal
	id, err = data.save(account)
	if _, ok := err.(*customErrors.AccountAlreadyUsed); ok {
		//TODO:
		id, err = data.updateAccountFromUser(account, user)
	}
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

// Method that retrieves all accounts of the user
func (data Database) SetUserAccounts(user *User) (err error) {
	principal, accounts, err := data.getAccountsByUser(user.UUID)
	if err != nil {
		return err
	}
	user.PrincipalAccount = principal
	user.Accounts = accounts
	return

}

// Method that finds an account by its id that belongs to a user
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

// Method that saves calendars to an account
func (data Database) AddCalendarsToAccount(user *User, account Account, values []string) (err error) {
	for _, value := range values {
		decodedToken, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			return err
		}
		query, err := url.QueryUnescape(string(decodedToken[:]))
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			log.Errorf("error unescaped url: %s", query)
			return err
		}
		info := strings.Split(query, ":::")

		err = data.addCalendar(account, Calendar{ID: info[0], Name: info[1]})
		if err != nil {
			log.Errorf("error adding calendar: %s", err.Error())
		}
	}
	return

}

// Method that removes a calendar from the db
func (data Database) DeleteCalendar(user *User, id string) (err error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		return
	}
	calendar := Calendar{UUID: uid}
	return data.deleteFromUser(calendar, user)
}

// Method that creates the relations between calendars
func (data Database) AddCalendarsRelation(user *User, parentCalendarUUID string, calendarIDs []string) (err error) {
	stmt, err := data.client.Prepare("update calendars set (parent_calendar_uuid) = ($1) from accounts as a where calendars.uuid = $2 and calendars.account_email = a.email and a.user_uuid = $3")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
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
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			log.Errorf("error executing query: %s", err.Error())
			return err
		}

		affect, err := res.RowsAffected()
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
			log.Errorf("error retrieving rows affected: %s", err.Error())
			return err
		}
		if affect != 1 {
			data.sentry.CaptureErrorAndWait(errors.New(fmt.Sprintf("could not create relations with parent: %s and child: %s", parentCalendarUUID, childID)), map[string]string{"database": "frontend"})
			return errors.New(fmt.Sprintf("could not create relations with parent: %s and child: %s", parentCalendarUUID, childID))
		}

	}
	return

}

// Method that updates a calendar
func (data Database) UpdateCalendar(user *User, calendarID string, parentID string) (err error) {
	return data.updateCalendarFromUser(user, calendarID, parentID)

}

// Method that looks for a user by its ID
func (data Database) findUserByID(id string) (user *User, err error) {
	var uid uuid.UUID
	var name string
	var email string
	err = data.client.QueryRow("SELECT users.uuid, users.name, users.email from users where users.uuid = $1;", id).Scan(&uid, &name, &email)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("No user with that id: %s.", id)}
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Debugf("No user with that id: %s.", id)
		return nil, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "frontend"})
		log.Errorf("error looking for user with id: %s", id)
		return
	}
	user = &User{UUID: uid, Name: name, Email: email}
	return

}
