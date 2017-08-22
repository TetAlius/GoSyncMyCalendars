package util

import (
	"encoding/base64"
	"encoding/json"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func MailFromToken(tokens []string) (email string, preferred bool, err error) {
	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1) + "=="
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

func CallAPIRoot(route string) string {
	client := http.Client{}
	root := os.Getenv("API_ROOT")
	req, err := http.NewRequest("GET", root+route, nil)

	if err != nil {
		log.Errorf("Error creating API request: %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error doing API request: %s", err.Error())
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body from API request: %s", err.Error())
	}

	return strings.Replace(string(contents[:]), "\"", "", -1)
}
