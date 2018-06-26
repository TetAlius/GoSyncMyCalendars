package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"encoding/json"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/google/uuid"
)

// Function that creates a new GoogleSubscription given specific info
func NewOutlookSubscription() (subscription *OutlookSubscription) {
	subscription = new(OutlookSubscription)
	subscription.NotificationURL = fmt.Sprintf("%s:8081/outlook/watcher", os.Getenv("ENDPOINT"))
	subscription.ChangeType = "Created,Deleted,Updated"
	subscription.Type = "#Microsoft.OutlookServices.PushSubscription"
	subscription.Uuid = uuid.New()
	return
}

// Function that returns a GoogleSubscription given specific info
func RetrieveOutlookSubscription(ID string, uid uuid.UUID, calendar CalendarManager, typ string) (subscription *OutlookSubscription) {
	subscription = new(OutlookSubscription)
	subscription.ID = ID
	subscription.Uuid = uid
	subscription.calendar = calendar.(*OutlookCalendar)
	subscription.Type = typ
	return
}

// Method that manages the data for a renewal
func manageRenewalData(subscription *OutlookSubscription) (data []byte, err error) {
	renewal := new(OutlookSubscription)
	renewal.Type = subscription.Type
	renewal.ExpirationDateTime = subscription.ExpirationDateTime
	data, err = json.Marshal(renewal)
	return
}

// Method that subscribes calendar for notifications
//
// POST https://outlook.office.com/api/v2.0/me/subscriptions
func (subscription *OutlookSubscription) Subscribe(calendar CalendarManager) (err error) {
	if err = subscription.setCalendar(calendar); err != nil {
		log.Errorf("kind of subscription and calender differs: %s", calendar.GetName())
		return err
	}
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

	log.Warningf("RESPONSE: %s", contents)

	err = createOutlookResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, subscription)
	subscription.setTime()

	return
}

// Method that renews subscription.
//
// PATCH https://outlook.office.com/api/v2.0/me/subscriptions/{subscriptionId}
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
	log.Warningf("RESPONSE: %s", contents)
	err = createOutlookResponseError(contents)
	subscription.setTime()

	return
}

// Method that deletes subscription
//
// DELETE https://outlook.office.com/api/v2.0/me/subscriptions('{subscriptionId}')
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
	log.Warningf("RESPONSE: %s", contents)
	if len(contents) != 0 {
		err = createOutlookResponseError(contents)
		return err
	}
	return
}

// Method that returns the ID of the subscription
func (subscription *OutlookSubscription) GetID() string {
	return subscription.ID
}

// Method that returns the UUID of the subscription
func (subscription *OutlookSubscription) GetUUID() uuid.UUID {
	return subscription.Uuid
}

// Method that returns the account of the subscription
func (subscription *OutlookSubscription) GetAccount() AccountManager {
	return subscription.calendar.account
}

// Method that returns the type of the subscription
func (subscription *OutlookSubscription) GetType() string {
	return subscription.Type
}

// Method that sets the expiration time to the subscription
func (subscription *OutlookSubscription) setTime() {
	expiration, err := time.Parse(time.RFC3339Nano, subscription.ExpirationDateTime)
	if err != nil {
		subscription.expirationDate = time.Now().Add(time.Hour * 24 * 7)
		return
	}
	subscription.expirationDate = expiration
}

// Method that sets the calendar to be watched by subscription
func (subscription *OutlookSubscription) setCalendar(calendar CalendarManager) (err error) {
	switch calendar.(type) {
	case *OutlookCalendar:
		subscription.calendar = calendar.(*OutlookCalendar)
	default:
		return &customErrors.WrongKindError{Mail: calendar.GetName()}
	}

	return
}

// Method that returns the expiration date of the subscription
func (subscription *OutlookSubscription) GetExpirationDate() time.Time {
	return subscription.expirationDate
}

// Method that returns the resourceID of the subscription
func (subscription *OutlookSubscription) GetResourceID() string {
	return ""
}
