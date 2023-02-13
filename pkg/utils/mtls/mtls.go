package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config returns a standardized mTLS configuration for all clients and servers that
// are using this package. The certificate is provided from the private provider and
// any additional providers can be added to extend the cert pool for verification.
func Config(chain *Provider, trusted ...*Provider) (_ *tls.Config, err error) {
	if !chain.IsPrivate() {
		return nil, ErrPrivateKeyRequired
	}

	var cert tls.Certificate
	if cert, err = chain.GetKeyPair(); err != nil {
		return nil, err
	}

	var pool *x509.CertPool
	if pool, err = chain.GetCertPool(); err != nil {
		return nil, err
	}

	// TODO: handle additional trust pool.

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}, nil
}

// ServerCreds returns a grpc.ServerOption to create a gRPC server with mTLS.
func ServerCreds(chain *Provider, trusted ...*Provider) (_ grpc.ServerOption, err error) {
	var conf *tls.Config
	if conf, err = Config(chain, trusted...); err != nil {
		return nil, err
	}
	return grpc.Creds(credentials.NewTLS(conf)), nil
}

// ClientCreds returns a grpc.DialOption to connect a gRPC client with mTLS.
func ClientCreds(endpoint string, chain *Provider, trusted ...*Provider) (_ grpc.DialOption, err error) {
	var conf *tls.Config
	if conf, err = Config(chain, trusted...); err != nil {
		return nil, err
	}

	var u *url.URL
	if u, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}

	// Update configuration for client instead of server
	conf.ServerName = u.Host
	conf.RootCAs = conf.ClientCAs
	conf.ClientCAs = nil

	return grpc.WithTransportCredentials(credentials.NewTLS(conf)), nil
}
