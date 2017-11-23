package outlook_test

import (
	"os"

	"testing"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
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
	account := new(outlook.Account)
	err := account.Refresh()
	logger.Debugln(err)
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
	logger.Debugln(err)

	if err == nil {
		t.Fail()
		t.Fatal("something went wrong. Expected an error found nil")
		return
	}

}

func setup() (account *outlook.Account) {
	account = &outlook.Account{
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
