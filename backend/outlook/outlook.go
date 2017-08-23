package outlook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"errors"
	"fmt"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

//Config TODO: improve this calls
var Config struct {
	outlookConfig `json:"outlook"`
}

// Config TODO
type outlookConfig struct {
	ID          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

// Requests TODO
var Requests struct {
	RootURI     string `json:"root_uri"`
	Version     string `json:"version"`
	UserContext string `json:"user_context"`
	Calendars   string `json:"calendars"`
	Events      string `json:"events"`
}

// Responses TODO: this will be change to type and not var when I store the access_token on the BD
var Responses struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	Scope             string `json:"scope"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

type Response struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	Scope             string `json:"scope"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

func NewResponse(contents []byte) (r *Response, err error) {
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling outlook response: %s", err.Error()))
	}

	email, preferred, err := util.MailFromToken(strings.Split(r.TokenID, "."))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving outlook mail: %s", err.Error()))
	}
	r.AnchorMailbox = email
	r.PreferredUsername = preferred
	return
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
		Config.LoginURI+Config.Version+"/token",
		strings.NewReader("grant_type=refresh_token"+
			"&client_id="+Config.ID+
			"&scope="+Config.Scope+
			"&refresh_token="+oldToken+
			"&client_secret="+Config.Secret))

	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing outlook request: %s", err.Error())
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading outlook response: %s", err.Error())
	}
	err = json.Unmarshal(contents, &Responses)
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
	return Requests.RootURI + Requests.Version + Requests.UserContext + Requests.Calendars + "/" + calendarID + Requests.Events
}

// https://outlook.office.com/api/v2.0/me/events/{eventID}
func eventURI(eventID string) (URI string) {
	return Requests.RootURI + Requests.Version + Requests.UserContext + Requests.Events + "/" + eventID
}

// https://outlook.office.com/api/v2.0/me/calendars/{calendarID}
func calendarsURI(calendarID string) (URI string) {
	return Requests.RootURI + Requests.Version + Requests.UserContext + Requests.Calendars + "/" + calendarID
}

func authorizationRequest() (auth string) {
	return Responses.TokenType + " " + Responses.AccessToken
}
