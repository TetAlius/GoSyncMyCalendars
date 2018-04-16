package backend

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/database/userdb"
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
	Name    string
	Surname string
	Email   string
	Pswd    string
}

// Save generates passwords and saves info to DB
func (u *User) Save() (err error) { //name string, surname string, email string, password string) (err error) {
	salt := make([]byte, PwSaltBytes)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Errorf("Error creating salt for password: %s", err.Error())
		return
	}

	hash, err := scrypt.Key([]byte(u.Pswd), salt, 1<<14, 8, 1, PwHashBytes)
	if err != nil {
		log.Errorf("Error creating password: %s", err.Error())
		return
	}

	p := base64.StdEncoding.EncodeToString(hash)
	s := base64.StdEncoding.EncodeToString(salt)
	u.Pswd = p

	err = userdb.SaveUser(nil, u.Name, u.Surname, u.Email, u.Pswd, s) // name, surname, email, p, s)
	return
}

// CheckInfo checks if the info of the user is valid info
// Passwords match and nothing is blank
// Returns error so that we can inform the user trying to register
func CheckInfo(name string, surname string, email string, pswd1 string, pswd2 string) (err error) {
	if len(strings.TrimSpace(name)) == 0 {
		log.Infoln("Name string is empty")
	}
	if len(strings.TrimSpace(surname)) == 0 {
		log.Infoln("Surname string is empty")
	}
	if len(strings.TrimSpace(email)) == 0 {
		log.Infoln("Email string is empty")
	}
	if len(strings.TrimSpace(pswd1)) == 0 {
		log.Infoln("Password 1 string is empty")
	}
	if len(strings.TrimSpace(pswd2)) == 0 {
		log.Infoln("Password 2 string is empty")
	}
	if pswd1 != pswd2 {
		log.Infoln("Passwords do not match")
	}

	return
}

// CheckExistingUser checks if the info introduced belongs to a user
// already in the system or if the info is correct
func CheckExistingUser(mail string, pswd string) (exists bool, err error) {
	salt, dpswd, err := userdb.CheckUser(nil, mail, pswd)
	if err != nil {
		log.Errorln(err)
	}

	s, err := base64.StdEncoding.DecodeString(salt)

	hash, err := scrypt.Key([]byte(pswd), s, 1<<14, 8, 1, PwHashBytes)
	if err != nil {
		log.Errorf("Error %s", err.Error())
	}

	log.Debugf("%s", dpswd)
	p := base64.StdEncoding.EncodeToString(hash)
	log.Debugf("%s", p)

	if dpswd != p {
		log.Debugln("NOT THE SAME :(")
	}

	return true, err
}
