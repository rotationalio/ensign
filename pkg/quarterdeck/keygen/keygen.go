/*
Package keygen provides functionality for generating API client IDs and secrets.
*/
package keygen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idxbits  = 6
	idxmask  = 1<<idxbits - 1
	idxmax   = 63 / idxbits
)

func Alpha(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for i, cache, remain := n-1, CryptoRandInt(), idxmax; i >= 0; {
		if remain == 0 {
			cache, remain = CryptoRandInt(), idxmax
		}

		if idx := int(cache & idxmask); idx < len(alphabet) {
			sb.WriteByte(alphabet[idx])
			i--
		}

		cache >>= idxbits
		remain--
	}

	return sb.String()
}

func AlphaNumeric(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for i, cache, remain := n-1, CryptoRandInt(), idxmax; i >= 0; {
		if remain == 0 {
			cache, remain = CryptoRandInt(), idxmax
		}

		if idx := int(cache & idxmask); idx < len(alphanum) {
			sb.WriteByte(alphanum[idx])
			i--
		}

		cache >>= idxbits
		remain--
	}

	return sb.String()
}

func CryptoRandInt() uint64 {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("cannot generate random number: %w", err))
	}
	return binary.BigEndian.Uint64(buf)
}
