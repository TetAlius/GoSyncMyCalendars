package api

import "fmt"

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
