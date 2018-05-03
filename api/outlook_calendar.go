package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (calendar *OutlookCalendar) Create(a AccountManager) (err error) {
	log.Debugln("createCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()
	contents, err := util.DoRequest(http.MethodPost,
		route,
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	log.Debugf("Response: %s", contents)

	return err
}

func (calendar *OutlookCalendar) Update(a AccountManager) error {
	log.Debugln("updateCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPatch,
		fmt.Sprintf(route, calendar.GetID()),
		bytes.NewBuffer(data),
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	return err

}

func (calendar *OutlookCalendar) Delete(a AccountManager) (err error) {
	log.Debugln("deleteCalendar outlook")
	if len(calendar.GetID()) == 0 {
		return errors.New("no ID for calendar was given")
	}

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodDelete,
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		headers, nil)

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)
	return

}

func (calendar *OutlookCalendar) GetAllEvents(a AccountManager) (events []EventManager, err error) {
	log.Debugln("getAllEvents outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()
	headers["Prefer"] = "outlook.timezone=UTC"

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventListResponse := new(OutlookEventListResponse)
	err = json.Unmarshal(contents, &eventListResponse)

	events = make([]EventManager, len(eventListResponse.Events))
	for i, s := range eventListResponse.Events {
		s.Calendar = calendar
		events[i] = s
	}
	return
}

// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func (calendar *OutlookCalendar) GetEvent(a AccountManager, ID string) (event EventManager, err error) {
	log.Debugln("getEvent outlook")
	if len(ID) == 0 {
		return nil, errors.New("an ID for the event must be given")
	}

	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()
	headers["Prefer"] = "outlook.timezone=UTC"

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, ID),
		nil,
		headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventResponse := new(OutlookEventResponse)
	err = json.Unmarshal(contents, &eventResponse)

	eventResponse.Calendar = calendar
	event = eventResponse.OutlookEvent

	return
}
func (calendar *OutlookCalendar) GetID() string {
	return calendar.ID
}
