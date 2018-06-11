package api_test

import (
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api"

	"encoding/json"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/convert"
)

func TestGoogleTime_JSON(t *testing.T) {
	var event api.GoogleEvent
	start := new(api.GoogleTime)
	end := new(api.GoogleTime)
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
	event = *new(api.GoogleEvent)
	err = json.Unmarshal(contents, &event)
	if err != nil {
		t.Fatalf("error unmarshalling json to event: %s", err.Error())
	}
	if event.Start.IsAllDay || event.End.IsAllDay {
		t.Fatalf("all day true when should be false")
	}
	if event.Start.DateTime.UTC().Format(time.RFC3339) != start.DateTime.UTC().Format(time.RFC3339) {
		t.Fatalf("start times does not match: %s vs json %s", start.DateTime.UTC().Format(time.RFC3339), event.Start.DateTime.UTC().Format(time.RFC3339))
	}
	if event.End.DateTime.UTC().Format(time.RFC3339) != end.DateTime.UTC().Format(time.RFC3339) {
		t.Fatalf("end times does not match: %s vs json %s", end.DateTime.UTC().Format(time.RFC3339), event.End.DateTime.UTC().Format(time.RFC3339))
	}

	event = *new(api.GoogleEvent)
	start = new(api.GoogleTime)
	end = new(api.GoogleTime)
	now = time.Now()
	more = now.Add(time.Hour * 2)
	start.Date = now
	start.IsAllDay = true
	end.Date = more
	end.IsAllDay = true
	event.Start = start
	event.End = end
	contents, err = json.Marshal(event)
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}
	event = *new(api.GoogleEvent)
	err = json.Unmarshal(contents, &event)
	if err != nil {
		t.Fatalf("error unmarshalling json to event: %s", err.Error())
	}
	if !event.Start.IsAllDay || !event.End.IsAllDay {
		t.Fatalf("all day false when should be true")
	}
	if event.Start.Date.UTC().Format("2006-01-02") != start.Date.UTC().Format("2006-01-02") {
		t.Fatalf("start times does not match: %s vs json %s", start.Date.UTC().Format("2006-01-02"), event.Start.Date.UTC().Format("2006-01-02"))
	}
	if event.End.Date.UTC().Format("2006-01-02") != end.Date.UTC().Format("2006-01-02") {
		t.Fatalf("end times does not match: %s vs json %s", end.Date.UTC().Format("2006-01-02"), event.End.Date.UTC().Format("2006-01-02"))
	}

}

func TestGoogleEventCalendar_EventLifeCycle(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var event api.GoogleEvent
	event.Subject = "Discuss the GoogleCalendar REST API"
	start := new(api.GoogleTime)
	end := new(api.GoogleTime)
	now := time.Now()
	more := now.Add(time.Hour * 2)
	start.DateTime = now
	end.DateTime = more
	event.Start = start
	event.End = end
	event.Description = "This is a description"
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
	ev.(*api.GoogleEvent).Subject = "Update"

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

	outlook := new(api.OutlookEvent)
	err = convert.Convert(ev, outlook)
	if err != nil {
		t.Fatalf("error converting from google to outlook: %s", err.Error())
	}
	google := new(api.GoogleEvent)
	err = convert.Convert(outlook, google)
	if err != nil {
		t.Fatalf("error converting from outlook to google: %s", err.Error())
	}
}
