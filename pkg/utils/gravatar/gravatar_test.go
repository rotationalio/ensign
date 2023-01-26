package gravatar_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/gravatar"
	"github.com/stretchr/testify/require"
)

func TestGravatar(t *testing.T) {
	email := "MyEmailAddress@example.com "
	url := gravatar.New(email, nil)
	require.Equal(t, "https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346?d=identicon&r=pg&s=80", url)
}

func TestHash(t *testing.T) {
	// Test case from: https://en.gravatar.com/site/implement/hash/
	input := "MyEmailAddress@example.com "
	expected := "0bc83cb571cd1c50ba6f3e8a78ef1346"
	require.Equal(t, expected, gravatar.Hash(input))
}
