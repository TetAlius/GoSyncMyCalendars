package api

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

func NewOutlookAccount(contents []byte) (a *OutlookAccount, err error) {
	err = json.Unmarshal(contents, &a)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error unmarshaling outlook response: %s", err.Error()))
	}

	email, preferred, err := util.MailFromToken(strings.Split(a.TokenID, "."), "=")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving outlook mail: %s", err.Error()))
	}
	a.AnchorMailbox = email
	a.PreferredUsername = preferred
	return
}

func RetrieveOutlookAccount(tokenType string, refreshToken string, email string, kind int, accessToken string, id int) (a *OutlookAccount) {
	a = new(OutlookAccount)
	a.TokenType = tokenType
	a.RefreshToken = refreshToken
	a.AnchorMailbox = email
	a.Kind = kind
	a.AccessToken = accessToken
	a.InternID = id
	return
}

func (a *OutlookAccount) Refresh() (err error) {
	client := http.Client{}
	//check if token is DEAD!!!

	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)

	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	params, err := util.CallAPIRoot("outlook/token/refresh-params")
	log.Debugf("Params: %s", fmt.Sprintf(params, a.RefreshToken))
	if err != nil {
		return errors.New(fmt.Sprintf("error generating params: %s", err.Error()))
	}

	req, err := http.NewRequest(http.MethodPost,
		route,
		strings.NewReader(fmt.Sprintf(params, a.RefreshToken)))

	if err != nil {
		return errors.New(fmt.Sprintf("error creating new request: %s", err.Error()))
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error doing outlook request: %s", err.Error()))
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading response body from outlook request: %s", err.Error()))
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

	log.Debugf("\nTokenType: %s\nExpiresIn: %d\nAccessToken: %s\nRefreshToken: %s\nTokenID: %s\nAnchorMailbox: %s\nPreferredUsername: %t",
		a.TokenType, a.ExpiresIn, a.AccessToken, a.RefreshToken, a.TokenID, a.AnchorMailbox, a.PreferredUsername)

	log.Debugf("%s\n", contents)
	err = json.Unmarshal(contents, &a)
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error with the outlook request: %s", err.Error()))
	}
	return
}

func (a *OutlookAccount) GetAllCalendars() (calendars []CalendarManager, err error) {
	log.Debugln("getAllCalendars outlook")

	route, err := util.CallAPIRoot("outlook/calendars")
	if err != nil {
		log.Errorf("%s", err.Error())
		return calendars, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()
	queryParams := map[string]string{"$filter": "CanEdit eq false"}

	contents, err := util.DoRequest(http.MethodGet,
		route,
		nil,
		headers, queryParams)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting all calendars for email %s. %s", a.AnchorMailbox, err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(OutlookCalendarListResponse)
	err = json.Unmarshal(contents, &calendarResponse)

	calendars = make([]CalendarManager, len(calendarResponse.Calendars))
	for i, s := range calendarResponse.Calendars {
		s.SetAccount(a)
		calendars[i] = s
	}
	return
}

func (a *OutlookAccount) GetCalendar(calendarID string) (calendar CalendarManager, err error) {
	if len(calendarID) == 0 {
		return calendar, errors.New("no ID for calendar was given")
	}
	log.Debugln("getCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/id")
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		return
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodGet,
		fmt.Sprintf(route, calendarID),
		nil,
		headers, nil)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting a calendar for email %s. %s", a.AnchorMailbox, err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(OutlookCalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)
	calendar = calendarResponse.OutlookCalendar
	calendar.SetAccount(a)
	return
}

func (a *OutlookAccount) GetPrimaryCalendar() (calendar CalendarManager, err error) {
	log.Debugln("getPrimaryCalendar outlook")

	route, err := util.CallAPIRoot("outlook/calendars/primary")
	if err != nil {
		log.Errorf("%s", err.Error())
		return calendar, errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = a.AuthorizationRequest()
	headers["X-AnchorMailbox"] = a.Mail()

	contents, err := util.DoRequest(http.MethodGet,
		route,
		nil,
		headers, nil)

	if err != nil {
		log.Errorf("%s", err.Error())
		return calendar, errors.New(fmt.Sprintf("error getting primary calendar for email %s. %s", a.AnchorMailbox, err.Error()))
	}
	err = createOutlookResponseError(contents)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s\n", contents)

	calendarResponse := new(OutlookCalendarResponse)
	err = json.Unmarshal(contents, &calendarResponse)
	calendar = calendarResponse.OutlookCalendar
	calendar.SetAccount(a)

	return
}

func (a *OutlookAccount) AuthorizationRequest() (auth string) {
	return fmt.Sprintf("%s %s", a.TokenType, a.AccessToken)
}

func (a *OutlookAccount) Mail() string {
	return a.AnchorMailbox
}

func (a *OutlookAccount) SetKind(kind int) {
	a.Kind = kind
}

func (a *OutlookAccount) GetTokenType() string {
	return a.TokenType
}

func (a *OutlookAccount) GetRefreshToken() string {
	return a.RefreshToken
}

func (a *OutlookAccount) GetKind() int {
	return a.Kind
}

func (a *OutlookAccount) GetAccessToken() string {
	return a.AccessToken
}

func (a *OutlookAccount) GetInternalID() int {
	return a.InternID
}

func (a *OutlookAccount) SetCalendars(calendars []CalendarManager) {
	a.calendars = calendars
}

func (a *OutlookAccount) GetSyncCalendars() []CalendarManager {
	return a.calendars
}
