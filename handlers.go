package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
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
		outlook.Outlook.LoginURI+outlook.Outlook.Version+
			"/authorize?client_id="+outlook.Outlook.ID+
			"&redirect_uri="+outlook.Outlook.RedirectURI+
			"&response_type=code&scope="+outlook.Outlook.Scope, 301)
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

	req, _ := http.NewRequest("POST",
		outlook.Outlook.LoginURI+outlook.Outlook.Version+
			"/token",
		strings.NewReader("grant_type=authorization_code"+
			"&code="+code+
			"&redirect_uri="+outlook.Outlook.RedirectURI+
			"&client_id="+outlook.Outlook.ID+
			"&client_secret="+outlook.Outlook.Secret))
	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", contents)
	err := json.Unmarshal(contents, &outlook.OutlookResp)
	//TODO save info
	if err != nil {
		log.Errorf("Error unmarshaling outlook response: %s", err.Error())
	}

	tokens := strings.Split(outlook.OutlookResp.IDToken, ".")

	//According to Outlook example, this replaces must be done
	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1)

	//TODO create evaluation of last two ==
	//Go must have the == at the end of base64 decode
	//in order to decode it without errors
	encodedToken = encodedToken + "=="
	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		log.Errorf("Error decoding outlook token: %s", err.Error())
	}

	var f interface{}
	err = json.Unmarshal(decodedToken, &f)
	if err != nil {
		log.Errorf("Error unmarshaling outlook decoded token: %s", err.Error())
	}
	m := f.(map[string]interface{})

	// TODO: El email petará si no recibo eso en el JSON
	if m["email"] != nil {
		log.Infoln("Got email on outlook")
		outlook.OutlookResp.AnchorMailbox = m["email"].(string)
		outlook.OutlookResp.PreferredUsername = false
	} else {
		log.Infoln("Got preferred username on outlook")
		outlook.OutlookResp.AnchorMailbox = m["preferred_username"].(string)
		outlook.OutlookResp.PreferredUsername = true
	}

	//TODO remove this call!
	outlook.TokenRefresh(outlook.OutlookResp.RefreshToken)

	http.Redirect(w, r, "/", 301)

}
func googleSignInHandler(w http.ResponseWriter, r *http.Request) {
	google.Requests.State = google.GenerateRandomState()
	log.Debugf("Random google state: %s", google.Requests.State)
	_ = google.GetDiscoveryDocument()

	//	fmt.Printf("%s\n", google.Requests.State)

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
	req, _ := http.NewRequest("POST",
		google.Config.TokenEndpoint,
		strings.NewReader("code="+code+
			"&client_id="+google.Config.ID+
			"&client_secret="+google.Config.Secret+
			"&redirect_uri="+google.Config.RedirectURI+
			"&grant_type=authorization_code"))

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(contents, &google.Responses)
	if err != nil {
		log.Errorf("Error unmarshaling google responses: %s", err.Error())
	}

	log.Debugf("%s", contents)

	tokens := strings.Split(google.Responses.TokenID, ".")

	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1)

	encodedToken = encodedToken + "=="
	_, err = base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		log.Errorf("Error decoding google token: %s", err.Error())
	}

	//fmt.Printf("%s\n", decoded)

	//TODO remove tests
	google.TokenRefresh(google.Responses.RefreshToken)

	//This is so that users cannot read the response
	http.Redirect(w, r, "/", 301)

}
