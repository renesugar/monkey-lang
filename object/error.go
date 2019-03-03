package object

import (
	"errors"
)

// Error is the error type and used to hold a message denoting the details of
// error encountered. This object is trakced through the evaluator and when
// encountered stops evaulation of the program or body of a function.
type Error struct {
	Message string
}

func (e *Error) Error() error {
	return errors.New(e.Message)
}

func (e *Error) Bool() bool {
	return false
}

func (e *Error) String() string {
	return e.Message
}

// Clone creates a new copy
func (e *Error) Clone() Object {
	return &Error{Message: e.Message}
}

// Type returns the type of the object
func (e *Error) Type() Type { return ERROR }

// Inspect returns a stringified version of the object for debugging
func (e *Error) Inspect() string { return "ERROR: " + e.Message }
