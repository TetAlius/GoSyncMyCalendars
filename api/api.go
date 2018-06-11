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
	Created = iota + 1
	Updated
	Deleted

	GOOGLE  = 1
	OUTLOOK = 2

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

type AccountManager interface {
	Refresh() error

	GetAllCalendars() ([]CalendarManager, error)
	GetCalendar(string) (CalendarManager, error)
	GetPrimaryCalendar() (CalendarManager, error)
	AuthorizationRequest() string
	Mail() string

	SetKind(int)
	GetTokenType() string
	GetRefreshToken() string
	GetKind() int
	GetAccessToken() string

	GetInternalID() int
	SetCalendars([]CalendarManager)
	GetSyncCalendars() []CalendarManager
}

type CalendarManager interface {
	SetAccount(AccountManager) error
	SetCalendars([]CalendarManager)
	GetCalendars() []CalendarManager

	Update() error
	Delete() error
	Create() error

	GetAllEvents() ([]EventManager, error)
	GetEvent(string) (EventManager, error)

	GetID() string
	GetQueryID() string
	GetName() string
	GetAccount() AccountManager
	SetName(string)
	GetUUID() string
	SetUUID(string)
	CreateEmptyEvent(string) EventManager
}

type EventManager interface {
	SetCalendar(CalendarManager) error
	SetRelations([]EventManager)

	Create() error
	Update() error
	Delete() error
	GetID() string

	GetCalendar() CalendarManager

	GetRelations() []EventManager

	GetUpdatedAt() (time.Time, error)
	MarkWrong()
	GetState() int
	SetState(int)
	CanProcessAgain() bool
	IncrementBackoff()
	SetInternalID(int)
	GetInternalID() int

	setAllDay()
}

type SubscriptionManager interface {
	Subscribe(CalendarManager) error
	Renew() error
	Delete() error
	GetID() string
	GetUUID() uuid.UUID
	GetAccount() AccountManager
	GetType() string
	GetExpirationDate() time.Time
	GetResourceID() string
}

type RefreshError struct {
	Code    string `json:"error,omitempty"`
	Message string `json:"error_description,omitempty"`
}

func (err RefreshError) Error() string {
	return fmt.Sprintf("code: %s. message: %s", err.Code, err.Message)
}

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

func sentryClient() (sentry *raven.Client) {
	sentry, _ = raven.New(os.Getenv("SENTRY_DSN"))
	sentry.SetEnvironment(os.Getenv("ENVIRONMENT"))
	sentry.SetRelease(os.Getenv("RELEASE"))
	return
}
