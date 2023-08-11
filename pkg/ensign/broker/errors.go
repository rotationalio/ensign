package broker

import "errors"

var (
	ErrBrokerNotRunning = errors.New("operation could not be completed: broker is not running")
	ErrUnknownID        = errors.New("no publisher or subscriber registered with specified id")
)
