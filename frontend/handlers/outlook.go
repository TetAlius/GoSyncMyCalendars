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
	log.Debugln("Starting outlook petition")
	route, err := util.CallAPIRoot("outlook/login")
	log.Debugln(route)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}
	http.Redirect(w, r, route, 302)
}
