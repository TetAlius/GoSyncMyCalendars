package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"encoding/json"

	"time"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/google/uuid"
)

func NewOutlookSubscription(ID string, notificationURL string, changeType string) (subscription *OutlookSubscription) {
	subscription = new(OutlookSubscription)
	subscription.NotificationURL = notificationURL
	subscription.ChangeType = changeType
	subscription.ID = ID
	subscription.Type = "#Microsoft.OutlookServices.PushSubscription"
	subscription.Uuid = uuid.New()
	return
}
func RetrieveOutlookSubscription(ID string, uid uuid.UUID, calendar CalendarManager, typ string) (subscription *OutlookSubscription) {
	subscription = new(OutlookSubscription)
	subscription.ID = ID
	subscription.Uuid = uid
	subscription.calendar = calendar.(*OutlookCalendar)
	subscription.Type = typ
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

	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, subscription)
	subscription.setTime()

	return
}

//PATCH https://outlook.office.com/api/v2.0/me/subscriptions/{subscriptionId}
func (subscription *OutlookSubscription) Renew() (err error) {
	a := subscription.calendar.GetAccount()

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
	err = createOutlookResponseError(contents)

	return
}

//DELETE https://outlook.office.com/api/v2.0/me/subscriptions('{subscriptionId}')
func (subscription *OutlookSubscription) Delete() (err error) {
	a := subscription.calendar.GetAccount()
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
	if len(contents) != 0 {
		err = createOutlookResponseError(contents)
		return err
	}
	return
}

func (subscription *OutlookSubscription) GetID() string {
	return subscription.ID
}

func (subscription *OutlookSubscription) GetUUID() uuid.UUID {
	return subscription.Uuid
}

func (subscription *OutlookSubscription) GetAccount() AccountManager {
	return subscription.calendar.account
}

func (subscription *OutlookSubscription) GetType() string {
	return subscription.Type
}
func (subscription *OutlookSubscription) SetCalendar(calendar *OutlookCalendar) {
	subscription.calendar = calendar
}

func (subscription *OutlookSubscription) setTime() {
	expiration, err := time.Parse(time.RFC3339Nano, subscription.ExpirationDateTime)
	if err != nil {
		subscription.expirationDate = time.Now().Add(time.Hour * 24 * 7)
		return
	}
	subscription.expirationDate = expiration
}

func (subscription *OutlookSubscription) GetExpirationDate() time.Time {
	return subscription.expirationDate
}
