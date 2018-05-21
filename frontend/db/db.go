package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}
