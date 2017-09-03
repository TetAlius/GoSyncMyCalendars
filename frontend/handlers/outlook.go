package handlers

import (
	//"encoding/json"
	//"fmt"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	//"io/ioutil"
	"net/http"
	//"strings"
)

type Outlook struct {
}

func NewOutlookHandler() (outlook *Outlook) {
	outlook = &Outlook{}
	return outlook
}

func (o *Outlook) SignInHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("outlook/login")
	code := 302
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		code = 500
		route = "http://localhost:8080"
	}
	http.Redirect(w, r, route, code)
}
