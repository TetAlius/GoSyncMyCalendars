package outlook_test

import (
	"os"
	"testing"
)

func TestOutlookAccount_GetPrimaryCalendar(t *testing.T) {
	setupApiRoot()
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	err := account.GetPrimaryCalendar()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	os.Setenv("API_ROOT", "")
	// Bad calling to GetPrimaryCalendar
	err = account.GetPrimaryCalendar()
	if err == nil {
		t.Fail()
	}

}
