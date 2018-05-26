package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func main() {
	missing := false
	user := os.Getenv("DB_USER")
	if len(user) <= 0 {
		log.Errorf("missing DB_USER variable")
		missing = true
	}
	password := os.Getenv("DB_PASSWORD")
	if len(password) <= 0 {
		log.Errorf("missing DB_PASSWORD variable")
		missing = true
	}
	name := os.Getenv("DB_NAME")
	if len(name) <= 0 {
		log.Errorf("missing DB_NAME variable")
		missing = true
	}
	host := os.Getenv("DB_HOST")
	if len(host) <= 0 {
		log.Errorf("missing DB_USER variable")
		missing = true
	}
	if len(os.Getenv("DNS_NAME")) <= 0 {
		log.Errorf("missing DNS_NAME variable")
		missing = true
	}
	if missing {
		os.Exit(1)
	}

	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, name)
	frontendDB, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Errorf("error opening frontend database: %s", err.Error())
		os.Exit(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	err = frontendDB.Ping()
	if err != nil {
		log.Errorf("error ping frontend database: %s", err.Error())
		os.Exit(1)
	}

	backendDB, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Errorf("error opening backend database: %s", err.Error())
		os.Exit(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	err = backendDB.Ping()
	if err != nil {
		log.Errorf("error ping backend database: %s", err.Error())
		os.Exit(1)
	}

	f := frontend.NewServer("127.0.0.1", 8080, "./frontend/resources", frontendDB)
	maxWorker := 15
	b := backend.NewServer("127.0.0.1", 8081, maxWorker, backendDB)

	// Control + C interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	//signal.Notify(c, syscall.SIGKILL)
	//
	//signal.Notify(c, syscall.SIGINT)
	//signal.Notify(c, syscall.SIGTERM)

	go func() {
		for range c {
			err := f.Stop()
			exit := 0
			if err != nil {
				log.Errorf("not finished frontend correctly: %s", err.Error())
				exit = 1
			}
			err = b.Stop()
			if err != nil {
				log.Errorf("not finished backend correctly: %s", err.Error())
				exit = 1
			}
			os.Exit(exit)
		}
	}()
	f.Start()
	b.Start()
}
