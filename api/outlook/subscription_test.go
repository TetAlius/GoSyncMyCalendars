package outlook_test

import (
	"net/http"
	"testing"

	"io/ioutil"

	"encoding/json"

	"os"
)

type TunnelJSON struct {
	Tunnels []Tunnel `json:"tunnels"`
}
type Tunnel struct {
	PublicURL string `json:"public_url"`
	Proto     string `json:"proto"`
}

func TestNGrokConfig(t *testing.T) {
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
