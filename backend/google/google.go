package google

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"
	"fmt"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/pkg/errors"
)

// Config TODO doc
var Config struct {
	googleConfig `json:"google"`
}

type googleConfig struct {
	ID            string `json:"client_id"`
	Secret        string `json:"client_secret"`
	RedirectURI   string `json:"redirect_uri"`
	Endpoint      string `json:"authorization_endpoint"`
	TokenEndpoint string `json:"token_endpoint"`
	Scope         string `json:"scope"`
}

// Requests TODO doc
var Requests struct {
	State        string
	RootURI      string `json:"root_uri"`
	CalendarAPI  string `json:"calendarAPI"`
	Version      string `json:"version"`
	Context      string `json:"user_context"`
	CalendarList string `json:"calendar_list"`
	Calendars    string `json:"calendars"`
	Events       string `json:"events"`
}

// Responses TODO doc
var Responses struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
	Email        string
}

type Response struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
	Email        string
}

func NewResponse(contents []byte) (r *Response, err error) {
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling google responses: %s", err.Error()))
	}

	log.Debugf("%s", contents)

	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(r.TokenID, "."))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving google mail: %s", err.Error()))
	}

	r.Email = email
	return
}

//GenerateRandomState TODO doc
func GenerateRandomState() (rs string) {
	size := 32

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		log.Errorf("Error creating random numbers: %s", err.Error())
	}

	rs = base64.URLEncoding.EncodeToString(rb)

	return
}

//TokenRefresh TODO doc
func TokenRefresh(oldToken string) {
	client := http.Client{}

	req, err := http.NewRequest("POST",
		Config.TokenEndpoint,
		strings.NewReader("client_id="+Config.ID+
			"&client_secret="+Config.Secret+
			"&refresh_token="+oldToken+
			"&grant_type=refresh_token"))

	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
	}

	//log.Debugf("%s\n", req.Body)

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing google request: %s", err.Error())
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from google request: %s", err.Error())
	}

	log.Debugf("%s\n", contents)

	//TODO CRUD events
	//getAllEvents("primary") //TESTED
	//createEvent("primary", nil) //TESTED
	//updateEvent("primary", "eventID", nil)//TESTED
	//deleteEvent("primary", "eventID")//TESTED
	//getEvent("primary", "eventID") //TESTED

	//TODO CRUD calendars
	//getAllCalendars() //TESTED
	//getCalendar("ID") //TESTED
	//updateCalendar("ID", []byte(`"Hola":"Adios"`)) //TESTED
	//deleteCalendar("ID") //TESTED
	//createCalendar([]byte(`"Hola":"Adios"`)) //TESTED
}

// https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func eventsURI(calendarID string, eventID string) (URI string) {
	return Requests.RootURI + Requests.CalendarAPI + Requests.Version + Requests.Calendars + "/" + calendarID + Requests.Events + "/" + eventID
}

// https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func calendarListURI(calendarID string) (URI string) {
	return Requests.RootURI + Requests.CalendarAPI + Requests.Version + Requests.Context + Requests.CalendarList + "/" + calendarID
}

// https://www.googleapis.com/calendar/v3/calendars/{calendarID}
func calendarsURI(calendarID string) (URI string) {
	return Requests.RootURI + Requests.CalendarAPI + Requests.Version + Requests.Calendars
}

func authorizationRequest() (auth string) {
	return Responses.TokenType + " " + Responses.AccessToken
}
