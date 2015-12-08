package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var outlook Outlook

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

type OutlookResp struct {
	TokenType        string `json:"token_type"`
	ExpiresIn        string `json:"expires_in"`
	Scope            string `json:"scope"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IdToken          string `json:"id_token"`
	IdTokenExpiresIn string `json:"id_token_expires_in"`
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
	var outlookResp OutlookResp

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(contents, &outlookResp)
	//TODO save info
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Token:%s\n", string(contents))

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
	var outlookResp OutlookResp
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &outlookResp)
	if err != nil {
		fmt.Println(err)
	}
	//TODO save info

}

func main() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &outlook)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/SignInOutlook", outlookSignInHandler)
	http.HandleFunc("/outlook", outlookTokenHandler)
	http.ListenAndServe(":8080", nil)
}
