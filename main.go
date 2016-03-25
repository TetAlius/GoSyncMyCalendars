package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
)

func main() {
	//Parse configuration of outlook and google
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &google.Config)
	if err != nil {
		fmt.Println(err)
	}

	//Parse all requests for google
	file, err = ioutil.ReadFile("./google.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &google.Requests)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/SignInWithGoogle", googleSignInHandler)
	http.HandleFunc("/google", googleTokenHandler)

	http.ListenAndServe(":8080", nil)
}
