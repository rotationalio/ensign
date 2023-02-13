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
	"encoding/pem"
	"fmt"
)

// Provider wraps a PEM-encoded certificate chain, which can optionally include private
// keys and an additional pool of valid certificates. Providers with keys (referred to
// as private providers) are used to instantiate mTLS servers and to make secure
// connections from clients to servers.
type Provider struct {
	chain tls.Certificate
	pool  *x509.CertPool
	key   interface{}
}

// Decode PEM blocks and adds them to the provider. Certificates are appended to the
// Provider chain and Private Keys are Unmarshalled from PKCS8. All other block types
// return an error and stop processing the block or chain. Only the private key is
// verified for correctness, certificates are unvalidated.
func (p *Provider) Decode(in []byte) (err error) {
	var block *pem.Block
	for {
		block, in = pem.Decode(in)
		if block == nil {
			break
		}

		switch block.Type {
		case BlockCertificate:
			p.chain.Certificate = append(p.chain.Certificate, block.Bytes)
		case BlockPrivateKey, BlockECPrivateKey, BlockRSAPrivateKey:
			if p.key, err = ParsePrivateKey(block); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled block type %q", block.Type)
		}
	}
	return nil
}

// Encode Provider in PCKS12 PEM format for serialization. Certificates are written to
// the array first. If the private key exists, it is written as the last PEM block.
func (p *Provider) Encode() (_ []byte, err error) {
	var b bytes.Buffer
	var block []byte

	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return nil, fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if block, err = PEMEncodeCertificate(crt); err != nil {
			return nil, fmt.Errorf("could not encode certificate %d: %s", i, err)
		}

		b.Write(block)
	}

	if p.key != nil {
		if block, err = PEMEncodePrivateKey(p.key); err != nil {
			return nil, fmt.Errorf("could not encode private key: %s", err)
		}

		b.Write(block)
	}

	return b.Bytes(), nil
}

// GetCertPool returns the x509.CertPool certificate set representing the root,
// intermediate, and leaf certificates of the Provider. This pool is provider-specific
// and does not include system certificates.
func (p *Provider) GetCertPool() (_ *x509.CertPool, err error) {
	if p.pool == nil {
		p.pool = x509.NewCertPool()
		for _, c := range p.chain.Certificate {
			var x509Cert *x509.Certificate
			if x509Cert, err = x509.ParseCertificate(c); err != nil {
				return nil, err
			}
			p.pool.AddCert(x509Cert)
		}
	}
	return p.pool, nil
}

// GetKeyPair returns a tls.Certificate parsed from the PEM encoded data maintained by
// the provider. This method uses tls.X509KeyPair to ensure that the public/private key
// pair are suitable for use with an HTTP Server.
func (p *Provider) GetKeyPair() (_ tls.Certificate, err error) {
	if p.key == nil {
		return tls.Certificate{}, ErrMissingKey
	}

	var block []byte
	var certs bytes.Buffer
	for i, asn1Data := range p.chain.Certificate {
		var crt *x509.Certificate
		if crt, err = x509.ParseCertificate(asn1Data); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not parse certificate %d: %s", i, err)
		}

		if block, err = PEMEncodeCertificate(crt); err != nil {
			return tls.Certificate{}, fmt.Errorf("could not encode certificate %d: %s", i, err)
		}

		certs.Write(block)
	}

	var key []byte
	if key, err = PEMEncodePrivateKey(p.key); err != nil {
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
