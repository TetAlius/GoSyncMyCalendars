package db

import (
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

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
	return
}
func findUserByID(id string) (user *User, err error) {
	return
}
func findUserByMail(email string) (user *User, err error) {
	return
}

func creteUser(user *User) (err error) {
	return
}
