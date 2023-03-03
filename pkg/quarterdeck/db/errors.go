package db

import "errors"

var (
	ErrTokenMissingEmail = errors.New("token is missing email address")
	ErrTokenExpired      = errors.New("token has expired")
	ErrInvalidSecret     = errors.New("invalid secret for token verification")
	ErrTokenInvalid      = errors.New("token has invalid signature")
)
