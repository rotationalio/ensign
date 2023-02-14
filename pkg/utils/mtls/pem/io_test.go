package pem_test

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/mtls/pem"
	"github.com/stretchr/testify/require"
)

func TestWriterReader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "certs.pem")

	f, err := os.Create(path)
	require.NoError(t, err, "could not open tmp file")
	defer f.Close()

	writer := pem.NewWriter(f)

	// Attempt to write a private key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	err = writer.EncodePrivateKey(priv)
	require.NoError(t, err)

	// Attempt to write a public key
	pub := &priv.PublicKey
	err = writer.EncodePublicKey(pub)
	require.NoError(t, err)

	// Attempt to write a Certificate
	cert, err := cert()
	require.NoError(t, err)

	err = writer.EncodeCertificate(cert)
	require.NoError(t, err)

	// Attempt to write a CSR
	csr, err := csr()
	require.NoError(t, err)

	err = writer.EncodeCSR(csr)
	require.NoError(t, err)

	// Close the writer
	require.NoError(t, writer.Close())

	// Create a reader
	fr, err := os.Open(path)
	require.NoError(t, err)
	defer fr.Close()

	reader, err := pem.NewReader(fr)
	require.NoError(t, err)

	// Should be able to decode blocks in the order they were written
	b := 0
	for reader.Next() {
		block := reader.Decode()
		require.NotNil(t, block)

		switch b {
		case 0:
			key, err := block.DecodePrivateKey()
			require.NoError(t, err)
			require.Equal(t, priv, key)
		case 1:
			key, err := block.DecodePublicKey()
			require.NoError(t, err)
			require.Equal(t, pub, key)
		case 2:
			c, err := block.DecodeCertificate()
			require.NoError(t, err)
			require.Equal(t, cert, c)
		case 3:
			c, err := block.DecodeCSR()
			require.NoError(t, err)
			require.Equal(t, csr, c)
		default:
			panic("???")
		}

		b++
	}

	// Close the reader
	require.NoError(t, reader.Close())
}
