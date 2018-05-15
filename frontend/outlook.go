package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (s *Server) outlookSignInHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("outlook/login")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}
	http.Redirect(w, r, route, http.StatusFound)
}

func (s *Server) OutlookTokenHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		return
	}
	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}
	params, err := util.CallAPIRoot("outlook/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w, err)
		return
	}

	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))
	if err != nil {
		log.Errorf("error creating new outlook request: %s", err.Error())
		serverError(w, err)
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error doing outlook request: %s", err.Error())
		serverError(w, err)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body from outlook request: %s", err.Error())
		serverError(w, err)
		return
	}
	//TODO: DB to implement
	account, err := api.NewGoogleAccount(contents)
	account.SetKind(api.OUTLOOK)
	err = currentUser.AddAccount(account)
	if err != nil {
		serverError(w, err)
		return
	}

	http.Redirect(w, r, "http://localhost:8080", http.StatusPermanentRedirect)
}
