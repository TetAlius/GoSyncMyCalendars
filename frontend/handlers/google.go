package handlers

import (
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"net/http"
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
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	http.Redirect(w, r, route, 302)
}
