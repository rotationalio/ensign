// Package tlstest provides a TLS server that can be used for testing in lieu of
// httptest since we need to get access to a TLSConfig before starting the server.
package tlstest

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func Config() *tls.Config {
	cert, err := tls.X509KeyPair(LocalhostCert, LocalhostKey)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{cert},
	}
}

func Client() *http.Client {
	cert, err := tls.X509KeyPair(LocalhostCert, LocalhostKey)
	if err != nil {
		panic(err)
	}

	certificate, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	certpool := x509.NewCertPool()
	certpool.AddCert(certificate)

	client := &http.Client{
		CheckRedirect: nil,
		Timeout:       30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certpool,
			},
			ForceAttemptHTTP2: false,
		},
	}

	if client.Jar, err = cookiejar.New(nil); err != nil {
		panic(err)
	}
	return client
}
