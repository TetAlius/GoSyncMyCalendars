package handlers_test

import (
	"net/http"
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/frontend"
)

//TestGoogle_SignInHandlerSignInHandler test the SignInHandler method
func TestOutlookSignInHandler(t *testing.T) {
	f := frontend.NewServer("127.0.0.1", 8080)
	f.Start()
	defer f.Stop()
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
