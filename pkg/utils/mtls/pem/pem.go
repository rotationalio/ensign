package pem

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"
)

// PEM Block types
const (
	BlockPublicKey          = "PUBLIC KEY"
	BlockPrivateKey         = "PRIVATE KEY"
	BlockRSAPublicKey       = "RSA PUBLIC KEY"
	BlockRSAPrivateKey      = "RSA PRIVATE KEY"
	BlockECPrivateKey       = "EC PRIVATE KEY"
	BlockCertificate        = "CERTIFICATE"
	BlockCertificateRequest = "CERTIFICATE REQUEST"
)

// EncodePrivateKey as a PKCS8 ASN.1 DER key and write a PEM block with type "PRIVATE KEY"
func EncodePrivateKey(key interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := encodePrivateKeyTo(key, &b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodePrivateKeyTo(key interface{}, out io.Writer) (err error) {
	block := &pem.Block{Type: BlockPrivateKey}
	if block.Bytes, err = x509.MarshalPKCS8PrivateKey(key); err != nil {
		return err
	}
	return pem.Encode(out, block)
}

// DecodePrivateKey from a PEM encoded block. If the block type is "EC PRIVATE KEY",
// then the block is parsed as an EC private key in SEC 1, ASN.1 DER form. If the block
// is "RSA PRIVATE KEY" then it is decoded as a PKCS 1, ASN.1 DER form. If the block
// type is "PRIVATE KEY", the block is decoded as a PKCS 8 ASN.1 DER key, if that fails,
// then the PKCS 1 and EC parsers are tried in that order, before returning an error.
func DecodePrivateKey(in []byte) (interface{}, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, ErrDecodePrivateKey
	}
	return ParsePrivateKey(block)
}

// ParsePrivateKey from PEM block. May return an *ecdsa.PrivateKey, *rsa.PrivateKey, or
// ed25519.PrivateKey depending on the block type and the x509 parsing method.
func ParsePrivateKey(block *pem.Block) (interface{}, error) {
	// EC PRIVATE KEY specific handling
	if block.Type == BlockECPrivateKey {
		return x509.ParseECPrivateKey(block.Bytes)
	}

	// RSA PRIVATE KEY specific handling
	if block.Type == BlockRSAPrivateKey {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	// Expect PRIVATE KEY if not EC or RSA at this point
	if block.Type != BlockPrivateKey {
		return nil, ErrDecodePrivateKey
	}

	// Try parsing private key using PKCS8, PKCS1, then EC
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Could not parse the private key
	return nil, ErrDecodePrivateKey
}

// EncodePublicKey as a PKIX ASN1.1 DER key and write a PEM block with type "PUBLIC KEY"
func EncodePublicKey(key interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := encodePublicKeyTo(key, &b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodePublicKeyTo(key interface{}, out io.Writer) (err error) {
	block := &pem.Block{Type: BlockPublicKey}
	if block.Bytes, err = x509.MarshalPKIXPublicKey(key); err != nil {
		return err
	}
	return pem.Encode(out, block)
}

// DecodePublicKey from a PEM encoded block. If the block type is "RSA PUBLIC KEY",
// then it is deocded as a PKCS 1, ASN.1 DER form. If the block is "PUBLIC KEY", then it
// is decoded from PKIX ASN1.1 DER form.
func DecodePublicKey(in []byte) (interface{}, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, ErrDecodePublicKey
	}
	return ParsePublicKey(block)
}

func ParsePublicKey(block *pem.Block) (interface{}, error) {
	if block.Type == BlockRSAPublicKey {
		return x509.ParsePKCS1PublicKey(block.Bytes)
	}

	if block.Type != BlockPublicKey {
		return nil, ErrDecodePublicKey
	}
	return x509.ParsePKIXPublicKey(block.Bytes)
}

// EncodeCertificate and write a PEM block with type "CERTIFICATE"
func EncodeCertificate(c *x509.Certificate) ([]byte, error) {
	var b bytes.Buffer
	if err := encodeCertificateTo(c, &b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodeCertificateTo(c *x509.Certificate, out io.Writer) (err error) {
	block := &pem.Block{Type: BlockCertificate, Bytes: c.Raw}
	return pem.Encode(out, block)
}

// DecodeCertificate from PEM encoded block with type "CERTIFICATE"
func DecodeCertificate(in []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(in)
	if block == nil || block.Type != BlockCertificate {
		return nil, ErrDecodeCertificate
	}
	return x509.ParseCertificate(block.Bytes)
}

// EncodeCSR and write a PEM block with type "CERTIFICATE REQUEST"
func EncodeCSR(c *x509.CertificateRequest) ([]byte, error) {
	var b bytes.Buffer
	if err := encodeCSRTo(c, &b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodeCSRTo(c *x509.CertificateRequest, out io.Writer) error {
	block := &pem.Block{Type: BlockCertificateRequest, Bytes: c.Raw}
	return pem.Encode(out, block)
}

// DecodeCSR from PEM encoded block with type "CERTIFICATE REQUEST"
func DecodeCSR(in []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(in)
	if block == nil || block.Type != BlockCertificateRequest {
		return nil, ErrDecodeCSR
	}
	return x509.ParseCertificateRequest(block.Bytes)
}
