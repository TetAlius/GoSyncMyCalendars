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

const (
	USER     = "postgres"
	PASSWORD = "postgres"
	NAME     = "postgres"
)

func main() {

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		USER, PASSWORD, NAME)
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
	f := frontend.NewServer("127.0.0.1", 8080, "./frontend/resources", frontendDB)
	maxWorker := 15
	b := backend.NewServer("127.0.0.1", 8081, maxWorker)

	// Control + C interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
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
