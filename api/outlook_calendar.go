package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/convert"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// Method that returns a OutlookCalendar given specific info
func RetrieveOutlookCalendar(ID string, uid string, account *OutlookAccount) *OutlookCalendar {
	cal := new(OutlookCalendar)
	cal.ID = ID
	cal.account = account
	cal.uuid = uid
	return cal
}

// Method that creates the calendar
//
// POST https://outlook.office.com/api/v2.0/me/calendars
func (calendar *OutlookCalendar) Create() (err error) {
	log.Debugln("createCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	headers["X-AnchorMailbox"] = calendar.GetAccount().Mail()

	contents, err := util.DoRequest(http.MethodPost,
		route,
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	return err
}

// Method that updates the calendar
//
// PUT https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (calendar *OutlookCalendar) Update() error {
	log.Debugln("updateCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	headers["X-AnchorMailbox"] = calendar.GetAccount().Mail()

	contents, err := util.DoRequest(http.MethodPatch,
		fmt.Sprintf(route, calendar.GetID()),
		bytes.NewBuffer(data),
		headers, nil)

	log.Debugf("contents: %s", contents)
	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	// default outlook calendar cannot be renamed, so ignore this kind of error as the request is valid.
	if err != nil && !strings.Contains(err.Error(), "default calendar cannot be renamed") {
		return err
	}
	if err != nil && strings.Contains(err.Error(), "default calendar cannot be renamed") {
		cal, err := calendar.GetAccount().GetCalendar(calendar.GetID())
		if err != nil {
			return err
		}
		convert.Convert(cal, calendar)
	}

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	return err

}

// Method that deletes the calendar
//
// DELETE https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (calendar *OutlookCalendar) Delete() (err error) {
	log.Debugln("deleteCalendar outlook")
	if len(calendar.GetID()) == 0 {
		return errors.New("no ID for calendar was given")
	}

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	headers["X-AnchorMailbox"] = calendar.GetAccount().Mail()

	contents, err := util.DoRequest(http.MethodDelete,
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}

	if len(contents) != 0 {
		err = createOutlookResponseError(contents)
		return err
	}
	return

}

// Method that returns all events inside the calendar
//
// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func (calendar *OutlookCalendar) GetAllEvents() (events []EventManager, err error) {
	log.Debugln("getAllEvents outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	headers["X-AnchorMailbox"] = calendar.GetAccount().Mail()
	headers["Prefer"] = "outlook.timezone=UTC, outlook.body-content-type=text"

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}

	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}
	eventListResponse := new(OutlookEventListResponse)
	err = json.Unmarshal(contents, &eventListResponse)

	for _, s := range eventListResponse.Events {
		s.SetCalendar(calendar)
		if err != nil {
			return
		}
		s.setAllDay()
		events = append(events, s)
	}
	return
}

// Method that returns a single event given the ID
//
// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func (calendar *OutlookCalendar) GetEvent(ID string) (event EventManager, err error) {
	log.Debugln("getEvent outlook")
	if len(ID) == 0 {
		return nil, errors.New("an ID for the event must be given")
	}

	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	headers["X-AnchorMailbox"] = calendar.GetAccount().Mail()
	headers["Prefer"] = "outlook.timezone=UTC,outlook.body-content-type=text"

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, ID),
		nil,
		headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return
	}

	eventResponse := new(OutlookEventResponse)
	err = json.Unmarshal(contents, &eventResponse)
	if err != nil {
		return
	}

	err = eventResponse.SetCalendar(calendar)
	if err != nil {
		return
	}
	e := eventResponse.OutlookEvent
	e.setAllDay()
	return e, nil
}

// Method that sets the account which the calendar belongs
func (calendar *OutlookCalendar) SetAccount(a AccountManager) (err error) {
	switch x := a.(type) {
	case *OutlookAccount:
		calendar.account = x
	default:
		return errors.New(fmt.Sprintf("type of account not valid for outlook: %T", x))
	}
	return
}

// Method that returns the ID of the calendar
func (calendar *OutlookCalendar) GetID() string {
	return calendar.ID
}

// Method that returns the ID formatted for a query request
func (calendar *OutlookCalendar) GetQueryID() string {
	return calendar.ID
}

// Method that returns the account
func (calendar *OutlookCalendar) GetAccount() AccountManager {
	return calendar.account
}

// Method that returns the name of the calendar
func (calendar *OutlookCalendar) GetName() string {
	return calendar.Name
}

// Method that returns the internal UUID given to the calendar
func (calendar *OutlookCalendar) GetUUID() string {
	return calendar.uuid
}

// Method that sets the internal UUID for the calendar
func (calendar *OutlookCalendar) SetUUID(id string) {
	calendar.uuid = id
}

// Method that sets the synced calendars
func (calendar *OutlookCalendar) SetCalendars(calendars []CalendarManager) {
	calendar.calendars = calendars

}

// Method that returns the synced calendar
func (calendar *OutlookCalendar) GetCalendars() []CalendarManager {
	return calendar.calendars
}

// Method that creates an empty event
func (calendar *OutlookCalendar) CreateEmptyEvent(ID string) EventManager {
	return &OutlookEvent{ID: ID, calendar: calendar}
}
