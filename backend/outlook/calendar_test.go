package outlook_test

import (
	"os"
	"testing"

	"encoding/json"

	"fmt"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestOutlookAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	log.Debugln("Started")
	err = account.GetPrimaryCalendar()
	if err != nil {
		log.Infoln(err.Error())
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
	}

}

func TestAccount_GetAllCalendars(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	err := account.GetAllCalendars()
	if err != nil {
		t.Fail()
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	err = account.GetAllCalendars()
	if err == nil {
		t.Fail()
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

	calendar, err := account.CreateCalendar(calendarJSON)

	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
	}

	os.Setenv("OUTLOOK_CALENDAR_ID", calendar.ID)
	log.Infof("CalendarID: %s", calendar.ID)
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
		log.Errorln(err.Error())
		t.Fail()
	}

	if calendarName != calendar.Name {
		log.Errorf("expected %s got %s", calendarName, calendar.Name)
	}
	log.Debugln(calendarID)
}

func TestAccount_UpdateCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	calendarID := os.Getenv("OUTLOOK_CALENDAR_ID")

	oldCalendar, err := account.GetCalendar(calendarID)

	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
	}
	oldCalendar.Name = fmt.Sprintf("TravisRenamed%s", calendarID)

	calendarJSON, err := json.Marshal(oldCalendar)

	if err != nil {
		log.Errorln(err)
		t.Fail()
	}

	calendar, err := account.UpdateCalendar(calendarID, calendarJSON)

	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
	}

	if oldCalendar.Name != calendar.Name {
		log.Errorf("expected %s got %s", oldCalendar.Name, calendar.Name)
		t.Fail()
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
