package customErrors

import "fmt"

// AccError provides error for when an account given
// to the database is not supported.
type AccError struct {
	Account string
}

func (err AccError) Error() string {
	return fmt.Sprintf("Account %s not supported", err.Account)
}

// RowsError provides error for when the rows affected on an db.Exec
// differ from the ones you expect
type RowsError struct {
	Expected int64
	Received int64
	Message  string
}

func (err RowsError) Error() string {
	return fmt.Sprintf("Expected: %d Received: %d. %s", err.Expected, err.Received, err.Message)
}
