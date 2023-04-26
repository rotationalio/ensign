package emails

import "errors"

var (
	ErrMissingSubject   = errors.New("missing email subject")
	ErrMissingSender    = errors.New("missing email sender")
	ErrMissingRecipient = errors.New("missing email recipient")
	ErrUnparsable       = errors.New("could not parse email address")
)
