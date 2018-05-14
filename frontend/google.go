package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (s *Server) googleSignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Starting google petition")
	route, err := util.CallAPIRoot("google/login")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}
	http.Redirect(w, r, route, http.StatusFound)
}

func (s *Server) googleTokenHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		return
	}
	route, err := util.CallAPIRoot("google/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}
	params, err := util.CallAPIRoot("google/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}
	query := r.URL.Query()
	// TODO: Know how to send state
	//state := query.Get("state")

	code := query.Get("code")

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))

	if err != nil {
		log.Errorf("error creating new google request: %s", err.Error())
		serverError(w, err)
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error doing google request: %s", err.Error())
		serverError(w, err)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body from google request: %s", err.Error())
		serverError(w, err)
		return
	}

	//TODO: DB to implement
	account, err := api.NewGoogleAccount(contents)
	account.SetKind(api.GOOGLE)
	err = currentUser.AddAccount(account)
	if err != nil {
		serverError(w, err)
		return
	}

	//go func(account *api.GoogleAccount) {
	//	log.Debugln(account)
	//	account.GetAllCalendars()
	//	account.Refresh()
	//}(account)

	//This is so that users cannot read the response
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}