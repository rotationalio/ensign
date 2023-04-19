package tokens_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestResourceToken(t *testing.T) {
	id := ulids.New()

	// Should be able to create the token from an ID
	token, err := tokens.NewConfirmation(id)
	require.NoError(t, err, "could not create token")

	// Should be able to decode the token
	decoded := &tokens.Confirmation{}
	require.NoError(t, decoded.Decode(token), "could not decode token")
	require.Equal(t, id, decoded.ID, "decoded ID does not match original ID")
	require.False(t, decoded.IsExpired(), "token should not be expired")
	require.NotEmpty(t, decoded.Secret, "token secret should not be empty")

	// Tokens should have unique secrets
	other, err := tokens.NewConfirmation(id)
	require.NoError(t, err, "could not create token")
	require.NotEqual(t, token, other, "tokens should not be equal")

	otherDecoded := &tokens.Confirmation{}
	require.NoError(t, otherDecoded.Decode(other), "could not decode token")
	require.NotEqual(t, decoded.Secret, otherDecoded.Secret, "secrets should not be equal")
}
