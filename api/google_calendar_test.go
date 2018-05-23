package api_test

import (
	"encoding/json"
	"testing"

	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

func TestGoogleCalendar_CalendarLifeCycle(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong refreshing. Expected nil found error: %s", err.Error())
		return
	}

	var calendar api.GoogleCalendar
	calendar.SetAccount(account)
	var calendarJSON = []byte(`{
  		"summary": "Travis"
	}`)
	err = json.Unmarshal(calendarJSON, &calendar)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// good call to create calendar
	err = calendar.Create()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	//	good call to get calendar
	_, err = account.GetCalendar(calendar.GetID())
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	//	good call to update calendar
	calendar.Name = fmt.Sprintf("TravisRenamed%s", calendar.GetID())
	err = calendar.Update()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	//	good call to delete calendar
	err = calendar.Delete()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	var calendarWrong api.GoogleCalendar
	calendarWrong.SetAccount(account)
	// wrong call to create calendar
	err = calendarWrong.Create()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}
	// wrong call to get calendar
	_, err = account.GetCalendar("asdasd")
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}
	//	wrong call to update calendar
	calendarWrong.Name = fmt.Sprintf("TravisRenamed%s", calendarWrong.ID)
	err = calendarWrong.Update()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}
	// wrong call to delete calendar
	err = calendarWrong.Delete()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}
}

func TestGoogleCalendar_GetAllEvents(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendarWrong api.GoogleCalendar
	calendarWrong.SetAccount(account)

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// wrong call to get all events from a calendar
	_, err = calendarWrong.GetAllEvents()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	// good call to get all events from a calendar
	_, err = calendar.GetAllEvents()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

}

func TestGoogleCalendar_GetEvent(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()
	var allEvents []api.EventManager

	calendars, err := account.GetAllCalendars()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	for _, calendar := range calendars {
		events, err := calendar.GetAllEvents()
		if err != nil {
			t.Fail()
			t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s",
				calendar.GetName(), err.Error())
			return
		}
		allEvents = append(allEvents, events...)

	}

	event := allEvents[0].(*api.GoogleEvent)
	calendar := event.GetCalendar()

	_, err = calendar.GetEvent("asdasd")
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	//good call to get event
	_, err = calendar.GetEvent(event.ID)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s", calendar.GetName(), err.Error())
		return
	}

}
