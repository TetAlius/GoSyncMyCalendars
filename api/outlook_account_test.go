package api_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestNewOutlookAccount(t *testing.T) {
	setupApiRoot()
	// Bad Info inside Json
	b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	_, err := api.NewOutlookAccount(b)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Bad formatted JSON
	b = []byte(`{"id_token":"ASD.ASD.ASD","Food":"Pickle"`)
	_, err = api.NewOutlookAccount(b)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Correct information given
	account, _ := setup()
	b, err = json.Marshal(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}
	_, err = api.NewOutlookAccount(b)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}
}

func TestOutlookAccount_Refresh(t *testing.T) {
	setupApiRoot()
	//Empty initialized info account
	account := new(api.OutlookAccount)
	err := account.Refresh()
	log.Debugln(err)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Good info account
	account, _ = setup()
	err = account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}

	os.Setenv("API_ROOT", "")
	err = account.Refresh()
	log.Debugln(err)

	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func TestOutlookAccount_GetAllCalendars(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	_, err := account.GetAllCalendars()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetAllCalendars()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected error found nil")
		return
	}

}

func TestOutlookAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	log.Debugln("Started")
	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	os.Setenv("OUTLOOK_CALENDAR_ID", calendar.GetID())
	os.Setenv("OUTLOOK_CALENDAR_NAME", calendar.(*api.OutlookCalendar).Name)

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func TestOutlookAccount_GetCalendar(t *testing.T) {
	setupApiRoot()
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("OUTLOOK_CALENDAR_ID")
	calendarName := os.Getenv("OUTLOOK_CALENDAR_NAME")

	calendar, err := account.GetCalendar(calendarID)

	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	if calendarName != calendar.(*api.OutlookCalendar).Name {
		t.Fail()
		t.Fatalf("something went wrong. Expected %s got %s", calendarName, calendar.(*api.OutlookCalendar).Name)
		return
	}

}
