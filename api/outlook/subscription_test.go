package outlook_test

import (
	"net/http"
	"testing"

	"io/ioutil"

	"encoding/json"

	"os"

	"github.com/TetAlius/GoSyncMyCalendars/api/outlook"
	"github.com/TetAlius/GoSyncMyCalendars/backend"

	"fmt"
)

type TunnelJSON struct {
	Tunnels []Tunnel `json:"tunnels"`
}
type Tunnel struct {
	PublicURL string `json:"public_url"`
	Proto     string `json:"proto"`
}

func setupNgrok(t *testing.T) {
	resp, err := http.Get("http://localhost:4040/api/tunnels")
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong: %s", err.Error())
		return
	}

	response := new(TunnelJSON)
	err = json.Unmarshal(contents, &response)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong: %s", err.Error())
		return
	}
	ngrokURI := response.Tunnels[0].PublicURL

	os.Setenv("NGROK_URI", ngrokURI)
	if len(os.Getenv("NGROK_URI")) < 0 {
		t.Fail()
		t.Fatal("NGROK was incorrectly used")
		return
	}
	t.Log(os.Getenv("NGROK_URI"))
}

func TestOutlookSubscription_SubscriptionLifeCycle(t *testing.T) {
	setupApiRoot()
	setupNgrok(t)
	ngrokURL := fmt.Sprintf("%s/outlook/watcher", os.Getenv("NGROK_URI"))
	account := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	subscription := outlook.NewSubscription("", ngrokURL, "Created,Deleted,Updated")

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	b := backend.NewServer("127.0.0.1", 8081)
	go func() {
		b.Start()
	}()

	err = subscription.Subscribe(account, calendar)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	err = subscription.Renew(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	err = subscription.Delete(account)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}

	// Wrong calls to subscription
	subs := outlook.NewSubscription("", "localhost:8081", "Created,Deleted,Updated")
	err = subs.Subscribe(account, calendar)
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
	err = subscription.Renew(account)
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
	err = subscription.Delete(account)
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
	b.Stop()
}
