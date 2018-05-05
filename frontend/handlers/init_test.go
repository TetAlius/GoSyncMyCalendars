package handlers_test

import (
	"os"

	"github.com/TetAlius/GoSyncMyCalendars/frontend"
)

func init() {
	f := frontend.NewServer("127.0.0.1", 8080)
	f.Start()
}

func setupApiRoot() {
	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}
