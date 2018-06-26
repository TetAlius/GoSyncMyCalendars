package api

import (
	"bytes"

	"fmt"

	"encoding/json"

	"net/url"

	"net/http"

	"errors"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// Method that returns a GoogleCalendar given specific info
func RetrieveGoogleCalendar(ID string, uid string, account *GoogleAccount) *GoogleCalendar {
	cal := new(GoogleCalendar)
	cal.ID = ID
	cal.account = account
	cal.uuid = uid
	return cal
}

// Method that updates the calendar
//
// PUT https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *GoogleCalendar) Update() (err error) {
	log.Debugln("updateCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	contents, err :=
		util.DoRequest(
			http.MethodPut,
			fmt.Sprintf(route, calendar.GetQueryID()),
			bytes.NewBuffer(data),
			headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}

	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &calendar)

	return
}

// Method that deletes the calendar
//
// DELETE https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *GoogleCalendar) Delete() (err error) {
	log.Debugln("Delete calendar")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()
	contents, err := util.DoRequest(
		http.MethodDelete,
		fmt.Sprintf(route, calendar.GetQueryID()),
		nil,
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}

	if len(contents) != 0 {
		err = createGoogleResponseError(contents)
		return err
	}

	return
}

// Method that creates the calendar
//
// POST https://www.googleapis.com/calendar/v3/calendars
func (calendar *GoogleCalendar) Create() (err error) {
	log.Debugln("createCalendar google")
	route, err := util.CallAPIRoot("google/calendars")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugln(route)
	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()

	contents, err :=
		util.DoRequest(
			http.MethodPost,
			route,
			bytes.NewBuffer(data),
			headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &calendar)
	return
}

// Method that returns all events inside the calendar
//
// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (calendar *GoogleCalendar) GetAllEvents() (events []EventManager, err error) {
	log.Debugln("getAllEvents google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()

	queryParams := map[string]string{"timeZone": "UTC"}

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, calendar.GetQueryID()),
		nil,
		headers, queryParams)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of g calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}
	eventList := new(GoogleEventList)
	err = json.Unmarshal(contents, &eventList)

	//events = new([]EventManager)
	for _, event := range eventList.Events {
		event.SetCalendar(calendar)
		// ignore cancelled events
		if event.Status != "cancelled" {
			event.setAllDay()
			//TODO: this status
			events = append(events, event)
		}
	}
	return events, err
}

// Method that returns a single event given the ID
//
// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (calendar *GoogleCalendar) GetEvent(eventID string) (event EventManager, err error) {
	log.Debugln("getEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = calendar.GetAccount().AuthorizationRequest()

	queryParams := map[string]string{"timeZone": "UTC"}

	contents, err := util.DoRequest(
		http.MethodGet,
		fmt.Sprintf(route, calendar.GetQueryID(), eventID),
		nil,
		headers, queryParams)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of g calendar for email %s. %s", calendar.GetAccount().Mail(), err.Error()))
	}

	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	eventResponse := new(GoogleEvent)
	err = json.Unmarshal(contents, &eventResponse)
	if err != nil {
		return
	}
	log.Warningf("GOOGLE EVENT: %s", contents)
	//TODO: this part
	if eventResponse.Status != "cancelled" {
		eventResponse.SetCalendar(calendar)
		event = eventResponse
		event.setAllDay()
	} else {
		return nil, &customErrors.NotFoundError{Message: fmt.Sprintf("event with id: %s not found", eventID)}
	}

	return
}

// Method that sets the account which the calendar belongs
func (calendar *GoogleCalendar) SetAccount(a AccountManager) (err error) {
	switch x := a.(type) {
	case *GoogleAccount:
		calendar.account = x
	default:
		return errors.New(fmt.Sprintf("type of account not valid for google: %T", x))
	}
	return
}

// Method that returns the ID formatted for a query request
func (calendar *GoogleCalendar) GetQueryID() string {
	return url.QueryEscape(calendar.GetID())
}

// Method that returns the ID of the calendar
func (calendar *GoogleCalendar) GetID() string {
	return calendar.ID
}

// Method that returns the name of the calendar
func (calendar *GoogleCalendar) GetName() string {
	return calendar.Name
}

// Method that returns the account
func (calendar *GoogleCalendar) GetAccount() AccountManager {
	return calendar.account
}

// Method that returns the internal UUID given to the calendar
func (calendar *GoogleCalendar) GetUUID() string {
	return calendar.uuid
}

// Method that sets the internal UUID for the calendar
func (calendar *GoogleCalendar) SetUUID(id string) {
	calendar.uuid = id
}

// Method that sets the synced calendars
func (calendar *GoogleCalendar) SetCalendars(calendars []CalendarManager) {
	calendar.calendars = calendars

}

// Method that returns the synced calendar
func (calendar *GoogleCalendar) GetCalendars() []CalendarManager {
	return calendar.calendars
}

// Method that creates an empty event
func (calendar *GoogleCalendar) CreateEmptyEvent(ID string) EventManager {
	return &GoogleEvent{ID: ID, calendar: calendar}
}
