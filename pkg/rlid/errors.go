package rlid

import "errors"

var (
	// Returned when constructing an RLID with a timestamp that is larger than maxTime
	ErrOverTime = errors.New("timestamp is too far in the future to encode")
)
