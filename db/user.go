package db

import (
	"fmt"

	"database/sql"
	"errors"
	"net/http"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	USER     = "postgres"
	PASSWORD = "postgres"
	NAME     = "postgres"
)

// User has the info for the user
type User struct {
	ID           uuid.UUID
	Name         string
	Surname      string
	Email        string
	LoginAccount api.AccountManager
	Accounts     []api.AccountManager
	Calendars    []api.CalendarManager
}

type NotFoundError struct {
	Code int
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("")
}

func RetrieveUser(uuid string) (user *User, err error) {
	return findUserByID(uuid)
}

func GetUserFromToken(token string) (user *User, err error) {
	email := token
	user, err = findUserByMail(email)
	if _, ok := err.(*NotFoundError); ok {
		log.Debugf("no user found with email %s", email)
		user = &User{ID: uuid.New(), Name: "TESTING", Email: email, Surname: "asasd"}
		err = creteUser(user)
	}
	if err != nil {
		log.Errorf("error retrieving email: %s", email)
		return nil, err

	}
	log.Infof("user with email %s successfully retrieve from DB", user.Email)

	return
}
func findUserByID(id string) (user *User, err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("error on db: %s", err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * from users where users.uuid = $1;", id)
	if err != nil {
		log.Errorf("error querying: %s", err.Error())
	}
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
		user = &User{ID: uid, Name: name, Surname: surname, Email: email}
	} else {
		return nil, &NotFoundError{Code: http.StatusNotFound}
	}

	return

}

func findUserByMail(email string) (user *User, err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("error on db: %s", err.Error())
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * from users where users.email = $1;", email)
	if err != nil {
		log.Errorf("error querying: %s", err.Error())
		return
	}
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
		user = &User{ID: uid, Name: name, Surname: surname, Email: email}
	} else {
		return nil, &NotFoundError{Code: http.StatusNotFound}
	}

	return
}

func creteUser(user *User) (err error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("error on db: %s", err.Error())
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into users(uuid,email,name,surname) values ($1,$2,$3,$4);")
	if err != nil {
		log.Errorf("error preparing query: %s", err.Error())
		return
	}

	res, err := stmt.Exec(user.ID, user.Email, user.Name, user.Surname)
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
