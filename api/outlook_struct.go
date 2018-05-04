package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type OutlookError struct {
	OutlookConcreteError `json:"error,omitempty"`
}
type OutlookConcreteError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type OutlookAccount struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

type OutlookCalendarResponse struct {
	OdataContext string `json:"@odata.context"`
	*OutlookCalendar
}

type OutlookCalendarListResponse struct {
	OdataContext string             `json:"@odata.context"`
	Calendars    []*OutlookCalendar `json:"value"`
}

// CalendarInfo TODO
type OutlookCalendar struct {
	OdataID string `json:"@odata.id,omitempty"`

	CalendarView        []OutlookEvent       `json:"CalendarView,omitempty"`
	CanEdit             bool                 `json:"CanEdit,omitempty"`
	CanShare            bool                 `json:"CanShare,omitempty"`
	CanViewPrivateItems bool                 `json:"CanViewPrivateItems,omitempty"`
	ChangeKey           string               `json:"ChangeKey,omitempty"`
	Color               OutlookCalendarColor `json:"Color,omitempty"`
	Events              []OutlookEvent       `json:"Events,omitempty"`
	ID                  string               `json:"Id"`
	Name                string               `json:"Name,omitempty"`
	Owner               OutlookEmailAddress  `json:"Owner,omitempty"`

	//	IsDefaultCalendar             bool         `json:"IsDefaultCalendar,omitempty"`
	//	IsShared                      bool         `json:"IsShared,omitempty"`
	//	IsSharedWithMe                bool         `json:"IsSharedWithMe,omitempty"`
	//	MultiValueExtendedProperties  []Properties `json:"MultiValueExtendedProperties"`
	//	SingleValueExtendedProperties []Properties `json:"SingleValueExtendedProperties"`
}

// Specifies the color theme to distinguish the calendar from other calendars in a UI.
// The property values are: LightBlue=0, LightGreen=1, LightOrange=2, LightGray=3,
// LightYellow=4, LightTeal=5, LightPink=6, LightBrown=7, LightRed=8, MaxColor=9, Auto=-1
type OutlookCalendarColor string

type OutlookEmailAddress struct {
	Address string `json:"Address,omitempty"`
	Name    string `json:"Name,omitempty"`
}

type OutlookEventResponse struct {
	OdataContext string `json:"@odata.context"`
	*OutlookEvent
}

type OutlookEventListResponse struct {
	OdataContext string          `json:"@odata.context"`
	Events       []*OutlookEvent `json:"value"`
}

type OutlookEvent struct {
	OdataID string `json:"@odata.id,omitempty"`
	//OdataEtag string `json:"@odata.etag,omitempty"`
	Relations   []string  `json:"-"`
	StartsAt    time.Time `json:"-"`
	EndsAt      time.Time `json:"-"`
	Description string    `json:"BodyPreview,omitempty"`

	Calendar *OutlookCalendar `json:"-"`

	Attachments []OutlookAttachment `json:"Attachments,omitempty"`
	Attendees   []OutlookAttendee   `json:"Attendees,omitempty"`
	Body        *OutlookItemBody    `json:"Body,omitempty"`

	Categories      []string                 `json:"Categories,omitempty"`
	ChangeKey       string                   `json:"ChangeKey,omitempty"`
	CreatedDateTime string                   `json:"CreatedDateTime,omitempty"` //"2014-10-19T23:13:47.3959685Z"
	End             *OutlookDateTimeTimeZone `json:"End,omitempty"`
	HasAttachments  bool                     `json:"HasAttachments,omitempty"`
	//ICalUID                    string               `json:"iCalUID,omitempty"`
	ID                         string                      `json:"Id"`
	Importance                 OutlookImportance           `json:"Importance,omitempty"`
	Instances                  []OutlookEvent              `json:"Instances,omitempty"`
	IsAllDay                   bool                        `json:"IsAllday,omitempty"`
	IsCancelled                bool                        `json:"IsCancelled,omitempty"`
	IsOrganizer                bool                        `json:"IsOrganizer,omitempty"`
	IsReminderOn               bool                        `json:"IsReminderOn,omitempty"`
	LastModifiedDateTime       string                      `json:"LastModifiedDateTime,omitempty"` //"2014-10-19T23:13:47.6772234Z"
	Location                   *OutlookLocation            `json:"Location,omitempty"`
	OnlineMeetingUrl           string                      `json:"OnlineMeetingUrl,omitempty"`
	Organizer                  *OutlookRecipient           `json:"Organizer,omitempty"`
	OriginalStartTimeZone      string                      `json:"OriginalStartTimeZone,omitempty"`
	OriginalEndTimeZone        string                      `json:"OriginalEndTimeZone,omitempty"`
	Recurrence                 *OutlookPatternedRecurrence `json:"Recurrence,omitempty"`
	ReminderMinutesBeforeStart int32                       `json:"ReminderMinutesBeforeStart,omitempty"`
	ResponseRequested          bool                        `json:"ResponseRequested,omitempty"`
	ResponseStatus             *OutlookResponseStatus      `json:"ResponseStatus,omitempty"`
	Sensitivity                OutlookSensitivity          `json:"Sensitivity,omitempty"`
	SeriesMasterID             string                      `json:"SeriesMasterId,omitempty"`
	ShowAs                     OutlookFreeBusyStatus       `json:"ShowAs,omitempty"`
	Start                      *OutlookDateTimeTimeZone    `json:"Start,omitempty"`
	Subject                    string                      `json:"Subject,omitempty"`
	Type                       OutlookEventType            `json:"Type,omitempty"`
	WebLink                    string                      `json:"WebLink,omitempty"`

	//Extensions                 []Extension         `json:"Extensions"`
}

type OutlookAttachment struct {
	ContentType          string `json:"ContentType,omitempty"`
	IsInline             bool   `json:"IsInline,omitempty"`
	LastModifiedDateTime string `json:"LastModifiedDateTime,omitempty"`
	Name                 string `json:"Name,omitempty"`
	Size                 int32  `json:"Size,omitempty"`
}

type OutlookAttendee struct {
	Recipient *OutlookRecipient `json:"EmailAddress,omitempty"`
	Status    *OutlookStatus    `json:"Status,omitempty"`
	Type      string            `json:"Type,omitempty"`
}

type OutlookRecipient struct {
	EmailAddress *OutlookEmailAddress `json:"EmailAddress,omitempty"`
}

type OutlookStatus struct {
	Response string `json:"Response,omitempty"`
	Time     string `json:"Time,omitempty"`
}

type OutlookItemBody struct {
	ContentType string `json:"ContentType,omitempty"`
	Content     string `json:"Content,omitempty"`
}

type OutlookDateTimeTimeZone struct {
	DateTime string `json:"DateTime,omitempty"`
	TimeZone string `json:"TimeZone,omitempty"`
}

// The importance of the event: Low, Normal, High.
type OutlookImportance string

type OutlookLocation struct {
	Address              OutlookPhysicalAddress `json:"Address,omitempty"`
	Coordinates          OutlookGeoCoordinates  `json:"Coordinates,omitempty"`
	DisplayName          string                 `json:"DisplayName,omitempty"`
	LocationEmailAddress string                 `json:"LocationEmailAddress,omitempty"`
}

type OutlookPhysicalAddress struct {
	Street          string `json:"Street,omitempty"`
	City            string `json:"City,omitempty"`
	State           string `json:"State,omitempty"`
	CountryOrRegion string `json:"CountryOrRegion,omitempty"`
	PostalCode      string `json:"PostalCode,omitempty"`
}

type OutlookGeoCoordinates struct {
	Altitude         float64 `json:"Altitude,omitempty"`
	Latitude         float64 `json:"Latitude,omitempty"`
	Longitude        float64 `json:"Longitude,omitempty"`
	Accuracy         float64 `json:"Accuracy,omitempty"`
	AltitudeAccuracy float64 `json:"AltitudeAccuracy,omitempty"`
}

// The type of location: Default, ConferenceRoom, HomeAddress, BusinessAddress,OutlookGeoCoordinates, StreetAddress, Hotel, Restaurant, LocalBusiness, PostalAddress.
type OutlookLocationType int32

type OutlookPatternedRecurrence struct {
	Pattern            OutlookRecurrencePattern `json:"Pattern,omitempty"`
	RecurrenceTimeZone string                   `json:"RecurrenceTimeZone,omitempty"`
	Range              OutlookRecurrenceRange   `json:"Range,omitempty"`
}

type OutlookRecurrencePattern struct {
	Type           OutlookRecurrencePatternType `json:"Type,omitempty"`
	Interval       int32                        `json:"Interval,omitempty"`
	DayOfMonth     int32                        `json:"DayOfMonth,omitempty"`
	Month          int32                        `json:"Month,omitempty"`
	DaysOfWeek     []OutlookDayOfWeek           `json:"DaysOfWeek,omitempty"`
	FirstDayOfWeek OutlookDayOfWeek             `json:"DayOfWeek,omitempty"`
	Index          OutlookWeekIndex             `json:"Index,omitempty"`
}

// The recurrence pattern type: Daily = 0, Weekly = 1, AbsoluteMonthly = 2, RelativeMonthly = 3, AbsoluteYearly = 4, RelativeYearly = 5.
type OutlookRecurrencePatternType int32

// The day of the week: Sunday = 0, Monday = 1, Tuesday = 2, Wednesday = 3, Thursday = 4, Friday = 5, Saturday = 6.
type OutlookDayOfWeek int32

// The week index: First = 0, Second = 1, Third = 2, Fourth = 3, Last = 4.
type OutlookWeekIndex int32

type OutlookRecurrenceRange struct {
	Type                OutlookRecurrenceRangeType `json:"Type,omitempty"`
	StartDate           string                     `json:"StartDate,omitempty"` //"2014-10-19T23:13:47.3959685Z" TODO
	EndDate             string                     `json:"EndDate,omitempty"`   //"2014-10-19T23:13:47.3959685Z" TODO
	NumberOfOccurrences int32                      `json:"NumberOfOccurrences,omitempty"`
}

// The recurrence range: EndDate = 0, NoEnd = 1, Numbered = 2.
type OutlookRecurrenceRangeType int32

type OutlookResponseStatus struct {
	Response OutlookResponseType `json:"Response,omitempty"`
	Time     string              `json:"Time,omitempty"`
}

// The response type: None, Organizer, TentativelyAccepted, Accepted, Declined, NotResponded.
type OutlookResponseType string

// Indicates the level of privacy for the event: Normal, Personal, Private, Confidential.
type OutlookSensitivity string

// The status to show: Free, Tentative, Busy, Oof, WorkingElsewhere, Unknown.
type OutlookFreeBusyStatus string

// The event type: SingleInstance, Occurrence, Exception, SeriesMaster.
type OutlookEventType string

type OutlookSubscription struct {
	Type               string `json:"@odata.type,omitempty"`
	Resource           string `json:"Resource,omitempty"`
	NotificationURL    string `json:"NotificationURL,omitempty"`
	ChangeType         string `json:"ChangeType,omitempty"`
	ID                 string `json:"id,omitempty"`
	ClientState        string `json:"ClientState,omitempty"`
	ExpirationDateTime string `json:"SubscriptionExpirationDateTime,omitempty"`
}

type OutlookNotification struct {
	Subscriptions []OutlookSubscriptionNotification `json:"value"`
}

type OutlookSubscriptionNotification struct {
	SubscriptionID                 string              `json:"SubscriptionId"`
	SubscriptionExpirationDateTime string              `json:"SubscriptionExpirationDateTime"`
	SequenceNumber                 int32               `json:"SequenceNumber"`
	Date                           OutlookResourceData `json:"ResourceData"`
}

type OutlookResourceData struct {
	ID string `json:"Id"`
}

func createOutlookResponseError(contents []byte) (err error) {
	e := new(OutlookError)
	err = json.Unmarshal(contents, &e)
	if len(e.Code) != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %s. message: %s", e.Code, e.Message))
	}
	return nil
}