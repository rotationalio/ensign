package pem

import "errors"

var (
	ErrDecodePrivateKey  = errors.New("could not decode PEM private key")
	ErrDecodePublicKey   = errors.New("could not decode PEM public key")
	ErrDecodeCertificate = errors.New("could not decode PEM certificate")
	ErrDecodeCSR         = errors.New("could not decode PEM certificate request")
)
