package handlers

import (
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"net/http"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"encoding/json"
	"io/ioutil"
)

//Google //TODO Document
type Google struct {
	googleConfig `json:"google"`
}

type googleConfig struct {
	ID            string `json:"client_id"`
	Secret        string `json:"client_secret"`
	RedirectURI   string `json:"redirect_uri"`
	Endpoint      string `json:"authorization_endpoint"`
	TokenEndpoint string `json:"token_endpoint"`
	Scope         string `json:"scope"`
}

func NewGoogleHandler() *Google {

	google := Google{}

	//Parse configuration of google
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Error reading config.json: %s", err.Error())
	}
	err = json.Unmarshal(file, &google)
	if err != nil {
		log.Fatalf("Error unmarshalling google config: %s", err.Error())
	}
	return &google
}

//SingInHandler Google SingIn handler
func (g *Google) SignInHandler(w http.ResponseWriter, r *http.Request) {
	//google.Requests.State = google.GenerateRandomState()
	//log.Debugf("Random google state: %s", google.Requests.State)

	http.Redirect(w, r, g.Endpoint+
		"?client_id="+g.ID+
		"&access_type=offline&response_type=code"+
		"&scope="+g.Scope+
		"&redirect_uri="+g.RedirectURI+
		//"&state="+g.State+ //TODO Uncomment
		"&prompt=consent&include_granted_scopes=true",
		301)
}
