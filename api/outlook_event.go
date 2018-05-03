package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/pkg/errors"
)

func (event *OutlookEvent) Create(a AccountManager) (err error) {
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
		fmt.Sprintf(route, event.Calendar.GetID()),
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

	log.Debugf("Response: %s", contents)
	return
}

func (event *OutlookEvent) Update(a AccountManager) (err error) {
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

	log.Debugf("Response: %s", contents)
	return
}

// DELETE https://outlook.office.com/api/v2.0/me/events/{eventID}
func (event *OutlookEvent) Delete(a AccountManager) (err error) {
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

	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	log.Debugf("%s\n", contents)
	return
}

func (event *OutlookEvent) GetID() string {
	return event.ID
}

func (event *OutlookEvent) GetCalendar() CalendarManager {
	return event.Calendar
}

func (event *OutlookEvent) PrepareTime() (err error) {
	panic("IMPLEMENT ME")
	return
}
