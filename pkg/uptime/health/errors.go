package health

import "errors"

var (
	ErrTooManyRedirects = errors.New("too many redirects, reporting unhealthy server")
	ErrNoContent        = errors.New("api response did not contain any content")
	ErrNoStatusResponse = errors.New("api response did not contain a status")
	ErrUnparsableStatus = errors.New("api response status is unparsable")
)
