package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound             = errors.New("object not found in the database")
	ErrInvalidOrganization  = errors.New("organization model is not correctly populated")
	ErrUserOrganization     = errors.New("user is not associated with the organization")
	ErrUserOrgExists        = errors.New("user is already associated with the organization")
	ErrInvalidUser          = errors.New("user model is not correctly populated")
	ErrInvalidPassword      = errors.New("user password should be stored as an argon2 derived key")
	ErrMissingModelID       = errors.New("model does not have an ID assigned")
	ErrMissingKeyMaterial   = errors.New("apikey model requires client id and secret")
	ErrInvalidSecret        = errors.New("apikey secrets should be stored as argon2 derived keys")
	ErrMissingRole          = errors.New("missing role")
	ErrInvalidRole          = errors.New("invalid role")
	ErrNoOwnerRole          = errors.New("organization is missing an owner")
	ErrOwnerRoleConstraint  = errors.New("organization must have at least one owner")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrMissingOrgID         = errors.New("model does not have an organization ID assigned")
	ErrMissingProjectID     = errors.New("model requires project id")
	ErrInvalidProjectID     = errors.New("invalid project id for apikey")
	ErrMissingKeyName       = errors.New("apikey model requires name")
	ErrMissingCreatedBy     = errors.New("apikey model requires created by")
	ErrNoPermissions        = errors.New("apikey model requires permissions")
	ErrInvalidPermission    = errors.New("invalid permission specified for apikey")
	ErrModifyPermissions    = errors.New("cannot modify permissions on an existing APIKey object")
	ErrMissingPageSize      = errors.New("cannot list database without a page size")
	ErrInvalidCursor        = errors.New("could not compute the next page of results")
	ErrDuplicate            = errors.New("unique constraint violated on model")
	ErrMissingRelation      = errors.New("foreign key relation violated on model")
	ErrNotNull              = errors.New("not null constraint violated on model")
	ErrConstraint           = errors.New("database constraint violated")
	ErrRequiresConfirmation = errors.New("row update requires confirmation")
)

type ValidationError struct {
	err error
}

func invalid(err error) *ValidationError {
	return &ValidationError{err}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.err)
}

func (e *ValidationError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

// ConstraintError attempts to parse a sqlite3.ErrConstraint error into a model error.
type ConstraintError struct {
	err   error
	dberr sqlite3.Error
}

func constraint(dberr sqlite3.Error) *ConstraintError {
	// String parsing seems to be the only way to deal with error handling for sqlite3
	errs := dberr.Error()
	switch {
	case strings.HasPrefix(errs, "UNIQUE"):
		return &ConstraintError{err: ErrDuplicate, dberr: dberr}
	case strings.HasPrefix(errs, "FOREIGN KEY"):
		return &ConstraintError{err: ErrMissingRelation, dberr: dberr}
	case strings.HasPrefix(errs, "NOT NULL"):
		return &ConstraintError{err: ErrNotNull, dberr: dberr}
	default:
		return &ConstraintError{err: ErrConstraint, dberr: dberr}
	}
}

func (e *ConstraintError) Error() string {
	if e.dberr.Code == sqlite3.ErrConstraint {
		return e.err.Error()
	}
	return e.dberr.Error()
}

func (e *ConstraintError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *ConstraintError) Unwrap() error {
	return e.err
}
