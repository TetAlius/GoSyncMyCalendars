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
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	event.Calendar = calendar.(*api.GoogleCalendar)

	// good call to create event
	err = event.Create(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// good call to get event
	ev, err := calendar.GetEvent(account, event.ID)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	ev.(*api.GoogleEvent).Subject = "Update"

	// good call to update event
	err = ev.Update(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// good call to delete event
	err = event.Delete(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

}
