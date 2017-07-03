package handlers

import (
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"net/http"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"encoding/json"
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

func NewOutlookHandler() *Outlook {

	outlook := Outlook{}

	//Parse configuration of google
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Error reading config.json: %s", err.Error())
	}
	err = json.Unmarshal(file, &outlook)
	if err != nil {
		log.Fatalf("Error unmarshalling google config: %s", err.Error())
	}
	return &outlook
}

func (o *Outlook) SignInHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		o.LoginURI+o.Version+
			"/authorize?client_id="+o.ID+
			"&redirect_uri="+o.RedirectURI+
			"&response_type=code&scope="+o.Scope, 301)
}
