package api

import (
	"errors"
	"fmt"
	"reflect"

	"time"

	"os"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
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
	Principal() bool
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

	SetSyncToken(string)
	GetSyncToken() string
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
	PrepareFields()
	CanProcessAgain() bool
	IncrementBackoff()
	SetInternalID(int)
	GetInternalID() int
}

type SubscriptionManager interface {
	Subscribe(CalendarManager) error
	Renew() error
	Delete() error
	GetID() string
	GetUUID() uuid.UUID
	GetAccount() AccountManager
	GetType() string
	setTime()
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

func Convert(from EventManager, to EventManager) (err error) {
	err = convert(from, to)
	if err != nil {
		return errors.New(fmt.Sprintf("could not convert events: %s", err.Error()))
	}
	to.PrepareFields()
	return
}

func convert(in interface{}, out interface{}) (err error) {
	log.Debugln("Converting...")
	tag := "sync"

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("conver only accepts structs, got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" && tagv != "-" {
			err := setField(out, tagv, v.Field(i).Interface())
			if err != nil {
				return err
			}
		}
	}
	return
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New(fmt.Sprintf("no such field: %s in obj", name))
	}

	if !structFieldValue.CanSet() {
		return errors.New(fmt.Sprintf("cannot set %s field value", name))
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New(fmt.Sprintf("provided value type didn't match obj field type"))
	}
	//TODO: error here
	sentry := sentryClient()
	sentry.CapturePanicAndWait(func() { structFieldValue.Set(val) }, map[string]string{"api": "setField"})

	return nil
}

func PrepareSync(calendar CalendarManager) (err error) {
	err = calendar.GetAccount().Refresh()
	if err != nil {
		log.Errorf("error refreshing account: %s", err.Error())
		return
	}

	cal, err := calendar.GetAccount().GetCalendar(calendar.GetID())
	convert(cal, calendar)
	for _, calen := range calendar.GetCalendars() {
		err := convert(calendar, calen)
		if err != nil {
			log.Errorf("error converting info: %s", err.Error())
			return err
		}
		log.Debugf("Name1: %s Name2: %s", calendar.GetName(), calen.GetName())
		err = calen.GetAccount().Refresh()
		if err != nil {
			log.Errorf("error refreshing account calendar: %s error: %s", calen.GetID(), err.Error())
			return err
		}
		err = calen.Update()

		if err != nil {
			log.Errorf("error updating calendar: %s error: %s", calen.GetID(), err.Error())
			return err
		}
	}
	return
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
