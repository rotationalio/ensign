package db_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestVerificationToken(t *testing.T) {
	// Test that the verification token is created correctly
	token, err := db.NewVerificationToken("leopold.wentzel@gmail.com")
	require.NoError(t, err, "could not create verification token")
	require.Equal(t, "leopold.wentzel@gmail.com", token.Email)
	require.True(t, token.ExpiresAt.After(time.Now()))
	require.Len(t, token.Nonce, 64)

	// Test signing a token
	signature, secret, err := token.Sign()
	require.NoError(t, err, "failed to sign token")
	require.NotEmpty(t, signature)
	require.Len(t, secret, 128)
	require.True(t, bytes.HasPrefix(secret, token.Nonce))

	// Signing again should produce a different signature
	differentSig, differentSecret, err := token.Sign()
	require.NoError(t, err, "failed to sign token")
	require.NotEqual(t, signature, differentSig, "expected different signatures")
	require.NotEqual(t, secret, differentSecret, "expected different secrets")

	// Verification should fail if the token is missing an email address
	verify := &db.VerificationToken{
		SigningInfo: db.SigningInfo{
			ExpiresAt: time.Now().AddDate(0, 0, 7),
		},
	}
	require.ErrorIs(t, verify.Verify(signature, secret), db.ErrTokenMissingEmail, "expected error when token is missing email address")

	// Verification should fail if the token is expired
	verify.Email = "leopold.wentzel@gmail.com"
	verify.ExpiresAt = time.Now().AddDate(0, 0, -1)
	require.ErrorIs(t, verify.Verify(signature, secret), db.ErrTokenExpired, "expected error when token is expired")

	// Verification should fail if the email is different
	verify.Email = "wes.anderson@gmail.com"
	verify.ExpiresAt = token.ExpiresAt
	require.ErrorIs(t, verify.Verify(signature, secret), db.ErrTokenInvalid, "expected error when email is different")

	// Verification should fail if the signature is not decodable
	verify.Email = "leopold.wentzel@gmail.com"
	require.Error(t, verify.Verify("^&**(", secret), "expected error when signature is not decodable")

	// Verification should fail if the signature was created with a different secret
	require.ErrorIs(t, verify.Verify(differentSig, secret), db.ErrTokenInvalid, "expected error when signature was created with a different secret")

	// Should error if the secret has the wrong length
	require.ErrorIs(t, verify.Verify(signature, nil), db.ErrInvalidSecret, "expected error when secret is nil")
	require.ErrorIs(t, verify.Verify(signature, []byte("wronglength")), db.ErrInvalidSecret, "expected error when secret is the wrong length")

	// Verification should fail if the wrong secret is used
	require.ErrorIs(t, verify.Verify(signature, differentSecret), db.ErrTokenInvalid, "expected error when wrong secret is used")

	// Successful verification
	require.NoError(t, verify.Verify(signature, secret), "expected successful verification")
}

func TestResetToken(t *testing.T) {
	t.Run("Valid Reset Token", func(t *testing.T) {
		// Test that the reset token is created correctly
		id := ulids.New()
		token, err := db.NewResetToken(id)
		require.NoError(t, err, "could not create reset token")

		// Test signing a token
		signature, secret, err := token.Sign()
		require.NoError(t, err, "failed to sign token")

		// Signing again should produce a different signature
		differentSig, differentSecret, err := token.Sign()
		require.NoError(t, err, "failed to sign token")
		require.NotEqual(t, signature, differentSig, "expected different signatures")
		require.NotEqual(t, secret, differentSecret, "expected different secrets")

		// Should be able to verify the token
		require.NoError(t, token.Verify(signature, secret), "expected successful verification")
	})

	t.Run("Missing ID", func(t *testing.T) {
		// Should fail to create a token without an ID
		_, err := db.NewResetToken(ulids.Null)
		require.ErrorIs(t, err, db.ErrMissingUserID, "expected error when token is missing ID")
	})

	t.Run("Token Missing User ID", func(t *testing.T) {
		// Token with missing user ID should be an error
		token := &db.ResetToken{}
		require.ErrorIs(t, token.Verify("", nil), db.ErrTokenMissingUserID, "expected error when token is missing ID")
	})

	t.Run("Token Expired", func(t *testing.T) {
		// Token that is expired should be an error
		token := &db.ResetToken{
			SigningInfo: db.SigningInfo{
				ExpiresAt: time.Now().AddDate(0, 0, -1),
			},
			UserID: ulids.New(),
		}
		require.ErrorIs(t, token.Verify("", nil), db.ErrTokenExpired, "expected error when token is expired")
	})

	t.Run("Wrong User ID", func(t *testing.T) {
		// Sign a valid token
		token, err := db.NewResetToken(ulids.New())
		require.NoError(t, err, "could not create reset token")
		signature, secret, err := token.Sign()
		require.NoError(t, err, "failed to sign token")

		// Verification should fail if the user ID is different
		token.UserID = ulids.New()
		require.ErrorIs(t, token.Verify(signature, secret), db.ErrTokenInvalid, "expected error when user ID is different")
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		// Sign a valid token
		token, err := db.NewResetToken(ulids.New())
		require.NoError(t, err, "could not create reset token")
		_, secret, err := token.Sign()
		require.NoError(t, err, "failed to sign token")

		// Verification should fail if the signature is not decodable
		require.Error(t, token.Verify("^&**(", secret), "expected error when signature is not decodable")

		// Verification should fail if the signature was created with a different secret
		otherToken, err := db.NewResetToken(token.UserID)
		require.NoError(t, err, "could not create reset token")
		otherSig, _, err := otherToken.Sign()
		require.NoError(t, err, "failed to sign token")
		require.ErrorIs(t, token.Verify(otherSig, secret), db.ErrTokenInvalid, "expected error when signature was created with a different secret")
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		// Sign a valid token
		token, err := db.NewResetToken(ulids.New())
		require.NoError(t, err, "could not create reset token")
		signature, _, err := token.Sign()
		require.NoError(t, err, "failed to sign token")

		// Should error if the secret has the wrong length
		require.ErrorIs(t, token.Verify(signature, nil), db.ErrInvalidSecret, "expected error when secret is nil")
		require.ErrorIs(t, token.Verify(signature, []byte("wronglength")), db.ErrInvalidSecret, "expected error when secret is the wrong length")

		fmt.Println("signature", signature)

		// Verification should fail if the wrong secret is used
		otherToken, err := db.NewResetToken(token.UserID)
		require.NoError(t, err, "could not create reset token")
		_, otherSecret, err := otherToken.Sign()
		require.NoError(t, err, "failed to sign token")
		require.ErrorIs(t, token.Verify(signature, otherSecret), db.ErrTokenInvalid, "expected error when wrong secret is used")
	})
}
