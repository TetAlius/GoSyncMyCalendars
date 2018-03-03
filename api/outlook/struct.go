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

