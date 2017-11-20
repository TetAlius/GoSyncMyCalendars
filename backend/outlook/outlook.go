package outlook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

type Account struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

func NewAccount(contents []byte) (r *Account, err error) {
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling outlook response: %s", err.Error()))
	}

	email, preferred, err := util.MailFromToken(strings.Split(r.TokenID, "."), "=")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving outlook mail: %s", err.Error()))
	}
	r.AnchorMailbox = email
	r.PreferredUsername = preferred
	return
}

type outlookEvent struct {
	ID                         string    `json:"Id"`
	OriginalStartTimeZone      string    `json:"OriginalStartTimeZone"`
	OriginalEndTimeZone        string    `json:"OriginalEndTimeZone"`
	ReminderMinutesBeforeStart string    `json:"ReminderMinutesBeforeStart"`
	IsReminderOn               bool      `json:"IsReminderOn"`
	HasAttachments             bool      `json:"HasAttachments"`
	Subject                    string    `json:"Subject"`
	Body                       body      `json:"Body"`
	BodyPreview                string    `json:"BodyPreview"`
	Importance                 string    `json:"Importance"`
	Sensitivity                string    `json:"Sensitivity"`
	Start                      eventDate `json:"Start"`
	End                        eventDate `json:"End"`
}

type body struct {
	ContentType string `json:"ContentType"`
	Body        string `json:"Content"`
}
type eventDate struct {
	DateTime string `json:"DateTime"`
	TimeZone string `json:"TimeZone"`
}

//CalendarResponse TODO
type CalendarResponse struct {
	OdataContext string         `json:"@odata.context"`
	Value        []CalendarInfo `json:"value"`
}

// CalendarInfo TODO
type CalendarInfo struct {
	OdataID   string `json:"@odata.id"`
	ID        string `json:"Id"`
	Name      string `json:"Name"`
	Color     string `json:"Color"`
	ChangeKey string `json:"ChangeKey"`
}

var calendar2 = []byte(`{
  "Name": "Social"contents
}`)

// TokenRefresh TODO
func (o *Account) Refresh() (err error) {
	client := http.Client{}
	//check if token is DEAD!!!

	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	params, err := util.CallAPIRoot("outlook/token/refresh-params")
	log.Debugf("Params: %s", fmt.Sprintf(params, o.RefreshToken))
	if err != nil {
		log.Errorf("Error generating params: %s", err.Error())
		return
	}

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(fmt.Sprintf(params, o.RefreshToken)))

	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing outlook request: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from outlook request: %s", err.Error())
		return
	}

	log.Debugf("\nTokenType: %s\nExpiresIn: %d\nAccessToken: %s\nRefreshToken: %s\nTokenID: %s\nAnchorMailbox: %s\nPreferredUsername: %t",
		o.TokenType, o.ExpiresIn, o.AccessToken, o.RefreshToken, o.TokenID, o.AnchorMailbox, o.PreferredUsername)

	log.Debugf("%s\n", contents)
	err = json.Unmarshal(contents, &o)
	if err != nil {
		log.Errorf("There was an error with the outlook request: %s", err.Error())
		return
	}

	log.Debugf("\nTokenType: %s\nExpiresIn: %d\nAccessToken: %s\nRefreshToken: %s\nTokenID: %s\nAnchorMailbox: %s\nPreferredUsername: %t",
		o.TokenType, o.ExpiresIn, o.AccessToken, o.RefreshToken, o.TokenID, o.AnchorMailbox, o.PreferredUsername)
	return
}

func (o *Account) authorizationRequest() (auth string) {
	return o.TokenType + " " + o.AccessToken
}
