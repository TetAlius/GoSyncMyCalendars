package api

import (
	"fmt"
	"strings"

	"time"

	"os"

	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
)

const (
	// iota+1 is used to distinguish from one of this to a just initialized int
	// Constants used to know which kind of state the event is in
	Created = iota + 1
	Updated
	Deleted

	// Different kinds of accounts used
	GOOGLE  = 1
	OUTLOOK = 2

	// maximum number of wrong requests in synchronization
	maxBackoff = 5
)

// tagOptions is the string following a comma in a struct field's
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Interface for account that defines the needed method to work inside the project
type AccountManager interface {
	// Method to refresh the access to the account
	Refresh() error

	// Method that retrieves all calendars from account
	GetAllCalendars() ([]CalendarManager, error)
	// Method that retrieves one calendar given an ID
	GetCalendar(string) (CalendarManager, error)
	// Method that returns the principal calendar from the account
	GetPrimaryCalendar() (CalendarManager, error)
	// Method that format the authorization request
	AuthorizationRequest() string
	// Method that returns the mail associated with the account
	Mail() string

	// Method that sets which kind of account is
	SetKind(int)
	// Method that returns the token type
	GetTokenType() string
	// Method that returns the refresh token
	GetRefreshToken() string
	// Method that returns the kind of the account
	GetKind() int
	// Method that returns the access token
	GetAccessToken() string

	// Method that returns the internal ID given to the account on DB
	GetInternalID() int
	// Method that sets all synced calendars associated with the account
	SetCalendars([]CalendarManager)
	// Method that returns all synced calendars associated with the account
	GetSyncCalendars() []CalendarManager
}

// Interface for calendar that defines the needed method to work inside the project
type CalendarManager interface {
	// Method that sets the account which the calendar belongs
	SetAccount(AccountManager) error
	// Method that sets the synced calendars
	SetCalendars([]CalendarManager)
	// Method that returns the synced calendar
	GetCalendars() []CalendarManager

	// Method that updates the calendar
	Update() error
	// Method that deletes the calendar
	Delete() error
	// Method that creates the calendar
	Create() error

	// Method that returns all events inside the calendar
	GetAllEvents() ([]EventManager, error)
	// Method that returns a single event given the ID
	GetEvent(string) (EventManager, error)

	// Method that returns the ID of the calendar
	GetID() string
	// Method that returns the ID formatted for a query request
	GetQueryID() string
	// Method that returns the name of the calendar
	GetName() string
	// Method that returns the account
	GetAccount() AccountManager
	// Method that returns the internal UUID given to the calendar
	GetUUID() string
	// Method that sets the internal UUID for the calendar
	SetUUID(string)
	// Method that creates an empty event
	CreateEmptyEvent(string) EventManager
}

// Interface for event that defines the needed method to work inside the project
type EventManager interface {
	// Method that sets the calendar which have the event
	SetCalendar(CalendarManager) error
	// Method that sets the events syncing with this
	SetRelations([]EventManager)

	// Method that creates the event
	Create() error
	// Method that updates the event
	Update() error
	// Method that deletes the event
	Delete() error
	// Method that returns the ID of the event
	GetID() string

	// Method that returns the calendar which have this event
	GetCalendar() CalendarManager

	// Method that returns the syncing events with this
	GetRelations() []EventManager

	// Method that returns the last update date
	GetUpdatedAt() (time.Time, error)
	// Method that returns the state of the event
	GetState() int
	// Method that sets the state of the event
	SetState(int)
	// Method that checks if the event can try sync again
	CanProcessAgain() bool
	// Method that increments the number of failed attempts to sync
	IncrementBackoff()
	// Method that sets the internal ID generated on db
	SetInternalID(int)
	// Method that gets the internal ID of the event
	GetInternalID() int

	// Method that sets all day to the necessary attributes
	setAllDay()
}

// Interface for subscription that defines the needed method to work inside the project
type SubscriptionManager interface {
	// Method that subscribes calendar for notifications
	Subscribe(CalendarManager) error
	// Method that renews subscription
	Renew() error
	// Method that deletes subscription
	Delete() error
	// Method that returns the ID of the subscription
	GetID() string
	// Method that returns the UUID of the subscription
	GetUUID() uuid.UUID
	// Method that returns the account of the subscription
	GetAccount() AccountManager
	// Method that returns the type of the subscription
	GetType() string
	// Method that returns the expiration date of the subscription
	GetExpirationDate() time.Time
	// Method that returns the resourceID of the subscription
	GetResourceID() string
}

// Specific error for a refresh try
type RefreshError struct {
	Code    string `json:"error,omitempty"`
	Message string `json:"error_description,omitempty"`
}

// Method implementing error interface
func (err RefreshError) Error() string {
	return fmt.Sprintf("code: %s. message: %s", err.Code, err.Message)
}

// Function to know in which state the event is
func GetChangeType(onCloud bool, onDB bool) int {
	if onCloud && !onDB {
		return Created
	} else if onCloud && onDB {
		return Updated
	} else if !onCloud && onDB {
		return Deleted
	}
	return 0
}

// Function that returns a sentry client prepared to report
func sentryClient() (sentry *raven.Client) {
	sentry, _ = raven.New(os.Getenv("SENTRY_DSN"))
	sentry.SetEnvironment(os.Getenv("ENVIRONMENT"))
	sentry.SetRelease(os.Getenv("RELEASE"))
	return
}
