package api_test

import (
	"encoding/json"
	"os"
	"testing"

	outlook "github.com/TetAlius/GoSyncMyCalendars/api/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestNewAccount(t *testing.T) {
	setupApiRoot()
	// Bad Info inside Json
	b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	_, err := outlook.NewAccount(b)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Bad formatted JSON
	b = []byte(`{"id_token":"ASD.ASD.ASD","Food":"Pickle"`)
	_, err = outlook.NewAccount(b)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Correct information given
	account := setup()
	b, err = json.Marshal(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}
	_, err = outlook.NewAccount(b)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}
}

func TestOutlookAccount_Refresh(t *testing.T) {
	setupApiRoot()
	//Empty initialized info account
	account := new(outlook.OutlookAccount)
	err := account.Refresh()
	log.Debugln(err)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Good info account
	account = setup()
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
	account := setup()
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
	account := setup()
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

	os.Setenv("OUTLOOK_CALENDAR_ID", calendar.(*outlook.OutlookCalendar).ID)
	os.Setenv("OUTLOOK_CALENDAR_NAME", calendar.(*outlook.OutlookCalendar).Name)

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
	account := setup()
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

	if calendarName != calendar.(*outlook.OutlookCalendar).Name {
		t.Fail()
		t.Fatalf("something went wrong. Expected %s got %s", calendarName, calendar.(*outlook.OutlookCalendar).Name)
		return
	}

}

func setup() (account *outlook.OutlookAccount) {
	account = &outlook.OutlookAccount{
		TokenType:         os.Getenv("OUTLOOK_TOKEN_TYPE"),
		ExpiresIn:         3600,
		AccessToken:       os.Getenv("OUTLOOK_ACCESS_TOKEN"),
		RefreshToken:      os.Getenv("OUTLOOK_REFRESH_TOKEN"),
		TokenID:           os.Getenv("OUTLOOK_TOKEN_ID"),
		AnchorMailbox:     os.Getenv("OUTLOOK_ANCHOR_MAILBOX"),
		PreferredUsername: false,
	}
	return
}

func setupApiRoot() {
	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}
