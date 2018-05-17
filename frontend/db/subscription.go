package db

import "github.com/google/uuid"

type Subscription struct {
	UUID         uuid.UUID
	CalendarUUID uuid.UUID
}
