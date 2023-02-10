package service

import "errors"

var (
	ErrNoServiceRegistered = errors.New("no service has been registered with the server")
)
