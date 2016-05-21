package google

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// Config TODO doc
var Config struct {
	googleConfig `json:"google"`
}

type googleConfig struct {
	ID                string `json:"client_id"`
	Secret            string `json:"client_secret"`
	RedirectURI       string `json:"redirect_uri"`
	DiscoveryDocument string `json:"discovery_document"`
	Endpoint          string `json:"authorization_endpoint"`
	TokenEndpoint     string `json:"token_endpoint"`
	LoginURI          string `json:"login_uri"`
	Version           string `json:"version"`
	Scope             string `json:"scope"`
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

// GetDiscoveryDocument TODO doc
func GetDiscoveryDocument() (document []byte) {
	client := http.Client{}
	req, err := http.NewRequest("GET", Config.DiscoveryDocument, nil)
	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing request: %s", err.Error())
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response: %s", err.Error())
	}
	//fmt.Printf("%s\n", contents)

	return contents
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

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)

	//fmt.Printf("%s\n", contents)
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
