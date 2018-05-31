package api_test

import (
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api"

	"time"
)

func TestGoogleEventCalendar_EventLifeCycle(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var event api.GoogleEvent
	event.Subject = "Discuss the GoogleCalendar REST API"
	event.Start = new(api.GoogleTime)
	event.Start.DateTime = time.Now().Format(time.RFC3339)
	event.End = new(api.GoogleTime)
	event.End.DateTime = time.Now().Add(time.Hour * time.Duration(2)).Format(time.RFC3339)

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

}
