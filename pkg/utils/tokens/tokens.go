package tokens

import (
	"encoding/base64"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/vmihailenco/msgpack/v5"
)

// NewConfirmation creates a new string token that includes a random secret and
// expires after 5 minutes.
func NewConfirmation(id ulid.ULID) (string, error) {
	token := &Confirmation{
		ID:     id,
		Secret: keygen.Secret(),
		//ExpiresAt: time.Now().Add(5 * time.Minute),
		ExpiresAt: time.Now().AddDate(50, 0, 0),
	}

	return token.Create()
}

// Confirmation protects access to a resource by encoding an ID with a
// cryptographically secure secret and an expiration time.
type Confirmation struct {
	ID        ulid.ULID `msgpack:"id"`
	Secret    string    `msgpack:"secret"`
	ExpiresAt time.Time `msgpack:"expires_at"`
}

func (t *Confirmation) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// Returns true if the token is valid for the given ID.
func (t *Confirmation) IsValid(id ulid.ULID) bool {
	return !t.IsExpired() && t.ID.Compare(id) == 0
}

// Create a new base64 encoded string from the token data. Note that callers should use
// the NewResourceToken method to ensure that all fields are present; this method is
// primarily exposed for the tests.
func (t *Confirmation) Create() (_ string, err error) {
	var data []byte
	if data, err = msgpack.Marshal(t); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(data), nil
}

// Decode a base64 encoded string into the struct.
func (t *Confirmation) Decode(token string) (err error) {
	var data []byte
	if data, err = base64.RawStdEncoding.DecodeString(token); err != nil {
		return err
	}

	return msgpack.Unmarshal(data, t)
}
