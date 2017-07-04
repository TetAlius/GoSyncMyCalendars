package handlers

import (
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"net/http"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"encoding/json"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"io/ioutil"
)

//Config TODO: improve this calls
type Outlook struct {
	outlookConfig `json:"outlook"`
}

// Config TODO
type outlookConfig struct {
	ID          string `json:"client_id"`
	Secret      string `json:"client_secret"`
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
