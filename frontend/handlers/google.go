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
	http.Redirect(w, r, util.CallAPIRoot("google/login"), 302)
}
