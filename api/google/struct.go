package google

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	ConcreteError `json:"error,omitempty"`
}
type ConcreteError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func createResponseError(contents []byte) (err error) {
	e := new(Error)
	err = json.Unmarshal(contents, &e)
	if e.Code != 0 && len(e.Message) != 0 {
		return errors.New(fmt.Sprintf("code: %d. message: %s", e.Code, e.Message))
	}
	return nil
}
