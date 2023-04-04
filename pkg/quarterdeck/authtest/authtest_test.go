package authtest_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/stretchr/testify/require"
)

// This test generates an example token with fake RSA keys for use in examples,
// documentation and other tests that don't need a valid token (since it will expire).
func TestGenerateToken(t *testing.T) {
	t.Skip("comment the skip out if you want to generate a token")

	srv, err := authtest.NewServer()
	require.NoError(t, err, "could not start authtest server")
	defer srv.Close()

	claims := &tokens.Claims{
		Name:        "John Doe",
		Email:       "jdoe@example.com",
		OrgID:       "123",
		ProjectID:   "abc",
		Permissions: []string{"read:data", "write:data"},
	}

	accessToken, refreshToken, err := srv.CreateTokenPair(claims)
	require.NoError(t, err, "could not generate access token")

	// Log the tokens then fail the test so the tokens are printed out.
	t.Logf("access token: %s", accessToken)
	t.Logf("refresh token: %s", refreshToken)
	t.FailNow()
}
