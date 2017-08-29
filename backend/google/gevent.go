package google

import (
	"bytes"
	"fmt"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"time"
)

type eventResource struct {
	end         timeCalendar
	start       timeCalendar
	attendees   []atendee
	attachments []attachment
	reminders   reminder
}

type reminder struct {
	UseDefault bool       `json:"useDefault"`
	Overrides  []override `json:"overrides"`
}

type override struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

type atendee struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	DisplayName      string `json:"displayName"`
	Organizer        bool   `json:"organizer"`
	Self             bool   `json:"self"`
	Resource         bool   `json:"resource"`
	Optional         bool   `json:"optional"`
	ResponseStatus   string `json:"responseStatus"`
	Comment          string `json:"comment"`
	AdditionalGuests int    `json:"additionalGuests"`
}

type attachment struct {
	FileURL  string `json:"fileUrl"`
	Title    string `json:"title"`
	MimeType string `json:"mimeType"`
	IconLink string `json:"iconLink"`
	FileID   string `json:"fileID"`
}

type person struct {
	id          string
	email       string
	displayName string
	self        bool
}
type timeCalendar struct {
	date     time.Time
	dateTime time.Time
	timeZone string
}

var event = []byte(`{
"summary":"Prueba desde Go",
  "end":
  {
    "dateTime": "2016-04-12T17:00:00-07:00",
    "timeZone": "America/Los_Angeles"
  },
  "start":
  {
    "dateTime": "2016-04-12T09:00:00-07:00",
    "timeZone": "America/Los_Angeles"
  }
}`)

var eventUpdated = []byte(`{
"summary":"Update desde Go",
  "end":
  {
    "dateTime": "2016-04-11T17:00:00-07:00",
    "timeZone": "America/Los_Angeles"
  },
  "start":
  {
    "dateTime": "2016-04-11T09:00:00-07:00",
    "timeZone": "America/Los_Angeles"
  }
}`)

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (g *GoogleAccount) GetAllEventsFromCalendar(calendarID string) {
	log.Debugln("getAllEvents google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("GET",
		fmt.Sprintf(route, calendarID),
		nil,
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error getting all events of g calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (g *GoogleAccount) CreateEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("POST",
		fmt.Sprintf(route, calendarID),
		bytes.NewBuffer(event),
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error creating event in g calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

//PUT https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (g *GoogleAccount) UpdateEvent(calendarID string, eventID string, eventData []byte) {
	log.Debugln("updateEvent google")

	//Meter en los header el etag
	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest("PUT",
		fmt.Sprintf(route, calendarID, eventID),
		bytes.NewBuffer(eventUpdated),
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error updating event of g calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

//DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (g *GoogleAccount) DeleteEvent(calendarID string, eventID string) {
	log.Debugln("deleteEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest(
		"DELETE",
		fmt.Sprintf(route, calendarID, eventID),
		nil,
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error deleting event of g calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (g *GoogleAccount) GetEvent(calendarID string, eventID string) {
	log.Debugln("getEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest(
		"GET",
		fmt.Sprintf(route, calendarID, eventID),
		nil,
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error getting an event of g calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}
