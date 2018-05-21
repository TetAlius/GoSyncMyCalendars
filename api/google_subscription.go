package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/google/uuid"
)

func NewGoogleSubscription(ID string, notificationURL string) (subscription *GoogleSubscription) {
	subscription = new(GoogleSubscription)
	subscription.NotificationURL = notificationURL
	subscription.Type = "web_hook"
	subscription.ID = ID
	subscription.Uuid = uuid.New()
	return
}

//TODO:
//func manageRenewalData(subscription *GoogleSubscription) (data []byte, err error) {
//	renewal := new(GoogleSubscription)
//	renewal.Type = subscription.Type
//	renewal.ExpirationDateTime = subscription.ExpirationDateTime
//	data, err = json.Marshal(renewal)
//	return
//}

// POST https://www.googleapis.com/apiName/apiVersion/resourcePath/watch
func (subscription *GoogleSubscription) Subscribe(calendar CalendarManager) (err error) {
	a := calendar.GetAccount()
	log.Debugln("subscribe calendar google")

	route, err := util.CallAPIRoot("google/subscription")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	data, err := json.Marshal(subscription)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPost,
		fmt.Sprintf(route, calendar.GetID()),
		bytes.NewBuffer(data),
		headers, nil)

	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, subscription)
	//TODO manage Expiration
	return
}

//Google does not let subscription be renewed
//A new subscription must be request
func (subscription *GoogleSubscription) Renew(a AccountManager) (err error) {
	log.Debugln("Renew google subscription")
	return subscription.Subscribe(subscription.calendar)
}

//POST https://www.googleapis.com/calendar/v3/channels/stop
func (subscription *GoogleSubscription) Delete(a AccountManager) (err error) {
	log.Debugln("Delete google subscription")
	//TODO: this URL
	route, err := util.CallAPIRoot("google/subscription/stop")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()
	data, err := json.Marshal(subscription)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}

	contents, err := util.DoRequest(http.MethodPost,
		route,
		bytes.NewBuffer(data),
		headers, nil)
	err = createGoogleResponseError(contents)
	return
}

func (subscription *GoogleSubscription) GetID() string {
	return subscription.ID
}
func (subscription *GoogleSubscription) GetUUID() uuid.UUID {
	return subscription.Uuid
}
