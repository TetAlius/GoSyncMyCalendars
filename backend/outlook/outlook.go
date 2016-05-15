package outlook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Outlook TODO: improve this calls
var Outlook struct {
	Config `json:"outlook"`
}

// Config TODO
type Config struct {
	ID          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

// OutlookRequests TODO
var OutlookRequests struct {
	RootURI     string `json:"root_uri"`
	Version     string `json:"version"`
	UserContext string `json:"user_context"`
	Calendars   string `json:"calendars"`
	Events      string `json:"events"`
}

// OutlookResp TODO: this will be change to type and not var when I store the access_token on the BD
var OutlookResp struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	Scope             string `json:"scope"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	IDToken           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

type outlookEvent struct {
	ID                         string    `json:"Id"`
	OriginalStartTimeZone      string    `json:"OriginalStartTimeZone"`
	OriginalEndTimeZone        string    `json:"OriginalEndTimeZone"`
	ReminderMinutesBeforeStart string    `json:"ReminderMinutesBeforeStart"`
	IsReminderOn               bool      `json:"IsReminderOn"`
	HasAttachments             bool      `json:"HasAttachments"`
	Subject                    string    `json:"Subject"`
	Body                       body      `json:"Body"`
	BodyPreview                string    `json:"BodyPreview"`
	Importance                 string    `json:"Importance"`
	Sensitivity                string    `json:"Sensitivity"`
	Start                      eventDate `json:"Start"`
	End                        eventDate `json:"End"`
}

type body struct {
	ContentType string `json:"ContentType"`
	Body        string `json:"Content"`
}
type eventDate struct {
	DateTime string `json:"DateTime"`
	TimeZone string `json:"TimeZone"`
}

//CalendarResponse TODO
type CalendarResponse struct {
	OdataContext string         `json:"@odata.context"`
	Value        []CalendarInfo `json:"value"`
}

// CalendarInfo TODO
type CalendarInfo struct {
	OdataID   string `json:"@odata.id"`
	ID        string `json:"Id"`
	Name      string `json:"Name"`
	Color     string `json:"Color"`
	ChangeKey string `json:"ChangeKey"`
}

var calendar = []byte(`{
  "Name": "Social events"
}`)

var calendar2 = []byte(`{
  "Name": "Social"
}`)

// TokenRefresh TODO
func TokenRefresh(oldToken string) {
	client := http.Client{}
	//check if token is DEAD!!!

	req, err := http.NewRequest("POST",
		Outlook.LoginURI+Outlook.Version+"/token",
		strings.NewReader("grant_type=refresh_token"+
			"&client_id="+Outlook.ID+
			"&scope="+Outlook.Scope+
			"&refresh_token="+oldToken+
			"&client_secret="+Outlook.Secret))

	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &OutlookResp)
	if err != nil {
		log.Errorf("Error unmarshaling outlook response: %s", err.Error())
	}
	//TODO save info

	//TODO CRUD events
	//getAllEvents() TESTED
	//createEvent("", nil)
	//updateEvent("", nil) TESTED
	//deleteEvent("") TESTED
	//getEvent("") TESTED

	//TODO CRUD calendars
	//getAllCalendars() TESTED
	//getCalendar() TESTED
	//updateCalendar() TESTED
	//deleteCalendar() TESTED
	//createCalendar() TESTED
}

// https://outlook.office.com/api/v2.0/me/calendars/{calendarID}/events
func eventsFromCalendarURI(calendarID string) (URI string) {
	return OutlookRequests.RootURI + OutlookRequests.Version + OutlookRequests.UserContext + OutlookRequests.Calendars + "/" + calendarID + OutlookRequests.Events
}

// https://outlook.office.com/api/v2.0/me/events/{eventID}
func eventURI(eventID string) (URI string) {
	return OutlookRequests.RootURI + OutlookRequests.Version + OutlookRequests.UserContext + OutlookRequests.Events + "/" + eventID
}

// https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func calendarsURI(calendarID string) (URI string) {
	return OutlookRequests.RootURI + OutlookRequests.Version + OutlookRequests.UserContext + OutlookRequests.Calendars + "/" + calendarID
}

func authorizationRequest() (auth string) {
	return OutlookResp.TokenType + " " + OutlookResp.AccessToken
}
