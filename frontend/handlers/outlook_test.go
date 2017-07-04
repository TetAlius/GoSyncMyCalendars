package handlers_test

import (
	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	"net/http"
	"testing"
)

//TestGoogle_SignInHandlerSignInHandler test the SignInHandler method
func TestOutlook_SignInHandler(t *testing.T) {
	setup()
	//defer server.Close()
	// Set up the HTTP request
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/SignInWithOutlook", nil)

	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		t.Fatal("Failed with status", resp.Status)
	}

	t.Log(resp.Header.Get("Location"))
}

func TestOutlook_NewOutlookHandler(t *testing.T) {
	configNameErr := "configgggg.json"
	_, err := handlers.NewOutlookHandler(&configNameErr)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}

	configNameErr = "./badConfig.json"
	_, err = handlers.NewOutlookHandler(&configNameErr)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}
	configName := "../../config.json"
	_, err = handlers.NewOutlookHandler(&configName)
	if err != nil {
		t.Fatalf("Error should be nil but it is: %s", err.Error())
	}

	notThisConfig := "../../google.json"
	_, err = handlers.NewOutlookHandler(&notThisConfig)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}
}
