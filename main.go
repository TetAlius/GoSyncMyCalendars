package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TetAlius/GoSyncMyCalendars/backend/google"
	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/backend/user"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	user.Save("Marta", "Jiménez", "tete", "password")

	log.Fatalln("Meh")
	//Parse configuration of outlook and google
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Error reading config.json: %s", err.Error())
	}
	err = json.Unmarshal(file, &outlook.Config)
	if err != nil {
		log.Fatalf("Error unmarshalling outlook config: %s", err.Error())
	}
	err = json.Unmarshal(file, &google.Config)
	if err != nil {
		log.Fatalf("Error unmarshalling google config: %s", err.Error())
	}

	//Parse all requests for outlook
	file, err = ioutil.ReadFile("./outlook.json")
	if err != nil {
		log.Fatalf("Error reading outlook.json: %s", err.Error())
	}
	err = json.Unmarshal(file, &outlook.Requests)
	if err != nil {
		log.Fatalf("Error unmarshalling outlook requests: %s", err.Error())
	}

	//Parse all requests for google
	file, err = ioutil.ReadFile("./google.json")
	if err != nil {
		log.Fatalf("Error reading google.json: %s", err.Error())
	}
	err = json.Unmarshal(file, &google.Requests)
	if err != nil {
		log.Fatalf("Error unmarshalling google requests: %s", err.Error())
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/signInWithOutlook", outlookSignInHandler)
	http.HandleFunc("/outlook", outlookTokenHandler)
	http.HandleFunc("/calendars", listCalendarsHandler)
	http.HandleFunc("/SignInWithGoogle", googleSignInHandler)
	http.HandleFunc("/google", googleTokenHandler)
	http.HandleFunc("/signUp", singUpHandler)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
