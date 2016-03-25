package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/welcome.html")
}

func googleSignInHandler(w http.ResponseWriter, r *http.Request) {
	google.Requests.State = google.GenerateRandomState()
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
		fmt.Println("state is not correct")
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
		fmt.Println(err)
	}

	fmt.Printf("%s,\n", contents)

	tokens := strings.Split(google.Responses.TokenID, ".")

	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1)

	encodedToken = encodedToken + "=="
	_, err = base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%s\n", decoded)

	//TODO remove tests
	google.TokenRefresh(google.Responses.RefreshToken)

	//This is so that users cannot read the response
	http.Redirect(w, r, "/", 301)

}
