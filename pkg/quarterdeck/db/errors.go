package db

import "errors"

var (
	ErrMissingEmail       = errors.New("email address is required")
	ErrMissingUserID      = errors.New("user id is required")
	ErrTokenMissingEmail  = errors.New("email verification token is missing email address")
	ErrTokenMissingUserID = errors.New("email verification token is missing user id")
	ErrTokenExpired       = errors.New("email verification token has expired")
	ErrInvalidSecret      = errors.New("invalid secret for email token verification")
	ErrTokenInvalid       = errors.New("email verification token has invalid signature")
	ErrSQLite3Conn        = errors.New("could not get sqlite3 connection for backups")
)
