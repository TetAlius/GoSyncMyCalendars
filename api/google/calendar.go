package google

import "github.com/TetAlius/GoSyncMyCalendars/api"

type Calendar struct {
	//	TODO
}

func (calendar *Calendar) Update() (err error) {
	return

}
func (calendar *Calendar) Delete() (err error) {
	return
}

func (calendar *Calendar) CreateEvent(api.EventManager) (err error) {
	return
}
func (calendar *Calendar) GetAllEvents() (events []api.EventManager, err error) {
	return
}
func (calendar *Calendar) GetEvent(string) (event api.EventManager, err error) {
	return
}
