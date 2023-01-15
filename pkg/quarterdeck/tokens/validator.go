package tokens

import (
	"errors"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Validators are able to verify that access and refresh tokens were issued by
// Quarterdeck and that their claims are valid (e.g. not expired).
type Validator interface {
	// Verify an access or a refresh token after parsing and return its claims.
	Verify(tks string) (claims *Claims, err error)

	// Parse an access or refresh token without verifying claims (e.g. to check an expired token).
	Parse(tks string) (claims *Claims, err error)
}

// JWKSValidator provides public verification that JWT tokens have been issued by the
// Quarterdeck authentication service by checking that the tokens have been signed using
// public keys from a JSON Web Key Set (JWKS). The validator then returns Quarterdeck
// specific claims if the token is in fact valid.
type JWKSValidator struct {
	keys     jwk.Set
	audience string
	issuer   string
}

func NewJWKSValidator(keys jwk.Set, audience, issuer string) *JWKSValidator {
	return &JWKSValidator{keys: keys, audience: audience, issuer: issuer}
}

func (v *JWKSValidator) Verify(tks string) (claims *Claims, err error) {
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(tks, &Claims{}, v.keyFunc); err != nil {
		return nil, err
	}

	var ok bool
	if claims, ok = token.Claims.(*Claims); ok && token.Valid {
		if !claims.VerifyAudience(v.audience, true) {
			return nil, fmt.Errorf("invalid audience %q", claims.Audience)
		}

		if !claims.VerifyIssuer(v.issuer, true) {
			return nil, fmt.Errorf("invalid issuer %q", claims.Issuer)
		}

		return claims, nil
	}

	return nil, fmt.Errorf("could not parse or verify claims from %T", token.Claims)
}

func (v *JWKSValidator) Parse(tks string) (claims *Claims, err error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	claims = &Claims{}
	if _, err = parser.ParseWithClaims(tks, claims, v.keyFunc); err != nil {
		return nil, err
	}
	return claims, nil
}

// keyFunc is an jwt.KeyFunc that selects the RSA public key from the list of managed
// internal keys based on the kid in the token header. If the kid does not exist an
// error is returned and the token will not be able to be verified.
func (v *JWKSValidator) keyFunc(token *jwt.Token) (publicKey interface{}, err error) {
	// Fetch the kid from the header
	kid, ok := token.Header["kid"]
	if !ok {
		return nil, errors.New("token does not have kid in header")
	}

	key, found := v.keys.LookupKeyID(kid.(string))
	if !found {
		return nil, errors.New("unknown signing key")
	}

	// Per JWT security notice: do not forget to validate alg is expected
	if token.Method.Alg() != key.Algorithm().String() {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Extract the raw public key from the key material and return it.
	if err = key.Raw(&publicKey); err != nil {
		return nil, fmt.Errorf("could not extract raw key: %w", err)
	}
	return publicKey, nil
}
