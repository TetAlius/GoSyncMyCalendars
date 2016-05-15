package google

import (
	"bytes"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
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
func getAllEvents(calendarID string) {
	log.Debugln("getAllEvents google")

	contents, _ := backend.NewRequest("GET",
		eventsURI(calendarID, ""),
		nil,
		authorizationRequest(),
		"")

	log.Debugf("Contents: %s", contents)

}

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func createEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent google")
	contents, _ := backend.NewRequest("POST",
		eventsURI(calendarID, ""),
		bytes.NewBuffer(event),
		authorizationRequest(),
		"")

	log.Debugf("Contents: %s", contents)

}

//PUT https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func updateEvent(calendarID string, eventID string, eventData []byte) {
	log.Debugln("updateEvent google")

	//Meter en los header el etag

	contents, _ := backend.NewRequest("PUT",
		eventsURI(calendarID, eventID),
		bytes.NewBuffer(eventUpdated),
		authorizationRequest(),
		"")

	log.Debugf("Contents: %s", contents)

}

//DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func deleteEvent(calendarID string, eventID string) {
	log.Debugln("deleteEvent google")

	contents, _ := backend.NewRequest(
		"DELETE",
		eventsURI(calendarID, eventID),
		nil,
		authorizationRequest(),
		"")

	log.Debugf("Contents: %s", contents)

}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func getEvent(calendarID string, eventID string) {
	log.Debugln("getEvent google")

	contents, _ := backend.NewRequest(
		"GET",
		eventsURI(calendarID, eventID),
		nil,
		authorizationRequest(),
		"")

	log.Debugf("Contents: %s", contents)

}
