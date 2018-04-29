package api_test

import (
	"testing"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api/google"
)

func TestEventCalendar_EventLifeCycle(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var event google.GoogleEvent
	event.Summary = "Discuss the GoogleCalendar REST API"
	event.Start = new(google.GoogleTime)
	event.Start.Date = time.Now().Format("2006-01-02")
	event.End = event.Start

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	event.Calendar = calendar.(*google.GoogleCalendar)

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

	ev.(*google.GoogleEvent).Summary = "Update"

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
