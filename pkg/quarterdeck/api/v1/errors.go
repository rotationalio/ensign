package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	unsuccessful = Reply{Success: false}
	notFound     = Reply{Success: false, Error: "resource not found"}
	notAllowed   = Reply{Success: false, Error: "method not allowed"}
)

var (
	ErrInvalidCredentials = errors.New("quarterdeck credentials are missing or invalid")
	ErrExpiredCredentials = errors.New("quarterdeck credentials have expired")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrPasswordTooWeak    = errors.New("password is too weak: use a combination of upper and lower case letters, numbers, and special characters")
	ErrMissingField       = errors.New("missing required field")
	ErrInvalidField       = errors.New("invalid or unparsable field")
	ErrRestrictedField    = errors.New("field restricted for request")
	ErrModelIDMismatch    = errors.New("resource id does not match id of endpoint")
	ErrUserExists         = errors.New("user or organization already exists")
)

// Construct a new response for an error or simply return unsuccessful.
func ErrorResponse(err interface{}) Reply {
	if err == nil {
		return unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}

// NotFound returns a JSON 404 response for the API.
// NOTE: we know it's weird to put server-side handlers like NotFound and NotAllowed
// here in the client/api side package but it unifies where we keep our error handling
// mechanisms.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}

// FieldError provides a general mechanism for specifying errors with specific API
// object fields such as missing required field or invalid field and giving some
// feedback about which fields are the problem.
// TODO: allow multiple field errors to be specified in one response.
type FieldError struct {
	Field string `json:"field"`
	Err   error  `json:"error"`
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Field)
}

func (e *FieldError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

func MissingField(field string) error {
	return &FieldError{Field: field, Err: ErrMissingField}
}

func InvalidField(field string) error {
	return &FieldError{Field: field, Err: ErrInvalidField}
}

func RestrictedField(field string) error {
	return &FieldError{Field: field, Err: ErrRestrictedField}
}

// StatusError decodes an error response from Quarterdeck.
type StatusError struct {
	StatusCode int
	Reply      Reply
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Reply.Error)
}
