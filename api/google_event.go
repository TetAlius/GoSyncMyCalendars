package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"time"

	"reflect"

	conv "github.com/TetAlius/GoSyncMyCalendars/convert"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

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
		log.Errorf("error creating event: %s", err.Error())
		return err
	}

	err = json.Unmarshal(contents, &event)
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

func (event *GoogleEvent) SetState(stateInformed int) {
	event.state = stateInformed
}

func (event *GoogleEvent) GetState() int {
	return event.state
}

func (event *GoogleEvent) SetInternalID(internalID int) {
	event.internalID = internalID
}
func (event *GoogleEvent) GetInternalID() int {
	return event.internalID
}

func (event *GoogleEvent) GetUpdatedAt() (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, event.Updated)
	if err != nil {
		sentryClient().CaptureErrorAndWait(err, map[string]string{"api": "google"})
		return
	}

	return t.UTC(), nil
}

func (date *GoogleTime) UnmarshalJSON(b []byte) error {
	var s map[string]string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	for key, value := range s {
		switch key {
		case "date":
			date.IsAllDay = true
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return err
			}
			date.Date = t.UTC()
		case "dateTime":
			date.IsAllDay = false
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			date.DateTime = t.UTC()
		}
	}
	date.TimeZone = time.UTC

	return nil
}

func (date *GoogleTime) MarshalJSON() ([]byte, error) {
	if date.DateTime.IsZero() && date.Date.IsZero() {
		return bytes.NewBufferString("{}").Bytes(), nil
	}
	var jsonValue string
	var name string
	if date.IsAllDay {
		name = "Date"
		jsonValue = date.Date.UTC().Format("2006-01-02")
	} else {
		name = "DateTime"
		jsonValue = date.DateTime.UTC().Format(time.RFC3339)
	}
	field, ok := reflect.TypeOf(date).Elem().FieldByName(name)
	if !ok {
		return nil, fmt.Errorf("could not retrieve field %s", name)
	}
	tag, _ := parseTag(field.Tag.Get("json"))
	buffer := bytes.NewBufferString("{")
	_, err := buffer.WriteString(fmt.Sprintf(`"%s":"%s"}`, tag, jsonValue))
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (recurrences *GoogleRecurrence) MarshalJSON() ([]byte, error) {
	return bytes.NewBufferString("").Bytes(), nil
}

func (recurrences *GoogleRecurrence) UnmarshalJSON(b []byte) error {
	return nil
}

func (date *GoogleTime) Deconvert() interface{} {
	m := make(map[string]interface{})
	var value time.Time
	if date.IsAllDay {
		value = date.Date.UTC()
	} else {
		value = date.DateTime.UTC()
	}
	field, ok := reflect.TypeOf(date).Elem().FieldByName("DateTime")
	if !ok {
		return nil
	}
	tag, _ := parseTag(field.Tag.Get("convert"))
	m[tag] = value
	field, ok = reflect.TypeOf(date).Elem().FieldByName("IsAllDay")
	if !ok {
		return nil
	}
	tag, _ = parseTag(field.Tag.Get("convert"))
	m[tag] = date.IsAllDay

	field, ok = reflect.TypeOf(date).Elem().FieldByName("TimeZone")
	if !ok {
		return nil
	}
	tag, _ = parseTag(field.Tag.Get("convert"))
	m[tag] = date.TimeZone
	return m
}

func (*GoogleTime) Convert(m interface{}, tag string, opts string) (conv.Converter, error) {
	d := m.(map[string]interface{})

	dateTime, ok := d["dateTime"].(time.Time)
	if !ok {
		return nil, errors.New("incorrect type of field dateTime")
	}
	isAllDay, ok := d["isAllDay"].(bool)
	if !ok {
		return nil, errors.New("incorrect type of field isAllDay")
	}
	timeZone, ok := d["timeZone"].(*time.Location)
	if !ok {
		return nil, errors.New("incorrect type of field timeZone")
	}

	return &GoogleTime{DateTime: dateTime, Date: dateTime, TimeZone: timeZone, IsAllDay: isAllDay}, nil
}

func (event *GoogleEvent) setAllDay() {
	if event.Start == nil && event.End == nil {
		event.IsAllDay = false
		return
	}
	event.IsAllDay = event.Start.IsAllDay
}
