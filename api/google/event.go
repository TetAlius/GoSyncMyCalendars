package google

import (
	"bytes"
	"errors"
	"fmt"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

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

// POST https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events
func (event *Event) Create(a api.AccountManager) (err error) {
	log.Debugln("createEvent google")

	route, err := util.CallAPIRoot("google/calendars/id/events")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}

	calendar := event.Calendar
	event.Calendar = nil

	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	log.Debugln(data)
	event.Calendar = calendar

	contents, err := util.DoRequest("POST",
		fmt.Sprintf(route, event.Calendar.ID),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return errors.New(fmt.Sprintf("error creating event in g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)

	log.Debugf("Response: %s", contents)
	return
}

// PUT https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *Event) Update(a api.AccountManager) (err error) {
	log.Debugln("updateEvent google")
	//TODO: Test if ids are two given

	//Meter en los header el etag
	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		return errors.New(fmt.Sprintf("error generating URL: %s", err.Error()))
	}
	calendar := event.Calendar
	event.Calendar = nil
	data, err := json.Marshal(event)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling event data: %s", err.Error()))
	}
	event.Calendar = calendar

	contents, err := util.DoRequest("PUT",
		fmt.Sprintf(route, event.Calendar.ID, event.ID),
		bytes.NewBuffer(data),
		a.AuthorizationRequest(),
		"")

	if err != nil {
		return errors.New(fmt.Sprintf("error updating event of g calendar for email %s. %s", a.Mail(), err.Error()))
	}
	err = createResponseError(contents)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, &event)

	log.Debugf("Response: %s", contents)
	return
}

// DELETE https://www.googleapis.com/calendar/v3/calendars/{calendarID}/events/{eventID}
func (event *Event) Delete(a api.AccountManager) (err error) {
	log.Debugln("deleteEvent google")
	//TODO: Test if ids are two given

	route, err := util.CallAPIRoot("google/calendars/id/events/id")
	if err != nil {
		log.Errorf("Error generating URL: %s", err.Error())
		return
	}

	contents, err := util.DoRequest(
		"DELETE",
		fmt.Sprintf(route, event.Calendar.ID, event.ID),
		nil,
		a.AuthorizationRequest(),
		"")

	if err != nil {
		log.Errorf("Error deleting event of g calendar for email %s. %s", a.Mail(), err.Error())
	}

	log.Debugf("Contents: %s", contents)
	if len(contents) != 0 {
		return errors.New(fmt.Sprintf("error deleting google event %s: %s", event.ID, contents))
	}

	return
}
