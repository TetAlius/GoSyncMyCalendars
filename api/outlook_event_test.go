package api_test

import (
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

func TestOutlookEventCalendar_EventLifeCycle(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var event api.OutlookEvent
	event.Subject = "Discuss the OutlookCalendar REST API"

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
