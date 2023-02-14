package mtls

import "errors"

// Standard errors for error type checking
var (
	ErrPrivateKeyRequired = errors.New("provider must contain a private key to initialize TLS certs")
	ErrNoCertificates     = errors.New("provider does not contain any certificates")
	ErrMissingKey         = errors.New("provider does not contain a private key")
	ErrZipEmpty           = errors.New("zip archive contains no providers")
	ErrZipTooMany         = errors.New("multiple providers in zip, is this a provider pool?")
)
