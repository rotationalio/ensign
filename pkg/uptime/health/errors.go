package health

import "errors"

var (
	ErrTooManyRedirects = errors.New("too many redirects, reporting unhealthy server")
	ErrNoContent        = errors.New("api response did not contain any content")
	ErrNoStatusResponse = errors.New("api response did not contain a status")
	ErrUnparsableStatus = errors.New("api response status is unparsable")
	ErrNoServiceID      = errors.New("service id must be specified before status can be saved to db")
	ErrNoTimestamp      = errors.New("the service status does not have a checked at timestamp")
)
