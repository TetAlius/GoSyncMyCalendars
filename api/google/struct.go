package google

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	ConcreteError `json:"error,omitempty"`
}
type ConcreteError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func createResponseError(contents []byte) (err error) {
	e := new(Error)
	err = json.Unmarshal(contents, &e)
	if e.Code != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %d. message: %s", e.Code, e.Message))
	}
	return nil
}

type Account struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"id_token"`
	Email        string
}

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
	AccessRole           string     `json:"accessRole,omitempty"`
	DefaultReminders     []Reminder `json:"defaultReminders,omitempty"`
	Primary              bool       `json:"primary,omitempty"`
	Deleted              bool       `json:"deleted,omitempty"`
	ConferenceProperties `json:"conferenceProperties,omitempty"`
	NotificationSetting  `json:"notificationSettings,omitempty"`
}

type Reminder struct {
	Method  string `json:"method,omitempty"`
	Minutes int32  `json:"minutes,omitempty"`
}

type NotificationSetting struct {
	Notifications []Notification `json:"notifications,omitempty"`
}

type Notification struct {
	Type   string `json:"type,omitempty"`
	Method string `json:"method,omitempty"`
}

type ConferenceProperties struct {
	AllowedConferenceSolutionTypes []string `json:"allowedConferenceSolutionTypes,omitempty"`
}

type EventList struct {
	Events []*Event `json:"items"`
}
type Event struct {
	Calendar                *Calendar       `json:"calendar,omitempty"`
	Kind                    string          `json:"kind,omitempty"`
	Etag                    string          `json:"etag,omitempty"`
	ID                      string          `json:"id"`
	Status                  string          `json:"status,omitempty"`
	HTMLLink                string          `json:"htmlLink,omitempty"`
	Created                 string          `json:"created,omitempty"`
	Updated                 string          `json:"updated,omitempty"`
	Summary                 string          `json:"summary,omitempty"`
	Description             string          `json:"description,omitempty"`
	Location                string          `json:"location,omitempty"`
	ColorID                 string          `json:"colorId,omitempty"`
	Creator                 *Person         `json:"creator,omitempty"`
	Organizer               *Person         `json:"organizer,omitempty"`
	Start                   *Time           `json:"start,omitempty"`
	End                     *Time           `json:"end,omitempty"`
	EndTimeUnspecified      bool            `json:"endTimeUnspecified,omitempty"`
	Recurrence              []string        `json:"recurrence,omitempty"`
	RecurringEventId        string          `json:"recurringEventId,omitempty"`
	OriginalStartTime       *Time           `json:"originalStartTime,omitempty"`
	Transparency            string          `json:"transparency,omitempty"`
	Visibility              string          `json:"visibility,omitempty"`
	ICalUID                 string          `json:"iCalUID,omitempty"`
	Sequence                int32           `json:"sequence,omitempty"`
	Attendees               []Person        `json:"attendees,omitempty"`
	AttendeesOmitted        bool            `json:"attendeesOmitted,omitempty"`
	HangoutLink             string          `json:"hangoutLink,omitempty"`
	ConferenceData          *ConferenceData `json:"conferenceData,omitempty"`
	Gadget                  *Gadget         `json:"gadget,omitempty"`
	AnyoneCanAddSelf        bool            `json:"anyoneCanAddSelf,omitempty"`
	GuestsCanInviteOthers   bool            `json:"guestsCanInviteOthers,omitempty"`
	GuestsCanModify         bool            `json:"guestsCanModify,omitempty"`
	GuestsCanSeeOtherGuests bool            `json:"guestsCanSeeOtherGuests,omitempty"`
	PrivateCopy             bool            `json:"privateCopy,omitempty"`
	Locked                  bool            `json:"locked,omitempty"`
	Reminders               EventReminder   `json:"reminders,omitempty"`
	Source                  *Source         `json:"source,omitempty"`
	Attachments             []Attachment    `json:"attachments,omitempty"`
	//ExtendedProperties
}

type Person struct {
	ID               string `json:"id,omitempty"`
	Email            string `json:"email,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`
	Self             bool   `json:"self,omitempty"`
	Organizer        bool   `json:"organizer,omitempty"`
	Resource         bool   `json:"resource,omitempty"`
	Optional         bool   `json:"optional,omitempty"`
	ResponseStatus   string `json:"responseStatus,omitempty"`
	Comment          string `json:"comment,omitempty"`
	AdditionalGuests int32  `json:"additionalGuests,omitempty"`
}
type Time struct {
	Date     string `json:"date,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

type ConferenceData struct {
	CreateRequest *CreateRequest `json:"createRequest,omitempty"`
	EntryPoints   []EntryPoint   `json:"entryPoints,omitempty"`
	//ConferenceSolution ConferenceSolution `json:"conferenceSolution"`
	ConferenceID string `json:"conferenceId,omitempty"`
	Signature    string `json:"signature,omitempty"`
}

type CreateRequest struct {
	RequestID             string       `json:"requestId,omitempty"`
	ConferenceSolutionKey *SolutionKey `json:"conferenceSolutionKey,omitempty"`
	Status                *Status      `json:"status,omitempty"`
}

type SolutionKey struct {
	Type string `json:"type,omitempty"`
}

type Status struct {
	Code string `json:"statusCode,omitempty"`
}

type EntryPoint struct {
	EntryPointType string `json:"entryPointType,omitempty"`
	URI            string `json:"uri,omitempty"`
	Label          string `json:"label,omitempty"`
	Pin            string `json:"pin,omitempty"`
	AccessCode     string `json:"accessCode,omitempty"`
	MeetingCode    string `json:"meetingCode,omitempty"`
	Passcode       string `json:"passcode,omitempty"`
	Password       string `json:"password,omitempty"`
}
type Gadget struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
	IconLink string `json:"iconLink,omitempty"`
	Width    int32  `json:"width,omitempty"`
	Height   int32  `json:"height,omitempty"`
	Display  string `json:"display,omitempty"`
	//Preferences
}
type EventReminder struct {
	UseDefault bool       `json:"useDefault,omitempty"`
	Overrides  []Reminder `json:"overrides,omitempty"`
}
type Source struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

type Attachment struct {
	FileURL  string `json:"fileUrl,omitempty"`
	Title    string `json:"title,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	IconLink string `json:"iconLink,omitempty"`
	FileID   string `json:"fileId,omitempty"`
}
