package outlook_test

import (
	"os"
	"testing"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestOutlookAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	err := account.GetPrimaryCalendar()
	if err != nil {
		log.Infoln(err.Error())
		t.Fail()
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
	}

}

func TestAccount_GetAllCalendars(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	err := account.GetAllCalendars()
	if err != nil {
		log.Infoln(err.Error())
		t.Fail()
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	err = account.GetAllCalendars()
	if err == nil {
		t.Fail()
	}
}
