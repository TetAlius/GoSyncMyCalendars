package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type calendarInfo struct {
	account string
	names   []string
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "welcome.html")
}

func outlookSignInHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		outlook.LoginURI+outlook.Version+
			"/authorize?client_id="+outlook.Id+
			"&redirect_uri="+outlook.RedirectURI+
			"&response_type=code&scope="+outlook.Scope, 301)
}

func listCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./frontend/calendars.html")
	if err != nil {
		log.Fatal("Parse file error: ", err)
	}

	calendars := []calendarInfo{
		{"outlook@outlook.com", []string{"a", "b"}},
		{"outlook@outlook.com", []string{"a", "b"}},
	}
	fmt.Println(calendars)
	t.Execute(w, calendars)
}

//TODO handle errors
func outlookTokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, _ := http.NewRequest("POST",
		outlook.LoginURI+outlook.Version+
			"/token",
		strings.NewReader("grant_type=authorization_code"+
			"&code="+code+
			"&redirect_uri="+outlook.RedirectURI+
			"&client_id="+outlook.Id+
			"&client_secret="+outlook.Secret))
	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(contents, &outlookResp)
	//TODO save info
	if err != nil {
		fmt.Println(err)
	}

	tokens := strings.Split(outlookResp.IdToken, ".")

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
		fmt.Printf("Decoding token: %s\n", err)
	}

	var f interface{}
	err = json.Unmarshal(decodedToken, &f)
	m := f.(map[string]interface{})
	outlookResp.AnchorMailbox = m["preferred_username"].(string)

	//TODO remove this call!
	outlookTokenRefresh(outlookResp.RefreshToken)

	http.Redirect(w, r, "/", 301)

}
