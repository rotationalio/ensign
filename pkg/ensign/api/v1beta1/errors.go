package api

import "errors"

// Statically defined errors for error checking the type of error returned by a method
// or function in the api package.
var (
	ErrNoEvent   = errors.New("event wrapper contains no event")
	ErrNoKeys    = errors.New("no keys specified for key based hashing")
	ErrNoFields  = errors.New("no fields specified for field based hashing")
	ErrNoGroupID = errors.New("consumer group requires either id or name")
)
