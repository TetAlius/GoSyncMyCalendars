package google

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	State string
}

// Responses TODO doc
var Responses struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
}

//GenerateRandomState TODO doc
func GenerateRandomState() (rs string) {
	size := 32

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	rs = base64.URLEncoding.EncodeToString(rb)

	return
}

// GetDiscoveryDocument TODO doc
func GetDiscoveryDocument() (document []byte) {
	client := http.Client{}
	req, err := http.NewRequest("GET", Config.DiscoveryDocument, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	fmt.Printf("%s\n", req.Body)

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s\n", contents)

	//TODO CRUD events
	//getAllEvents()
	//createEvent("", nil)
	//updateEvent("", nil)
	//deleteEvent("")
	//getEvent("")

	//TODO CRUD calendars
	getAllCalendars()
	//getCalendar()
	//updateCalendar()
	//deleteCalendar()
	//createCalendar()
}
