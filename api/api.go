package api

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
