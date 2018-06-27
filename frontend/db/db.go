package db

import (
	"database/sql"

	"github.com/getsentry/raven-go"
	"github.com/lib/pq"
)

const uniqueViolationError = pq.ErrorCode("23505")

// Database object for the frontend
type Database struct {
	client *sql.DB
	sentry *raven.Client
}

// Function that creates a new database instance
func New(client *sql.DB, sentry *raven.Client) Database {
	return Database{
		client: client,
		sentry: sentry,
	}

}

// Method that closes the database client
func (data Database) Close() error {
	return data.client.Close()
}
