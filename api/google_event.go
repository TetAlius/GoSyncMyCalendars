package api

import (
	"bytes"
	"errors"
	"fmt"

	"encoding/json"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (event *GoogleEvent) Create(a AccountManager) (err error) {
	log.Debugln("createEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	calendar := event.Calendar
	event.Calendar = nil

	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)
	event.Calendar = calendar

	contents, err := util.DoRequest("POST",
		fmt.Sprintf(route, event.Calendar.GetQueryID()),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return errors.New(fmt.Sprintf("error creating event in g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)

	log.Debugf("Response: %s", contents)
	return
}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *GoogleEvent) Update(a AccountManager) (err error) {
	log.Debugln("updateEvent google")
	//TODO: Test if ids are two given

	//Meter en los header el etag
	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	calendar := event.Calendar
	event.Calendar = nil
	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	event.Calendar = calendar

	contents, err := util.DoRequest("PUT",
		fmt.Sprintf(route, event.Calendar.GetQueryID(), event.ID),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return errors.New(fmt.Sprintf("error updating event of g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)

	log.Debugf("Response: %s", contents)
	return
}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *GoogleEvent) Delete(a AccountManager) (err error) {
	log.Debugln("deleteEvent google")
	//TODO: Test if ids are two given

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest(
		"DELETE",
		fmt.Sprintf(route, event.Calendar.GetQueryID(), event.ID),
		nil,
		a.AuthorizationRequest(),
		"")

	if err != nil {
		log.Errorf("error deleting event of g calendar for email %s. %s", a.Mail(), err.Error())
	}

	log.Debugf("Contents: %s", contents)
	if len(contents) != 0 {
		return errors.New(fmt.Sprintf("error deleting google event %s: %s", event.ID, contents))
	}

	return
}
func (event *GoogleEvent) GetID() string {
	return event.ID
}

func (event *GoogleEvent) GetCalendar() CalendarManager {
	return event.Calendar
}

func (event *GoogleEvent) PrepareTime() (err error) {
	panic("IMPLEMENT ME")
	return
}
