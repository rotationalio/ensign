package db

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	nonceLength = 64
	keyLength   = 64
)

// NewVerificationToken creates a token struct from an email address that expires in 7
// days.
func NewVerificationToken(email string) (token *VerificationToken, err error) {
	if email == "" {
		return nil, errors.New("email address is required")
	}

	token = &VerificationToken{
		Email:     email,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		Nonce:     make([]byte, nonceLength),
	}

	if _, err = rand.Read(token.Nonce); err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}
	return token, nil
}

// VerificationToken packages an email address with random data and an expiration time
// so that it can be serialized and hashed into a token which can be sent to users.
type VerificationToken struct {
	Email     string    `msgpack:"email"`
	ExpiresAt time.Time `msgpack:"expires_at"`
	Nonce     []byte    `msgpack:"nonce"`
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
	key := make([]byte, keyLength)
	if _, err = rand.Read(key); err != nil {
		return "", nil, err
	}

	mac := hmac.New(sha256.New, key)
	if _, err = mac.Write(data); err != nil {
		return "", nil, err
	}

	// Include the nonce with the key so that the token can be reconstructed later
	secret = make([]byte, nonceLength+keyLength)
	copy(secret[0:nonceLength], t.Nonce)
	copy(secret[nonceLength:], key)

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), secret, nil
}

// Verify checks that a token was signed with the secret and is not expired.
func (t *VerificationToken) Verify(signature string, secret []byte) (err error) {
	if t.Email == "" {
		return ErrTokenMissingEmail
	}

	if t.IsExpired() {
		return ErrTokenExpired
	}

	if len(secret) != nonceLength+keyLength {
		return ErrInvalidSecret
	}

	// Serialize the struct with the nonce from the secret
	t.Nonce = secret[0:nonceLength]
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
		return ErrTokenInvalid
	}

	return nil
}
