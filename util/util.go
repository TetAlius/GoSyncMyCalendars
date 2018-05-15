package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"time"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func MailFromToken(tokens []string, addFinal string) (email string, preferred bool, err error) {
	if len(tokens) < 2 {
		return "", false, errors.New("TokenID was not parsed correctly")
	}
	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1) + addFinal
	log.Debugf("EncodedToken: %s", encodedToken)

	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)

	if err != nil {
		log.Errorf("Error decoding token: %s", err.Error())
	}
	log.Debugf("DecodedToken: %s", decodedToken)
	var f interface{}
	err = json.Unmarshal(decodedToken, &f)
	if err != nil {
		log.Errorf("Error unmarshaling decoded token: %s", err.Error())
	}

	if reflect.TypeOf(f) != nil && reflect.TypeOf(f).Kind() == reflect.Map {
		m := f.(map[string]interface{})

		if m["email"] != nil {
			log.Debugf("Got email %s", m["email"].(string))
			email = m["email"].(string)
			preferred = false
		} else if m["preferred_username"] != nil {
			log.Debugf("Got preferred email %s", m["preferred_username"].(string))
			email = m["preferred_username"].(string)
			preferred = true
		} else {
			err = customErrors.DecodedError{Message: "Not email nor preferred_username on token"}
		}
	} else {
		err = customErrors.DecodedError{Message: "Decoded token is not a map"}
	}
	return
}

func CallAPIRoot(route string) (apiRoute string, err error) {
	root := os.Getenv("API_ROOT")
	if len(root) == 0 {
		return "", errors.New("API_ROOT not given on environment")
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", root+route, nil)

	if err != nil {
		return "", errors.New(fmt.Sprintf("Error creating API request: %s", err.Error()))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error doing API request: %s", err.Error()))
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response body from API request: %s", err.Error()))
	}

	return strings.Replace(string(contents[:]), "\"", "", -1), nil
}

//DoRequest TODO Creates and executes the request for all petitions
//TODO return responseCode
//and returns the JSON so that it can be parsed into the correct struct
func DoRequest(method string, url string, body io.Reader, headers map[string]string, params map[string]string) (contents []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return contents, errors.New(fmt.Sprintf("error creating new request: %s", err.Error()))
	}

	//Set all headers for request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// If body is given, has to put a content-Type json on the header
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return contents, errors.New(fmt.Sprintf("error doing request: %s", err.Error()))
	}

	defer resp.Body.Close()
	//TODO parse errors and content
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return contents, errors.New(fmt.Sprintf("error reading response body: %s", err.Error()))
	}

	return
}
