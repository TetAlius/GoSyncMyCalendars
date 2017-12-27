package google

import "github.com/TetAlius/GoSyncMyCalendars/api"

type Account struct {
	//	TODO
}

func (a *Account) Refresh() error {
	return nil
}

func (a *Account) GetAllCalendars() ([]api.CalendarManager, error) {
	return nil, nil
}

func (a *Account) CreateCalendar(calendar *api.CalendarManager) error {
	return nil
}

func (a *Account) GetCalendar(id string) (api.CalendarManager, error) {
	return nil, nil
}

func (a *Account) GetPrimaryCalendar() (api.CalendarManager, error) {
	return nil, nil
}
