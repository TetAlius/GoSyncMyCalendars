package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Outlook struct {
	OutlookConfig `json:"outlook"`
}

type OutlookConfig struct {
	Id          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

type OutlookRequests struct {
	RootUri     string `json:"root_uri"`
	Version     string `json:"version"`
	UserContext string `json:"user_context"`
	Calendars   string `json:"calendars"`
	Events      string `json:"events"`
}

type OutlookResp struct {
	TokenType     string `json:"token_type"`
	ExpiresIn     string `json:"expires_in"`
	Scope         string `json:"scope"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	IdToken       string `json:"id_token"`
	AnchorMailbox string
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
type OutlookCalendarResponse struct {
	OdataContext string                `json:"@odata.context"`
	value        []OutlookCalendarInfo `json:"value"`
}

type OutlookCalendarInfo struct {
	OdataId   string `json:"@odata.id"`
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Color     string `json:"Color"`
	ChangeKey string `json:"ChangeKey"`
}

func outlookTokenRefresh(oldToken string) {
	client := http.Client{}
	//check if token is DEAD!!!

	req, err := http.NewRequest("POST",
		outlook.LoginURI+outlook.Version+"/token",
		strings.NewReader("grant_type=refresh_token"+
			"&client_id="+outlook.Id+
			"&scope="+outlook.Scope+
			"&refresh_token="+oldToken+
			"&client_secret="+outlook.Secret))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &outlookResp)
	if err != nil {
		fmt.Println(err)
	}
	//TODO save info
	//getAllCalendars()
	//getAllEvents()
	//createEvent("", nil)
	//updateEvent("", nil)
	deleteEvent("")
}

func getAllCalendars() {
	fmt.Println("All Calendars")

	contents := NewRequest(
		"GET",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			outlookRequests.Calendars,
		nil,
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)
	fmt.Printf("%s\n", contents)

}

func getAllEvents() {
	fmt.Println("All Events")

	contents := NewRequest("GET",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			outlookRequests.Events,
		nil,
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)

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
	fmt.Println("Create event")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents := NewRequest("POST",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			//	outlookRequests.Calendars+ TODO: Uncomment this two
			//	calendarID+
			outlookRequests.Events,
		bytes.NewBuffer(event),
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

var update = []byte(`{
  "Location": {
    "DisplayName": "Your office"
  }
}`)

func updateEvent(eventID string, eventData []byte) {
	fmt.Println("Update event")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents := NewRequest("PATCH",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			outlookRequests.Events+"/AAMkADIyZTVhZWUzLTZkNDUtNDM0Mi04MmVkLTA3YTM1NmZjZmRhMABGAAAAAADBz_m20_ARTLPlrdxoDR-VBwAM8uNeNO_IS54Z_auRX3ZoAAAAHLWyAACVymTsMw3zQK86n81a2jLeAAIs-oJBAAA=",
		bytes.NewBuffer(update),
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}

func deleteEvent(eventID string) {
	fmt.Println("Delete event")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents := NewRequest("DELETE",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			outlookRequests.Events+"/AAMkADIyZTVhZWUzLTZkNDUtNDM0Mi04MmVkLTA3YTM1NmZjZmRhMABGAAAAAADBz_m20_ARTLPlrdxoDR-VBwAM8uNeNO_IS54Z_auRX3ZoAAAAHLWyAACVymTsMw3zQK86n81a2jLeAAIs-oJBAAA=",
		nil,
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)
}

func getEvent(eventID string) {
	fmt.Println("Get event")
	//POST https://outlook.office.com/api/v2.0/me/calendars/{calendar_id}/events
	contents := NewRequest("GET",
		outlookRequests.RootUri+
			outlookRequests.Version+
			outlookRequests.UserContext+
			outlookRequests.Events+"/AAMkADIyZTVhZWUzLTZkNDUtNDM0Mi04MmVkLTA3YTM1NmZjZmRhMABGAAAAAADBz_m20_ARTLPlrdxoDR-VBwAM8uNeNO_IS54Z_auRX3ZoAAAAHLWyAACVymTsMw3zQK86n81a2jLeAAIs-oJBAAA=",
		nil,
		outlookResp.TokenType+" "+
			outlookResp.AccessToken,
		outlookResp.AnchorMailbox)

	fmt.Printf("%s\n", contents)

}
