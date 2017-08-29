package main

import (
	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
	_ "github.com/jackc/pgx/stdlib"
	"os"
	"os/signal"
)

func main() {

	f := frontend.NewFrontend("127.0.0.1", 8080)
	b := backend.NewServer("127.0.0.1", 8081)

	// Control + C interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			err := f.Stop()
			if err != nil {
				os.Exit(1)
			}
			err = b.Stop()
			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}()
	f.Start()
	b.Start()

	/*
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
		http.HandleFunc("/signUp", signUpHandler)
		http.HandleFunc("/signIn", signInHandler)
		http.HandleFunc("/cookies", cookiesHandlerTest)

		log.Fatalln(http.ListenAndServe(":8080", nil))
	*/
}
