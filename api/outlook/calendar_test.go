package outlook_test

import (
	"encoding/json"
	"testing"

	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/api/outlook"
)

func TestOutlookCalendar_CalendarLifeCycle(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendar outlook.Calendar
	var calendarWrong outlook.Calendar
	var calendarJSON = []byte(`{
  		"Name": "Travis"
	}`)

	// wrong call to create calendar
	err := calendarWrong.Create(account)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	t.Fail()
	t.Fatal()
	return

	err = json.Unmarshal(calendarJSON, &calendar)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// good call to create calendar
	err = calendar.Create(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// wrong call to get calendar
	_, err = account.GetCalendar("")
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	//	good call to get calendar
	_, err = account.GetCalendar(calendar.ID)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	//	wrong call to update calendar
	calendarWrong.Name = fmt.Sprintf("TravisRenamed%s", calendarWrong.ID)
	err = calendarWrong.Update(account)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	//	good call to update calendar
	calendar.Name = fmt.Sprintf("TravisRenamed%s", calendar.ID)
	err = calendar.Update(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// wrong call to delete calendar
	err = calendarWrong.Delete(account)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}
	//	good call to delete calendar
	err = calendar.Delete(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

}

func TestOutlookCalendar_GetAllEvents(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendarWrong outlook.Calendar

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// wrong call to get all events from a calendar
	_, err = calendarWrong.GetAllEvents(account)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	// good call to get all events from a calendar
	_, err = calendar.GetAllEvents(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

}

func TestOutlookCalendar_GetEvent(t *testing.T) {
	setupApiRoot()
	account := setup()
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
		events, err := calendar.GetAllEvents(account)
		if err != nil {
			t.Fail()
			t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s",
				calendar.(*outlook.Calendar).Name, err.Error())
			return
		}
		allEvents = append(allEvents, events...)

	}

	event := allEvents[0].(*outlook.Event)
	calendar := event.Calendar

	_, err = calendar.GetEvent(account, "")
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

	//good call to get event
	_, err = calendar.GetEvent(account, event.ID)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s", calendar.Name, err.Error())
		return
	}

}
