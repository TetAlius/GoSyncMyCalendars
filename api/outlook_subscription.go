package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func NewOutlookSubscription(ID string, notificationURL string, changeType string) (subscription *OutlookSubscription) {
	subscription = new(OutlookSubscription)
	subscription.NotificationURL = notificationURL
	subscription.ChangeType = changeType
	subscription.ID = ID
	subscription.Type = "#Microsoft.OutlookServices.PushSubscription"
	return
}

func manageRenewalData(subscription *OutlookSubscription) (data []byte, err error) {
	renewal := new(OutlookSubscription)
	renewal.Type = subscription.Type
	renewal.ExpirationDateTime = subscription.ExpirationDateTime
	data, err = json.Marshal(renewal)
	return
}

// POST https://outlook.office.com/api/v2.0/me/subscriptions
func (subscription *OutlookSubscription) Subscribe(calendar CalendarManager) (err error) {
	a := calendar.GetAccount()
	log.Debugln("subscribe calendar outlook")

	route, err := util.CallAPIRoot("outlook/subscription")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	resource, err := util.CallAPIRoot("outlook/calendars/id/events")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	resource = fmt.Sprintf(resource, calendar.GetID())
	subscription.Resource = resource
	data, err := json.Marshal(subscription)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPost,
		route,
		bytes.NewBuffer(data),
		headers, nil)

	log.Debugf("%s\n", contents)
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, subscription)

	return
}

//PATCH https://outlook.office.com/api/v2.0/me/subscriptions/{subscriptionId}
func (subscription *OutlookSubscription) Renew(a AccountManager) (err error) {
	log.Debugln("subscribe calendar outlook")

	route, err := util.CallAPIRoot("outlook/subscription")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	route = fmt.Sprintf("%s/%s", route, subscription.GetID())

	data, err := manageRenewalData(subscription)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodPatch,
		route,
		bytes.NewBuffer(data),
		headers, nil)
	log.Debugf("%s\n", contents)
	err = createOutlookResponseError(contents)

	return
}

//DELETE https://outlook.office.com/api/v2.0/me/subscriptions('{subscriptionId}')
func (subscription *OutlookSubscription) Delete(a AccountManager) (err error) {
	log.Debugln("Delete outlook subscription")
	route, err := util.CallAPIRoot("outlook/subscription")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	route = fmt.Sprintf("%s('%s')", route, subscription.GetID())

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodDelete,
		route,
		nil,
		headers, nil)
	log.Debugf("%s\n", contents)
	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}
	return
}

func (subscription *OutlookSubscription) GetID() string {
	return subscription.ID
}
