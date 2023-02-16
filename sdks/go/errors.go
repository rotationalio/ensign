package ensign

import "errors"

var (
	ErrMissingEndpoint     = errors.New("endpoint is required")
	ErrMissingClientID     = errors.New("client ID is required")
	ErrMissingClientSecret = errors.New("client secret is required")
)

type Errorer interface {
	Err() error
}
