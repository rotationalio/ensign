package contexts

import (
	"context"

	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

// Ensign-specific context keys for passing values to concurrent requests
type contextKey uint8

// Allocate context keys to simplify context key usage in Ensign
const (
	KeyUnknown contextKey = iota
	KeyClaims
)

// WithClaims returns a copy of the parent context with the access claims stored as a
// value on the new context. Users can fetch the claims using the ClaimsFrom function.
func WithClaims(parent context.Context, claims *tokens.Claims) context.Context {
	return context.WithValue(parent, KeyClaims, claims)
}

// ClaimsFrom returns the claims from the context if they exist or false if not.
func ClaimsFrom(ctx context.Context) (*tokens.Claims, bool) {
	claims, ok := ctx.Value(KeyClaims).(*tokens.Claims)
	return claims, ok
}

// Authorize reduces a multistep process into a single step; fetching the claims from
// the context and checking that the claims have the required permission. If there are
// no claims in the context or the permission is invalid, then an error is returned.
func Authorize(ctx context.Context, permission string) (*tokens.Claims, error) {
	claims, ok := ClaimsFrom(ctx)
	if !ok {
		return nil, ErrNoClaimsInContext
	}

	if !claims.HasPermission(permission) {
		return nil, ErrNotAuthorized
	}

	return claims, nil
}

var contextKeyNames = []string{"unknown", "claims"}

// String returns a human readable representation of the context key for easier debugging.
func (c contextKey) String() string {
	if int(c) < len(contextKeyNames) {
		return contextKeyNames[c]
	}
	return contextKeyNames[0]
}
