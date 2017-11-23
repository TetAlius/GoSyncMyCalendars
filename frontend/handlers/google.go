package handlers

import (
	"net/http"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type Google struct {
}

func NewGoogleHandler() (google *Google) {
	google = &Google{}
	return google
}

//SingInHandler Google SingIn handler
func (g *Google) SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln("Starting google petition")
	route, err := util.CallAPIRoot("google/login")
	code := 302
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		code = 500
		route = "http://localhost:8080"
	}
	http.Redirect(w, r, route, code)
}
