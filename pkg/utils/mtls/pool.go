package mtls

import "crypto/x509"

// CertPool returns an x509.CertPool, a collection of certificates with public keys
// that are used to verify clients or servers when connecting over mTLS. The CertPool
// is constructed from the providers passed in (whether they are private or not) and
// the collected CertPool is returned. This pool is provider-specific and does not
// include system certificates.
func CertPool(providers ...*Provider) (pool *x509.CertPool, err error) {
	pool = x509.NewCertPool()
	for _, provider := range providers {
		for _, asn1der := range provider.chain.Certificate {
			var cert *x509.Certificate
			if cert, err = x509.ParseCertificate(asn1der); err != nil {
				return nil, err
			}
			pool.AddCert(cert)
		}
	}
	return pool, nil
}
