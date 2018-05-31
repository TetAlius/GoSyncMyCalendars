package frontend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/frontend/db"
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
	query := r.URL.Query()
	if len(query.Get("error")) > 0 {
		log.Errorf("google authorization with error: %s", query.Get("error"))
		http.Redirect(w, r, "/accounts", http.StatusPermanentRedirect)
		return
	}
	currentUser, ok := s.manageSession(w, r)
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
	var objmap map[string]interface{}
	err = json.Unmarshal(contents, &objmap)
	if err != nil {
		serverError(w, err)
		return
	}

	//MailFromToken(strings.Split(objmap["id_token"].(string), "."), "=")
	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(objmap["id_token"].(string), "."))
	if err != nil {
		serverError(w, err)
		return
	}
	acc := db.Account{
		User:         currentUser,
		TokenType:    objmap["token_type"].(string),
		RefreshToken: objmap["refresh_token"].(string),
		Email:        email,
		AccessToken:  objmap["access_token"].(string),
		Kind:         api.GOOGLE,
	}
	err = s.database.AddAccount(currentUser, acc)
	if err != nil {
		serverError(w, err)
		return
	}

	//This is so that users cannot read the response
	http.Redirect(w, r, "/accounts", http.StatusPermanentRedirect)
}
