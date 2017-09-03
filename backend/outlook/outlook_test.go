package outlook_test

import (
	"os"

	"testing"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

func TestNewAccount(t *testing.T) {
	// Bad Info inside Json
	b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
	_, err := outlook.NewAccount(b)
	if err == nil {
		t.Fail()
	}

	//Bad formatted JSON
	b = []byte(`{"id_token":"ASD.ASD.ASD","Food":"Pickle"`)
	_, err = outlook.NewAccount(b)
	if err == nil {
		t.Fail()
	}

	//Correct information given
	account := setup()
	b, err = json.Marshal(account)
	if err != nil {
		t.Fail()
	}
	_, err = outlook.NewAccount(b)
	if err != nil {
		t.Fail()
	}
}

func TestOutlookAccount_Refresh(t *testing.T) {
	//Empty initialized info account
	account := new(outlook.OutlookAccount)
	err := account.Refresh()
	logger.Debugln(err)
	if err == nil {
		t.Fail()
	}

	//Good info account
	account = setup()
	err = account.Refresh()
	if err != nil {
		t.Fail()
	}

	os.Setenv("API_ROOT", "")

	err = account.Refresh()
	logger.Debugln(err)

	if err == nil {
		t.Fail()
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
