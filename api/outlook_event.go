package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	"time"

	"reflect"

	conv "github.com/TetAlius/GoSyncMyCalendars/convert"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// Method that creates the event
//
// POST https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
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
	if err != nil {
		return err
	}
	event.setAllDay()
	return
}

// Method that updates the event
//
// PATCH https://outlook.office.com/api/v2.0/me/events/{eventID}
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
	if err != nil {
		return err
	}
	event.setAllDay()
	return
}

// Method that deletes the event
//
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

// Method that returns the ID of the event
func (event *OutlookEvent) GetID() string {
	return event.ID
}

// Method that returns the calendar which have this event
func (event *OutlookEvent) GetCalendar() CalendarManager {
	return event.calendar
}

// Method that returns the syncing events with this
func (event *OutlookEvent) GetRelations() []EventManager {
	return event.relations
}

// Method that returns the calendar which have this event
func (event *OutlookEvent) SetCalendar(calendar CalendarManager) (err error) {
	switch x := calendar.(type) {
	case *OutlookCalendar:
		event.calendar = x
	default:
		return errors.New(fmt.Sprintf("type of calendar not valid for outlook: %T", x))
	}
	return
}

// Method that checks if the event can try sync again
func (event *OutlookEvent) CanProcessAgain() bool {
	return event.exponentialBackoff < maxBackoff
}

// Method that sets the events syncing with this
func (event *OutlookEvent) SetRelations(relations []EventManager) {
	event.relations = relations
}

// Method that increments the number of failed attempts to sync
func (event *OutlookEvent) IncrementBackoff() {
	event.exponentialBackoff += 1
}

// Method that sets the state of the event
func (event *OutlookEvent) SetState(stateInformed int) {
	event.state = stateInformed
}

// Method that sets the internal ID generated on db
func (event *OutlookEvent) SetInternalID(internalID int) {
	event.internalID = internalID
}

// Method that gets the internal ID of the event
func (event *OutlookEvent) GetInternalID() int {
	return event.internalID
}

// Method that returns the state of the event
func (event *OutlookEvent) GetState() int {
	return event.state
}

// Method that returns the last update date
func (event *OutlookEvent) GetUpdatedAt() (t time.Time, err error) {
	//format := time.RFC3339Nano
	//format := "2006-01-02T15:04:05.999999999"
	t, err = time.Parse(time.RFC3339, event.LastModifiedDateTime)
	if err != nil {
		sentryClient().CaptureErrorAndWait(err, map[string]string{"api": "outlook"})
		return
	}
	return t.UTC(), nil
}

// Method that converts from a JSON to a OutlookDateTimeTimeZone struct.
// This method implements Unmarshaler interface
func (date *OutlookDateTimeTimeZone) UnmarshalJSON(b []byte) error {
	var s map[string]string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	field, ok := reflect.TypeOf(date).Elem().FieldByName("TimeZone")
	if !ok {
		return fmt.Errorf("could not retrieve field TimeZone")
	}
	tag, _ := parseTag(field.Tag.Get("json"))

	location, err := time.LoadLocation(s[tag])
	date.TimeZone = location
	if err != nil {
		log.Errorf("error getting location: %s")
		return err
	}

	field, ok = reflect.TypeOf(date).Elem().FieldByName("DateTime")
	if !ok {
		return fmt.Errorf("could not retrieve field DateTime")
	}
	tag, _ = parseTag(field.Tag.Get("json"))
	t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999", s[tag], location)
	if err != nil {
		log.Errorf("error parsing time: %s", err.Error())
	}
	date.DateTime = t.UTC()
	return nil
}

// Method that converts from OutlookDateTimeTimeZone struct to a json.
// This method implements Marshaler interface
func (date *OutlookDateTimeTimeZone) MarshalJSON() (b []byte, err error) {
	if date.DateTime.IsZero() {
		return bytes.NewBufferString("{}").Bytes(), nil
	}
	field, ok := reflect.TypeOf(date).Elem().FieldByName("DateTime")
	if !ok {
		return nil, fmt.Errorf("could not retrieve field DateTime")
	}
	tag, _ := parseTag(field.Tag.Get("json"))
	buffer := bytes.NewBufferString("{")
	//RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	_, err = buffer.WriteString(fmt.Sprintf(`"%s":"%s"`, tag, date.DateTime.UTC().Format("2006-01-02T15:04:05.999999999")))
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(",")
	if err != nil {
		return nil, err
	}

	field, ok = reflect.TypeOf(date).Elem().FieldByName("TimeZone")
	if !ok {
		return nil, fmt.Errorf("could not retrieve field TimeZone")
	}
	tag, _ = parseTag(field.Tag.Get("json"))
	_, err = buffer.WriteString(fmt.Sprintf(`"%s":"%s"`, tag, date.TimeZone))
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteString("}")
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Method that converts a OutlookDateTimeTimeZone struct to a interface{}.
// This method implements Deconverter interface
func (date *OutlookDateTimeTimeZone) Deconvert() interface{} {
	m := make(map[string]interface{})
	field, ok := reflect.TypeOf(date).Elem().FieldByName("DateTime")
	if !ok {
		return nil
	}
	tag, _ := parseTag(field.Tag.Get("convert"))
	m[tag] = date.DateTime.UTC()

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

// Method that converts a OutlookItemBody struct to a interface{}.
// This method implements Deconverter interface
func (body *OutlookItemBody) Deconvert() interface{} {
	return body.Description
}

// Method that converts an interface{} ti a OutlookItemBody struct.
// This method implements Converter interface
func (*OutlookItemBody) Convert(m interface{}, tag string, opts string) (conv.Converter, error) {
	desc, ok := m.(string)
	if !ok {
		return nil, errors.New("incorrect type of field description")
	}
	return &OutlookItemBody{Description: desc}, nil
}

// Method that converts an interface{} ti a OutlookDateTimeTimeZone struct.
// This method implements Converter interface
func (*OutlookDateTimeTimeZone) Convert(m interface{}, tag string, opts string) (conv.Converter, error) {
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

	return &OutlookDateTimeTimeZone{DateTime: dateTime, TimeZone: timeZone, IsAllDay: isAllDay}, nil
}

// Method that sets all day to the necessary attributes
func (event *OutlookEvent) setAllDay() {
	event.Start.IsAllDay = event.IsAllDay
	event.End.IsAllDay = event.IsAllDay
}
