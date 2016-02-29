package outlook

import (
	"bytes"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
)

func getAllCalendars() {
	fmt.Println("All Calendars")

	contents := backend.NewRequest("GET",
		OutlookRequests.RootUri+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}

//GET https://outlook.office.com/api/v2.0/me/calendar
func getPrimaryCalendar() {
	//TODO
}

//GET https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}
func getCalendar(calendarID string) {
	fmt.Println("Get Calendar given")
	contents := backend.NewRequest("GET",
		OutlookRequests.RootUri+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars+"/"+
			calendarID,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

//POST https://outlook.office.com/api/v2.0/me/calendars
func createCalendar(calendarData []byte) {
	fmt.Println("Create calendar")

	contents := backend.NewRequest("POST",
		OutlookRequests.RootUri+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars,
		bytes.NewBuffer(calendarData),
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}

//PATCH https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}
func updateCalendar(calendarID string, calendarData []byte) {
	fmt.Println("Update calendar")

	contents := backend.NewRequest("PATCH",
		OutlookRequests.RootUri+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars+"/"+
			calendarID,
		bytes.NewBuffer(calendarData),
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

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
//DELETE https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}
//Does not return json if OK, only status 204
func deleteCalendar(calendarID string) {
	fmt.Println("Delete calendar")

	contents := backend.NewRequest("DELETE",
		OutlookRequests.RootUri+
			OutlookRequests.Version+
			OutlookRequests.UserContext+
			OutlookRequests.Calendars+"/"+
			calendarID,
		nil,
		OutlookResp.TokenType+" "+
			OutlookResp.AccessToken,
		OutlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}
