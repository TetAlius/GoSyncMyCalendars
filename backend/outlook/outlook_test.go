package outlook_test

import (
	"os"

	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/backend/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/logger"
)

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
	logger.Debugln(err)

	if err != nil {
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
