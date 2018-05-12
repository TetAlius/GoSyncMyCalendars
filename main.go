package main

import (
	"os"
	"os/signal"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func main() {

	f := frontend.NewServer("127.0.0.1", 8080, "./frontend/resources")
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
