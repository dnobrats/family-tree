// Package domain holds the core entities and the errors that cross layer
// boundaries. It intentionally imports nothing from the rest of the project
// (no SQL, no HTTP) so it stays easy to read and reason about.
package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound means the requested record does not exist.
	ErrNotFound = errors.New("not found")
	// ErrConflict means a uniqueness/duplication rule was violated.
	ErrConflict = errors.New("conflict")
	// ErrValidation is the parent of every ValidationError.
	ErrValidation = errors.New("validation error")
)

// ValidationError is a user-facing error: its message is safe to show to the
// person who submitted the form. Handlers check errors.Is(err, ErrValidation)
// and only then surface err.Error() to the client.
type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }
func (e *ValidationError) Unwrap() error { return ErrValidation }

// Invalid builds a ValidationError with a formatted, user-safe message.
func Invalid(format string, args ...any) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}
