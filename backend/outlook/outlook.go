package outlook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Outlook TODO: improve this calls
var Outlook struct {
	OutlookConfig `json:"outlook"`
}

// OutlookConfig TODO
type OutlookConfig struct {
	Id          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

// OutlookRequests TODO
var OutlookRequests struct {
	RootUri     string `json:"root_uri"`
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
	IdToken           string `json:"id_token"`
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

//OutlookCalendarResponse TODO
type OutlookCalendarResponse struct {
	OdataContext string                `json:"@odata.context"`
	Value        []OutlookCalendarInfo `json:"value"`
}

//
type OutlookCalendarInfo struct {
	OdataId   string `json:"@odata.id"`
	Id        string `json:"Id"`
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

// OutlookTokenRefresh TODO
func OutlookTokenRefresh(oldToken string) {
	client := http.Client{}
	//check if token is DEAD!!!

	req, err := http.NewRequest("POST",
		Outlook.LoginURI+Outlook.Version+"/token",
		strings.NewReader("grant_type=refresh_token"+
			"&client_id="+Outlook.Id+
			"&scope="+Outlook.Scope+
			"&refresh_token="+oldToken+
			"&client_secret="+Outlook.Secret))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &OutlookResp)
	if err != nil {
		fmt.Println(err)
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
