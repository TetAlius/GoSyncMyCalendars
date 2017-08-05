package handlers

import (
	"encoding/json"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//Config TODO: improve this calls
type Outlook struct {
	ID          string `json:"client_id"`
	Secret      string
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

func NewOutlookHandler(fileName *string) (outlook *Outlook, err error) {
	outlook = &Outlook{}
	//Parse configuration of outlook
	//file, err := ioutil.ReadFile("./config.json")
	file, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Errorf("Error reading config.json: %s", err.Error())
		return
	}
	err = json.Unmarshal(file, &outlook)
	if err != nil {
		log.Errorf("Error unmarshalling google config: %s", err.Error())
		return
	}
	env := os.Getenv("outlook_secret")
	outlook.Secret = env
	if len(outlook.ID) == 0 || len(outlook.Secret) == 0 || len(outlook.RedirectURI) == 0 || len(outlook.LoginURI) == 0 || len(outlook.Version) == 0 || len(outlook.Scope) == 0 {
		err = customErrors.ConfigNotChargedCorrectlyError{Message: "Outlook Config has not loaded properly"}
		log.Errorf("%s", err.Error())
		return
	}
	return outlook, err
}

func (o *Outlook) SignInHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		o.LoginURI+o.Version+
			"/authorize?client_id="+o.ID+
			"&redirect_uri="+o.RedirectURI+
			"&response_type=code&scope="+o.Scope, 302)
}

//TODO handle errors
func (o *Outlook) TokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		o.LoginURI+o.Version+
			"/token",
		strings.NewReader("grant_type=authorization_code"+
			"&code="+code+
			"&redirect_uri="+o.RedirectURI+
			"&client_id="+o.ID+
			"&client_secret="+o.Secret))
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
