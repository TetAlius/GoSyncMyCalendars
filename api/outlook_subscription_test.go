package api_test

/*
import (
	"net/http"
	"testing"

	"io/ioutil"

	"encoding/json"

	"os"

	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/backend"
)

type Tunnel struct {
	PublicURL string `json:"public_url"`
	Proto     string `json:"proto"`
}

func setupNgrok(t *testing.T) {
	resp, err := http.Get("http://localhost:4040/api/tunnels/gosync")
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

	response := new(Tunnel)
	err = json.Unmarshal(contents, &response)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong: %s", err.Error())
		return
	}
	ngrokURI := response.PublicURL

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
	account, _ := setup()
	//Refresh previous petition in order to have tokens updated
	account.Refresh()

	subscription := api.NewOutlookSubscription("", ngrokURL, "Created,Deleted,Updated")

	calendar, err := account.GetPrimaryCalendar()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s", err.Error())
		return
	}
	b := backend.NewServer("127.0.0.1", 8081, 10, nil)
	go func() {
		b.Start()
	}()

	err = subscription.Subscribe(calendar)
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s. ID: %s", err.Error(), subscription.ID)
		return
	}

	subscription = api.RetrieveOutlookSubscription(subscription.ID, subscription.Uuid, calendar, subscription.Type)

	err = subscription.Renew()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s. ID: %s", err.Error(), subscription.ID)
		return
	}
	err = subscription.Delete()
	if err != nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected nil found error: %s. ID: %s", err.Error(), subscription.ID)
		return
	}

	// Wrong calls to subscription
	subs := api.NewOutlookSubscription("", "localhost:8081", "Created,Deleted,Updated")
	err = subs.Subscribe(calendar)
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
	subs = api.RetrieveOutlookSubscription(subs.ID, subs.Uuid, calendar, subs.Type)
	err = subscription.Renew()
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
	err = subscription.Delete()
	if err == nil {
		t.Fail()
		t.Fatalf("something went wrong. Expected error found nil")
		return
	}
}
*/
