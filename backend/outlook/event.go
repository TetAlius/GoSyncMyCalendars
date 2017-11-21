package outlook

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

// GET https://outlook.office.com/api/v2.0/me/events
// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func (o *Account) GetAllEventsFromCalendar(calendarID string) (events []EventInfo, err error) {
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

	err = json.Unmarshal(contents, &events)
	return
}

// POST https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func (o *Account) CreateEvent(calendarID string, eventData []byte) (event EventInfo, err error) {
	log.Debugln("createEvent outlook")
	route, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("POST",
		fmt.Sprintf(route, calendarID),
		bytes.NewBuffer(eventData),
		o.authorizationRequest(),
		o.AnchorMailbox)

	if err != nil {
		log.Errorf("Error creating event in a calendar for email %s. %s", o.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

	err = json.Unmarshal(contents, &event)
	return
}

var update = []byte(`{
  "Location": {
    "DisplayName": "Your office"
  }
}`)

// PATCH https://outlook.office.com/api/v2.0/me/events/{eventID}
func (o *Account) UpdateEvent(eventData []byte, ids ...string) {
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
func (o *Account) DeleteEvent(ids ...string) {
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
func (o *Account) GetEvent(ids ...string) {
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
