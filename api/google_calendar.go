package api

import (
	"bytes"

	"fmt"

	"encoding/json"

	"net/url"

	"net/http"

	"errors"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

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

// DELETE https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *GoogleCalendar) Delete() (err error) {
	//return
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

	events = make([]EventManager, len(eventList.Events))
	for i, s := range eventList.Events {
		s.SetCalendar(calendar)
		err := s.extractTime()
		if err != nil {
			return nil, err
		}
		events[i] = s
	}
	return
}

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
	err = eventResponse.extractTime()
	if err != nil {
		return
	}

	eventResponse.SetCalendar(calendar)
	event = eventResponse

	return
}

func (calendar *GoogleCalendar) Subscribe(a AccountManager) (err error) {
	log.Debugln("subscribe calendar google")

	route, err := util.CallAPIRoot("google/calendars/subscription")
	log.Debugln(route)
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data := []byte(`{
	  "id": "01234567-89ab-cdef-0123456789ab",
	  "type": "web_hook",
	  "address": "https://mcbjyngjgh.execute-api.eu-west-1.amazonaws.com/prod/testing"
	}`)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPost,
		route,
		bytes.NewBuffer(data),
		headers, nil)

	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	return
}

func (calendar *GoogleCalendar) SetAccount(a AccountManager) (err error) {
	switch x := a.(type) {
	case *GoogleAccount:
		calendar.account = x
	default:
		return errors.New(fmt.Sprintf("type of account not valid for google: %T", x))
	}
	return
}

// There's no explicit way to renew subscription.
// One new must be created
func (calendar *GoogleCalendar) RenewSubscription(a AccountManager, subscriptionID string) (err error) {
	panic("NOT YET IMPLEMENTED")
	return
}

func (calendar *GoogleCalendar) DeleteSubscription(a AccountManager, subscriptionID string) (err error) {
	panic("NOT YET IMPLEMENTED")
	return
}

func (calendar *GoogleCalendar) GetQueryID() string {
	return url.QueryEscape(calendar.GetID())
}

func (calendar *GoogleCalendar) GetID() string {
	return calendar.ID
}

func (calendar *GoogleCalendar) GetName() string {
	return calendar.Name
}

func (calendar *GoogleCalendar) GetAccount() AccountManager {
	return calendar.account
}

func (calendar *GoogleCalendar) GetUUID() string {
	return calendar.uuid
}

func (calendar *GoogleCalendar) SetUUID(id string) {
	calendar.uuid = id
}

func (calendar *GoogleCalendar) SetCalendars(calendars []CalendarManager) {
	calendar.calendars = calendars

}
func (calendar *GoogleCalendar) GetCalendars() []CalendarManager {
	return calendar.calendars
}
