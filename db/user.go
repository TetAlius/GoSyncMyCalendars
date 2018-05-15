package db

import (
	"fmt"

	"errors"
	"net/http"

	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// User has the info for the user
type User struct {
	UUID     uuid.UUID
	Name     string
	Surname  string
	Email    string
	Accounts []api.AccountManager
}

type NotFoundError struct {
	Code int
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("")
}

func (user *User) AddAccount(account api.AccountManager) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	err = saveAccount(db, account, *user)
	if err != nil {
		log.Errorf("could not save account: %s", err.Error())
		return
	}
	user.Accounts = append(user.Accounts, account)
	return
}

func (user *User) FindAccount(internalID int) (account api.AccountManager, err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	account, err = findAccountFromUser(db, user, internalID)
	err = account.Refresh()
	if err != nil {
		log.Errorf("error refreshing account: %s", account.Mail())
	}
	updateAccount(db, account)
	return
}

func (user *User) RetrieveCalendarsFromAccount(account api.AccountManager) (calendars []api.CalendarManager, err error) {
	//TODO: Change this
	return account.GetAllCalendars()
}

func (user *User) AddCalendarsToAccount(account api.AccountManager, ids []string) (err error) {
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	for _, id := range ids {
		calendar, err := account.GetCalendar(id)
		if err != nil {
			return err
		}
		err = addCalendarToAccount(db, account, calendar)
		if err != nil {
			log.Errorf("error adding calendar: %s", err.Error())
		} else {
			account.SetCalendars(append(account.GetSyncCalendars(), calendar))
		}
	}
	return
}

func (user *User) DeleteCalendar(uuid string) (err error) {
	//account, err := getAccountFromCalendarID(user, uuid)
	//if err != nil {
	//	log.Errorf("could not get account associated: %s", err.Error())
	//	return
	//}
	//calendar, err := findSubscriptionFromCalendar(user, uuid)
	//if err != nil {
	//	log.Errorf("could not get calendar: %s", err.Error())
	//	return
	//}
	//	TODO: subscrition delete
	db, err := connect()
	if err != nil {
		log.Errorf("db could not load: %s", err.Error())
	}
	defer db.Close()
	return deleteCalendarFromUser(db, user, uuid)
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
	err = user.setAccounts(db)
	if err != nil {
		log.Errorf("error retrieving accounts: %s", err.Error())
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
	if _, ok := err.(*NotFoundError); ok {
		log.Debugf("no user found with email %s", email)
		user = &User{UUID: uuid.New(), Name: "TESTING", Email: email, Surname: "asasd"}
		err = creteUser(db, user)
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
		return nil, &NotFoundError{Code: http.StatusNotFound}
	}

	return
}

func creteUser(db *sql.DB, user *User) (err error) {
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

func (user *User) setAccounts(db *sql.DB) (err error) {
	user.Accounts, err = getAccountsByUser(db, user.UUID)
	return
}
