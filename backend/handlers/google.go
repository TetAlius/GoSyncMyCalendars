package handlers

import (
	//"encoding/json"
	"fmt"
	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"io/ioutil"
	"net/http"
	"strings"
)

type Google struct {
}

func NewGoogleHandler() (google *Google) {
	google = &Google{}
	return google
}

func (g *Google) TokenHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// TODO: Know how to send state
	//state := query.Get("state")

	//if strings.Compare(google.Requests.State, state) != 0 {
	//	log.Errorf("State is not correct, expected %s got %s", google.Requests.State, state)
	//}

	code := query.Get("code")

	client := http.Client{}
	req, err := http.NewRequest("POST",
		util.CallAPIRoot("google/token/uri"),
		strings.NewReader(
			fmt.Sprintf(util.CallAPIRoot("google/token/request-params"), code)))

	if err != nil {
		log.Errorf("Error creating new google request: %s", err.Error())
	}

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

	//TODO: DB to implement
	_, err = google.NewResponse(contents)

	//This is so that users cannot read the response
	http.Redirect(w, r, "http://localhost:8080", 301)
}
