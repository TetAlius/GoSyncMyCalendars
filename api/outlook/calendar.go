package outlook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type CalendarResponse struct {
	OdataContext string `json:"@odata.context"`
	*Calendar
}

type CalendarListResponse struct {
	OdataContext string      `json:"@odata.context"`
	Calendars    []*Calendar `json:"value"`
}

// CalendarInfo TODO
type Calendar struct {
	OdataID string `json:"@odata.id,omitempty"`

	CalendarView        []Event       `json:"CalendarView,omitempty"`
	CanEdit             bool          `json:"CanEdit,omitempty"`
	CanShare            bool          `json:"CanShare,omitempty"`
	CanViewPrivateItems bool          `json:"CanViewPrivateItems,omitempty"`
	ChangeKey           string        `json:"ChangeKey,omitempty"`
	Color               CalendarColor `json:"Color,omitempty"`
	Events              []Event       `json:"Events,omitempty"`
	ID                  string        `json:"Id"`
	Name                string        `json:"Name,omitempty"`
	Owner               EmailAddress  `json:"Owner,omitempty"`

	//	IsDefaultCalendar             bool         `json:"IsDefaultCalendar,omitempty"`
	//	IsShared                      bool         `json:"IsShared,omitempty"`
	//	IsSharedWithMe                bool         `json:"IsSharedWithMe,omitempty"`
	//	MultiValueExtendedProperties  []Properties `json:"MultiValueExtendedProperties"`
	//	SingleValueExtendedProperties []Properties `json:"SingleValueExtendedProperties"`
}

// Specifies the color theme to distinguish the calendar from other calendars in a UI.
// The property values are: LightBlue=0, LightGreen=1, LightOrange=2, LightGray=3,
// LightYellow=4, LightTeal=5, LightPink=6, LightBrown=7, LightRed=8, MaxColor=9, Auto=-1
type CalendarColor string

type EmailAddress struct {
	Address string `json:"Address,omitempty"`
	Name    string `json:"Name,omitempty"`
}

func (calendar *Calendar) Create(a api.AccountManager) (err error) {
	log.Debugln("createCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	contents, err := util.DoRequest("POST",
		route,
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("%s\n", contents)

	calendarResponse := CalendarResponse{OdataContext: "", Calendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	log.Debugf("Response: %s", contents)

	return err
}

func (calendar *Calendar) Update(a api.AccountManager) error {
	log.Debugln("updateCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("Error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	contents, err := util.DoRequest("PATCH",
		fmt.Sprintf(route, calendar.ID),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("%s\n", contents)

	calendarResponse := CalendarResponse{OdataContext: "", Calendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	return err

}

func (calendar *Calendar) Delete(a api.AccountManager) (err error) {
	log.Debugln("deleteCalendar outlook")
	if len(calendar.ID) == 0 {
		return errors.New("no ID for calendar was given")
	}

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("DELETE",
		fmt.Sprintf(route, calendar.ID),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("%s\n", contents)
	return

}

func (calendar *Calendar) GetAllEvents(a api.AccountManager) (events []api.EventManager, err error) {
	log.Debugln("getAllEvents outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, calendar.ID),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("%s\n", contents)
	eventListResponse := new(EventListResponse)
	err = json.Unmarshal(contents, &eventListResponse)

	events = make([]api.EventManager, len(eventListResponse.Events))
	for i, s := range eventListResponse.Events {
		s.Calendar = calendar
		events[i] = s
	}
	return
}

// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func (calendar *Calendar) GetEvent(a api.AccountManager, ID string) (event api.EventManager, err error) {
	log.Debugln("getEvent outlook")
	if len(ID) == 0 {
		return nil, errors.New("an ID for the event must be given")
	}

	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, ID),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("%s\n", contents)
	eventResponse := new(EventResponse)
	err = json.Unmarshal(contents, &eventResponse)

	eventResponse.Calendar = calendar
	event = eventResponse.Event

	return
}
func (calendar *Calendar) GetID() string {
	return calendar.ID
}
