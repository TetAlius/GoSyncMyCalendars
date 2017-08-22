package handlers

import (
	"encoding/json"
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

func (o *Outlook) SignInHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, util.CallAPIRoot("outlook/login"), 302)
}

//TODO handle errors
func (o *Outlook) TokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		util.CallAPIRoot("outlook/token/uri"),
		strings.NewReader(
			fmt.Sprintf(util.CallAPIRoot("outlook/token/request"), code)))

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	if err != nil {
		log.Errorf("Error creating new outlook request: %s", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing outlook request: %s", err.Error())
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from outlook request: %s", err.Error())
	}
	log.Debugf("%s\n", contents)
	err = json.Unmarshal(contents, &outlook.Responses)
	//TODO save info
	if err != nil {
		log.Errorf("Error unmarshaling outlook response: %s", err.Error())
	}

	email, preferred, err := util.MailFromToken(strings.Split(outlook.Responses.IDToken, "."))
	if err != nil {
		log.Errorf("Error retrieving outlook mail: %s", err.Error())
	} else {
		outlook.Responses.AnchorMailbox = email
		outlook.Responses.PreferredUsername = preferred
	}

	//TODO remove this call!
	outlook.TokenRefresh(outlook.Responses.RefreshToken)

	http.Redirect(w, r, "/", 301)

}
