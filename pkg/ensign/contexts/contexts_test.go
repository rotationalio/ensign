package contexts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/stretchr/testify/require"
)

func TestClaimsContext(t *testing.T) {
	claims := &tokens.Claims{
		Name:  "Barbara Testly",
		Email: "btest@testing.io",
	}

	parent, cancel := context.WithCancel(context.Background())
	ctx := contexts.WithClaims(parent, claims)

	cmpt, ok := contexts.ClaimsFrom(ctx)
	require.True(t, ok)
	require.Same(t, claims, cmpt)

	cancel()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}

func TestKeyString(t *testing.T) {
	testCases := []struct {
		key      fmt.Stringer
		expected string
	}{
		{contexts.KeyUnknown, "unknown"},
		{contexts.KeyClaims, "claims"},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.key.String())
	}
}
