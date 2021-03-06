package api_test

import (
	"testing"

	"fmt"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

func TestOutlookCalendar_CalendarLifeCycle(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendar api.OutlookCalendar
	calendar.Name = fmt.Sprintf("Travis%d", time.Now().UnixNano())
	calendar.SetAccount(account)
	var calendarWrong api.OutlookCalendar
	calendarWrong.SetAccount(account)

	// wrong call to create calendar
	err := calendarWrong.Create()
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

	// good call to create calendar
	err = calendar.Create()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	// wrong call to get calendar
	_, err = account.GetCalendar("")
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

	//	good call to get calendar
	_, err = account.GetCalendar(calendar.GetID())
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	//	wrong call to update calendar
	calendarWrong.Name = fmt.Sprintf("TravisRenamed%s", calendarWrong.ID)
	err = calendarWrong.Update()
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

	//	good call to update calendar
	calendar.Name = fmt.Sprintf("TravisRenamed%s", calendar.GetID())
	err = calendar.Update()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	// wrong call to delete calendar
	err = calendarWrong.Delete()
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}
	//	good call to delete calendar
	err = calendar.Delete()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

}

func TestOutlookCalendar_GetAllEvents(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendarWrong api.OutlookCalendar
	calendarWrong.SetAccount(account)

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	// wrong call to get all events from a calendar
	_, err = calendarWrong.GetAllEvents()
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

	// good call to get all events from a calendar
	_, err = calendar.GetAllEvents()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

}

func TestOutlookCalendar_GetEvent(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()
	var allEvents []api.EventManager

	calendars, err := account.GetAllCalendars()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
	}

	for _, calendar := range calendars {
		events, err := calendar.GetAllEvents()
		if err != nil {
			t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s",
				calendar.GetName(), err.Error())
			break
		}
		allEvents = append(allEvents, events...)

	}

	event := allEvents[0].(*api.OutlookEvent)
	calendar := event.GetCalendar()

	_, err = calendar.GetEvent("")
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

	//good call to get event
	_, err = calendar.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("something went wrong with calendar: %s. Expected nil found error: %s", calendar.GetName(), err.Error())
	}

}
