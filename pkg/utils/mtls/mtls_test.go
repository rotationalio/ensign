package mtls_test

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Test that Config returns a valid TLS config.
func TestConfig(t *testing.T) {

	// Public provider should return an error
	public := private.Public()
	_, err = mtls.Config(public, mtls.NewPool())
	require.Error(t, err)

	// Valid config
	cfg, err := mtls.Config(private, mtls.NewPool())
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Len(t, cfg.Certificates, 1)
	require.Equal(t, uint16(tls.VersionTLS12), cfg.MinVersion)
	require.NotEmpty(t, cfg.CurvePreferences)
	require.NotEmpty(t, cfg.CipherSuites)
	require.Equal(t, tls.RequireAndVerifyClientCert, cfg.ClientAuth)
	require.NotNil(t, cfg.ClientCAs)
}

// Test that ServerCreds returns a grpc.ServerOption for mtls.
func TestServerCreds(t *testing.T) {
	// Public provider should return an error
	public := private.Public()
	_, err = mtls.ServerCreds(public, mtls.NewPool())
	require.Error(t, err)

	// Succesfully retuning a grpc.ServerOption
	opt, err := mtls.ServerCreds(private, mtls.NewPool())
	require.NoError(t, err)
	require.Implements(t, (*grpc.ServerOption)(nil), opt)
}

// Test that ClientCreds returns a grpc.DialOption for mtls.
func TestClientCreds(t *testing.T) {
	// Public provider should return an error
	public := private.Public()
	_, err = mtls.ClientCreds("https://localhost:12345", public, mtls.NewPool())
	require.Error(t, err)

	// Successfully returning a grpc.DialOption
	opt, err := mtls.ClientCreds("https://localhost:12345", private, mtls.NewPool())
	require.NoError(t, err)
	require.Implements(t, (*grpc.DialOption)(nil), opt)
}
