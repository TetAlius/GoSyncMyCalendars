package outlook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (calendar *OutlookCalendar) Create(a api.AccountManager) (err error) {
	log.Debugln("createCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	contents, err := util.DoRequest("POST",
		route,
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error creating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	log.Debugf("Response: %s", contents)

	return err
}

func (calendar *OutlookCalendar) Update(a api.AccountManager) error {
	log.Debugln("updateCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling calendar data: %s", err.Error()))
	}

	contents, err := util.DoRequest("PATCH",
		fmt.Sprintf(route, calendar.GetID()),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error updating a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := OutlookCalendarResponse{OdataContext: "", OutlookCalendar: calendar}
	err = json.Unmarshal(contents, &calendarResponse)

	return err

}

func (calendar *OutlookCalendar) Delete(a api.AccountManager) (err error) {
	log.Debugln("deleteCalendar outlook")
	if len(calendar.GetID()) == 0 {
		return errors.New("no ID for calendar was given")
	}

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("DELETE",
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return errors.New(fmt.Sprintf("error deleting a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)
	return

}

func (calendar *OutlookCalendar) GetAllEvents(a api.AccountManager) (events []api.EventManager, err error) {
	log.Debugln("getAllEvents outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, calendar.GetID()),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all events of a calendar for email %s. %s", a.Mail(), err.Error()))
	}

	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)
	eventListResponse := new(OutlookEventListResponse)
	err = json.Unmarshal(contents, &eventListResponse)

	events = make([]api.EventManager, len(eventListResponse.Events))
	for i, s := range eventListResponse.Events {
		s.Calendar = calendar
		events[i] = s
	}
	return
}

// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func (calendar *OutlookCalendar) GetEvent(a api.AccountManager, ID string) (event api.EventManager, err error) {
	log.Debugln("getEvent outlook")
	if len(ID) == 0 {
		return nil, errors.New("an ID for the event must be given")
	}

	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, ID),
		nil,
		a.AuthorizationRequest(),
		a.Mail())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting an event of a calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
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
