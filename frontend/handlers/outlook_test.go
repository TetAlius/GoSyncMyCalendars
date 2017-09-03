package handlers_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/TetAlius/GoSyncMyCalendars/frontend"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func init() {
	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}

//TestGoogle_SignInHandlerSignInHandler test the SignInHandler method
func TestOutlookSignInHandler(t *testing.T) {
	f := frontend.NewServer("127.0.0.1", 8080)
	f.Start()
	defer f.Stop()
	//defer server.Close()
	// Correct request
	resp := requestOutlookSignIn(t)
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		log.Errorf("Failed with status %d %s", resp.StatusCode, resp.Status)
		t.Fail()
	}
	log.Infoln(resp.Header.Get("Location"))

	os.Setenv("API_ROOT", "")
	// Bad request
	resp = requestOutlookSignIn(t)
	log.Infoln(resp.Header.Get("Location"))
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 500 && resp.Header.Get("Location") != "htpp://localhost:8080/" {
		log.Errorf("Failed with status %d %s", resp.StatusCode, resp.Status)
		t.Fail()
	}

	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}
func requestOutlookSignIn(t *testing.T) *http.Response {
	// Set up the HTTP request
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/SignInWithOutlook", nil)
	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
	}

	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)

	if err != nil {
		log.Errorln(err.Error())
		t.Fail()
	}

	return resp
}
