package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type Outlook struct {
}

func NewOutlookHandler() (outlook *Outlook) {
	outlook = &Outlook{}
	return outlook
}

func (o *Outlook) TokenHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	params, err := util.CallAPIRoot("outlook/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))
	if err != nil {
		log.Errorf("Error creating new outlook request: %s", err.Error())
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
		log.Errorf("Error reading response body from outlook request: %s", err.Error())
	}
	log.Debugln(contents)
	//TODO: DB to implement
	account, err := outlook.NewAccount(contents)
	if err != nil {
		log.Errorf("Error creating new account request: %s", err.Error())
	}
	go func(account *outlook.Account) {
		log.Debugln(account)
		account.GetAllCalendars()
		account.Refresh()
	}(account)

	http.Redirect(w, r, "http://localhost:8080", 301)
}
