package api_test

import (
	"os"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

func setup() (outAcc *api.OutlookAccount, gooAcc *api.GoogleAccount) {
	outAcc = &api.OutlookAccount{
		TokenType:         os.Getenv("OUTLOOK_TOKEN_TYPE"),
		ExpiresIn:         3600,
		AccessToken:       os.Getenv("OUTLOOK_ACCESS_TOKEN"),
		RefreshToken:      os.Getenv("OUTLOOK_REFRESH_TOKEN"),
		TokenID:           os.Getenv("OUTLOOK_TOKEN_ID"),
		AnchorMailbox:     os.Getenv("OUTLOOK_ANCHOR_MAILBOX"),
		PreferredUsername: false,
	}

	gooAcc = &api.GoogleAccount{
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
