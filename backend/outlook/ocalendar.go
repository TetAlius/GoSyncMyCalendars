package outlook

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// GET https://outlook.office.com/api/v2.0/me/calendars
func getAllCalendars() {
	log.Debugln("getAllCalendars outlook")

	contents, err := backend.NewRequest("GET",
		calendarsURI(""),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	if err != nil {
		log.Errorf("Error getting all calendars for email %s. %s", Responses.AnchorMailbox, err.Error())
	}

	log.Debugf("%s\n", contents)

}

//GET https://outlook.office.com/api/v2.0/me/calendar
func getPrimaryCalendar() {
	//TODO
}

// GET https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func getCalendar(calendarID string) {
	log.Debugln("getCalendar outlook")
	contents, err := backend.NewRequest("GET",
		calendarsURI(calendarID),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	if err != nil {
		log.Errorf("Error getting a calendar for email %s. %s", Responses.AnchorMailbox, err.Error())
	}

	fmt.Printf("%s\n", contents)
}

// POST https://outlook.office.com/api/v2.0/me/calendars
func createCalendar(calendarData []byte) {
	log.Debugln("createCalendars outlook")

	contents, err := backend.NewRequest("POST",
		calendarsURI(""),
		bytes.NewBuffer(calendarData),
		authorizationRequest(),
		Responses.AnchorMailbox)

	if err != nil {
		log.Errorf("Error creating a calendar for email %s. %s", Responses.AnchorMailbox, err.Error())
	}

	fmt.Printf("%s\n", contents)

}

// PATCH https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func updateCalendar(calendarID string, calendarData []byte) {
	log.Debugln("updateCalendar outlook")

	contents, err := backend.NewRequest("PATCH",
		calendarsURI(calendarID),
		bytes.NewBuffer(calendarData),
		authorizationRequest(),
		Responses.AnchorMailbox)

	if err != nil {
		log.Errorf("Error updateing a calendar for email %s. %s", Responses.AnchorMailbox, err.Error())
	}

	fmt.Printf("%s\n", contents)
}

//TODO check if calendar is primary or birthdays if it is, the following error is send
/*
{
	"error": {
		"code": "ErrorInvalidRequest",
		"message": "Your request can't be completed. The default calendar cannot be deleted."
	}
}
*/
// DELETE https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
//Does not return json if OK, only status 204
func deleteCalendar(calendarID string) {
	log.Debugln("deleteCalendar outlook")

	contents, err := backend.NewRequest("DELETE",
		calendarsURI(calendarID),
		nil,
		authorizationRequest(),
		Responses.AnchorMailbox)

	if err != nil {
		log.Errorf("Error deleting a calendar for email %s. %s", Responses.AnchorMailbox, err.Error())
	}

	fmt.Printf("%s\n", contents)
}