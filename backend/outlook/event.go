package outlook

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func getAllEvents() {
	log.Debugln("getAllEvents outlook")

	contents, _ := backend.NewRequest("GET",
		OutlookRequests.RootURI+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Events,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

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

func createEvent(calendarID string, eventData []byte) {
	log.Debugln("createEvent outlook")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents, _ := backend.NewRequest("POST",
		OutlookRequests.RootURI+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars+"/"+
			calendarID+
			OutlookRequests.Events,
		bytes.NewBuffer(event),
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

var update = []byte(`{
  "Location": {
    "DisplayName": "Your office"
  }
}`)

func updateEvent(eventID string, eventData []byte) {
	log.Debugln("updateEvent outlook")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents, _ := backend.NewRequest("PATCH",
		OutlookRequests.RootURI+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Events+"/"+
			eventID,
		bytes.NewBuffer(update),
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}

func deleteEvent(eventID string) {
	log.Debugln("deleteEvent outlook")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents, _ := backend.NewRequest("DELETE",
		OutlookRequests.RootURI+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Events+"/"+
			eventID,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

func getEvent(eventID string) {
	log.Debugln("getEvent outlook")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents, _ := backend.NewRequest("GET",
		OutlookRequests.RootURI+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Events+"/"+
			eventID,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}
