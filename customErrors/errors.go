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

// DecodedError TODO
type DecodedError struct {
	Message string
}

func (err DecodedError) Error() string {
	return fmt.Sprintln(err.Message)
}

// EmptyValueError TODO
type EmptyValueError struct {
}

func (err EmptyValueError) Error() string {
	return fmt.Sprintln("Values cannot be empty")
}

type ConfigNotChargedCorrectlyError struct {
	Message string
}

func (err ConfigNotChargedCorrectlyError) Error() string {
	return fmt.Sprintln(err.Message)
}

type WrongKindError struct {
	Mail string
}

func (err *WrongKindError) Error() string {
	return fmt.Sprintf("wrong kind of account %s", err.Mail)
}

type NotFoundError struct {
	Message string
}

func (err *NotFoundError) Error() string {
	return err.Message
}
