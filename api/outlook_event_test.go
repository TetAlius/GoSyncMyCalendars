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
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	event.Calendar = calendar.(*api.OutlookCalendar)

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

	ev.(*api.OutlookEvent).Subject = "Update"

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