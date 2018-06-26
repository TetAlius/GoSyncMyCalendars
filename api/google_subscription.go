package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/google/uuid"
)

// Function that creates a new GoogleSubscription given specific info
func NewGoogleSubscription(ID string) (subscription *GoogleSubscription) {
	subscription = new(GoogleSubscription)
	subscription.NotificationURL = fmt.Sprintf("%s:8081/google/watcher", os.Getenv("ENDPOINT"))
	subscription.Type = "web_hook"
	subscription.ID = ID
	subscription.Uuid = uuid.New()
	return
}

// Function that returns a GoogleSubscription given specific info
func RetrieveGoogleSubscription(ID string, uid uuid.UUID, calendar CalendarManager, resourceID string) (subscription *GoogleSubscription) {
	subscription = new(GoogleSubscription)
	subscription.ID = ID
	subscription.Uuid = uid
	subscription.calendar = calendar.(*GoogleCalendar)
	subscription.ResourceID = resourceID
	return
}

// Method that manages the data for a renewal
func (subscription *GoogleSubscription) manageRenewalData() {
	subscription.ID = uuid.New().String()
	subscription.NotificationURL = fmt.Sprintf("%s:8081/google/watcher", os.Getenv("ENDPOINT"))
	subscription.Type = "web_hook"
}

// Method that subscribes calendar for notifications
//
// POST https://www.googleapis.com/calendar/v3/calendars/calendarId/events/watch
func (subscription *GoogleSubscription) Subscribe(calendar CalendarManager) (err error) {
	if err = subscription.setCalendar(calendar); err != nil {
		log.Errorf("kind of subscription and calender differs: %s", calendar.GetName())
		return err
	}
	a := calendar.GetAccount()
	log.Debugln("subscribe calendar google")

	route, err := util.CallAPIRoot("google/calendars/subscription")
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
	log.Warningf("RESPONSE: %s", contents)

	err = createGoogleResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, subscription)
	subscription.setTime()
	//TODO manage Expiration
	return
}

// Method that renews subscription.
// Google does not let a subscription be renewed so
// a new subscription must be request
func (subscription *GoogleSubscription) Renew() (err error) {
	log.Debugln("Renew google subscription")
	subscription.manageRenewalData()
	return subscription.Subscribe(subscription.calendar)
}

// Method that deletes subscription
//
// POST https://www.googleapis.com/calendar/v3/channels/stop
func (subscription *GoogleSubscription) Delete() (err error) {
	a := subscription.calendar.GetAccount()
	log.Debugln("Delete google subscription")

	route, err := util.CallAPIRoot("google/calendars/subscription/stop")
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
	log.Warningf("RESPONSE: %s", contents)

	if len(contents) != 0 {
		err = createGoogleResponseError(contents)
		log.Errorf("error deleting subscription: %s", err.Error())
	}
	return
}

// Method that returns the ID of the subscription
func (subscription *GoogleSubscription) GetID() string {
	return subscription.ID
}

// Method that returns the UUID of the subscription
func (subscription *GoogleSubscription) GetUUID() uuid.UUID {
	return subscription.Uuid
}

// Method that returns the account of the subscription
func (subscription *GoogleSubscription) GetAccount() AccountManager {
	return subscription.calendar.account
}

// Method that returns the type of the subscription
func (subscription *GoogleSubscription) GetType() string {
	return subscription.Type
}

// Method that sets the expiration time to the subscription
func (subscription *GoogleSubscription) setTime() {
	subscription.expirationDate = time.Unix(subscription.Expiration/1000, 0)
}

// Method that sets the calendar to be watched by subscription
func (subscription *GoogleSubscription) setCalendar(calendar CalendarManager) (err error) {
	switch calendar.(type) {
	case *GoogleCalendar:
		subscription.calendar = calendar.(*GoogleCalendar)
	default:
		return &customErrors.WrongKindError{Mail: calendar.GetName()}
	}

	return
}

// Method that returns the expiration date of the subscription
func (subscription *GoogleSubscription) GetExpirationDate() time.Time {
	return subscription.expirationDate
}

// Method that returns the resourceID of the subscription
func (subscription *GoogleSubscription) GetResourceID() string {
	return subscription.ResourceID
}
