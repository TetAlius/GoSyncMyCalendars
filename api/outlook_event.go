package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"time"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (event *OutlookEvent) Create() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("createEvent outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
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
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPost,
		fmt.Sprintf(route, event.GetCalendar().GetID()),
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating event in a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	eventResponse := OutlookEventResponse{OdataContext: "", OutlookEvent: event}
	err = json.Unmarshal(contents, &eventResponse)

	err = event.extractTime()
	return
}

func (event *OutlookEvent) Update() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("updateEvent outlook")

	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugln(route)
	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPatch,
		fmt.Sprintf(route, event.ID),
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error updating event of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	eventResponse := OutlookEventResponse{OdataContext: "", OutlookEvent: event}
	err = json.Unmarshal(contents, &eventResponse)

	err = event.extractTime()
	return
}

// DELETE https://outlook.office.com/api/v2.0/me/events/{eventID}
func (event *OutlookEvent) Delete() (err error) {
	a := event.GetCalendar().GetAccount()
	log.Debugln("deleteEvent outlook")

	route, err := util.CallAPIRoot("outlook/events/id")

	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugln(route)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodDelete,
		fmt.Sprintf(route, event.ID),
		nil,
		headers, nil)

	if err != nil {
		log.Errorf("error deleting event of a calendar for email %s. %s", a.Mail(), err.Error())
		return errors.New(fmt.Sprintf("error deleting event of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	if len(contents) != 0 {
		err = createOutlookResponseError(contents)
		return err
	}

	return
}

func (event *OutlookEvent) GetID() string {
	return event.ID
}

func (event *OutlookEvent) GetCalendar() CalendarManager {
	return event.calendar
}

func (event *OutlookEvent) GetRelations() []EventManager {
	return event.relations
}

func (event *OutlookEvent) PrepareFields() {
	event.Start = &OutlookDateTimeTimeZone{event.StartsAt.Format(time.RFC3339Nano), "UTC"}
	event.End = &OutlookDateTimeTimeZone{event.EndsAt.Format(time.RFC3339Nano), "UTC"}
	event.Body = &OutlookItemBody{"Text", event.Description}
	return
}

func (event *OutlookEvent) SetCalendar(calendar CalendarManager) (err error) {
	switch x := calendar.(type) {
	case *OutlookCalendar:
		event.calendar = x
	default:
		return errors.New(fmt.Sprintf("type of calendar not valid for google: %T", x))
	}
	return
}

func (event *OutlookEvent) CanProcessAgain() bool {
	return event.exponentialBackoff < maxBackoff
}

func (event *OutlookEvent) SetRelations(relations []EventManager) {
	event.relations = relations
}

func (event *OutlookEvent) MarkWrong() {
	//	TODO: Implement marking to db
	log.Fatalf("not implemented yet. ID: %s", event.GetID())
}

func (event *OutlookEvent) IncrementBackoff() {
	event.exponentialBackoff += 1
}

func (event *OutlookEvent) SetState(stateInformed string) (err error) {
	state := states[stateInformed]
	if state == 0 {
		return errors.New(fmt.Sprintf("state %s not supported", stateInformed))
	}
	return
}

func (event *OutlookEvent) GetState() int {
	return event.state
}

func (event *OutlookEvent) extractTime() (err error) {
	format := "2006-01-02T15:04:05.999999999"
	location, err := time.LoadLocation(event.Start.TimeZone)
	if err != nil {
		log.Errorf("error getting start location: %s", event.Start.TimeZone)
		return err
	}
	date, err := time.ParseInLocation(format, event.Start.DateTime, location)
	if err != nil {
		log.Errorf("error parsing start time: %s %s", event.Start.DateTime, err.Error())
		return
	}
	event.StartsAt = date.UTC()

	location, err = time.LoadLocation(event.End.TimeZone)
	if err != nil {
		log.Errorf("error getting end location: %s", event.End.TimeZone)
		return err
	}
	event.EndsAt, err = time.ParseInLocation(format, event.Start.DateTime, location)
	if err != nil {
		log.Errorf("error parsing end time: %s %s", event.Start.DateTime, err.Error())
		return
	}
	event.EndsAt = date.UTC()
	return
}
