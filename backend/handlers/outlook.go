package handlers

import (
	//"encoding/json"
	"fmt"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"io/ioutil"
	"net/http"
	"strings"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/google"
)

type Outlook struct {
}

func NewOutlookHandler() (outlook *Outlook) {
	outlook = &Outlook{}
	return outlook
}

func (o *Outlook) TokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		util.CallAPIRoot("outlook/token/uri"),
		strings.NewReader(
			fmt.Sprintf(util.CallAPIRoot("outlook/token/request-params"), code)))
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
	_, err = outlook.NewResponse(contents)

	//TODO remove this call!
	outlook.TokenRefresh(outlook.Responses.RefreshToken)

	http.Redirect(w, r, "/", 301)
}
