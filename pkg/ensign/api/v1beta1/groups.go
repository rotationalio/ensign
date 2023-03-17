package api

import (
	"errors"

	"github.com/twmb/murmur3"
)

var ErrNoGroupID = errors.New("consumer group requires either id or name")

func (c *ConsumerGroup) Key() ([16]byte, error) {
	key := [16]byte{}

	// If the ID is already 16 bytes, use it directly without hashing.
	if len(c.Id) == 16 {
		copy(key[:], c.Id)
		return key, nil
	}

	// Hash the ID or the name to get the key
	hash := murmur3.New128()
	switch {
	case len(c.Id) > 0:
		hash.Write(c.Id)
	case c.Name != "":
		hash.Write([]byte(c.Name))
	default:
		return key, ErrNoGroupID
	}

	copy(key[:], hash.Sum(nil))
	return key, nil
}
