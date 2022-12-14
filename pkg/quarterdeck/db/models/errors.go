package models

import "errors"

var (
	ErrNotFound      = errors.New("object not found in the database")
	ErrInvalidUser   = errors.New("user model is not correctly populated")
	ErrMissingUserID = errors.New("user model does not have an ID assigned")
)
