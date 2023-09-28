package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

var (
	unsuccessful              = Reply{Success: false}
	notFound                  = Reply{Success: false, Error: "resource not found"}
	notAllowed                = Reply{Success: false, Error: "method not allowed"}
	ErrOrganizationIDRequired = errors.New("organization id is required for this endpoint")
	ErrAPIKeyIDRequired       = errors.New("apikey id is required for this endpoint")
	ErrMemberIDRequired       = errors.New("member id is required for this endpoint")
	ErrProjectIDRequired      = errors.New("project id is required for this endpoint")
	ErrTenantIDRequired       = errors.New("tenant id is required for this endpoint")
	ErrTopicIDRequired        = errors.New("topic id is required for this endpoint")
	ErrTokenRequired          = errors.New("token is required for this endpoint")
	ErrPasswordRequired       = errors.New("password is required for this endpoint")
	ErrPasswordMismatch       = errors.New("passwords do not match")
	ErrInvalidTenantField     = errors.New("invalid tenant field")
	ErrMissingQueryField      = errors.New("missing query field")
	ErrQueryTooLong           = errors.New("query string is too long, please use the SDKs for complex queries")
	ErrInvalidUserClaims      = errors.New("user claims invalid or unavailable")
	ErrUnparsable             = errors.New("could not parse request")
	ErrNoCookies              = errors.New("no cookies available")
	ErrNoRefreshToken         = errors.New("refresh token not found in cookies")
	ErrNoAccessToken          = errors.New("access token not found in cookies")
)

// FieldValidationError represents a validation error for a specific field, when the
// frontend needs to know which field is in error, and includes the name of the field
// and the index if the field is a list along with the error string.
type FieldValidationError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
	Index int    `json:"index"`
}

func (e *FieldValidationError) Error() string {
	switch {
	case e.Field == "":
		return e.Err
	case e.Index > -1:
		return fmt.Sprintf("%q at index %d: %s", e.Field, e.Index, e.Err)
	default:
		return fmt.Sprintf("%q: %s", e.Field, e.Err)
	}
}

type FieldValidationErrors []*FieldValidationError

func (e FieldValidationErrors) Error() string {
	errs := make([]string, 0, len(e))
	for _, err := range e {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("%d validation errors occurred:\n  -%s", len(e), strings.Join(errs, "\n  -"))
}

func FieldTypeError(field string, t string) error {
	return fmt.Errorf("invalid type for field %q, expected %q", field, t)
}

func InvalidFieldError(field string) error {
	return fmt.Errorf("invalid field %q", field)
}

// Constructs a new response for an error or returns unsuccessful.
func ErrorResponse(err interface{}) Reply {
	if err == nil {
		return unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case FieldValidationErrors:
		rep.ValidationErrors = err
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

// NotFound returns a JSON response for the API.
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

// ReplyQuarterdeckError returns a JSON response for a Quarterdeck error by attempting
// to decode a generic error into a StatusError. If the error is not a StatusError,
// then a JSON 500 response is returned.
// TODO: Does this need to have more user-friendly error messaging? :point-down:
func ReplyQuarterdeckError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusOK, Reply{Success: true})
		return
	}

	if serr, ok := err.(*qd.StatusError); ok {
		if serr.StatusCode == 0 {
			serr.StatusCode = http.StatusInternalServerError
		}
		c.JSON(serr.StatusCode, serr.Reply)
	} else {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}
