package db

import "github.com/google/uuid"

type Event struct {
	UUID     uuid.UUID
	Calendar Calendar
}
