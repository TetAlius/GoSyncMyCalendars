package api_test

import (
	"testing"

	"encoding/json"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

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
