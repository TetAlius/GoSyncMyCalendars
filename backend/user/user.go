package user

import (
	"crypto/rand"
	"io"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"golang.org/x/crypto/scrypt"
)

const (
	//PwSaltBytes number of bytes for salt
	PwSaltBytes = 32
	//PwHashBytes number of bytes for hash
	PwHashBytes = 64
)

// User has the info for the user
type User struct {
	name     string
	surname  string
	email    string
	pswdHash string
	pswdSalt string
}

// Save generates passwords and saves info to DB
func Save(name string, surname string, email string, password string) (err error) {
	salt := make([]byte, PwSaltBytes)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Errorf("Error creating salt for password: %s", err.Error())
		return
	}

	hash, err := scrypt.Key([]byte("ashdasd"), salt, 1<<14, 8, 1, PwHashBytes)
	if err != nil {
		log.Errorf("Error creating password: %s", err.Error())
		return
	}
	log.Debugf("%x", hash)
	return
}

// CheckInfo checks if the info of the user is valid info
// Passwords match and nothing is blank
// Returns error so that we can inform the user trying to register
func CheckInfo(name string, surname string, email string, pswd1 string, pswd2 string) (err error) {
	return
}
