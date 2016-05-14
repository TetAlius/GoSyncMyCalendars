package google

import (
	"bytes"
	"fmt"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/logs/logger"
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

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events
func getAllEvents(calendarID string) {
	log.Debugln("getAllEvents google")

	contents, _ := backend.NewRequest("GET",
		Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/"+calendarID+Requests.Events,
		nil,
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)

}

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events
func createEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent google")
	contents, _ := backend.NewRequest("POST",
		//Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/"+calendarID+Requests.Events,
		"https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events",
		bytes.NewBuffer(event),
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)

}

//PUT https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events/{eventId}
func updateEvent(calendarID string, eventID string, eventData []byte) {
	log.Debugln("updateEvent google")

	//Meter en los header el etag

	contents, _ := backend.NewRequest("PUT",
		//Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/"+calendarID+Requests.Events,
		"https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events/"+eventID,
		bytes.NewBuffer(eventUpdated),
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)

}

//DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events/{eventId}
func deleteEvent(calendarID string, eventID string) {
	log.Debugln("deleteEvent google")

	contents, _ := backend.NewRequest(
		"DELETE",
		"https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events/"+eventID,
		nil,
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)

}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events/{eventId}
func getEvent(calendarID string, eventID string) {
	log.Debugln("getEvent google")

	contents, _ := backend.NewRequest(
		"GET",
		"https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events/"+eventID,
		nil,
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)

}
