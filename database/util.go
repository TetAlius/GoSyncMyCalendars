package database

import (
	"database/sql"
	"os"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func getDBConfig() (connection string) {
	host := os.Getenv("DBHOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DBUSER")
	if user == "" {
		user = "marta"
	}

	pass := os.Getenv("DBPASS")
	if pass == "" {
		pass = "1234"
	}

	data := os.Getenv("DATABASE")
	if data == "" {
		data = "marta"
	}
	port := os.Getenv("DBPORT")
	if port == "" {
		port = "5432"
	}
	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + data
}

// OpenDB TODO
func OpenDB() *sql.DB {
	connection := getDBConfig()
	db, err := sql.Open("pgx", connection)
	if err != nil {
		log.Fatalf("sql.Open failed: %v", err)
	}

	return db
}

const (
	//GOOGLE is the string for the account type in the BD
	GOOGLE = "google"
	//OUTLOOK is the string for the account type in the BD
	OUTLOOK = "outlook"
)

// HandleCalendarRowsForAccount TODO
func HandleCalendarRowsForAccount(account string, rows *sql.Rows) (v interface{}, err error) {
	return
}
