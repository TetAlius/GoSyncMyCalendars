package db

import (
	"database/sql"

	"github.com/getsentry/raven-go"
	_ "github.com/lib/pq"
)

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
