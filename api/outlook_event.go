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

func (event *OutlookEvent) SetState(stateInformed int) {
	event.state = stateInformed
}

func (event *OutlookEvent) SetInternalID(internalID int) {
	event.internalID = internalID
}

func (event *OutlookEvent) GetInternalID() int {
	return event.internalID
}

func (event *OutlookEvent) GetState() int {
	return event.state
}

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

func (body *OutlookItemBody) Deconvert() interface{} {
	return body.Description
}

func (*OutlookItemBody) Convert(m interface{}, tag string, opts string) (conv.Converter, error) {
	desc, ok := m.(string)
	if !ok {
		return nil, errors.New("incorrect type of field description")
	}
	return &OutlookItemBody{Description: desc}, nil
}

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

func (event *OutlookEvent) setAllDay() {
	event.Start.IsAllDay = event.IsAllDay
	event.End.IsAllDay = event.IsAllDay
}
