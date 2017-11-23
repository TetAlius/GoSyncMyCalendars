package outlook_test

import (
	"testing"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestOutlookAccount_CreateEvent(t *testing.T) {
	setupApiRoot()
	account := setup()

	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		log.Infoln(err.Error())
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	//TODO: delete this
	var createEvent = []byte(`{
  "Subject": "Discuss the Calendar REST API",
  "Body": {
    "ContentType": "HTML",
    "Content": "I think it will meet our requirements!"
  },
  "Start": {
      "DateTime": "2016-02-02T18:00:00",
      "TimeZone": "Pacific Standard Time"
  },
  "End": {
      "DateTime": "2016-02-02T19:00:00",
      "TimeZone": "Pacific Standard Time"
  },
	"ReminderMinutesBeforeStart": "30",
  "IsReminderOn": "false"
}`)
	_, err = account.CreateEvent(calendar.ID, createEvent)
	if err != nil {
		t.Fail()
		t.Fatalf("error creating new event: %s", err.Error())
	}

	t.Log("Not implemented yet")
}

func TestOutlookAccount_GetAllEventsFromCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()

	err := account.Refresh()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		log.Infoln(err.Error())
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found %s", err.Error())
	}
	_, err = account.GetAllEventsFromCalendar(calendar.ID)
	if err != nil {
		t.Fail()
		t.Fatalf("error creating new event: %s", err.Error())
	}

	t.Log("Not implemented yet")
}

func TestOutlookAccount_GetEvent(t *testing.T) {
	t.Log("Not Implemented yet")
}

func TestOutlookAccount_UpdateEvent(t *testing.T) {
	t.Log("Not Implemented yet")
}

func TestOutlookAccount_DeleteEvent(t *testing.T) {
	t.Log("Not Implemented yet")
}
