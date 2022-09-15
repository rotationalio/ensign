package rlid

import "errors"

var (
	// Returned when constructing an RLID with a timestamp that is larger than maxTime
	ErrOverTime = errors.New("timestamp is too far in the future to encode")

	// Returned when marshaling RLIDs to a buffer that is the incorrect size.
	ErrBufferSize = errors.New("bad buffer size: cannot encode the id")

	// Returned when unmarshaling RLIDs from a string that is the incorrect length.
	ErrDataSize = errors.New("bad data size: cannot decode the id from string")

	// Returned when unmarshaling RLIDs in strict mode if the string contains bad chars
	ErrInvalidCharacters = errors.New("invalid characters: cannot decode the id from string")
)
