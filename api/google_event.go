package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"time"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

var recoverableGoogleErrors = map[string]bool{}

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (event *GoogleEvent) Create() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("createEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err := util.DoRequest(http.MethodPost,
		fmt.Sprintf(route, event.GetCalendar().GetQueryID()),
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating event in g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)
	err = event.extractTime()
	return
}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *GoogleEvent) Update() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("updateEvent google")
	//TODO: Test if ids are two given

	//Meter en los header el etag
	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err := util.DoRequest(http.MethodPut,
		fmt.Sprintf(route, event.GetCalendar().GetQueryID(), event.ID),
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error updating event of g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)
	err = event.extractTime()
	return
}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *GoogleEvent) Delete() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("deleteEvent google")
	//TODO: Test if ids are two given

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		return
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err := util.DoRequest(
		http.MethodDelete,
		fmt.Sprintf(route, event.GetCalendar().GetQueryID(), event.ID),
		nil,
		headers, nil)

	if err != nil {
		log.Errorf("error deleting event of g calendar for email %s. %s", a.Mail(), err.Error())
	}

	if len(contents) != 0 {
		return errors.New(fmt.Sprintf("error deleting google event %s: %s", event.ID, contents))
	}

	return
}

func (event *GoogleEvent) GetID() string {
	return event.ID
}

func (event *GoogleEvent) GetCalendar() CalendarManager {
	return event.calendar
}

func (event *GoogleEvent) GetRelations() []EventManager {
	return event.relations
}

func (event *GoogleEvent) PrepareFields() {
	var startDate, endDate string
	if event.IsAllDay {
		startDate = event.StartsAt.Format("2006-01-02")
		endDate = event.EndsAt.Format("2006-01-02")
	}

	event.Start = &GoogleTime{startDate, event.StartsAt.Format(time.RFC3339), "UTC"}
	event.End = &GoogleTime{endDate, event.EndsAt.Format(time.RFC3339), "UTC"}
	return
}

func (event *GoogleEvent) CanProcessAgain() bool {
	return event.exponentialBackoff < maxBackoff
}

func (event *GoogleEvent) SetCalendar(calendar CalendarManager) (err error) {
	switch x := calendar.(type) {
	case *GoogleCalendar:
		event.calendar = x
	default:
		return errors.New(fmt.Sprintf("type of calendar not valid for google: %T", x))
	}
	return
}

func (event *GoogleEvent) SetRelations(relations []EventManager) {
	event.relations = relations
}

func (event *GoogleEvent) MarkWrong() {
	//	TODO: Implement marking to db
	log.Fatalf("not implemented yet. ID: %s", event.GetID())
}

func (event *GoogleEvent) IncrementBackoff() {
	event.exponentialBackoff += 1
}

func (event *GoogleEvent) SetState(stateInformed string) (err error) {
	state := states[stateInformed]
	if state == 0 {
		return errors.New(fmt.Sprintf("state %s not supported", stateInformed))
	}
	return
}

func (event *GoogleEvent) GetState() int {
	return event.state
}

func (event *GoogleEvent) extractTime() (err error) {
	var start, end, format string
	if len(event.Start.Date) != 0 && len(event.End.Date) != 0 {
		event.IsAllDay = true
		start = event.Start.Date
		end = event.End.Date
		format = "2006-01-02"

	} else {
		event.IsAllDay = false
		start = event.Start.DateTime
		end = event.End.DateTime
		format = time.RFC3339
	}

	event.StartsAt, err = time.Parse(format, start)
	if err != nil {
		return errors.New(fmt.Sprintf("error parsing start time: %s %s", start, err.Error()))
	}
	event.EndsAt, err = time.Parse(format, end)
	if err != nil {
		return errors.New(fmt.Sprintf("error parsing end time: %s %s", end, err.Error()))
	}
	return
}
