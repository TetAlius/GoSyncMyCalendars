package handlers

import (
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"net/http"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"encoding/json"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"html/template"
	"io/ioutil"
	"os"
)

//Google //TODO Document
type Google struct {
	ID            string `json:"client_id"`
	Secret        string
	RedirectURI   string `json:"redirect_uri"`
	Endpoint      string `json:"authorization_endpoint"`
	TokenEndpoint string `json:"token_endpoint"`
	Scope         string `json:"scope"`
	//googleConfig `json:"google"`
}

func NewGoogleHandler(fileName *string) (google *Google, err error) {
	google = &Google{}
	//Parse configuration of google
	//file, err := ioutil.ReadFile("./config.json")
	file, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Errorf("Error reading config.json: %s", err.Error())
		return
	}
	err = json.Unmarshal(file, &google)
	if err != nil {
		log.Errorf("Error unmarshalling google config: %s", err.Error())
		return
	}
	env := os.Getenv("google_secret")
	google.Secret = env

	if len(google.ID) == 0 || len(google.Secret) == 0 || len(google.RedirectURI) == 0 || len(google.Endpoint) == 0 || len(google.TokenEndpoint) == 0 || len(google.Scope) == 0 {
		err = customErrors.ConfigNotChargedCorrectlyError{Message: "Google Config has not loaded properly"}
		log.Errorf("%s", err.Error())
		return
	}
	return google, err
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
		302)
}

func (g *Google) TokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Has entered token handler")
	t, err := template.ParseFiles("./frontend/resources/html/welcome.html")
	if err != nil {
		log.Errorln("Error reading config.json: %s", err.Error())
	}
	err = t.Execute(w, nil) //No template at this moment
	if err != nil {
		log.Errorln(err)
	}
}
