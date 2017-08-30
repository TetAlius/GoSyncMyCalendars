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

type GoogleAccount struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
	Email        string
}

func NewAccount(contents []byte) (r *GoogleAccount, err error) {
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

func (g *GoogleAccount) authorizationRequest() string {
	return fmt.Sprintf("%s %s", g.TokenType, g.AccessToken)
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
func (g *GoogleAccount) Refresh() {
	client := http.Client{}

	route, err := util.CallAPIRoot("google/token/uri")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	params, err := util.CallAPIRoot("google/token/refresh-params")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, g.RefreshToken)))

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
