package redlock

import (
	"errors"
	"fmt"
	"strings"
)

var ErrAcquireLockFailed = errors.New("redlock: failed to acquire lock")

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Errors []error
}

func (e *Error) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	errors := make([]string, 0, len(e.Errors))
	for i, err := range e.Errors {
		errors = append(errors, fmt.Sprintf("#%d %s", i, err.Error()))
	}

	return strings.Join(errors, ",")
}

// Append will append more errors onto an Error in order to create a larger multi-error.
func (e *Error) Append(errs ...error) {
	for _, err := range errs {
		e.Errors = append(e.Errors, err)
	}
}
