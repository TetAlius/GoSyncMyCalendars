package db

import "github.com/google/uuid"

type Calendar struct {
	UUID         uuid.UUID
	Account      Account
	Name         string
	ID           string
	ParentUUID   uuid.UUID
	Events       []Event
	Subscription Subscription
}
