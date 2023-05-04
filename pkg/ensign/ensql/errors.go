package ensql

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyQuery             = errors.New("empty query is invalid")
	ErrMissingTopic           = errors.New("topic name cannot be empty")
	ErrUnhandledStep          = errors.New("parser has reached an unhandled state")
	ErrNoFieldsSelected       = errors.New("SELECT requires field projection or *")
	ErrInvalidSelectAllFields = errors.New("cannot select * and specify fields")
	ErrNonNumeric             = errors.New("cannot parse non-numeric token as a number")
)

type SyntaxError struct {
	position int
	near     string
	message  string
}

func (e *SyntaxError) Error() string {
	if e.message != "" {
		return fmt.Sprintf("syntax error at position %d near %q: %s", e.position, e.near, e.message)
	}
	return fmt.Sprintf("syntax error at position %d near %q", e.position, e.near)
}

func Error(pos int, near, msg string) *SyntaxError {
	return &SyntaxError{pos, near, msg}
}

func Errorf(pos int, near, format string, args ...interface{}) *SyntaxError {
	return Error(pos, near, fmt.Sprintf(format, args...))
}

// InvalidState is a developer error; it means that the parser proceeded to a step but
// modified the underlying state of the parsing incorrectly. This is not a user error,
// e.g. due to syntax, so invalid state usually panics.
func InvalidState(expected, actual string) error {
	return fmt.Errorf("%w: expected %q but received %q", ErrUnhandledStep, expected, actual)
}
