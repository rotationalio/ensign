package models

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("object not found in the database")
	ErrInvalidUser        = errors.New("user model is not correctly populated")
	ErrInvalidPassword    = errors.New("user password should be stored as an argon2 derived key")
	ErrMissingModelID     = errors.New("model does not have an ID assigned")
	ErrMissingKeyMaterial = errors.New("apikey model requires client id and secret")
	ErrInvalidSecret      = errors.New("apikey secrets should be stored as argon2 derived keys")
	ErrMissingOrgID       = errors.New("model does not have an organization ID assigned")
	ErrMissingProjectID   = errors.New("apikey model requires project id")
	ErrMissingKeyName     = errors.New("apikey model requires name")
	ErrNoPermissions      = errors.New("apikey model requires permissions")
	ErrModifyPermissions  = errors.New("cannot modify permissions on an existing APIKey object")
	ErrMissingPageSize    = errors.New("cannot list database without a page size")
	ErrInvalidCursor      = errors.New("could not compute the next page of results")
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
