package db

import "errors"

var (
	ErrTokenMissingEmail = errors.New("email verification token is missing email address")
	ErrTokenExpired      = errors.New("email verification token has expired")
	ErrInvalidSecret     = errors.New("invalid secret for email token verification")
	ErrTokenInvalid      = errors.New("email verification token has invalid signature")
)
