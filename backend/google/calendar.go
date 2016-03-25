package google

import (
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
)

//GET https://apidata.googleusercontent.com/caldav/v2/calid/user
func getAllCalendars() {
	fmt.Println("All calendars")

	contents :=
		backend.NewRequest(
			"GET",
			"https://www.googleapis.com/calendar/v3/users/me/calendarList",
			nil,
			Responses.TokenType+" "+Responses.AccessToken,
			"")

	fmt.Printf("%s\n", contents)
}

//TODO
func getPrimaryCalendar() {

}
func getCalendar(calendarID string) {

}

func createCalendar(calendarData []byte) {

}
func updateCalendar(calendarID string, calendarData []byte) {

}

func deleteCalendar(calendarID string) {

}
