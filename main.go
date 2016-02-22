package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var outlook Outlook

// TODO: this will be removed when I store the access_token on the BD
var outlookResp OutlookResp

var outlookRequests OutlookRequests

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

	//TODO: here will be all requests for gmail

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/signInOutlook", outlookSignInHandler)
	http.HandleFunc("/outlook", outlookTokenHandler)
	http.HandleFunc("/calendars", listCalendarsHandler)
	http.ListenAndServe(":8080", nil)
}
