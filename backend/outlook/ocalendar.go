package outlook

import (
	"bytes"
	"fmt"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// GET https://outlook.office.com/api/v2.0/me/calendars
func (o *OutlookAccount) GetAllCalendars() {
	log.Debugln("getAllCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
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

}

//GET https://outlook.office.com/api/v2.0/me/calendar
func (o *OutlookAccount) GetPrimaryCalendar() {
	log.Debugln("getPrimaryCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/primary")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
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
}

// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (o *OutlookAccount) GetCalendar(calendarID string) {
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
}

// POST https://outlook.office.com/api/v2.0/me/calendars
func (o *OutlookAccount) CreateCalendar(calendarData []byte) {
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

}

// PATCH https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func (o *OutlookAccount) UpdateCalendar(calendarID string, calendarData []byte) {
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
func (o *OutlookAccount) DeleteCalendar(calendarID string) {
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
