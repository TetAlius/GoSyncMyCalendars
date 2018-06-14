package api_test

import (
	"testing"

	"encoding/json"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

var contOutlook = []byte(` {
  "CreatedDateTime": "2018-06-11T18:05:01.2183293Z",
  "LastModifiedDateTime": "2018-06-11T18:05:01.2495814Z",
  "ChangeKey": "7Tv+UgXJaEmg/mpJO213QQACVRFjeg==",
  "Categories": [],
  "ReminderMinutesBeforeStart": 15,
  "IsReminderOn": true,
  "HasAttachments": false,
  "Subject": "RECURRENCE OUTLOOK",
  "BodyPreview": "",
  "Importance": "Normal",
  "Sensitivity": "Normal",
  "IsAllDay": false,
  "IsCancelled": false,
  "IsOrganizer": true,
  "ResponseRequested": true,
  "ResponseStatus": {
    "Response": "Organizer",
    "Time": "0001-01-01T00:00:00Z"
  },
  "Body": {
    "ContentType": "Text",
    "Content": "\r\n"
  },
  "Start": {
    "DateTime": "2018-06-14T15:00:00.0000000",
    "TimeZone": "UTC"
  },
  "End": {
    "DateTime": "2018-06-14T15:30:00.0000000",
    "TimeZone": "UTC"
  },
  "Recurrence": {
    "Pattern": {
      "Type": "Weekly",
      "Interval": 1,
      "Month": 0,
      "DayOfMonth": 0,
      "DaysOfWeek": [
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday"
      ],
      "FirstDayOfWeek": "Monday",
      "Index": "First"
    },
    "Range": {
      "Type": "EndDate",
      "StartDate": "2018-06-15",
      "EndDate": "2018-11-30",
      "RecurrenceTimeZone": "Romance Standard Time",
      "NumberOfOccurrences": 0
    }
  }
}`)

func TestOutlookTime_JSON(t *testing.T) {
	var event api.OutlookEvent
	start := new(api.OutlookDateTimeTimeZone)
	end := new(api.OutlookDateTimeTimeZone)
	now := time.Now()
	more := now.Add(time.Hour * 2)
	start.DateTime = now
	end.DateTime = more
	event.Start = start
	event.End = end
	contents, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	event = *new(api.OutlookEvent)
	err = json.Unmarshal(contents, &event)
	if err != nil {
		t.Fatalf("error unmarshalling json to event: %s", err.Error())
	}
	//TODO:
	//if event.Start.IsAllDay || event.End.IsAllDay {
	//	t.Fatalf("all day true when should be false")
	//}
	if event.Start.DateTime.UTC().Format(time.RFC3339Nano) != start.DateTime.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("start times does not match: %s vs json %s", start.DateTime.UTC().Format(time.RFC3339Nano), event.Start.DateTime.UTC().Format(time.RFC3339Nano))
	}
	if event.End.DateTime.UTC().Format(time.RFC3339Nano) != end.DateTime.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("end times does not match: %s vs json %s", end.DateTime.UTC().Format(time.RFC3339Nano), event.End.DateTime.UTC().Format(time.RFC3339Nano))
	}
	event = *new(api.OutlookEvent)
	event.Start = new(api.OutlookDateTimeTimeZone)
	event.End = new(api.OutlookDateTimeTimeZone)
	contents, err = json.Marshal(event)
	if err != nil {
		t.Fatalf("error marshaling empty dates: %s", err.Error())
	}

	eventUn := new(api.OutlookEvent)
	err = json.Unmarshal(contOutlook, &eventUn)
	if err != nil {
		t.Fatalf("something happened: %s", err.Error())
	}

	contents, err = json.Marshal(eventUn)
	if err != nil {
		t.Fatalf("something happened: %s", err.Error())
	}
	logger.Debugf("%s", contents)
}

func TestOutlookEventCalendar_EventLifeCycle(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var event api.OutlookEvent
	event.Subject = "Discuss the OutlookCalendar REST API"
	start := new(api.OutlookDateTimeTimeZone)
	end := new(api.OutlookDateTimeTimeZone)
	now := time.Now()
	more := now.Add(time.Hour * 2)
	start.DateTime = now
	end.DateTime = more
	event.Start = start
	event.End = end

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}
	event.SetCalendar(calendar)

	// good call to create event
	err = event.Create()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	// good call to get event
	ev, err := calendar.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	ev.(*api.OutlookEvent).Subject = "Update"

	// good call to update event
	err = ev.Update()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())

	}

	// good call to delete event
	err = event.Delete()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

}
