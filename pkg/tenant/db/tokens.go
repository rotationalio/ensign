package db

import (
	"encoding/base64"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/vmihailenco/msgpack/v5"
)

// NewResourceToken creates a new string token that includes a random secret and
// expires after 5 minutes.
func NewResourceToken(id ulid.ULID) (string, error) {
	token := &ResourceToken{
		ID:        id,
		Secret:    keygen.Secret(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return token.create()
}

// ResourceToken protects access to a resource by encoding an ID with a
// cryptographically secure secret and an expiration time.
type ResourceToken struct {
	ID        ulid.ULID `msgpack:"id"`
	Secret    string    `msgpack:"token"`
	ExpiresAt time.Time `msgpack:"expires_at"`
}

func (t *ResourceToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t *ResourceToken) create() (_ string, err error) {
	var data []byte
	if data, err = msgpack.Marshal(t); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(data), nil
}

// Decode a base64 encoded string into the struct.
func (t *ResourceToken) Decode(token string) (err error) {
	var data []byte
	if data, err = base64.RawStdEncoding.DecodeString(token); err != nil {
		return err
	}

	return msgpack.Unmarshal(data, t)
}
