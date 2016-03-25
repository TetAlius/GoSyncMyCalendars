package backend

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//NewRequest TODO Creates and executes the request for all google petitions
//and returns the JSON so that it can be parsed into the correct struct
func NewRequest(
	method string,
	url string,
	body io.Reader,
	authorization string,
	anchorMailbox string) (contents []byte) {
	client := http.Client{}
	req, err := http.NewRequest(method,
		url,
		body)

	if err != nil {
		fmt.Println(err)
	}
	//Add the authorization to the header
	req.Header.Set("Authorization",
		authorization)

	//Add the anchorMailbox to the header
	req.Header.Set("X-AnchorMailbox", anchorMailbox)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return contents
}
