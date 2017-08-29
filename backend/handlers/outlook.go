package handlers

import (
	"fmt"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"io/ioutil"
	"net/http"
	"strings"
)

type Outlook struct {
}

func NewOutlookHandler() (outlook *Outlook) {
	outlook = &Outlook{}
	return outlook
}

func (o *Outlook) TokenHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("outlook/token/uri")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	params, err := util.CallAPIRoot("outlook/token/request-params")
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

	//TODO: DB to implement
	account, err := outlook.NewAccount(contents)

	go func(account *outlook.OutlookAccount) {
		log.Debugln(account)
		account.GetAllCalendars()
	}(account)

	http.Redirect(w, r, "/", 301)
}
