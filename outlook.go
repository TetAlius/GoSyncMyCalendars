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
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "welcome.html")
}

func outlookSignInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hola\n")
	http.Redirect(w, r,
		outlook.LoginURI+outlook.Version+
			"/authorize?client_id="+outlook.Id+
			"&redirect_uri="+outlook.RedirectURI+
			"&response_type=code&scope=https://outlook.office.com/Calendars.ReadWrite", 301)
}

func outlookTokenHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	code := r.URL.Query().Get("code")
	fmt.Printf("%s\n", code)

	req, _ := http.NewRequest("POST",
		outlook.LoginURI+outlook.Version+
			"/token",
		strings.NewReader("grant_type=authorization_code&code="+code+
			"&redirect_uri="+outlook.RedirectURI+
			"&client_id="+outlook.Id+
			"&client_secret="+outlook.Secret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(contents))

	http.Redirect(w, r, "/", 301)

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
