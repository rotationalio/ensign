package tokens

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	ulid "github.com/oklog/ulid/v2"
)

// Token time constraint constants.
// TODO: move to configuration file.
const (
	accessTokenDuration  = 1 * time.Hour
	refreshTokenDuration = 2 * time.Hour
	accessRefreshOverlap = -15 * time.Minute
)

// Global variables that should really not be changed except between major versions.
// NOTE: the signing method should match the value returned by the JWKS
var (
	signingMethod = jwt.SigningMethodRS256
	nilID         = ulid.ULID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

// TokenManager handles the creation and verification of RSA signed JWT tokens. To
// facilitate signing key rollover, TokenManager can accept multiple keys identified by
// a ulid. JWT tokens generated by token managers include a kid in the header that
// allows the token manager to verify the key with the specified signature. To sign keys
// the token manager will always use the latest private key by ulid.
//
// When the TokenManager creates tokens it will use JWT standard claims as well as
// extended claims based on Quarterdeck usage. The standard claims included are exp, nbf
// aud, and sub. On token verification, the exp, nbf, iss and aud claims are validated.
// TODO: Create automatic key rotation mechanism rather than loading keys.
type TokenManager struct {
	audience     string
	issuer       string
	currentKeyID ulid.ULID
	currentKey   *rsa.PrivateKey
	keys         map[ulid.ULID]*rsa.PublicKey
}

// New creates a TokenManager with the specified keys which should be a mapping of ULID
// strings to paths to files that contain PEM encoded RSA private keys. This input is
// specifically designed for the config environment variable so that keys can be loaded
// from k8s or vault secrets that are mounted as files on disk.
func New(keys map[string]string, audience, issuer string) (tm *TokenManager, err error) {
	tm = &TokenManager{
		keys:     make(map[ulid.ULID]*rsa.PublicKey),
		audience: audience,
		issuer:   issuer,
	}

	for kid, path := range keys {
		// Parse the key id
		var keyID ulid.ULID
		if keyID, err = ulid.Parse(kid); err != nil {
			return nil, fmt.Errorf("could not parse kid %q for path %s: %s", kid, path, err)
		}

		// Load the keys from disk
		var data []byte
		if data, err = os.ReadFile(path); err != nil {
			return nil, fmt.Errorf("could not read kid %s from %s: %s", kid, path, err)
		}

		var key *rsa.PrivateKey
		if key, err = jwt.ParseRSAPrivateKeyFromPEM(data); err != nil {
			return nil, fmt.Errorf("could not parse RSA private key kid %s from %s: %s", kid, path, err)
		}

		// Add the key to the key map
		tm.keys[keyID] = &key.PublicKey

		// Set the current key if it is the latest key
		if tm.currentKey == nil || keyID.Time() > tm.currentKeyID.Time() {
			tm.currentKey = key
			tm.currentKeyID = keyID
		}
	}

	return tm, nil
}

// Verify an access or a refresh token after parsing and return its claims.
func (tm *TokenManager) Verify(tks string) (claims *Claims, err error) {
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(tks, &Claims{}, tm.keyFunc); err != nil {
		return nil, err
	}

	var ok bool
	if claims, ok = token.Claims.(*Claims); ok && token.Valid {
		if !claims.VerifyAudience(tm.audience, true) {
			return nil, fmt.Errorf("invalid audience %q", claims.Audience)
		}

		if !claims.VerifyIssuer(tm.issuer, true) {
			return nil, fmt.Errorf("invalid issuer %q", claims.Issuer)
		}

		return claims, nil
	}

	return nil, fmt.Errorf("could not parse or verify claims from %T", token.Claims)
}

// Parse an access or refresh token verifying its signature but without verifying its
// claims. This ensures that valid JWT tokens are still accepted but claims can be
// handled on a case-by-case basis; for example by validating an expired access token
// during reauthentication.
func (tm *TokenManager) Parse(tks string) (claims *Claims, err error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	claims = &Claims{}
	if _, err = parser.ParseWithClaims(tks, claims, tm.keyFunc); err != nil {
		return nil, err
	}
	return claims, nil
}

// Sign an access or refresh token and return the token string.
func (tm *TokenManager) Sign(token *jwt.Token) (tks string, err error) {
	// Sanity check to prevent nil panics.
	if tm.currentKey == nil || tm.currentKeyID.Compare(nilID) == 0 {
		return "", errors.New("token manager not initialized with signing keys")
	}

	// Add the kid (key id - this is the standard 3 letter JWT name) to the header.
	token.Header["kid"] = tm.currentKeyID.String()

	// Return the signed string
	return token.SignedString(tm.currentKey)
}

// CreateTokenPair returns signed access and refresh tokens for the specified claims in
// one step (since normally users want both an access and a refresh token)!
func (tm *TokenManager) CreateTokenPair(claims *Claims) (accessToken, refreshToken string, err error) {
	var atk, rtk *jwt.Token
	if atk, err = tm.CreateAccessToken(claims); err != nil {
		return "", "", fmt.Errorf("could not create access token: %w", err)
	}

	if rtk, err = tm.CreateRefreshToken(atk); err != nil {
		return "", "", fmt.Errorf("could not create refresh token: %w", err)
	}

	if accessToken, err = tm.Sign(atk); err != nil {
		return "", "", fmt.Errorf("could not sign access token: %w", err)
	}

	if refreshToken, err = tm.Sign(rtk); err != nil {
		return "", "", fmt.Errorf("could not sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// CreateAccessToken from the credential payload or from an previous token if the
// access token is being reauthorized from previous credentials. Note that the returned
// token only contains the claims and is unsigned.
func (tm *TokenManager) CreateAccessToken(claims *Claims) (_ *jwt.Token, err error) {
	// Create the claims for the access token, using access token defaults.
	now := time.Now()
	sub := claims.RegisteredClaims.Subject

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ID:        strings.ToLower(ulid.Make().String()), // ID is randomly generated and shared between access and refresh tokens.
		Subject:   sub,
		Audience:  jwt.ClaimStrings{tm.audience},
		Issuer:    tm.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
	}
	return jwt.NewWithClaims(signingMethod, claims), nil
}

// CreateRefreshToken from the Access token claims with predefined expiration. Note that
// the returned token only contains the claims and is unsigned.
func (tm *TokenManager) CreateRefreshToken(accessToken *jwt.Token) (refreshToken *jwt.Token, err error) {
	accessClaims, ok := accessToken.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not retrieve claims from access token")
	}

	// Create claims for the refresh token from the access token defaults.
	// TODO: should we make this a refresh-specific audience or subject?
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessClaims.ID, // ID is randomly generated and shared between access and refresh tokens.
			Audience:  accessClaims.Audience,
			Issuer:    accessClaims.Issuer,
			Subject:   accessClaims.Subject,
			IssuedAt:  accessClaims.IssuedAt,
			NotBefore: jwt.NewNumericDate(accessClaims.ExpiresAt.Add(accessRefreshOverlap)),
			ExpiresAt: jwt.NewNumericDate(accessClaims.IssuedAt.Add(refreshTokenDuration)),
		},
	}

	return jwt.NewWithClaims(signingMethod, claims), nil
}

// Keys returns the map of ulid to public key for use externally.
func (tm *TokenManager) Keys() map[ulid.ULID]*rsa.PublicKey {
	return tm.keys
}

// CurrentKey returns the ulid of the current key being used to sign tokens.
func (tm *TokenManager) CurrentKey() ulid.ULID {
	return tm.currentKeyID
}

// keyFunc is an jwt.KeyFunc that selects the RSA public key from the list of managed
// internal keys based on the kid in the token header. If the kid does not exist an
// error is returned and the token will not be able to be verified.
func (tm *TokenManager) keyFunc(token *jwt.Token) (key interface{}, err error) {
	// Per JWT security notice: do not forget to validate alg is expected
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Fetch the kid from the header
	kid, ok := token.Header["kid"]
	if !ok {
		return nil, errors.New("token does not have kid in header")
	}

	// Parse the kid
	var keyID ulid.ULID
	if keyID, err = ulid.Parse(kid.(string)); err != nil {
		return nil, fmt.Errorf("could not parse kid: %s", err)
	}

	// Fetch the key from the list of managed keys
	if key, ok = tm.keys[keyID]; !ok {
		return nil, errors.New("unknown signing key")
	}
	return key, nil
}
