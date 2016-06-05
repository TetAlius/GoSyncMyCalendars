package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type calendarInfo struct {
	account string
	names   []string
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/welcome.html")
}

func outlookSignInHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		outlook.Config.LoginURI+outlook.Config.Version+
			"/authorize?client_id="+outlook.Config.ID+
			"&redirect_uri="+outlook.Config.RedirectURI+
			"&response_type=code&scope="+outlook.Config.Scope, 301)
}

func listCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./frontend/calendars.html")
	if err != nil {
		log.Fatalf("Parse file error: %s", err.Error())
	}

	calendars := []calendarInfo{
		{"outlook@outlook.com", []string{"a", "b"}},
		{"outlook@outlook.com", []string{"a", "b"}},
	}
	log.Debugln(calendars)

	t.Execute(w, calendars)
}

//TODO handle errors
func outlookTokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		outlook.Config.LoginURI+outlook.Config.Version+
			"/token",
		strings.NewReader("grant_type=authorization_code"+
			"&code="+code+
			"&redirect_uri="+outlook.Config.RedirectURI+
			"&client_id="+outlook.Config.ID+
			"&client_secret="+outlook.Config.Secret))
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

	email, preferred, err := util.MailFromToken(strings.Split(google.Responses.TokenID, "."))
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
func googleSignInHandler(w http.ResponseWriter, r *http.Request) {
	google.Requests.State = google.GenerateRandomState()
	log.Debugf("Random google state: %s", google.Requests.State)

	http.Redirect(w, r, google.Config.Endpoint+
		"?client_id="+google.Config.ID+
		"&access_type=offline&response_type=code"+
		"&scope="+google.Config.Scope+
		"&redirect_uri="+google.Config.RedirectURI+
		"&state="+google.Requests.State+
		"&prompt=consent&include_granted_scopes=true",
		301)
}

func googleTokenHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query.Get("state")

	if strings.Compare(google.Requests.State, state) != 0 {
		log.Errorf("State is not correct, expected %s got %s", google.Requests.State, state)
	}

	code := query.Get("code")

	client := http.Client{}
	req, err := http.NewRequest("POST",
		google.Config.TokenEndpoint,
		strings.NewReader("code="+code+
			"&client_id="+google.Config.ID+
			"&client_secret="+google.Config.Secret+
			"&redirect_uri="+google.Config.RedirectURI+
			"&grant_type=authorization_code"))

	if err != nil {
		log.Errorf("Error creating new google request: %s", err.Error())
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing google request: %s", err.Error())
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from google request: %s", err.Error())
	}

	err = json.Unmarshal(contents, &google.Responses)
	if err != nil {
		log.Errorf("Error unmarshaling google responses: %s", err.Error())
	}

	log.Debugf("%s", contents)

	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(google.Responses.TokenID, "."))
	if err != nil {
		log.Errorf("Error retrieving google mail: %s", err.Error())
	} else {
		google.Responses.Email = email
	}

	//TODO remove tests
	google.TokenRefresh(google.Responses.RefreshToken)

	//This is so that users cannot read the response
	http.Redirect(w, r, "/", 301)

}
