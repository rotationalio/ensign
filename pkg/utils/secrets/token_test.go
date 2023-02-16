package secrets_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/secrets"
	"github.com/stretchr/testify/require"
)

// Test that the CreateToken function creates a valid token of the given length.
func TestCreateToken(t *testing.T) {
	// Negative or zero length should return an empty string
	require.Equal(t, "", secrets.CreateToken(-1))
	require.Equal(t, "", secrets.CreateToken(0))

	// Returned token should not contain any unexpected characters
	token := secrets.CreateToken(20)
	require.True(t, secrets.ValidateToken(token), "token %s contains invalid characters", token)

	// Successive calls should return different tokens
	nextToken := secrets.CreateToken(20)
	require.True(t, secrets.ValidateToken(nextToken), "token %s contains invalid characters", nextToken)
	require.NotEqual(t, token, nextToken, "CreateToken returned the same token twice: %s", token)
}

// Test that the ValidateToken function can determine if a token contains invalid
// characters.
func TestValidateToken(t *testing.T) {
	// Valid tokens
	require.True(t, secrets.ValidateToken(""))
	require.True(t, secrets.ValidateToken("abcABC12345"))
	require.True(t, secrets.ValidateToken("abc123#+-{}"))
	require.True(t, secrets.ValidateToken("()&*[]#"))
	require.False(t, secrets.ValidateToken("abc12345!"))

	// Invalid tokens
	require.False(t, secrets.ValidateToken("abc12345\n"))
	require.False(t, secrets.ValidateToken("`~\t"))
	require.False(t, secrets.ValidateToken(",."))
}
