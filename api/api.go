package api

import (
	"errors"
	"fmt"
	"reflect"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

const (
	// iota+1 is used to distinguish from one of this to a just initialized int
	Created = iota + 1
	Updated
	Deleted

	UpdatedText = "Updated"
	CreatedText = "Created"
	DeletedText = "Deleted"

	GOOGLE  = 1
	OUTLOOK = 2

	maxBackoff = 5
)

var states = map[string]int{
	"Created": Created,
	"Updated": Updated,
	"Deleted": Deleted}

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

	MarkWrong()
	GetState() int
	SetState(string) error
	PrepareFields()
	CanProcessAgain() bool
	IncrementBackoff()
}

type SubscriptionManager interface {
	Subscribe(CalendarManager) error
	Renew() error
	Delete() error
	GetID() string
	GetUUID() uuid.UUID
	GetAccount() AccountManager
	GetType() string
	GetExpiration() string
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
	//m := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("conver only accepts structs, got %T", v))
	}

	typ := v.Type()
	//logger.Debugf("%d\n", v.NumField())
	for i := 0; i < v.NumField(); i++ {
		//logger.Debugln("Looping...")
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" && tagv != "-" {
			log.Debugf("tag: %s, value: %s", tagv, v.Field(i).Interface())
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

	structFieldValue.Set(val)
	return nil
}

func StartSync(calendar CalendarManager) (err error) {
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
	//TODO: create subscriptions for all calendars
	//TODO synchronize all events from principal to the other accounts

	return
}
