package backend

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/TetAlius/logs/logger"
)

//NewRequest TODO Creates and executes the request for all petitions
//and returns the JSON so that it can be parsed into the correct struct
func NewRequest(method string, url string, body io.Reader, authorization string, anchorMailbox string) (contents []byte, err error) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Errorf("Error creating new request: %s", err.Error())
	}

	//Add the authorization to the header
	req.Header.Set("Authorization", authorization)

	//Add the anchorMailbox to the header
	req.Header.Set("X-AnchorMailbox", anchorMailbox)

	// If body is given, has to put a content-Type json on the header
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing request: %s", err.Error())
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err.Error())
	}
	return
}
