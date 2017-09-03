package handlers

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type Google struct {
}

func NewGoogleHandler() (google *Google) {
	google = &Google{}
	return google
}

func (g *Google) TokenHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("google/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	params, err := util.CallAPIRoot("google/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	query := r.URL.Query()
	// TODO: Know how to send state
	//state := query.Get("state")

	code := query.Get("code")

	client := http.Client{}
	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))

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
	account, err := google.NewAccount(contents)

	go func(account *google.Account) {
		log.Debugln(account)
		account.GetAllCalendars()
		account.Refresh()
	}(account)

	//This is so that users cannot read the response
	http.Redirect(w, r, "http://localhost:8080", 301)
}
