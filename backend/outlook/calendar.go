package outlook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

var calendar2 = []byte(`{
  "Name": "Social"contents
}`)

// GET https://outlook.office.com/api/v2.0/me/calendars
func (o *Account) GetAllCalendars() (calendars []CalendarInfo, err error) {
	log.Debugln("getAllCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		log.Errorf("%s", err.Error())
		return calendars, errors.New(fmt.Sprintf("Error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		route,
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error getting all calendars for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarListResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	return calendarResponse.Calendars, err

}

//GET https://outlook.office.com/api/v2.0/me/calendar
func (o *Account) GetPrimaryCalendar() (calendar CalendarInfo, err error) {
	log.Debugln("getPrimaryCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/primary")
	if err != nil {
		log.Errorf("%s", err.Error())
		return calendar, errors.New(fmt.Sprintf("Error generating URL: %s", err.Error()))
	}

	contents, err := util.DoRequest("GET",
		route,
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("%s", err.Error())
		return calendar, errors.New(fmt.Sprintf("Error getting primary calendar for email %s. %s", o.AnchorMailbox, err.Error()))
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	return calendarResponse.CalendarInfo, err
}

// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (o *Account) GetCalendar(calendarID string) (calendar CalendarInfo, err error) {
	log.Debugln("getCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, calendarID),
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error getting a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	return calendarResponse.CalendarInfo, err
}

// POST https://outlook.office.com/api/v2.0/me/calendars
func (o *Account) CreateCalendar(calendarData []byte) (calendar CalendarInfo, err error) {
	log.Debugln("createCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("POST",
		route,
		bytes.NewBuffer(calendarData),
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error creating a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	return calendarResponse.CalendarInfo, err

}

// PATCH https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (o *Account) UpdateCalendar(calendarID string, calendarData []byte) (calendar CalendarInfo, err error) {
	log.Debugln("updateCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("PATCH",
		fmt.Sprintf(route, calendarID),
		bytes.NewBuffer(calendarData),
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error updateing a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	return calendarResponse.CalendarInfo, err
}

//TODO check if calendar is primary or birthdays if it is, the following error is send
//
//{
//	"error": {
//		"code": "ErrorInvalidRequest",
//		"message": "Your request can't be completed. The default calendar cannot be deleted."
//	}
//}

// DELETE https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
//Does not return json if OK, only status 204
func (o *Account) DeleteCalendar(calendarID string) {
	log.Debugln("deleteCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("DELETE",
		fmt.Sprintf(route, calendarID),
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error deleting a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)
}
