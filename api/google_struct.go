package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type GoogleError struct {
	GoogleConcreteError `json:"error,omitempty"`
}

type GoogleConcreteError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type GoogleAccount struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
	Email        string
}

type GoogleCalendarListResponse struct {
	Kind          string            `json:"kind"`
	Etag          string            `json:"etag"`
	NextPageToken string            `json:"nextPageToken"`
	NextSyncToken string            `json:"nextSyncToken"`
	Calendars     []*GoogleCalendar `json:"items"`
}

type GoogleCalendar struct {
	account *GoogleAccount
	//From CalendarLIST resource
	ID              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description,omitempty"`
	Location        string `json:"location,omitempty"`
	TmeZone         string `json:"timeZone,omitempty"`
	SummaryOverride string `json:"summaryOverride,omitempty"`
	ColorId         string `json:"colorId,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	ForegroundColor string `json:"foregroundColor,omitempty"`
	Hidden          bool   `json:"hidden,omitempty"`
	Selected        bool   `json:"selected,omitempty"`
	// Only valid accessRoles with 'writer' or 'owner'
	AccessRole                 string           `json:"accessRole,omitempty"`
	DefaultReminders           []GoogleReminder `json:"defaultReminders,omitempty"`
	Primary                    bool             `json:"primary,omitempty"`
	Deleted                    bool             `json:"deleted,omitempty"`
	GoogleConferenceProperties `json:"conferenceProperties,omitempty"`
	GoogleNotificationSetting  `json:"notificationSettings,omitempty"`
}

type GoogleReminder struct {
	Method  string `json:"method,omitempty"`
	Minutes int32  `json:"minutes,omitempty"`
}

type GoogleNotificationSetting struct {
	Notifications []GoogleNotification `json:"notifications,omitempty"`
}

type GoogleNotification struct {
	Type   string `json:"type,omitempty"`
	Method string `json:"method,omitempty"`
}

type GoogleConferenceProperties struct {
	AllowedConferenceSolutionTypes []string `json:"allowedConferenceSolutionTypes,omitempty"`
}

type GoogleEventList struct {
	Events []*GoogleEvent `json:"items"`
}

type GoogleEvent struct {
	Calendar  *GoogleCalendar `json:"-"`
	Relations []string        `json:"-"`

	ID string `json:"id"`

	Subject     string      `json:"summary,omitempty" sync:"Subject"`
	Description string      `json:"description,omitempty" sync:"Description"`
	StartsAt    time.Time   `json:"-" sync:"StartsAt"`
	EndsAt      time.Time   `json:"-" sync:"EndsAt"`
	IsAllDay    bool        `json:"-" sync:"IsAllDay"`
	Start       *GoogleTime `json:"start,omitempty"`
	End         *GoogleTime `json:"end,omitempty"`

	Status             string   `json:"status,omitempty"`
	ColorID            string   `json:"colorId,omitempty"`
	EndTimeUnspecified bool     `json:"endTimeUnspecified,omitempty"`
	Recurrence         []string `json:"recurrence,omitempty"`
	RecurringEventId   string   `json:"recurringEventId,omitempty"`
	Transparency       string   `json:"transparency,omitempty"`
	Visibility         string   `json:"visibility,omitempty"`
	ICalUID            string   `json:"iCalUID,omitempty"`
	Sequence           int32    `json:"sequence,omitempty"`
	HangoutLink        string   `json:"hangoutLink,omitempty"`
	Locked             bool     `json:"locked,omitempty"`

	OriginalStartTime *GoogleTime           `json:"originalStartTime,omitempty"`
	Attendees         []GooglePerson        `json:"attendees,omitempty"`
	Gadget            *GoogleGadget         `json:"gadget,omitempty"`
	ConferenceData    *GoogleConferenceData `json:"conferenceData,omitempty"`
	Reminders         GoogleEventReminder   `json:"reminders,omitempty"`
	Source            *GoogleSource         `json:"source,omitempty"`
	Attachments       []GoogleAttachment    `json:"attachments,omitempty"`
	Organizer         *GooglePerson         `json:"organizer,omitempty"`

	//Not to sync
	Link                    string `json:"htmlLink,omitempty"`
	Created                 string `json:"created,omitempty"`
	Updated                 string `json:"updated,omitempty"`
	Location                string `json:"location,omitempty"`
	AttendeesOmitted        bool   `json:"attendeesOmitted,omitempty"`
	AnyoneCanAddSelf        bool   `json:"anyoneCanAddSelf,omitempty"`
	GuestsCanInviteOthers   bool   `json:"guestsCanInviteOthers,omitempty"`
	GuestsCanModify         bool   `json:"guestsCanModify,omitempty"`
	GuestsCanSeeOtherGuests bool   `json:"guestsCanSeeOtherGuests,omitempty"`
	PrivateCopy             bool   `json:"privateCopy,omitempty"`
}

type GooglePerson struct {
	ID               string `json:"id,omitempty"`
	Email            string `json:"email,omitempty"`
	Name             string `json:"displayName,omitempty"`
	Self             bool   `json:"self,omitempty"`
	Organizer        bool   `json:"organizer,omitempty"`
	Resource         bool   `json:"resource,omitempty"`
	Optional         bool   `json:"optional,omitempty"`
	ResponseStatus   string `json:"responseStatus,omitempty"`
	Comment          string `json:"comment,omitempty"`
	AdditionalGuests int32  `json:"additionalGuests,omitempty"`
}

type GoogleTime struct {
	Date string `json:"date,omitempty"`
	//time.RFC3339 gives TimeZone inside string
	DateTime string `json:"dateTime,omitempty"`
	//Ignore TimeZone as the json returns the original TimeZon
	//Although it is always asked in UTC it may cause confusion
	TimeZone string `json:"-"` //`json:"timeZone,omitempty"`
}

type GoogleConferenceData struct {
	CreateRequest *GoogleCreateRequest `json:"createRequest,omitempty"`
	EntryPoints   []GoogleEntryPoint   `json:"entryPoints,omitempty"`
	//ConferenceSolution ConferenceSolution `json:"conferenceSolution"`
	ConferenceID string `json:"conferenceId,omitempty"`
	Signature    string `json:"signature,omitempty"`
}

type GoogleCreateRequest struct {
	RequestID             string             `json:"requestId,omitempty"`
	ConferenceSolutionKey *GoogleSolutionKey `json:"conferenceSolutionKey,omitempty"`
	Status                *GoogleStatus      `json:"status,omitempty"`
}

type GoogleSolutionKey struct {
	Type string `json:"type,omitempty"`
}

type GoogleStatus struct {
	Code string `json:"statusCode,omitempty"`
}

type GoogleEntryPoint struct {
	EntryPointType string `json:"entryPointType,omitempty"`
	URI            string `json:"uri,omitempty"`
	Label          string `json:"label,omitempty"`
	Pin            string `json:"pin,omitempty"`
	AccessCode     string `json:"accessCode,omitempty"`
	MeetingCode    string `json:"meetingCode,omitempty"`
	Passcode       string `json:"passcode,omitempty"`
	Password       string `json:"password,omitempty"`
}

type GoogleGadget struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
	IconLink string `json:"iconLink,omitempty"`
	Width    int32  `json:"width,omitempty"`
	Height   int32  `json:"height,omitempty"`
	Display  string `json:"display,omitempty"`
	//Preferences
}

type GoogleEventReminder struct {
	UseDefault bool             `json:"useDefault,omitempty"`
	Overrides  []GoogleReminder `json:"overrides,omitempty"`
}

type GoogleSource struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

type GoogleAttachment struct {
	FileURL  string `json:"fileUrl,omitempty"`
	Title    string `json:"title,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	IconLink string `json:"iconLink,omitempty"`
	FileID   string `json:"fileId,omitempty"`
}

func createGoogleResponseError(contents []byte) (err error) {
	e := new(GoogleError)
	err = json.Unmarshal(contents, &e)
	if e.Code != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %d. message: %s", e.Code, e.Message))
	}
	return nil
}
