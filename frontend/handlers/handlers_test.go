package handlers_test

import (
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
)

func setup() {
	configNameGood := "../../config.json"
	f := frontend.NewFrontend("127.0.0.1", 8080, &configNameGood)
	f.Start()
	defer f.Stop()

}
