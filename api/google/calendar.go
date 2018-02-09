package google

import "github.com/TetAlius/GoSyncMyCalendars/api"

type CalendarListResponse struct {
	Kind          string      `json:"kind"`
	Etag          string      `json:"etag"`
	NextPageToken string      `json:"nextPageToken"`
	NextSyncToken string      `json:"nextSyncToken"`
	Calendars     []*Calendar `json:"items"`
}

type Calendar struct {
	//From CalendarLIST resource
	ID              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description"`
	Location        string `json:"location"`
	TmeZone         string `json:"timeZone"`
	SummaryOverride string `json:"summaryOverride"`
	ColorId         string `json:"colorId"`
	BackgroundColor string `json:"backgroundColor"`
	ForegroundColor string `json:"foregroundColor"`
	Hidden          bool   `json:"hidden"`
	Selected        bool   `json:"selected"`
	// Only valid accessRoles with 'writer' or 'owner'
	AccessRole           string     `json:"accessRole"`
	DefaultReminders     []Reminder `json:"defaultReminders"`
	Primary              bool       `json:"primary"`
	Deleted              bool       `json:"deleted"`
	ConferenceProperties `json:"conferenceProperties"`
	NotificationSetting  `json:"notificationSettings"`
}

type Reminder struct {
	Method  string `json:"method"`
	Minutes int32  `json:"minutes"`
}

type NotificationSetting struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	Type   string `json:"type"`
	Method string `json:"method"`
}

type ConferenceProperties struct {
	AllowedConferenceSolutionTypes []string `json:"allowedConferenceSolutionTypes"`
}

// PUT https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (Calendar) Update(api.AccountManager) error {
	panic("implement me")
}

// DELETE https://www.googleapis.com/calendar/v3/users/me/calendarList/{calendarId}
func (Calendar) Delete(api.AccountManager) error {
	panic("implement me")
}

// POST https://www.googleapis.com/calendar/v3/users/me/calendarList
func (Calendar) Create(api.AccountManager) error {
	panic("implement me")
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (Calendar) GetAllEvents(api.AccountManager) ([]api.EventManager, error) {
	panic("implement me")
}

// GET https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (Calendar) GetEvent(api.AccountManager, string) (api.EventManager, error) {
	panic("implement me")
}

func (calendar Calendar) GetID() string {
	return calendar.ID
}
