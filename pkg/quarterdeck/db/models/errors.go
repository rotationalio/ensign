package models

import "errors"

var (
	ErrNotFound           = errors.New("object not found in the database")
	ErrInvalidUser        = errors.New("user model is not correctly populated")
	ErrInvalidPassword    = errors.New("user password should be stored as an argon2 derived key")
	ErrMissingUserID      = errors.New("user model does not have an ID assigned")
	ErrMissingKeyMaterial = errors.New("apikey model requires client id and secret")
	ErrInvalidSecret      = errors.New("apikey secrets should be stored as argon2 derived keys")
	ErrMissingProjectID   = errors.New("apikey model requires project id")
	ErrMissingKeyName     = errors.New("apikey model requires name")
)
