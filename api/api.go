package api

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

type AccountManager interface {
	Refresh() error

	GetAllCalendars() ([]CalendarManager, error)
	GetCalendar(string) (CalendarManager, error)
	GetPrimaryCalendar() (CalendarManager, error)
	AuthorizationRequest() string
	Mail() string
}

type CalendarManager interface {
	Update(AccountManager) error
	Delete(AccountManager) error
	Create(AccountManager) error

	GetAllEvents(AccountManager) ([]EventManager, error)
	GetEvent(AccountManager, string) (EventManager, error)

	GetID() string
}

type EventManager interface {
	Create(AccountManager) error
	Update(AccountManager) error
	Delete(AccountManager) error
	GetID() string
	GetCalendar() CalendarManager

	PrepareTime() error
}

type SubscriptionManager interface {
	Subscribe(AccountManager, CalendarManager) error
	Renew(AccountManager) error
	Delete(AccountManager) error
	GetID() string
}

type RefreshError struct {
	Code    string `json:"error,omitempty"`
	Message string `json:"error_description,omitempty"`
}

func (err RefreshError) Error() string {
	return fmt.Sprintf("code: %s. message: %s", err.Code, err.Message)
}

func Convert(from EventManager, to EventManager) (err error) {
	var tag string
	switch from.(type) {
	case *OutlookEvent:
		tag = "outlook"
	case *GoogleEvent:
		tag = "google"
	default:
		return errors.New(fmt.Sprintf("type: %s not suported", reflect.TypeOf(from)))
	}
	logger.Debugln(tag)
	return
}
