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

// Function that parses the JSON of the request to a GoogleAccount
func NewGoogleAccount(contents []byte) (a *GoogleAccount, err error) {
	err = json.Unmarshal(contents, &a)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshaling google responses: %s", err.Error()))
	}

	log.Debugf("%s", contents)

	// preferred is ignored on google
	email, _, err := util.MailFromToken(strings.Split(a.TokenID, "."))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error retrieving google mail: %s", err.Error()))
	}

	a.Email = email
	return
}

// Function that returns a GoogleAccount given specific info
func RetrieveGoogleAccount(tokenType string, refreshToken string, email string, kind int, accessToken string) (a *GoogleAccount) {
	a = new(GoogleAccount)
	a.TokenType = tokenType
	a.RefreshToken = refreshToken
	a.Email = email
	a.Kind = kind
	a.AccessToken = accessToken
	return
}

// Method to refresh the access to the google account
func (a *GoogleAccount) Refresh() (err error) {
	client := http.Client{}

	route, err := util.CallAPIRoot("google/token/uri")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugf(route)

	params, err := util.CallAPIRoot("google/token/refresh-params")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	log.Debugln(a.RefreshToken)
	log.Debugln(fmt.Sprintf(params, a.RefreshToken))
	req, err := http.NewRequest(http.MethodPost,
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

	err = json.Unmarshal(contents, &a)
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error with the google request: %s", err.Error()))
	}
	return

}

// Method that retrieves all calendars from account
//
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList
func (a *GoogleAccount) GetAllCalendars() (calendars []CalendarManager, err error) {

	log.Debugln("getAllCalendars google")
	route, err := util.CallAPIRoot("google/calendar-list")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	queryParams := map[string]string{"minAccessRole": "writer"}
	contents, err :=
		util.DoRequest(
			http.MethodGet,
			route,
			nil,
			headers, queryParams)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all calendars for email %s. %s", a.Mail(), err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(GoogleCalendarListResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	for _, s := range calendarResponse.Calendars {
		s.SetAccount(a)
		calendars = append(calendars, s)
	}
	return
}

// Method that retrieves one calendar given an ID
//
// GET https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarID}
func (a *GoogleAccount) GetCalendar(calendarID string) (calendar CalendarManager, err error) {
	log.Debugln("getCalendar google")
	route, err := util.CallAPIRoot("google/calendars/id")
	log.Debugln(route)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	contents, err :=
		util.DoRequest(
			http.MethodGet,
			fmt.Sprintf(route, url.QueryEscape(calendarID)),
			nil,
			headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting calendar for email %s. %s", a.Email, err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(GoogleCalendar)
	err = json.Unmarshal(contents, &calendarResponse)
	calendarResponse.SetAccount(a)

	return calendarResponse, err

}

// Method that returns the principal calendar from the account
//
// GET https://www.googleapis.com/calendar/v3/calendars/primary
func (a *GoogleAccount) GetPrimaryCalendar() (calendar CalendarManager, err error) {
	log.Debugln("getPrimaryCalendar google")
	route, err := util.CallAPIRoot("google/calendars/primary")
	if err != nil {
		return calendar, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()

	contents, err :=
		util.DoRequest(
			http.MethodGet,
			route,
			nil,
			headers, nil)

	if err != nil {
		return calendar, errors.New(fmt.Sprintf("error getting primary calendar for email %s. %s", a.Email, err.Error()))
	}
	err = createGoogleResponseError(contents)
	if err != nil {
		return nil, err
	}

	calendarResponse := new(GoogleCalendar)
	err = json.Unmarshal(contents, &calendarResponse)
	calendar = calendarResponse
	calendar.SetAccount(a)
	return
}

// Method that format the authorization request
func (a *GoogleAccount) AuthorizationRequest() string {
	return fmt.Sprintf("%s %s", a.TokenType, a.AccessToken)
}

// Method that returns the mail associated with the account
func (a *GoogleAccount) Mail() string {
	return a.Email
}

// Method that sets which kind of account is
func (a *GoogleAccount) SetKind(kind int) {
	a.Kind = kind
}

// Method that returns the token type
func (a *GoogleAccount) GetTokenType() string {
	return a.TokenType
}

// Method that returns the refresh token
func (a *GoogleAccount) GetRefreshToken() string {
	return a.RefreshToken
}

// Method that returns the kind of the account
func (a *GoogleAccount) GetKind() int {
	return a.Kind
}

// Method that returns the access token
func (a *GoogleAccount) GetAccessToken() string {
	return a.AccessToken
}

// Method that returns the internal ID given to the account on DB
func (a *GoogleAccount) GetInternalID() int {
	return a.InternID
}

// Method that sets all synced calendars associated with the account
func (a *GoogleAccount) SetCalendars(calendars []CalendarManager) {
	a.calendars = calendars
}

// Method that returns all synced calendars associated with the account
func (a *GoogleAccount) GetSyncCalendars() []CalendarManager {
	return a.calendars
}
