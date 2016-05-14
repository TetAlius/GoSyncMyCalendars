package google

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/logs/logger"
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

	contents, _ :=
		backend.NewRequest(
			"GET",
			Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Context+Requests.CalendarList,
			nil,
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)
}

// GET https://www.googleapis.com/calendar/v3/calendars/primary This is the one used
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/primary
func getPrimaryCalendar() {
	log.Debugln("getPrimaryCalendar google")
	contents, _ :=
		backend.NewRequest(
			"GET",
			Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/primary",
			nil,
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)

}

// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func getCalendar(calendarID string) {
	log.Debugln("getCalendar google")

	contents, _ :=
		backend.NewRequest(
			"GET",
			Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Context+Requests.CalendarList+"/"+calendarID,
			nil,
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)

}

var calendar = []byte(`{"summary":"CalendarGO"}`)

var calendarUpdate = []byte(`{"summary":"Updated CalendarGO"}`)

// POST https://www.googleapis.com/calendar/v3/calendars
func createCalendar(calendarData []byte) {
	log.Debugln("createCalendar google")

	contents, _ :=
		backend.NewRequest(
			"POST",
			Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars,
			bytes.NewBuffer(calendarData),
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)

}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func updateCalendar(calendarID string, calendarData []byte) {
	log.Debugln("updateCalendar google")

	contents, _ :=
		backend.NewRequest(
			"PUT",
			Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/"+calendarID,
			bytes.NewBuffer(calendarData),
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)

}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func deleteCalendar(calendarID string) {
	fmt.Println("Delete calendar")

	contents, _ := backend.NewRequest(
		"DELETE",
		Requests.RootURI+Requests.CalendarAPI+Requests.Version+Requests.Calendars+"/"+calendarID,
		nil,
		Responses.TokenType+" "+Responses.AccessToken,
		"")

	fmt.Printf("%s\n", contents)
}
