package ensql

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyQuery   = errors.New("empty query is invalid")
	ErrMissingTopic = errors.New("topic name cannot be empty")
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
