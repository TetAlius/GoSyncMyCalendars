package handlers_test

import (
	"net/http"
	"testing"

	"os"
)

//TestGoogleSignInHandler test the SignInHandler method
func TestGoogleSignInHandler(t *testing.T) {
	setupApiRoot()

	// Correct request
	resp := requestGoogleSignIn(t)
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		t.Fail()
		t.Fatalf("Failed with status %d %s", resp.StatusCode, resp.Status)
		return
	}

	os.Setenv("API_ROOT", "")
	// Bad request
	resp = requestGoogleSignIn(t)
	// Check if you received the status codes you expect. There may
	// status codes other than 200 which are acceptable.
	if resp.StatusCode != 500 && resp.Header.Get("Location") != "http://localhost:8080/" {
		t.Fail()
		t.Fatalf("Failed with status %d %s", resp.StatusCode, resp.Status)
		return
	}

	os.Setenv("API_ROOT", os.Getenv("API_ROOT_TEST"))
}

func requestGoogleSignIn(t *testing.T) (resp *http.Response) {
	// Set up the HTTP request
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/SignInWithGoogle", nil)
	if err != nil {
		t.Fail()
		t.Fatal(err.Error())
		return
	}

	transport := http.Transport{}
	resp, err = transport.RoundTrip(req)

	if err != nil {
		t.Fail()
		t.Fatal(err.Error())
		return
	}

	return resp
}
