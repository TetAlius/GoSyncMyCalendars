package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"net/url"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func NewAccount(contents []byte) (a *Account, err error) {
	err = json.Unmarshal(contents, &a)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling google responses: %s", err.Error()))
	}

	log.Debugf("%s", contents)

	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(a.TokenID, "."), "==")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving google mail: %s", err.Error()))
	}

	a.Email = email
	return
}

func (a *Account) Refresh() (err error) {
	client := http.Client{}

	route, err := util.CallAPIRoot("google/token/uri")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	params, err := util.CallAPIRoot("google/token/refresh-params")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugln(a.RefreshToken)
	log.Debugln(fmt.Sprintf(params, a.RefreshToken))
	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, a.RefreshToken)))

	if err != nil {
		return errors.New(fmt.Sprintf("error creating new request: %s", err.Error()))
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error doing google request: %s", err.Error()))
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading response body from google request: %s", err.Error()))
	}
	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		e := new(api.RefreshError)
		_ = json.Unmarshal(contents, &e)
		if len(e.Code) != 0 && len(e.Message) != 0 {
			log.Errorln(e.Code)
			log.Errorln(e.Message)
			return e
		}
	}
	log.Debugf("%s\n", contents)

	err = json.Unmarshal(contents, &a)
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error with the outlook request: %s", err.Error()))
	}
	return

}

//GET https://www.googleapis.com/calendar/v3/users/me/calendarList
func (a *Account) GetAllCalendars() (calendars []api.CalendarManager, err error) {
	log.Debugln("getAllCalendars google")
	route, err := util.CallAPIRoot("google/calendar-list")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err :=
		util.DoRequest(
			"GET",
			route,
			nil,
			a.AuthorizationRequest(),
			"")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all calendars for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(CalendarListResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	calendars = make([]api.CalendarManager, len(calendarResponse.Calendars))
	for i, s := range calendarResponse.Calendars {
		calendars[i] = s
	}
	return
}

// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func (a *Account) GetCalendar(calendarID string) (calendar api.CalendarManager, err error) {
	log.Debugln("getCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	log.Debugln(route)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err :=
		util.DoRequest(
			"GET",
			fmt.Sprintf(route, url.QueryEscape(calendarID)),
			nil,
			a.AuthorizationRequest(),
			"")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting calendar for email %s. %s", a.Email, err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(Calendar)
	err = json.Unmarshal(contents, &calendarResponse)
	log.Debugln(contents)

	return calendarResponse, err

}

// GET https://www.googleapis.com/calendar/v3/calendars/primary
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/primary This is the one used
func (a *Account) GetPrimaryCalendar() (calendar api.CalendarManager, err error) {
	log.Debugln("getPrimaryCalendar google")
	route, err := util.CallAPIRoot("google/calendars/primary")
	if err != nil {
		return calendar, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	contents, err :=
		util.DoRequest(
			"GET",
			route,
			nil,
			a.AuthorizationRequest(),
			"")

	if err != nil {
		return calendar, errors.New(fmt.Sprintf("error getting primary calendar for email %s. %s", a.Email, err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(Calendar)
	err = json.Unmarshal(contents, &calendarResponse)
	log.Debugln(contents)

	return calendarResponse, err
}

func (a *Account) AuthorizationRequest() string {
	return fmt.Sprintf("%s %s", a.TokenType, a.AccessToken)
}

func (a *Account) Mail() string {
	return a.Email
}
