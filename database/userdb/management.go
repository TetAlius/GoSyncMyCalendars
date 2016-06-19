package userdb

import (
	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"github.com/TetAlius/GoSyncMyCalendars/database"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// SaveUser saves the user info and the salt of the password
func SaveUser(db *sql.DB, name string, surname string, email string, digestedPwd string, salt string) (err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}
	sentenceUser := `INSERT INTO users(email, name, surname, dpswd) VALUES($1, $2, $3, $4)`
	sentenceSalt := `INSERT INTO salt(salt, user_email) VALUES($1, $2)`

	resultUser, err := db.Exec(sentenceUser, email, name, surname, digestedPwd)
	if err != nil {
		log.Errorf("User with email: %s unabled to save, db.Exec", email)
		log.Errorf("%s", err.Error())
	}

	nUser, err := resultUser.RowsAffected()
	if err != nil {
		log.Errorln("Something went wrong when looking number of rows on user")
	}
	if nUser != 1 {
		return customErrors.RowsError{Expected: 1, Received: nUser, Message: "On user.Save with email: " + email}
	}

	resultSalt, err := db.Exec(sentenceSalt, salt, email)
	if err != nil {
		log.Errorf("Salt for user with email: %s unabled to save, db.Exec", email)
		log.Errorf("%s", err.Error())
	}

	nSalt, err := resultSalt.RowsAffected()
	if err != nil {
		log.Errorln("Something went wrong when looking number of rows on user")
	}
	if nSalt != 1 {
		return customErrors.RowsError{Expected: 1, Received: nSalt, Message: "On user.Save with email: " + email}
	}
	return
}

//CheckUser checks if the user exists and if the password is correct
func CheckUser(db *sql.DB, mail string, pswd string) (salt string, dpswd string, err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}

	stmtUser, err := db.Prepare(`SELECT dpswd FROM users WHERE email = $1`)
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	rowUser := stmtUser.QueryRow(mail)
	err = rowUser.Scan(&dpswd)
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	log.Debugln(dpswd)

	stmtSalt, err := db.Prepare(`SELECT salt FROM salt WHERE user_email = $1`)
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	rowsalt := stmtSalt.QueryRow(mail)
	err = rowsalt.Scan(&salt)
	log.Debugln(salt)
	return
}
