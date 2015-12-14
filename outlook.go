package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var outlook Outlook

//TODO this will be removed when I store the access_token on the BD
var outlookResp OutlookResp

var outlookRequests OutlookRequests

type Outlook struct {
	OutlookConfig `json:"outlook"`
}

type OutlookConfig struct {
	Id          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	LoginURI    string `json:"login_uri"`
	Version     string `json:"version"`
	Scope       string `json:"scope"`
}

type OutlookRequests struct {
	RootUri         string `json:"root_uri"`
	Version         string `json:"version"`
	UserContext     string `json:"user_context"`
	GetAllCalendars string `json:"get_calendars"`
}

type OutlookResp struct {
	TokenType        string `json:"token_type"`
	ExpiresIn        string `json:"expires_in"`
	Scope            string `json:"scope"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IdToken          string `json:"id_token"`
	IdTokenExpiresIn string `json:"id_token_expires_in"`
	AnchorMailbox    string `json:"anchor_mailbox"`
}

type OutlookCalendarResponse struct {
	OdataContext string                `json:"@odata.context"`
	value        []OutlookCalendarInfo `json:"value"`
}

type OutlookCalendarInfo struct {
	OdataId   string `json:"@odata.id"`
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Color     string `json:"Color"`
	ChangeKey string `json:"ChangeKey"`
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

func outlookTokenRefresh(oldToken string) {
	client := http.Client{}
	//check if token is DEAD!!!

	req, err := http.NewRequest("POST",
		outlook.LoginURI+outlook.Version+"/token",
		strings.NewReader("grant_type=refresh_token"+
			"&client_id="+outlook.Id+
			"&scope="+outlook.Scope+
			"&refresh_token="+oldToken+
			"&client_secret="+outlook.Secret))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &outlookResp)
	if err != nil {
		fmt.Println(err)
	}
	//TODO save info
	getAllCalendars()

}

func getAllCalendars() {
	fmt.Println("All Calendars")
	client := http.Client{}

	req, err := http.NewRequest("GET",

		"https://outlook.office.com/api/v2.0/me/calendars",
		strings.NewReader(""))

	req.Header.Set("Authorization",
		outlookResp.TokenType+" "+outlookResp.AccessToken)
	//req.Header.Set("X-AnchorMailbox", outlookResp.AnchorMailbox)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", contents)

}

func main() {
	//Parse configuration of outlook and gmail
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &outlook)
	if err != nil {
		fmt.Println(err)
	}

	//Parse all requests for outlook
	file, err = ioutil.ReadFile("./outlook.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &outlookRequests)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/SignInOutlook", outlookSignInHandler)
	http.HandleFunc("/outlook", outlookTokenHandler)
	http.ListenAndServe(":8080", nil)
}
