/*
Package mtls has helper functionality for organizing and maintaining key material on
disk and establishing connections between servers and clients using mTLS cryptography.
*/
package mtls

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/rotationalio/ensign/pkg/utils/mtls/pem"
)

// Provider wraps a PEM-encoded certificate chain, which can optionally include private
// keys and an additional pool of valid certificates. Providers with keys (referred to
// as private providers) are used to instantiate mTLS servers and to make secure
// connections from clients to servers.
type Provider struct {
	chain tls.Certificate
	key   interface{}
}

// New creates a provider from PEM encoded data.
func New(chain []byte) (provider *Provider, err error) {
	var reader *pem.Reader
	if reader, err = pem.NewReader(bytes.NewBuffer(chain)); err != nil {
		return nil, err
	}

	provider = &Provider{}
	if err = provider.Decode(reader); err != nil {
		return nil, err
	}
	return provider, nil
}

// Load opens the certificate chain at the specified path and reads the data.
// TODO: handle compressed formats.
func Load(path string) (provider *Provider, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	var reader *pem.Reader
	if reader, err = pem.NewReader(f); err != nil {
		return nil, err
	}

	provider = &Provider{}
	if err = provider.Decode(reader); err != nil {
		return nil, err
	}
	return provider, nil
}

// Decode PEM blocks and adds them to the provider. Certificates are appended to the
// Provider chain and Private Keys are decoded from PKCS8. All other block types
// return an error and stop processing the block or chain. Only the private key is
// verified for correctness, certificates are unverified.
func (p *Provider) Decode(reader *pem.Reader) (err error) {
	for reader.Next() {
		block := reader.Decode()
		switch block.Type {
		case pem.BlockCertificate:
			p.chain.Certificate = append(p.chain.Certificate, block.Bytes)
		case pem.BlockPrivateKey, pem.BlockECPrivateKey, pem.BlockRSAPrivateKey:
			if p.key, err = block.DecodePrivateKey(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled block type %s", block.Type)
		}

	}
	return nil
}

// Encode Provider in PCKS12 PEM format for serialization. Certificates are written to
// the array first. If the private key exists, it is written as the last PEM block.
func (p *Provider) Encode(writer *pem.Writer) (err error) {
	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if err = writer.EncodeCertificate(crt); err != nil {
			return fmt.Errorf("could not encode certificate %d: %s", i, err)
		}
	}

	if p.key != nil {
		if err = writer.EncodePrivateKey(p.key); err != nil {
			return fmt.Errorf("could not encode private key: %s", err)
		}
	}
	return nil
}

// Dump the provider to the specified path for serialization.
func (p *Provider) Dump(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	writer := pem.NewWriter(f)
	if err = p.Encode(writer); err != nil {
		return err
	}
	return nil
}

// GetKeyPair returns a tls.Certificate parsed from the PEM encoded data maintained by
// the provider. This method uses tls.X509KeyPair to ensure that the public/private key
// pair are suitable for use with an HTTP Server.
func (p *Provider) GetKeyPair() (_ tls.Certificate, err error) {
	if p.key == nil {
		return tls.Certificate{}, ErrMissingKey
	}

	var certs bytes.Buffer
	writer := pem.NewWriter(&certs)

	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if err = writer.EncodeCertificate(crt); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not encode certificate %d: %s", i, err)
		}
	}

	var key []byte
	if key, err = pem.EncodePrivateKey(p.key); err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certs.Bytes(), key)
}

// GetLeafCertificate returns the parsed x509 leaf certificate if it exists, returning
// an error if there are no certificates or if there is a parse error.
func (p *Provider) GetLeafCertificate() (*x509.Certificate, error) {
	if p.chain.Leaf != nil {
		return p.chain.Leaf, nil
	}

	if len(p.chain.Certificate) == 0 {
		return nil, ErrNoCertificates
	}
	return x509.ParseCertificate(p.chain.Certificate[0])
}

// GetKey returns the private key, or nil if this is a public provider.
func (p *Provider) GetKey() interface{} {
	return p.key
}

// GetRSAKeys returns a fully constructed RSA PrivateKey that includes the public key
// material property. This method errors if the key is not an RSA key or does not exist.
func (p *Provider) GetRSAKey() (key *rsa.PrivateKey, err error) {
	if p.key == nil {
		return nil, ErrMissingKey
	}

	var ok bool
	if key, ok = p.key.(*rsa.PrivateKey); !ok {
		return nil, fmt.Errorf("private key is not RSA but is %T", p.key)
	}
	return key, nil
}

// IsPrivate returns true if the Provider contains a non-nil key.
func (p *Provider) IsPrivate() bool {
	return p.key != nil
}

// String returns the common name of the Provider from the leaf certificate.
func (p *Provider) String() string {
	cert, err := p.GetLeafCertificate()
	if err != nil {
		return ""
	}
	return cert.Subject.CommonName
}
