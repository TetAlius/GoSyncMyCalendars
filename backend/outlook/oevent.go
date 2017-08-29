package outlook

import (
	"bytes"
	"fmt"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// GET https://outlook.office.com/api/v2.0/me/events
// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func (o *OutlookAccount) GetAllEventsFromCalendar(calendarID string) {
	log.Debugln("getAllEvents outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
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
		log.Errorf("Error getting all events of a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)
}

//TODO: delete this
var event = []byte(`{
  "Subject": "Discuss the Calendar REST API",
  "Body": {
    "ContentType": "HTML",
    "Content": "I think it will meet our requirements!"
  },
  "Start": {
      "DateTime": "2016-02-02T18:00:00",
      "TimeZone": "Pacific Standard Time"
  },
  "End": {
      "DateTime": "2016-02-02T19:00:00",
      "TimeZone": "Pacific Standard Time"
  },
	"ReminderMinutesBeforeStart": "30",
  "IsReminderOn": "false"
}`)

// POST https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func (o *OutlookAccount) CreateEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("POST",
		fmt.Sprintf(route, calendarID),
		bytes.NewBuffer(event),
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error creating event in a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)
}

var update = []byte(`{
  "Location": {
    "DisplayName": "Your office"
  }
}`)

// PATCH https://outlook.office.com/api/v2.0/me/events/{eventID}
func (o *OutlookAccount) UpdateEvent(eventData []byte, ids ...string) {
	log.Debugln("updateEvent outlook")
	//TODO: Test if ids are two given
	route, err := util.CallAPIRoot("outlook/events/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	contents, err := util.DoRequest("PATCH",
		fmt.Sprintf(route, ids[0]),
		bytes.NewBuffer(update),
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error updating event of a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

}

// DELETE https://outlook.office.com/api/v2.0/me/events/{eventID}
func (o *OutlookAccount) DeleteEvent(ids ...string) {
	log.Debugln("deleteEvent outlook")
	//TODO: Test if ids are two given
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	contents, err := util.DoRequest("DELETE",
		fmt.Sprintf(route, ids[0]),
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error deleting event of a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)
}

// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func (o *OutlookAccount) GetEvent(ids ...string) {
	log.Debugln("getEvent outlook")
	//TODO: Test if ids are one given

	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, ids[0]),
		nil,
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error getting an event of a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

}
