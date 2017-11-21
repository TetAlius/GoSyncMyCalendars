package outlook

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

type Account struct {
	TokenType         string `json:"token_type"`
	ExpiresIn         int    `json:"expires_in"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	TokenID           string `json:"id_token"`
	AnchorMailbox     string
	PreferredUsername bool
}

// The type of location: Default, ConferenceRoom, HomeAddress, BusinessAddress,GeoCoordinates, StreetAddress, Hotel, Restaurant, LocalBusiness, PostalAddress.
type LocationType int32

// The recurrence pattern type: Daily = 0, Weekly = 1, AbsoluteMonthly = 2, RelativeMonthly = 3, AbsoluteYearly = 4, RelativeYearly = 5.
type RecurrencePatternType int32

// The day of the week: Sunday = 0, Monday = 1, Tuesday = 2, Wednesday = 3, Thursday = 4, Friday = 5, Saturday = 6.
type DayOfWeek int32

// The week index: First = 0, Second = 1, Third = 2, Fourth = 3, Last = 4.
type WeekIndex int32

// The recurrence range: EndDate = 0, NoEnd = 1, Numbered = 2.
type RecurrenceRangeType int32

// Indicates the level of privacy for the event: Normal, Personal, Private, Confidential.
type Sensitivity string

// Specifies the color theme to distinguish the calendar from other calendars in a UI.
// The property values are: LightBlue=0, LightGreen=1, LightOrange=2, LightGray=3, LightYellow=4, LightTeal=5,
// LightPink=6, LightBrown=7, LightRed=8, MaxColor=9, Auto=-1
type CalendarColor string

// The importance of the event: Low, Normal, High.
type Importance string

// The response type: None, Organizer, TentativelyAccepted, Accepted, Declined, NotResponded.
type ResponseType string

// The status to show: Free, Tentative, Busy, Oof, WorkingElsewhere, Unknown.
type FreeBusyStatus string

// The event type: SingleInstance, Occurrence, Exception, SeriesMaster.
type EventType string

type CalendarResponse struct {
	OdataContext string `json:"@odata.context"`
	CalendarInfo
}

type CalendarListResponse struct {
	OdataContext string         `json:"@odata.context"`
	Calendars    []CalendarInfo `json:"value"`
}

// CalendarInfo TODO
type CalendarInfo struct {
	OdataID string `json:"@odata.id,omitempty"`

	CalendarView        []EventInfo   `json:"CalendarView,omitempty"`
	CanEdit             bool          `json:"CanEdit,omitempty"`
	CanShare            bool          `json:"CanShare,omitempty"`
	CanViewPrivateItems bool          `json:"CanViewPrivateItems,omitempty"`
	ChangeKey           string        `json:"ChangeKey,omitempty"`
	Color               CalendarColor `json:"Color,omitempty"`
	Events              []EventInfo   `json:"Events,omitempty"`
	ID                  string        `json:"Id"`
	IsDefaultCalendar   bool          `json:"IsDefaultCalendar,omitempty"`
	IsShared            bool          `json:"IsShared,omitempty"`
	IsSharedWithMe      bool          `json:"IsSharedWithMe,omitempty"`
	Name                string        `json:"Name,omitempty"`
	Owner               Recipient     `json:"Owner,omitempty"`
	//MultiValueExtendedProperties  []Properties  `json:"MultiValueExtendedProperties"`
	//SingleValueExtendedProperties []Properties  `json:"SingleValueExtendedProperties"`
}

type EventInfo struct {
	OdataID   string `json:"@odata.id"`
	OdataEtag string `json:"@odata.etag"`

	Attachments                []Attachment        `json:"Attachments"`
	Attendees                  []Attendee          `json:"Attendees"`
	Body                       ItemBody            `json:"Body"`
	BodyPreview                string              `json:"BodyPreview"`
	Calendar                   CalendarInfo        `json:"Calendar"`
	Categories                 []string            `json:"Categories"`
	ChangeKey                  string              `json:"ChangeKey"`
	CreatedDateTime            string              `json:"CreatedDateTime"` //"2014-10-19T23:13:47.3959685Z"
	End                        DateTimeTimeZone    `json:"End"`
	HasAttachments             bool                `json:"HasAttachments"`
	ICalUID                    string              `json:"iCalUID"`
	ID                         string              `json:"Id"`
	Importance                 Importance          `json:"Importance"`
	Instances                  []EventInfo         `json:"Instances"`
	IsAllDay                   bool                `json:"IsAllday"`
	IsCancelled                bool                `json:"IsCancelled"`
	IsOrganizer                bool                `json:"IsOrganizer"`
	IsReminderOn               bool                `json:"IsReminderOn"`
	LastModifiedDateTime       string              `json:"LastModifiedDateTime"` //"2014-10-19T23:13:47.6772234Z"
	Location                   Location            `json:"Location"`
	Locations                  []Location          `json:"Location"`
	OnlineMeetingUrl           string              `json:"OnlineMeetingUrl"`
	Organizer                  Recipient           `json:"Organizer"`
	OriginalStartTimeZone      string              `json:"OriginalStartTimeZone"`
	OriginalEndTimeZone        string              `json:"OriginalEndTimeZone"`
	Recurrence                 PatternedRecurrence `json:"Recurrence"`
	ReminderMinutesBeforeStart int32               `json:"ReminderMinutesBeforeStart"`
	ResponseRequested          bool                `json:"ResponseRequested"`
	ResponseStatus             ResponseStatus      `json:"ResponseStatus"`
	Sensitivity                Sensitivity         `json:"Sensitivity"`
	SeriesMasterID             string              `json:"SeriesMasterId"`
	ShowAs                     FreeBusyStatus      `json:"ShowAs"`
	Start                      DateTimeTimeZone    `json:"Start"`
	Subject                    string              `json:"Subject"`
	Type                       EventType           `json:"Type"`
	WebLink                    string              `json:"WebLink"`

	//Extensions                 []Extension         `json:"Extensions"`
}

type ItemBody struct {
	ContentType string `json:"ContentType"`
	Content     string `json:"Content"`
}
type DateTimeTimeZone struct {
	DateTime string `json:"DateTime"`
	TimeZone string `json:"TimeZone"`
}
type Location struct {
	Address              PhysicalAddress `json:"Address"`
	Coordinates          GeoCoordinates  `json:"Coordinates"`
	DisplayName          string          `json:"DisplayName"`
	LocationEmailAddress string          `json:"LocationEmailAddress"`
	LocationUri          string          `json:"LocationUri"`
	// TODO
	LocationType LocationType `json:"LocationType"`
}

type PhysicalAddress struct {
	Street          string `json:"Street"`
	City            string `json:"City"`
	State           string `json:"State"`
	CountryOrRegion string `json:"CountryOrRegion"`
	PostalCode      string `json:"PostalCode"`
}
type GeoCoordinates struct {
	Altitude         float64 `json:"Altitude"`
	Latitude         float64 `json:"Latitude"`
	Longitude        float64 `json:"Longitude"`
	Accuracy         float64 `json:"Accuracy"`
	AltitudeAccuracy float64 `json:"AltitudeAccuracy"`
}

type PatternedRecurrence struct {
	Pattern            RecurrencePattern `json:"Pattern"`
	RecurrenceTimeZone string            `json:"RecurrenceTimeZone"`
	Range              RecurrenceRange   `json:"Range"`
}

type RecurrencePattern struct {
	Type           RecurrencePatternType `json:"Type"`
	Interval       int32                 `json:"Interval"`
	DayOfMonth     int32                 `json:"DayOfMonth"`
	Month          int32                 `json:"Month"`
	DaysOfWeek     []DayOfWeek           `json:"DaysOfWeek"`
	FirstDayOfWeek DayOfWeek             `json:"DayOfWeek"`
	Index          WeekIndex             `json:"Index"`
}

type RecurrenceRange struct {
	Type                RecurrenceRangeType `json:"Type"`
	StartDate           string              `json:"StartDate"` //"2014-10-19T23:13:47.3959685Z" TODO
	EndDate             string              `json:"EndDate"`   //"2014-10-19T23:13:47.3959685Z" TODO
	NumberOfOccurrences int32               `json:"NumberOfOccurrences"`
}

type Attendee struct {
	Recipient Recipient `json:"EmailAddress"`
	Status    Status    `json:"Status"`
	Type      string    `json:"Type"`
}

type Status struct {
	Response string `json:"Response"`
	Time     string `json:"Time"`
}

type Recipient struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
}

type ResponseStatus struct {
	Response ResponseType `json:"Response"`
	Time     string       `json:"Time"`
}

type Attachment struct {
	ContentType          string `json:"ContentType"`
	IsInline             bool   `json:"IsInline"`
	LastModifiedDateTime string `json:"LastModifiedDateTime"`
	Name                 string `json:"Name"`
	Size                 int32  `json:"Size"`
}

func NewAccount(contents []byte) (r *Account, err error) {
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshaling outlook response: %s", err.Error()))
	}

	email, preferred, err := util.MailFromToken(strings.Split(r.TokenID, "."), "=")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error retrieving outlook mail: %s", err.Error()))
	}
	r.AnchorMailbox = email
	r.PreferredUsername = preferred
	return
}

type RefreshError struct {
	Code    string `json:"error,omitempty"`
	Message string `json:"error_description,omitempty"`
}

// TokenRefresh TODO
func (o *Account) Refresh() (err error) {
	client := http.Client{}
	//check if token is DEAD!!!

	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	params, err := util.CallAPIRoot("outlook/token/refresh-params")
	log.Debugf("Params: %s", fmt.Sprintf(params, o.RefreshToken))
	if err != nil {
		log.Errorf("Error generating params: %s", err.Error())
		return
	}

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(fmt.Sprintf(params, o.RefreshToken)))

	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing outlook request: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from outlook request: %s", err.Error())
		return
	}

	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		e := new(RefreshError)
		err = json.Unmarshal(contents, &e)
		log.Errorln(e.Code)
		log.Errorln(e.Message)
		if len(e.Code) != 0 && len(e.Message) != 0 {
			return errors.New(fmt.Sprintf("code: %s. message: %s", e.Code, e.Message))
		}
	}

	log.Debugf("\nTokenType: %s\nExpiresIn: %d\nAccessToken: %s\nRefreshToken: %s\nTokenID: %s\nAnchorMailbox: %s\nPreferredUsername: %t",
		o.TokenType, o.ExpiresIn, o.AccessToken, o.RefreshToken, o.TokenID, o.AnchorMailbox, o.PreferredUsername)

	log.Debugf("%s\n", contents)
	err = json.Unmarshal(contents, &o)
	if err != nil {
		log.Errorf("There was an error with the outlook request: %s", err.Error())
		return
	}

	log.Debugf("\nTokenType: %s\nExpiresIn: %d\nAccessToken: %s\nRefreshToken: %s\nTokenID: %s\nAnchorMailbox: %s\nPreferredUsername: %t",
		o.TokenType, o.ExpiresIn, o.AccessToken, o.RefreshToken, o.TokenID, o.AnchorMailbox, o.PreferredUsername)
	return
}

func (o *Account) authorizationRequest() (auth string) {
	return o.TokenType + " " + o.AccessToken
}
