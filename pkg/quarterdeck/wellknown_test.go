package quarterdeck_test

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (s *quarterdeckTestSuite) TestJWKS() {
	// Fetch the JWK resource by the specified URL (in the same manner that clients will)
	require := s.Require()
	keys, err := jwk.Fetch(context.Background(), s.srv.URL()+"/.well-known/jwks.json")
	require.NoError(err, "could not fetch jwks key set")
	require.Equal(2, keys.Len(), "unexpected number of keys returned")

	expected := []string{"01GE6191AQTGMCJ9BN0QC3CCVG", "01GE62EXXR0X0561XD53RDFBQJ"}
	for _, kid := range expected {
		key, ok := keys.LookupKeyID(kid)
		require.True(ok, "could not find key with id %q", kid)
		require.Equal(jwk.ForSignature.String(), key.KeyUsage())
		require.Equal(kid, key.KeyID())
		require.Equal(jwa.RSA, key.KeyType())
		require.Equal(jwa.RS256, key.Algorithm())

		var pub rsa.PublicKey
		err = key.Raw(&pub)
		require.NoError(err, "could not parse raw public key")
		require.NotNil(pub, "could not extract public key from jwks")
	}

}

func (s *quarterdeckTestSuite) TestOpenIDConfiguration() {
	require := s.Require()

	// Create a basic HTTP request rather than use the Quarterdeck client to ensure the
	// headers and data returned are the expected values.
	req, err := http.NewRequest(http.MethodGet, s.srv.URL()+"/.well-known/openid-configuration", nil)
	require.NoError(err, "could not create basic http request")

	rep, err := http.DefaultClient.Do(req)
	require.NoError(err, "could not execute basic http request")
	defer rep.Body.Close()

	// The content-type must be application/json for the configuration
	require.Equal("application/json; charset=utf-8", rep.Header.Get("Content-Type"))

	// Parse the response
	openid := &api.OpenIDConfiguration{}
	err = json.NewDecoder(rep.Body).Decode(openid)
	require.NoError(err, "could not decode JSON response from server")

	require.Equal("http://quarterdeck.test/", openid.Issuer)
	require.Equal("http://quarterdeck.test/.well-known/jwks.json", openid.JWKSURI)
}

func (s *quarterdeckTestSuite) TestSecurityTxt() {
	require := s.Require()

	// Create a basic HTTP request rather than use the Quarterdeck client to ensure the
	// headers and data returned are the expected values.
	req, err := http.NewRequest(http.MethodGet, s.srv.URL()+"/.well-known/security.txt", nil)
	require.NoError(err, "could not create basic http request")

	rep, err := http.DefaultClient.Do(req)
	require.NoError(err, "could not execute basic http request")
	defer rep.Body.Close()

	// The content-type must be text/plain for the security.text file
	require.Equal("text/plain; charset=utf-8", rep.Header.Get("Content-Type"))

	// Ensure the content returned is as expected
	expected, err := os.ReadFile("testdata/security.txt")
	require.NoError(err, "could not read testdata/security.txt fixture")

	actual, err := io.ReadAll(rep.Body)
	require.NoError(err, "could not read the body response from the server")
	require.Equal(expected, actual, "the security.txt file does not match the testdata fixture")

	// TODO: check that the GPG signature is valid
	// This can be done on the command line with: gpg --auto-key-retrieve --verify --output -
	// Not sure how to do this in golang tests, however.

	// Error if the security.txt has expired - if this test fails, regenerate the security.txt file!
	lines := strings.Split(string(actual), "\n")
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "Expires:") {
			found = true
			parts := strings.Split(line, " ")
			require.Len(parts, 2, "could not split expires directive")

			expires, err := time.Parse(time.RFC3339, parts[1])
			require.NoError(err, "could not parse expires timestamp")

			require.True(time.Now().Before(expires), "the security.txt file has expired, regenerate it!")
		}
	}
	require.True(found, "could not find expires line")
}
