package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"os"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
	"github.com/google/uuid"
)

func NewGoogleSubscription(ID string) (subscription *GoogleSubscription) {
	subscription = new(GoogleSubscription)
	subscription.NotificationURL = fmt.Sprintf("%s:8081/google/watcher", os.Getenv("ENDPOINT"))
	subscription.Type = "web_hook"
	subscription.ID = ID
	subscription.Uuid = uuid.New()
	return
}

func RetrieveGoogleSubscription(ID string, uid uuid.UUID, calendar CalendarManager) (subscription *GoogleSubscription) {
	subscription = new(GoogleSubscription)
	subscription.ID = ID
	subscription.Uuid = uid
	subscription.calendar = calendar.(*GoogleCalendar)
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

//Google does not let subscription be renewed
//A new subscription must be request
func (subscription *GoogleSubscription) Renew() (err error) {
	log.Debugln("Renew google subscription")
	return subscription.Subscribe(subscription.calendar)
}

//POST https://www.googleapis.com/calendar/v3/channels/stop
func (subscription *GoogleSubscription) Delete() (err error) {
	a := subscription.calendar.GetAccount()
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
	log.Warningf("RESPONSE: %s", contents)

	if len(contents) != 0 {
		err = createGoogleResponseError(contents)
	}
	return
}

func (subscription *GoogleSubscription) GetID() string {
	return subscription.ID
}
func (subscription *GoogleSubscription) GetUUID() uuid.UUID {
	return subscription.Uuid
}
func (subscription *GoogleSubscription) GetAccount() AccountManager {
	return subscription.calendar.account
}

func (subscription *GoogleSubscription) GetType() string {
	return subscription.Type
}

func (subscription *GoogleSubscription) setTime() {
	subscription.expirationDate = time.Unix(subscription.Expiration/1000, 0)
}
func (subscription *GoogleSubscription) setCalendar(calendar CalendarManager) (err error) {
	switch calendar.(type) {
	case *GoogleCalendar:
		subscription.calendar = calendar.(*GoogleCalendar)
	default:
		return &customErrors.WrongKindError{Mail: calendar.GetName()}
	}

	return
}

func (subscription *GoogleSubscription) GetExpirationDate() time.Time {
	return subscription.expirationDate
}

func (subscription *GoogleSubscription) GetToken() string {
	return subscription.Token
}
