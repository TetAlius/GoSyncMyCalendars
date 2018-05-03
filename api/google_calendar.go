package api

import (
	"bytes"

	"fmt"

	"encoding/json"

	"net/url"

	"net/http"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/pkg/errors"
)

// PUT https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *GoogleCalendar) Update(a AccountManager) (err error) {
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
	headers["Authorization"] = a.AuthorizationRequest()
	contents, err :=
		util.DoRequest(
			http.MethodPut,
			fmt.Sprintf(route, calendar.GetQueryID()),
			bytes.NewBuffer(data),
			headers)

	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("Contents: %s", contents)
	err = json.Unmarshal(contents, &calendar)

	return
}

// DELETE https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *GoogleCalendar) Delete(a AccountManager) (err error) {
	log.Debugln("Delete calendar")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	contents, err := util.DoRequest(
		http.MethodDelete,
		fmt.Sprintf(route, calendar.GetQueryID()),
		nil,
		headers)

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	if len(contents) != 0 {
		return errors.New(fmt.Sprintf("error deleting google calendar %s: %s", calendar.GetID(), contents))
	}

	log.Debugf("Contents: %s", contents)

	return
}

// POST https://www.googleapis.com/calendar/v3/calendars
func (calendar *GoogleCalendar) Create(a AccountManager) (err error) {
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
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err :=
		util.DoRequest(
			http.MethodPost,
			route,
			bytes.NewBuffer(data),
			headers)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("Contents: %s", contents)
	err = json.Unmarshal(contents, &calendar)
	return
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (calendar *GoogleCalendar) GetAllEvents(a AccountManager) (events []EventManager, err error) {
	log.Debugln("getAllEvents google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, calendar.GetQueryID()),
		nil,
		headers)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventList := new(GoogleEventList)
	err = json.Unmarshal(contents, &eventList)

	events = make([]EventManager, len(eventList.Events))
	for i, s := range eventList.Events {
		s.Calendar = calendar
		events[i] = s
	}
	return
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (calendar *GoogleCalendar) GetEvent(a AccountManager, eventID string) (event EventManager, err error) {
	log.Debugln("getEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err := util.DoRequest(
		http.MethodGet,
		fmt.Sprintf(route, calendar.GetQueryID(), eventID),
		nil,
		headers)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of g calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("Contents: %s", contents)
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventResponse := new(GoogleEvent)
	err = json.Unmarshal(contents, &eventResponse)

	eventResponse.Calendar = calendar
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
		headers)

	log.Debugf("%s\n", contents)
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
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

func (calendar GoogleCalendar) GetQueryID() string {
	return url.QueryEscape(calendar.GetID())
}

func (calendar GoogleCalendar) GetID() string {
	return calendar.ID
}
