package main

import (
	"os"
	"os/signal"

	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {

	f := frontend.NewServer("127.0.0.1", 8080)
	b := backend.NewServer("127.0.0.1", 8081, maxWorker)

	// Control + C interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			err := f.Stop()
			exit := 0
			if err != nil {
				exit = 1
			}
			err = b.Stop()
			if err != nil {
				exit = 1
			}
			os.Exit(exit)
		}
	}()
	f.Start()
	b.Start()
}
