package outlook_test

import (
	"os"
	"testing"

	"encoding/json"

	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

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
	_, err = account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	_, err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func TestAccount_GetAllCalendars(t *testing.T) {
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

func TestAccount_CreateCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	var calendarJSON = []byte(`{
  		"Name": "Travis"
	}`)

	calendarData := new(outlook.CalendarInfo)
	err := json.Unmarshal(calendarJSON, &calendarData)

	calendar, err := account.CreateCalendar(calendarData)

	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
		return
	}

	os.Setenv("OUTLOOK_CALENDAR_ID", calendar.ID)
	os.Setenv("OUTLOOK_CALENDAR_NAME", calendar.Name)
}

func TestAccount_GetCalendar(t *testing.T) {
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

	if calendarName != calendar.Name {
		t.Fail()
		t.Fatalf("something went wrong. Expected %s got %s", calendarName, calendar.Name)
		return
	}
}

func TestAccount_UpdateCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("OUTLOOK_CALENDAR_ID")
	oldCalendar, err := account.GetCalendar(calendarID)

	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}
	oldCalendar.Name = fmt.Sprintf("TravisRenamed%s", calendarID)

	calendar, err := account.UpdateCalendar(oldCalendar)

	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found: %s", err.Error())
		return
	}

	if oldCalendar.Name != calendar.Name {
		log.Errorf("expected %s got %s", oldCalendar.Name, calendar.Name)
		t.Fail()
		t.Fatalf("something went wrong. Expected %s got %s. Error: %s", oldCalendar.Name, calendar.Name, err.Error())
		return
	}

	os.Setenv("OUTLOOK_CALENDAR_ID", calendar.ID)
	os.Setenv("OUTLOOK_CALENDAR_NAME", calendar.Name)
}

func TestAccount_DeleteCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("OUTLOOK_CALENDAR_ID")
	log.Debugln(calendarID)

	account.DeleteCalendar(calendarID)
}
