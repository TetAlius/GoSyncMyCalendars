package google

import (
	"bytes"

	"fmt"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/pkg/errors"
)

type CalendarListResponse struct {
	Kind          string      `json:"kind"`
	Etag          string      `json:"etag"`
	NextPageToken string      `json:"nextPageToken"`
	NextSyncToken string      `json:"nextSyncToken"`
	Calendars     []*Calendar `json:"items"`
}

type Calendar struct {
	//From CalendarLIST resource
	ID              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description,omitempty"`
	Location        string `json:"location,omitempty"`
	TmeZone         string `json:"timeZone,omitempty"`
	SummaryOverride string `json:"summaryOverride,omitempty"`
	ColorId         string `json:"colorId,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	ForegroundColor string `json:"foregroundColor,omitempty"`
	Hidden          bool   `json:"hidden,omitempty"`
	Selected        bool   `json:"selected,omitempty"`
	// Only valid accessRoles with 'writer' or 'owner'
	AccessRole           string     `json:"accessRole,omitempty"`
	DefaultReminders     []Reminder `json:"defaultReminders,omitempty"`
	Primary              bool       `json:"primary,omitempty"`
	Deleted              bool       `json:"deleted,omitempty"`
	ConferenceProperties `json:"conferenceProperties,omitempty"`
	NotificationSetting  `json:"notificationSettings,omitempty"`
}

type Reminder struct {
	Method  string `json:"method,omitempty"`
	Minutes int32  `json:"minutes,omitempty"`
}

type NotificationSetting struct {
	Notifications []Notification `json:"notifications,omitempty"`
}

type Notification struct {
	Type   string `json:"type,omitempty"`
	Method string `json:"method,omitempty"`
}

type ConferenceProperties struct {
	AllowedConferenceSolutionTypes []string `json:"allowedConferenceSolutionTypes,omitempty"`
}

// PUT https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *Calendar) Update(a api.AccountManager) (err error) {
	log.Debugln("updateCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	contents, err :=
		util.DoRequest(
			"PUT",
			fmt.Sprintf(route, calendar.ID),
			bytes.NewBuffer(data),
			a.AuthorizationRequest(),
			"")

	if err != nil {
		return errors.New(fmt.Sprintf("Error updating a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	err = createResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("Contents: %s", contents)
	err = json.Unmarshal(contents, &calendar)

	return
}

// DELETE https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (calendar *Calendar) Delete(a api.AccountManager) (err error) {
	log.Debugln("Delete calendar")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest(
		"DELETE",
		fmt.Sprintf(route, calendar.ID),
		nil,
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	if len(contents) != 0 {
		return errors.New(fmt.Sprintf("error deleting google calendar %s: %s", calendar.ID, contents))
	}

	log.Debugf("Contents: %s", contents)

	return
}

// POST https://www.googleapis.com/calendar/v3/calendars
func (calendar *Calendar) Create(a api.AccountManager) (err error) {
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
	contents, err :=
		util.DoRequest(
			"POST",
			route,
			bytes.NewBuffer(data),
			a.AuthorizationRequest(),
			"")

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("Contents: %s", contents)
	err = json.Unmarshal(contents, &calendar)
	return
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (calendar *Calendar) GetAllEvents(a api.AccountManager) (events []api.EventManager, err error) {
	log.Debugln("getAllEvents google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, calendar.ID),
		nil,
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting all events of g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventList := new(EventList)
	err = json.Unmarshal(contents, &eventList)

	events = make([]api.EventManager, len(eventList.Events))
	for i, s := range eventList.Events {
		s.Calendar = calendar
		events[i] = s
	}
	return
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (calendar *Calendar) GetEvent(a api.AccountManager, ID string) (event api.EventManager, err error) {
	log.Debugln("getEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest(
		"GET",
		fmt.Sprintf(route, calendar.ID, ID),
		nil,
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of g calendar for email %s. %s", a.Mail(), err.Error()))
	}

	log.Debugf("Contents: %s", contents)
	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventResponse := new(Event)
	err = json.Unmarshal(contents, &eventResponse)

	eventResponse.Calendar = calendar
	event = eventResponse

	return
}

func (calendar Calendar) GetID() string {
	return calendar.ID
}
