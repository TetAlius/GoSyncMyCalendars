package google_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/api/google"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestNewAccount(t *testing.T) {
	setupApiRoot()
	// Bad Info inside Json
	b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	_, err := google.NewAccount(b)
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

	//Bad formatted JSON
	b = []byte(`{"id_token":"ASD.ASD.ASD","Food":"Pickle"`)
	_, err = google.NewAccount(b)
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
	acc, err := google.NewAccount(b)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
		return
	}

	acc.GetAllCalendars()

}

func TestGoogleAccount_Refresh(t *testing.T) {
	setupApiRoot()
	//Empty initialized info account
	account := new(google.GoogleAccount)
	err := account.Refresh()
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

	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func TestGoogleAccount_GetAllCalendars(t *testing.T) {
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

func TestGoogleAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	logger.Debugln(calendar.(*google.GoogleCalendar).ID)
	logger.Debugln(calendar.(*google.GoogleCalendar).Name)
	os.Setenv("GOOGLE_CALENDAR_ID", calendar.(*google.GoogleCalendar).ID)
	os.Setenv("GOOGLE_CALENDAR_NAME", calendar.(*google.GoogleCalendar).Name)

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func TestGoogleAccount_GetCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("GOOGLE_CALENDAR_ID")
	calendarName := os.Getenv("GOOGLE_CALENDAR_NAME")

	logger.Debugln(calendarID)
	logger.Debugln(calendarName)

	calendar, err := account.GetCalendar(calendarID)

	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	if calendarName != calendar.(*google.GoogleCalendar).Name {
		t.Fail()
		t.Fatalf("something went wrong. Expected %s got %s", calendarName, calendar.(*google.GoogleCalendar).Name)
		return
	}

}

func setup() (account *google.GoogleAccount) {
	account = &google.GoogleAccount{
		TokenType:    os.Getenv("GOOGLE_TOKEN_TYPE"),
		ExpiresIn:    3600,
		AccessToken:  os.Getenv("GOOGLE_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
		TokenID:      os.Getenv("GOOGLE_TOKEN_ID"),
		Email:        os.Getenv("GOOGLE_EMAIL"),
	}
	return
}

func setupApiRoot() {
	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}
