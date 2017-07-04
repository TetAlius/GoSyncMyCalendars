package handlers_test

import (
	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	"net/http"
	"testing"
)

//TestGoogle_SignInHandler test the SignInHandler method
func TestGoogle_SignInHandler(t *testing.T) {
	setup()
	//defer server.Close()
	// Set up the HTTP request
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/SignInWithGoogle", nil)

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

func TestGoogle_NewGoogleHandler(t *testing.T) {
	configNameErr := "configgggg.json"
	_, err := handlers.NewGoogleHandler(&configNameErr)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}

	configNameErr = "./badConfig.json"
	_, err = handlers.NewGoogleHandler(&configNameErr)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}
	configName := "../../google.json"
	_, err = handlers.NewGoogleHandler(&configName)
	if err != nil {
		t.Fatalf("Error should be nil but it is: %s", err.Error())
	}

	notThisConfig := "../../outlook.json"
	_, err = handlers.NewGoogleHandler(&notThisConfig)
	if err == nil {
		t.Fatal("Error should not be nil but it is")
	}
}

//func TestMain(m *testing.M) {
//
//	//Seccion mala
//
//	configNameErr := "configgggg.json"
//	configNameGood := "../../config.json"
//
//	f := frontend.NewFrontend("127.0.0.1", 8080, &configNameErr)
//	f.Start()
//	m.Run()
//
//	f.Stop()
//
//	log.Debugln("Main")
//	f = frontend.NewFrontend("127.0.0.1", 8080, &configNameGood)
//	f.Start()
//	m.Run()
//	f.Stop()
//
//	//Seccion buena
//}
