package handlers

import (
	//"encoding/json"
	//"fmt"
	//"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	//log "github.com/TetAlius/GoSyncMyCalendars/logger"
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
	http.Redirect(w, r, util.CallAPIRoot("outlook/login"), 302)
}
