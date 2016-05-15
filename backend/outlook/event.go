package outlook

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// GET https://outlook.office.com/api/v2.0/me/events
// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func getAllEvents(calendarID string) {
	log.Debugln("getAllEvents outlook")
	contents, _ := backend.NewRequest("GET",
		eventsFromCalendarURI(calendarID),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	fmt.Printf("%s\n", contents)
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
func createEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent outlook")
	contents, _ := backend.NewRequest("POST",
		eventsFromCalendarURI(calendarID),
		bytes.NewBuffer(event),
		authorizationRequest(),
		Responses.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

var update = []byte(`{
  "Location": {
    "DisplayName": "Your office"
  }
}`)

// PATCH https://outlook.office.com/api/v2.0/me/events/{eventID}
func updateEvent(eventID string, eventData []byte) {
	log.Debugln("updateEvent outlook")
	contents, _ := backend.NewRequest("PATCH",
		eventURI(eventID),
		bytes.NewBuffer(update),
		authorizationRequest(),
		Responses.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}

// DELETE https://outlook.office.com/api/v2.0/me/events/{eventID}
func deleteEvent(eventID string) {
	log.Debugln("deleteEvent outlook")
	contents, _ := backend.NewRequest("DELETE",
		eventURI(eventID),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

// GET https://outlook.office.com/api/v2.0/me/events/{eventID}
func getEvent(eventID string) {
	log.Debugln("getEvent outlook")
	contents, _ := backend.NewRequest("GET",
		eventURI(eventID),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}
