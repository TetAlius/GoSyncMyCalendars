package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"net/url"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func NewGoogleAccount(contents []byte) (a *GoogleAccount, err error) {
	err = json.Unmarshal(contents, &a)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshaling google responses: %s", err.Error()))
	}

	log.Debugf("%s", contents)

	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(a.TokenID, "."), "==")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error retrieving google mail: %s", err.Error()))
	}

	a.Email = email
	return
}

func (a *GoogleAccount) Refresh() (err error) {
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
		e := new(RefreshError)
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
func (a *GoogleAccount) GetAllCalendars() (calendars []CalendarManager, err error) {
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
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(GoogleCalendarListResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	calendars = make([]CalendarManager, len(calendarResponse.Calendars))
	for i, s := range calendarResponse.Calendars {
		calendars[i] = s
	}
	return
}

// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func (a *GoogleAccount) GetCalendar(calendarID string) (calendar CalendarManager, err error) {
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
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(GoogleCalendar)
	err = json.Unmarshal(contents, &calendarResponse)
	log.Debugln(contents)

	return calendarResponse, err

}

// GET https://www.googleapis.com/calendar/v3/calendars/primary
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/primary This is the one used
func (a *GoogleAccount) GetPrimaryCalendar() (calendar CalendarManager, err error) {
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
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(GoogleCalendar)
	err = json.Unmarshal(contents, &calendarResponse)
	log.Debugln(contents)

	return calendarResponse, err
}

func (a *GoogleAccount) AuthorizationRequest() string {
	return fmt.Sprintf("%s %s", a.TokenType, a.AccessToken)
}

func (a *GoogleAccount) Mail() string {
	return a.Email
}
