package keygen_test

import (
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/stretchr/testify/require"
)

func TestAlpha(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	// Test creating different random strings at different lengths
	for i := 0; i < 10000; i++ {
		len := rand.Intn(512) + 1
		alpha := keygen.Alpha(len)
		require.Len(t, alpha, len)
		require.Regexp(t, regexp.MustCompile(`[a-zA-Z]+`), alpha)
	}

	vals := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		val := keygen.Alpha(16)
		vals[val] = struct{}{}
	}
	require.Len(t, vals, 10000, "there is a very low chance that a duplicate value was generated")
}

func TestAlphaNumeric(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	// Test creating different random strings at different lengths
	for i := 0; i < 10000; i++ {
		len := rand.Intn(512) + 1
		alpha := keygen.AlphaNumeric(len)
		require.Len(t, alpha, len)
		require.Regexp(t, regexp.MustCompile(`[a-zA-Z0-9]+`), alpha)
	}

	vals := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		val := keygen.AlphaNumeric(16)
		vals[val] = struct{}{}
	}
	require.Len(t, vals, 10000, "there is a very low chance that a duplicate value was generated")
}

func TestKeyGen(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	keyID := keygen.KeyID()
	require.Len(t, keyID, keygen.KeyIDLength)

	secret := keygen.Secret()
	require.Len(t, secret, keygen.SecretLength)
}

func TestCryptoRandInt(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	nums := make(map[uint64]struct{})
	for i := 0; i < 10000; i++ {
		val := keygen.CryptoRandInt()
		nums[val] = struct{}{}
	}
	require.Len(t, nums, 10000, "there is a very low chance that a duplicate value was generated")
}

func benchmarkAlpha(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		keygen.Alpha(i)
	}
}

func benchmarkAlphaNumeric(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		keygen.AlphaNumeric(i)
	}
}

func BenchmarkAlpha16(b *testing.B)  { benchmarkAlpha(16, b) }
func BenchmarkAlpha64(b *testing.B)  { benchmarkAlpha(64, b) }
func BenchmarkAlpha256(b *testing.B) { benchmarkAlpha(256, b) }

func BenchmarkAlphaNumeric16(b *testing.B)  { benchmarkAlphaNumeric(16, b) }
func BenchmarkAlphaNumeric64(b *testing.B)  { benchmarkAlphaNumeric(64, b) }
func BenchmarkAlphaNumeric256(b *testing.B) { benchmarkAlphaNumeric(256, b) }

func BenchmarkCryptoRandInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		keygen.CryptoRandInt()
	}
}

func BenchmarkRandInt(b *testing.B) {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for n := 0; n < b.N; n++ {
		source.Int63()
	}
}
