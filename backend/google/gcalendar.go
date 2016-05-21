package google

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

type calendarResp struct {
	Kind        string `json:"kind"`
	Etag        string `json:"etag"`
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	TimeZone    string `json:"timeZone"`
}

//GET https://www.googleapis.com/calendar/v3/users/me/calendarList
func getAllCalendars() {
	log.Debugln("getAllCalendars google")

	contents, err :=
		backend.NewRequest(
			"GET",
			calendarListURI(""),
			nil,
			authorizationRequest(),
			"")
	if err != nil {
		log.Errorf("Error getting all calendars for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)
}

// GET https://www.googleapis.com/calendar/v3/calendars/primary This is the one used
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/primary
func getPrimaryCalendar() {
	log.Debugln("getPrimaryCalendar google")
	contents, err :=
		backend.NewRequest(
			"GET",
			calendarsURI("primary"),
			nil,
			authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error getting primary calendar for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func getCalendar(calendarID string) {
	log.Debugln("getCalendar google")

	contents, err :=
		backend.NewRequest(
			"GET",
			calendarListURI(calendarID),
			nil,
			authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error getting a calendar for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

var calendar = []byte(`{"summary":"CalendarGO"}`)

var calendarUpdate = []byte(`{"summary":"Updated CalendarGO"}`)

// POST https://www.googleapis.com/calendar/v3/calendars
func createCalendar(calendarData []byte) {
	log.Debugln("createCalendar google")

	contents, err :=
		backend.NewRequest(
			"POST",
			calendarsURI(""),
			bytes.NewBuffer(calendarData),
			authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error creating a calendar for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func updateCalendar(calendarID string, calendarData []byte) {
	log.Debugln("updateCalendar google")

	contents, err :=
		backend.NewRequest(
			"PUT",
			calendarsURI(calendarID),
			bytes.NewBuffer(calendarData),
			authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error updateing a calendar for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func deleteCalendar(calendarID string) {
	fmt.Println("Delete calendar")

	contents, err := backend.NewRequest(
		"DELETE",
		calendarsURI(calendarID),
		nil,
		authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error deleting a calendar for email %s. %s", Responses.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)
}
