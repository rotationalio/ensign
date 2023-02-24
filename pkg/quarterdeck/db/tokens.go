package db

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/keygen"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	nonceLength = 64
	keyLength   = 64
)

// NewVerificationToken creates a token struct from an email address that expires in 7
// days.
func NewVerificationToken(email string) *VerificationToken {
	return &VerificationToken{
		Email:     email,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		Nonce:     keygen.AlphaNumeric(nonceLength),
	}
}

// VerificationToken packages an email address with random data and an expiration time
// so that it can be serialized and hashed into a token which can be sent to users.
type VerificationToken struct {
	Email     string    `msgpack:"email"`
	ExpiresAt time.Time `msgpack:"expires_at"`
	Nonce     string    `msgpack:"nonce"`
}

func (t *VerificationToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

// Sign creates a base64 encoded string from the token data so that it can be sent to
// users as part of a URL. The returned secret should be stored in the database so that
// the string can be recomputed when verifying a user provided token.
func (t *VerificationToken) Sign() (_ string, secret []byte, err error) {
	var data []byte
	if data, err = msgpack.Marshal(t); err != nil {
		return "", nil, err
	}

	// Compute hash with a random 64 byte key
	key := []byte(keygen.AlphaNumeric(keyLength))
	mac := hmac.New(sha256.New, key)
	if _, err = mac.Write(data); err != nil {
		return "", nil, err
	}

	// Include the nonce with the key so that the token can be reconstructed later
	secret = append([]byte(t.Nonce), key...)

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), secret, nil
}

// Verify checks that a token was signed with the secret and is not expired.
func (t *VerificationToken) Verify(signature string, secret []byte) (err error) {
	if t.Email == "" {
		return errors.New("token is missing email address")
	}

	if t.IsExpired() {
		return errors.New("token is expired")
	}

	if len(secret) != nonceLength+keyLength {
		return errors.New("invalid secret for token verification")
	}

	// Serialize the struct with the nonce from the secret
	t.Nonce = string(secret[:nonceLength])
	var data []byte
	if data, err = msgpack.Marshal(t); err != nil {
		return err
	}

	// Compute hash to verify the user token
	mac := hmac.New(sha256.New, secret[nonceLength:])
	if _, err = mac.Write(data); err != nil {
		return err
	}

	// Decode the user token
	var token []byte
	if token, err = base64.RawURLEncoding.DecodeString(signature); err != nil {
		return err
	}

	// Check if the recomputed token matches the user token
	if !hmac.Equal(mac.Sum(nil), token) {
		return errors.New("invalid token signature")
	}

	return nil
}
