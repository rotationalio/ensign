package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

var (
	ErrNotConnected   = errors.New("not connected to trtl database")
	ErrNotFound       = errors.New("object not found for the specified key")
	ErrUnavailable    = errors.New("trtl database service is unavailable")
	ErrListBreak      = errors.New("on list item has stopped iterating")
	ErrOrgNotVerified = errors.New("could not verify organization")

	// Missing fields
	ErrMissingID           = errors.New("object requires id for serialization")
	ErrMissingOrgID        = errors.New("object requires organization id for serialization")
	ErrMissingTenantID     = errors.New("object requires tenant id for serialization")
	ErrMissingProjectID    = errors.New("object requires project id for serialization")
	ErrMissingMemberEmail  = errors.New("member email is required")
	ErrMissingMemberRole   = errors.New("member role is required")
	ErrMissingMemberStatus = errors.New("member status is required")
	ErrMissingProjectName  = errors.New("project name is required")
	ErrMissingTenantName   = errors.New("tenant name is required")
	ErrMissingEnvType      = errors.New("tenant environment type is required")
	ErrMissingTopicName    = errors.New("topic name is required")
	ErrMissingPageSize     = errors.New("cannot list database without a page size")
	ErrMissingOwnerID      = errors.New("model is missing owner id")

	// Invalid fields
	ErrNameTooLong         = errors.New("name cannot be longer than 1024 characters")
	ErrOrganizationTooLong = errors.New("organization name cannot be longer than 1024 characters")
	ErrWorkspaceTooLong    = errors.New("workspace name cannot be longer than 1024 characters")
	ErrProfessionTooLong   = errors.New("profession segment cannot be longer than 1024 characters")
	ErrDeveloperTooLong    = errors.New("developer segment cannot be longer than 1024 characters")
	ErrInvalidWorkspace    = errors.New("workspace name must be at least 3 characters and cannot start with a number")
	ErrUnknownMemberRole   = errors.New("unknown member role")
	ErrInvalidProjectName  = errors.New("invalid project name")
	ErrInvalidTenantName   = errors.New("invalid tenant name")
	ErrInvalidTopicName    = errors.New("invalid topic name")
	ErrInvalidStorage      = errors.New("data storage must be greater than or equal to 0")

	// Database state errors
	ErrMemberExists        = errors.New("member already exists")
	ErrMemberEmailNotFound = errors.New("member does not exist")

	// Key errors
	ErrKeyNoID      = errors.New("key does not contain an id")
	ErrKeyWrongSize = errors.New("key is not the correct size")

	// Max-length errors
	ErrProjectDescriptionTooLong = errors.New("project description is too long")
	ErrTopicNameTooLong          = errors.New("topic name is too long")
)

// ValidationError represents a validation error for a specific field and may contain
// an index if the field is a slice.
type ValidationError struct {
	Field string
	Err   error
	Index int
}

func validationError(field string, err error) *ValidationError {
	return &ValidationError{
		Field: field,
		Err:   err,
		Index: -1,
	}
}

func (e *ValidationError) AtIndex(index int) *ValidationError {
	e.Index = index
	return e
}

func (e *ValidationError) Error() string {
	switch {
	case e.Field == "":
		return e.Err.Error()
	case e.Index > -1:
		return fmt.Sprintf("validation error for field %s at index %d: %s", e.Field, e.Index, e.Err.Error())
	default:
		return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Err.Error())
	}
}

func (e *ValidationError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type ValidationErrors []*ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, e := range v {
		errs = append(errs, e.Error())
	}
	return fmt.Sprintf("%d validation errors occurred:\n  -%s", len(v), strings.Join(errs, "\n  -"))
}

func (v ValidationErrors) ToAPI() api.FieldValidationErrors {
	errs := make(api.FieldValidationErrors, 0, len(v))
	for _, e := range v {
		errs = append(errs, &api.FieldValidationError{
			Field: e.Field,
			Err:   e.Err.Error(),
			Index: e.Index,
		})
	}
	return errs
}
