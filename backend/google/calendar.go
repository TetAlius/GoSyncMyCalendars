package google

import (
	"bytes"
	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
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
func (g *Account) GetAllCalendars() {
	log.Debugln("getAllCalendars google")
	route, err := util.CallAPIRoot("google/calendars")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err :=
		util.DoRequest(
			"GET",
			route,
			nil,
			g.authorizationRequest(),
			"")
	if err != nil {
		log.Errorf("Error getting all calendars for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)
}

// GET https://www.googleapis.com/calendar/v3/calendars/primary
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/primary This is the one used
func (g *Account) GetPrimaryCalendar() {
	log.Debugln("getPrimaryCalendar google")
	route, err := util.CallAPIRoot("google/calendars/primary")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err :=
		util.DoRequest(
			"GET",
			route,
			nil,
			g.authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error getting primary calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func (g *Account) GetCalendar(calendarID string) {
	log.Debugln("getCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err :=
		util.DoRequest(
			"GET",
			fmt.Sprintf(route, calendarID),
			nil,
			g.authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error getting a calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

var calendar = []byte(`{"summary":"CalendarGO"}`)

var calendarUpdate = []byte(`{"summary":"Updated CalendarGO"}`)

// POST https://www.googleapis.com/calendar/v3/calendars
func (g *Account) CreateCalendar(calendarData []byte) {
	log.Debugln("createCalendar google")
	route, err := util.CallAPIRoot("google/calendars")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err :=
		util.DoRequest(
			"POST",
			route,
			bytes.NewBuffer(calendarData),
			g.authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error creating a calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func (g *Account) UpdateCalendar(calendarID string, calendarData []byte) {
	log.Debugln("updateCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err :=
		util.DoRequest(
			"PUT",
			fmt.Sprintf(route, calendarID),
			bytes.NewBuffer(calendarData),
			g.authorizationRequest(),
			"")

	if err != nil {
		log.Errorf("Error updating a calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)

}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarId}
func (g *Account) DeleteCalendar(calendarID string) {
	log.Debugln("Delete calendar")
	route, err := util.CallAPIRoot("google/calendars/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest(
		"DELETE",
		fmt.Sprintf(route, calendarID),
		nil,
		g.authorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error deleting a calendar for email %s. %s", g.Email, err.Error())
	}

	log.Debugf("Contents: %s", contents)
}
