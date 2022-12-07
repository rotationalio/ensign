package raft

import "errors"

var (
	ErrCannotSetRunningState = errors.New("can only set the running state from the initialized state")
)
