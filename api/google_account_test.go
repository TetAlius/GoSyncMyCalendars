package api_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestNewGoogleAccount(t *testing.T) {
	setupApiRoot()
	// Bad Info inside Json
	b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	_, err := api.NewGoogleAccount(b)
	if err == nil {
		t.Fatal("something went wrong. Expected an error found nil")
	}

	//Bad formatted JSON
	b = []byte(`{"id_token":"ASD.ASD.ASD","Food":"Pickle"`)
	_, err = api.NewGoogleAccount(b)
	if err == nil {
		t.Fatal("something went wrong. Expected an error found nil")
	}

	//Correct information given
	_, account := setup()
	b, err = json.Marshal(account)
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}
	_, err = api.NewGoogleAccount(b)
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

}

func TestGoogleAccount_Refresh(t *testing.T) {
	setupApiRoot()
	//Empty initialized info account
	account := new(api.GoogleAccount)
	err := account.Refresh()
	if err == nil {
		t.Fatal("something went wrong. Expected an error found nil")
	}

	//Good info account
	_, account = setup()
	err = account.Refresh()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	os.Setenv("API_ROOT", "")
	err = account.Refresh()

	if err == nil {
		t.Fatal("something went wrong. Expected an error found nil")
	}

}

func TestGoogleAccount_GetAllCalendars(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	_, err := account.GetAllCalendars()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetAllCalendars()
	if err == nil {
		t.Fatal("something went wrong. Expected error found nil")
	}

}

func TestGoogleAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	err := account.Refresh()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
	}

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
	}
	os.Setenv("GOOGLE_CALENDAR_ID", calendar.GetID())
	os.Setenv("GOOGLE_CALENDAR_NAME", calendar.GetName())

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fatal("something went wrong. Expected an error found nil")
	}

}

func TestGoogleAccount_GetCalendar(t *testing.T) {
	setupApiRoot()
	_, account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("GOOGLE_CALENDAR_ID")
	calendarName := os.Getenv("GOOGLE_CALENDAR_NAME")

	logger.Debugln(calendarID)
	logger.Debugln(calendarName)

	calendar, err := account.GetCalendar(calendarID)

	if err != nil {
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
	}

	if calendarName != calendar.GetName() {
		t.Fatalf("something went wrong. Expected %s got %s", calendarName, calendar.GetName())

	}

}
