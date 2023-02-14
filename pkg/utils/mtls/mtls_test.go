package mtls_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/mtls"
	"github.com/rotationalio/ensign/pkg/utils/mtls/pem"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Test that Config returns a valid TLS config.
func TestConfig(t *testing.T) {
	err := checkFixtures()
	require.NoError(t, err, "could not create required fixtures")

	// Load a provider chain (astros) and trusted pool (banks) that does not have a private key
	var chain, trusted *mtls.Provider

	chain, err = mtls.Load("testdata/server.astros.com.pem")
	require.NoError(t, err)

	trusted, err = mtls.Load("testdata/banks.com.pool.pem")
	require.NoError(t, err)

	// Public provider should return an error
	_, err = mtls.Config(trusted)
	require.Error(t, err)

	// Valid config
	cfg, err := mtls.Config(chain, trusted)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Len(t, cfg.Certificates, 1)
	require.Equal(t, uint16(tls.VersionTLS13), cfg.MinVersion)
	require.NotEmpty(t, cfg.CurvePreferences)
	require.NotEmpty(t, cfg.CipherSuites)
	require.Equal(t, tls.RequireAndVerifyClientCert, cfg.ClientAuth)
	require.NotNil(t, cfg.ClientCAs)

	// Should be able to create a config without a trusted pool
	cfg2, err := mtls.Config(chain)
	require.NoError(t, err)
	require.Len(t, cfg2.Certificates, 1)
	require.NotNil(t, cfg2.ClientCAs)

	require.False(t, cfg.ClientCAs.Equal(cfg2.ClientCAs))

}

// Test that ServerCreds returns a grpc.ServerOption for mtls.
func TestServerCreds(t *testing.T) {
	err := checkFixtures()
	require.NoError(t, err, "could not create required fixtures")

	// Load a provider chain (astros) and trusted pool (banks) that does not have a private key
	var chain, trusted *mtls.Provider

	chain, err = mtls.Load("testdata/server.astros.com.pem")
	require.NoError(t, err)

	trusted, err = mtls.Load("testdata/banks.com.pool.pem")
	require.NoError(t, err)

	// Public provider should return an error
	_, err = mtls.ServerCreds(trusted)
	require.Error(t, err)

	// Succesfully retuning a grpc.ServerOption
	opt, err := mtls.ServerCreds(chain, trusted)
	require.NoError(t, err)
	require.Implements(t, (*grpc.ServerOption)(nil), opt)
}

// Test that ClientCreds returns a grpc.DialOption for mtls.
func TestClientCreds(t *testing.T) {
	err := checkFixtures()
	require.NoError(t, err, "could not create required fixtures")

	// Load a provider chain (astros) and trusted pool (banks) that does not have a private key
	var chain, trusted *mtls.Provider

	chain, err = mtls.Load("testdata/client.astros.com.pem")
	require.NoError(t, err)

	trusted, err = mtls.Load("testdata/banks.com.pool.pem")
	require.NoError(t, err)

	// Public provider should return an error
	_, err = mtls.ClientCreds("server.astros.com:4434", trusted)
	require.Error(t, err)

	// Successfully returning a grpc.DialOption
	opt, err := mtls.ClientCreds("server.astros.com:4434", chain, trusted)
	require.NoError(t, err)
	require.Implements(t, (*grpc.DialOption)(nil), opt)
}

// Helper function to check fixtures and if they don't exist to create them.
func checkFixtures() (err error) {
	requiredFixtures := []string{
		"astros.com.pool.pem", "banks.com.pool.pem",
		"server.astros.com.pem", "client.astros.com.pem",
		"server.banks.com.pem", "client.banks.com.pem",
	}

	// If all of the required fixtures exist, return nil
	if err = allFixturesExist(requiredFixtures...); err != nil {
		// Create fixtures since they do not exist
		if err = createGroupA(); err != nil {
			return err
		}
		if err = createGroupB(); err != nil {
			return err
		}
	}
	return nil
}

func allFixturesExist(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join("testdata", path)); errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func createGroupA() error {
	rsub := pkix.Name{
		CommonName:         "root.astros.com",
		Organization:       []string{"Astros, Inc."},
		OrganizationalUnit: []string{"Network Security"},
		StreetAddress:      []string{"Herzbachweg 2"},
		Locality:           []string{"Gelnhausen"},
		Province:           []string{"Hesse"},
		Country:            []string{"DE"},
	}

	isub := pkix.Name{
		CommonName:         "astros.com",
		Organization:       []string{"Astros, Inc."},
		OrganizationalUnit: []string{"Field Division"},
		StreetAddress:      []string{"Friedrich-von-Schletz-Strasse 30"},
		Locality:           []string{"Forchheim"},
		Province:           []string{"Bavaria"},
		Country:            []string{"DE"},
	}

	rootCA, intermediateCA, signKey, err := createCA(rsub, isub)
	if err != nil {
		return err
	}

	server := pkix.Name{
		CommonName:         "server.astros.com",
		Organization:       []string{"Astros, Inc."},
		OrganizationalUnit: []string{"Field Division"},
		StreetAddress:      []string{"Friedrich-von-Schletz-Strasse 30"},
		Locality:           []string{"Forchheim"},
		Province:           []string{"Bavaria"},
		Country:            []string{"DE"},
	}

	if err = createCerts(rootCA, intermediateCA, signKey, server); err != nil {
		return err
	}

	client := pkix.Name{
		CommonName:         "client.astros.com",
		Organization:       []string{"Astros, Inc."},
		OrganizationalUnit: []string{"Field Division"},
		StreetAddress:      []string{"Friedrich-von-Schletz-Strasse 30"},
		Locality:           []string{"Forchheim"},
		Province:           []string{"Bavaria"},
		Country:            []string{"DE"},
	}

	if err = createCerts(rootCA, intermediateCA, signKey, client); err != nil {
		return err
	}

	return nil
}

func createGroupB() error {
	rsub := pkix.Name{
		CommonName:         "root.banks.com",
		Organization:       []string{"Banking Better, PTE"},
		OrganizationalUnit: []string{"Engineering"},
		StreetAddress:      []string{"100 Norman Street"},
		Locality:           []string{"Waverley"},
		Province:           []string{"Dunedin"},
		PostalCode:         []string{"9013"},
		Country:            []string{"NZ"},
	}

	isub := pkix.Name{
		CommonName:         "banks.com",
		Organization:       []string{"Banking Better, PTE"},
		OrganizationalUnit: []string{"Software Engineering"},
		StreetAddress:      []string{"113 Hazel Terrace"},
		Locality:           []string{"Tauriko"},
		Province:           []string{"Tauranga"},
		PostalCode:         []string{"3110"},
		Country:            []string{"NZ"},
	}

	rootCA, intermediateCA, signKey, err := createCA(rsub, isub)
	if err != nil {
		return err
	}

	server := pkix.Name{
		CommonName:         "server.banks.com",
		Organization:       []string{"Banking Better, PTE"},
		OrganizationalUnit: []string{"Software Engineering"},
		StreetAddress:      []string{"113 Hazel Terrace"},
		Locality:           []string{"Tauriko"},
		Province:           []string{"Tauranga"},
		PostalCode:         []string{"3110"},
		Country:            []string{"NZ"},
	}

	if err = createCerts(rootCA, intermediateCA, signKey, server); err != nil {
		return err
	}

	client := pkix.Name{
		CommonName:         "client.banks.com",
		Organization:       []string{"Banking Better, PTE"},
		OrganizationalUnit: []string{"Software Engineering"},
		StreetAddress:      []string{"113 Hazel Terrace"},
		Locality:           []string{"Tauriko"},
		Province:           []string{"Tauranga"},
		PostalCode:         []string{"3110"},
		Country:            []string{"NZ"},
	}

	if err = createCerts(rootCA, intermediateCA, signKey, client); err != nil {
		return err
	}
	return nil
}

// Helper function to create mock certs and write them to testdata/common_name.pem
func createCerts(rootCA, intermediateCA *x509.Certificate, signKey interface{}, subject pkix.Name) (err error) {
	tmpl := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: ulid.Make().Bytes(),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{subject.CommonName},
	}

	// Set the serial number of the template to a random UUID.
	tmpl.SerialNumber.SetBytes(ulid.Make().Bytes())

	var priv *rsa.PrivateKey
	if priv, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return fmt.Errorf("could not generate rsa key: %w", err)
	}

	// Sign the certificate
	var signed []byte
	if signed, err = x509.CreateCertificate(rand.Reader, tmpl, intermediateCA, &priv.PublicKey, signKey); err != nil {
		return fmt.Errorf("could not sign certificate: %w", err)
	}

	var cert *x509.Certificate
	if cert, err = x509.ParseCertificate(signed); err != nil {
		return fmt.Errorf("could not parse signed certificate: %w", err)
	}

	var f *os.File
	if f, err = os.Create(filepath.Join("testdata", fmt.Sprintf("%s.pem", subject.CommonName))); err != nil {
		return fmt.Errorf("could not open fixture file: %w", err)
	}
	defer f.Close()

	writer := pem.NewWriter(f)
	defer writer.Close()

	// Write the certificate chain in order
	if err = writer.EncodeCertificate(cert); err != nil {
		return fmt.Errorf("could not encode leaf certificate: %w", err)
	}

	if err = writer.EncodeCertificate(intermediateCA); err != nil {
		return fmt.Errorf("could not encode intermediate ca: %w", err)
	}

	if err = writer.EncodeCertificate(rootCA); err != nil {
		return fmt.Errorf("could not encode root ca: %w", err)
	}

	if err = writer.EncodePrivateKey(priv); err != nil {
		return fmt.Errorf("could not encode private key: %w", err)
	}
	return nil
}

// Helper function to create a certificate signing chain for testing.
// The certificates of the pool are written to testdata.
func createCA(rootSubject, intermediateSubject pkix.Name) (rootCA, intermediateCA *x509.Certificate, signKey *rsa.PrivateKey, err error) {
	// Create self signed root certificate
	rootmpl := &x509.Certificate{
		SerialNumber:          &big.Int{},
		Subject:               rootSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		SubjectKeyId:          ulid.Make().Bytes(),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Set the serial number of the template to a random UUID.
	rootmpl.SerialNumber.SetBytes(ulid.Make().Bytes())

	var rootKey *rsa.PrivateKey
	if rootKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, nil, nil, fmt.Errorf("could not generate rsa root key: %w", err)
	}

	var data []byte
	if data, err = x509.CreateCertificate(rand.Reader, rootmpl, rootmpl, &rootKey.PublicKey, rootKey); err != nil {
		return nil, nil, nil, fmt.Errorf("could not create self-signed rootCA certificate: %w", err)
	}

	if rootCA, err = x509.ParseCertificate(data); err != nil {
		return nil, nil, nil, fmt.Errorf("could not parse self-signed rootCA certificate: %w", err)
	}

	intertmpl := &x509.Certificate{
		SerialNumber:          &big.Int{},
		Subject:               intermediateSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		SubjectKeyId:          ulid.Make().Bytes(),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}

	// Set the serial number of the template to a random UUID.
	intertmpl.SerialNumber.SetBytes(ulid.Make().Bytes())

	if signKey, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, nil, nil, fmt.Errorf("could not generate rsa signing key: %w", err)
	}

	if data, err = x509.CreateCertificate(rand.Reader, intertmpl, rootCA, &signKey.PublicKey, rootKey); err != nil {
		return nil, nil, nil, fmt.Errorf("could not create signed intermediateCA certificate: %w", err)
	}

	if intermediateCA, err = x509.ParseCertificate(data); err != nil {
		return nil, nil, nil, fmt.Errorf("could not parse signed intermediateCA certificate: %w", err)
	}

	var f *os.File
	if f, err = os.Create(filepath.Join("testdata", fmt.Sprintf("%s.pool.pem", intermediateSubject.CommonName))); err != nil {
		return nil, nil, nil, fmt.Errorf("could not open fixture file: %w", err)
	}
	defer f.Close()

	writer := pem.NewWriter(f)
	defer writer.Close()

	// Write the certificate chain in order
	if err = writer.EncodeCertificate(intermediateCA); err != nil {
		return nil, nil, nil, fmt.Errorf("could not encode intermediate ca: %w", err)
	}

	if err = writer.EncodeCertificate(rootCA); err != nil {
		return nil, nil, nil, fmt.Errorf("could not encode leaf certificate: %w", err)
	}

	return rootCA, intermediateCA, signKey, nil
}
