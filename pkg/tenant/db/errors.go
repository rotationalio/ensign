package db

import "errors"

var (
	ErrNotConnected = errors.New("not connected to trtl database")
	ErrNotFound     = errors.New("object not found for the specified key")
	ErrMissingID    = errors.New("object requires id for serialization")
)
