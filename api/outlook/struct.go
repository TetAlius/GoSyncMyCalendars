package outlook

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	ConcreteError `json:"error,omitempty"`
}
type ConcreteError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Account struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

type CalendarResponse struct {
	OdataContext string `json:"@odata.context"`
	*Calendar
}

type CalendarListResponse struct {
	OdataContext string      `json:"@odata.context"`
	Calendars    []*Calendar `json:"value"`
}

// CalendarInfo TODO
type Calendar struct {
	OdataID string `json:"@odata.id,omitempty"`

	CalendarView        []Event       `json:"CalendarView,omitempty"`
	CanEdit             bool          `json:"CanEdit,omitempty"`
	CanShare            bool          `json:"CanShare,omitempty"`
	CanViewPrivateItems bool          `json:"CanViewPrivateItems,omitempty"`
	ChangeKey           string        `json:"ChangeKey,omitempty"`
	Color               CalendarColor `json:"Color,omitempty"`
	Events              []Event       `json:"Events,omitempty"`
	ID                  string        `json:"Id"`
	Name                string        `json:"Name,omitempty"`
	Owner               EmailAddress  `json:"Owner,omitempty"`

	//	IsDefaultCalendar             bool         `json:"IsDefaultCalendar,omitempty"`
	//	IsShared                      bool         `json:"IsShared,omitempty"`
	//	IsSharedWithMe                bool         `json:"IsSharedWithMe,omitempty"`
	//	MultiValueExtendedProperties  []Properties `json:"MultiValueExtendedProperties"`
	//	SingleValueExtendedProperties []Properties `json:"SingleValueExtendedProperties"`
}

// Specifies the color theme to distinguish the calendar from other calendars in a UI.
// The property values are: LightBlue=0, LightGreen=1, LightOrange=2, LightGray=3,
// LightYellow=4, LightTeal=5, LightPink=6, LightBrown=7, LightRed=8, MaxColor=9, Auto=-1
type CalendarColor string

type EmailAddress struct {
	Address string `json:"Address,omitempty"`
	Name    string `json:"Name,omitempty"`
}

type EventResponse struct {
	OdataContext string `json:"@odata.context"`
	*Event
}

type EventListResponse struct {
	OdataContext string   `json:"@odata.context"`
	Events       []*Event `json:"value"`
}

type Event struct {
	OdataID string `json:"@odata.id,omitempty"`
	//OdataEtag string `json:"@odata.etag,omitempty"`

	Calendar *Calendar

	Attachments     []Attachment      `json:"Attachments,omitempty"`
	Attendees       []Attendee        `json:"Attendees,omitempty"`
	Body            *ItemBody         `json:"Body,omitempty"`
	BodyPreview     string            `json:"BodyPreview,omitempty"`
	Categories      []string          `json:"Categories,omitempty"`
	ChangeKey       string            `json:"ChangeKey,omitempty"`
	CreatedDateTime string            `json:"CreatedDateTime,omitempty"` //"2014-10-19T23:13:47.3959685Z"
	End             *DateTimeTimeZone `json:"End,omitempty"`
	HasAttachments  bool              `json:"HasAttachments,omitempty"`
	//ICalUID                    string               `json:"iCalUID,omitempty"`
	ID                         string               `json:"Id"`
	Importance                 Importance           `json:"Importance,omitempty"`
	Instances                  []Event              `json:"Instances,omitempty"`
	IsAllDay                   bool                 `json:"IsAllday,omitempty"`
	IsCancelled                bool                 `json:"IsCancelled,omitempty"`
	IsOrganizer                bool                 `json:"IsOrganizer,omitempty"`
	IsReminderOn               bool                 `json:"IsReminderOn,omitempty"`
	LastModifiedDateTime       string               `json:"LastModifiedDateTime,omitempty"` //"2014-10-19T23:13:47.6772234Z"
	Location                   *Location            `json:"Location,omitempty"`
	OnlineMeetingUrl           string               `json:"OnlineMeetingUrl,omitempty"`
	Organizer                  *Recipient           `json:"Organizer,omitempty"`
	OriginalStartTimeZone      string               `json:"OriginalStartTimeZone,omitempty"`
	OriginalEndTimeZone        string               `json:"OriginalEndTimeZone,omitempty"`
	Recurrence                 *PatternedRecurrence `json:"Recurrence,omitempty"`
	ReminderMinutesBeforeStart int32                `json:"ReminderMinutesBeforeStart,omitempty"`
	ResponseRequested          bool                 `json:"ResponseRequested,omitempty"`
	ResponseStatus             *ResponseStatus      `json:"ResponseStatus,omitempty"`
	Sensitivity                Sensitivity          `json:"Sensitivity,omitempty"`
	SeriesMasterID             string               `json:"SeriesMasterId,omitempty"`
	ShowAs                     FreeBusyStatus       `json:"ShowAs,omitempty"`
	Start                      *DateTimeTimeZone    `json:"Start,omitempty"`
	Subject                    string               `json:"Subject,omitempty"`
	Type                       EventType            `json:"Type,omitempty"`
	WebLink                    string               `json:"WebLink,omitempty"`

	//Extensions                 []Extension         `json:"Extensions"`
}

type Attachment struct {
	ContentType          string `json:"ContentType,omitempty"`
	IsInline             bool   `json:"IsInline,omitempty"`
	LastModifiedDateTime string `json:"LastModifiedDateTime,omitempty"`
	Name                 string `json:"Name,omitempty"`
	Size                 int32  `json:"Size,omitempty"`
}

type Attendee struct {
	Recipient *Recipient `json:"EmailAddress,omitempty"`
	Status    *Status    `json:"Status,omitempty"`
	Type      string     `json:"Type,omitempty"`
}

type Recipient struct {
	EmailAddress *EmailAddress `json:"EmailAddress,omitempty"`
}

type Status struct {
	Response string `json:"Response,omitempty"`
	Time     string `json:"Time,omitempty"`
}

type ItemBody struct {
	ContentType string `json:"ContentType,omitempty"`
	Content     string `json:"Content,omitempty"`
}

type DateTimeTimeZone struct {
	DateTime string `json:"DateTime,omitempty"`
	TimeZone string `json:"TimeZone,omitempty"`
}

// The importance of the event: Low, Normal, High.
type Importance string

type Location struct {
	Address              PhysicalAddress `json:"Address,omitempty"`
	Coordinates          GeoCoordinates  `json:"Coordinates,omitempty"`
	DisplayName          string          `json:"DisplayName,omitempty"`
	LocationEmailAddress string          `json:"LocationEmailAddress,omitempty"`
}

type PhysicalAddress struct {
	Street          string `json:"Street,omitempty"`
	City            string `json:"City,omitempty"`
	State           string `json:"State,omitempty"`
	CountryOrRegion string `json:"CountryOrRegion,omitempty"`
	PostalCode      string `json:"PostalCode,omitempty"`
}

type GeoCoordinates struct {
	Altitude         float64 `json:"Altitude,omitempty"`
	Latitude         float64 `json:"Latitude,omitempty"`
	Longitude        float64 `json:"Longitude,omitempty"`
	Accuracy         float64 `json:"Accuracy,omitempty"`
	AltitudeAccuracy float64 `json:"AltitudeAccuracy,omitempty"`
}

// The type of location: Default, ConferenceRoom, HomeAddress, BusinessAddress,GeoCoordinates, StreetAddress, Hotel, Restaurant, LocalBusiness, PostalAddress.
type LocationType int32

type PatternedRecurrence struct {
	Pattern            RecurrencePattern `json:"Pattern,omitempty"`
	RecurrenceTimeZone string            `json:"RecurrenceTimeZone,omitempty"`
	Range              RecurrenceRange   `json:"Range,omitempty"`
}

type RecurrencePattern struct {
	Type           RecurrencePatternType `json:"Type,omitempty"`
	Interval       int32                 `json:"Interval,omitempty"`
	DayOfMonth     int32                 `json:"DayOfMonth,omitempty"`
	Month          int32                 `json:"Month,omitempty"`
	DaysOfWeek     []DayOfWeek           `json:"DaysOfWeek,omitempty"`
	FirstDayOfWeek DayOfWeek             `json:"DayOfWeek,omitempty"`
	Index          WeekIndex             `json:"Index,omitempty"`
}

// The recurrence pattern type: Daily = 0, Weekly = 1, AbsoluteMonthly = 2, RelativeMonthly = 3, AbsoluteYearly = 4, RelativeYearly = 5.
type RecurrencePatternType int32

// The day of the week: Sunday = 0, Monday = 1, Tuesday = 2, Wednesday = 3, Thursday = 4, Friday = 5, Saturday = 6.
type DayOfWeek int32

// The week index: First = 0, Second = 1, Third = 2, Fourth = 3, Last = 4.
type WeekIndex int32

type RecurrenceRange struct {
	Type                RecurrenceRangeType `json:"Type,omitempty"`
	StartDate           string              `json:"StartDate,omitempty"` //"2014-10-19T23:13:47.3959685Z" TODO
	EndDate             string              `json:"EndDate,omitempty"`   //"2014-10-19T23:13:47.3959685Z" TODO
	NumberOfOccurrences int32               `json:"NumberOfOccurrences,omitempty"`
}

// The recurrence range: EndDate = 0, NoEnd = 1, Numbered = 2.
type RecurrenceRangeType int32

type ResponseStatus struct {
	Response ResponseType `json:"Response,omitempty"`
	Time     string       `json:"Time,omitempty"`
}

// The response type: None, Organizer, TentativelyAccepted, Accepted, Declined, NotResponded.
type ResponseType string

// Indicates the level of privacy for the event: Normal, Personal, Private, Confidential.
type Sensitivity string

// The status to show: Free, Tentative, Busy, Oof, WorkingElsewhere, Unknown.
type FreeBusyStatus string

// The event type: SingleInstance, Occurrence, Exception, SeriesMaster.
type EventType string

func createResponseError(contents []byte) (err error) {
	e := new(Error)
	err = json.Unmarshal(contents, &e)
	if len(e.Code) != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %s. message: %s", e.Code, e.Message))
	}
	return nil
}
