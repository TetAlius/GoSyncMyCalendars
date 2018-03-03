package outlook

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	ConcreteError `json:"error,omitempty"`
}
type ConcreteError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func createResponseError(contents []byte) (err error) {
	e := new(Error)
	err = json.Unmarshal(contents, &e)
	if len(e.Code) != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %s. message: %s", e.Code, e.Message))
	}
	return nil
}
