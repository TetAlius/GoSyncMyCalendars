package db

import (
	"database/sql"

	"github.com/getsentry/raven-go"
	"github.com/lib/pq"
)

const uniqueViolationError = pq.ErrorCode("23505")

type Database struct {
	client *sql.DB
	sentry *raven.Client
}

func New(client *sql.DB, sentry *raven.Client) Database {
	return Database{
		client: client,
		sentry: sentry,
	}

}

func (data Database) Close() error {
	return data.client.Close()
}
