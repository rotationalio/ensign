package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

const (
	bearer                    = "Bearer "
	authorization             = "Authorization"
	ContextUserClaims         = "user_claims"
	DefaultKeysURL            = "https://auth.ensign.app/.well-known/jwks.json"
	DefaultAudience           = "https://ensign.app"
	DefaultIssuer             = "https://auth.ensign.app"
	DefaultMinRefreshInterval = 5 * time.Minute
)

// Authenticate middleware ensures that the request has a valid Bearer JWT in the
// Authenticate header of the request otherwise it stops processing of the request and
// returns a 401 unauthorized error. A valid Bearer JWT means that the access token is
// supplied as the Bearer token, it is unexpired, and it was issued by Quarterdeck by
// checking with the Quarterdeck public keys.
//
// In order to perform authorization, this middleware fetches public JSON Web Key Sets
// (JWKS) from the authorizing Quarterdeck server and then caches them according to the
// Cache-Control or Expires headers in the response. As Quarterdeck keys are rotated,
// the cache must refresh the public keys in a background routine in order to correctly
// authorize incoming JWT tokens. Users can control how the JWKS are fetched and cached
// using AuthOptions (which are particularly helpful for tests).
func Authenticate(opts ...AuthOption) (_ gin.HandlerFunc, err error) {
	// Create the authorization options from the variadic arguments.
	conf := NewAuthOptions(opts...)

	// Create the JWK cache object using the context from the configuration
	// This configuration tells the cache we want to refresh the JWKs when it needs to
	// based on the Cache-Control or Expires header from the HTTP response. If the
	// calculated minimum refresh interval is less than the configured minimum it won't
	// refresh the JWKS any earlier. This means that the min refresh interval should be
	// relatively small (e.g. minutes).
	var validator tokens.Validator
	if validator, err = conf.Validator(); err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		var (
			err         error
			accessToken string
			claims      *tokens.Claims
		)

		// Get access token from the request
		if accessToken, err = GetAccessToken(c); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrAuthRequired))
			return
		}

		// Verify the access token is authorized for use with Quarterdeck and extract claims.
		if claims, err = validator.Verify(accessToken); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrAuthRequired))
			return
		}

		// Add claims to context for use in downstream processing and continue handlers
		c.Set(ContextUserClaims, claims)
		c.Next()
	}, nil
}

// Authorize is a middleware that requires specific permissions in an authenticated
// user's claims. If those permissions do not match or the request is unauthenticated
// the middleware returns a 401 response. The Authorize middleware must be chained
// following the Authenticate middleware.
func Authorize(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoAuthorization))
			return
		}

		if !claims.HasAllPermissions(permissions...) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(ErrNoPermission))
			return
		}

		c.Next()
	}
}

// GetAccessToken retrieves the bearer token from the authorization header and parses it
// to return only the JWT access token component of the header. If the header is missing
// or the token is not available, an error is returned.
// TODO: switch to regular expression based parsing.
func GetAccessToken(c *gin.Context) (tks string, err error) {
	header := c.GetHeader(authorization)
	if header != "" {
		parts := strings.Split(header, bearer)
		if len(parts) == 2 {
			tks = strings.TrimSpace(parts[1])
			if tks != "" {
				return tks, nil
			}
		}
		return "", ErrParseBearer
	}
	return "", ErrNoAuthorization
}

// GetClaims fetches and parses Quarterdeck claims from the gin context. Returns an
// error if no claims exist on the context; panics if the claims are not the correct
// type -- however the panic should be recovered by middleware.
func GetClaims(c *gin.Context) (*tokens.Claims, error) {
	claims, exists := c.Get(ContextUserClaims)
	if !exists {
		return nil, ErrNoClaims
	}
	return claims.(*tokens.Claims), nil
}

// AuthOption allows users to optionally supply configuration to the Authorization middleware.
type AuthOption func(opts *AuthOptions)

// AuthOptions is constructed from variadic AuthOption arguments with reasonable defaults.
type AuthOptions struct {
	KeysURL            string           // The URL endpoint to the JWKS public keys on the Quarterdeck server
	Audience           string           // The audience to verify on tokens
	Issuer             string           // The issuer to verify on tokens
	MinRefreshInterval time.Duration    // Minimum amount of time the JWKS public keys are cached
	Context            context.Context  // The context object to control the lifecycle of the background fetch routine
	validator          tokens.Validator // The validator constructed by the auth options (can be directly supplied by the user).
}

// NewAuthOptions creates an AuthOptions object with reasonable defaults and any user
// supplied input from the AuthOption variadic arguments.
func NewAuthOptions(opts ...AuthOption) (conf AuthOptions) {
	conf = AuthOptions{
		KeysURL:            DefaultKeysURL,
		Audience:           DefaultAudience,
		Issuer:             DefaultIssuer,
		MinRefreshInterval: DefaultMinRefreshInterval,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	// Create a context if one has not been supplied by the user.
	if conf.Context == nil && conf.validator == nil {
		conf.Context = context.Background()
	}
	return conf
}

// Validator returns the user supplied validator or constructs a new JWKS Cache
// Validator from the supplied options. If the options are invalid or the validator
// cannot be created an error is returned.
func (conf *AuthOptions) Validator() (_ tokens.Validator, err error) {
	if conf.validator == nil {
		cache := jwk.NewCache(conf.Context)
		cache.Register(conf.KeysURL, jwk.WithMinRefreshInterval(conf.MinRefreshInterval))

		if conf.validator, err = tokens.NewCachedJWKSValidator(conf.Context, cache, conf.KeysURL, conf.Audience, conf.Issuer); err != nil {
			return nil, err
		}
	}
	return conf.validator, nil
}

// WithAuthOptions allows the user to update the default auth options with an auth
// options struct to set many options values at once. Zero values are ignored, so if
// using this option, the defaults will still be preserved if not set on the input.
func WithAuthOptions(opts AuthOptions) AuthOption {
	return func(conf *AuthOptions) {
		if opts.KeysURL != "" {
			conf.KeysURL = opts.KeysURL
		}

		if opts.Audience != "" {
			conf.Audience = opts.Audience
		}

		if opts.Issuer != "" {
			conf.Issuer = opts.Issuer
		}

		if opts.MinRefreshInterval != 0 {
			conf.MinRefreshInterval = opts.MinRefreshInterval
		}

		if opts.Context != nil {
			conf.Context = opts.Context
		}
	}
}

// WithJWKSEndpoint allows the user to specify an alternative endpoint to fetch the JWKS
// public keys from. This is useful for testing or for different environments.
func WithJWKSEndpoint(url string) AuthOption {
	return func(opts *AuthOptions) {
		opts.KeysURL = url
	}
}

// WithAudience allows the user to specify an alternative audience.
func WithAudience(audience string) AuthOption {
	return func(opts *AuthOptions) {
		opts.Audience = audience
	}
}

// WithIssuer allows the user to specify an alternative issuer.
func WithIssuer(issuer string) AuthOption {
	return func(opts *AuthOptions) {
		opts.Issuer = issuer
	}
}

// WithMinRefreshInterval allows the user to specify an alternative minimum duration
// between cache refreshes to control refresh behavior for the JWKS public keys.
func WithMinRefreshInterval(interval time.Duration) AuthOption {
	return func(opts *AuthOptions) {
		opts.MinRefreshInterval = interval
	}
}

// WithContext allows the user to specify an external, cancelable context to control
// the background refresh behavior of the JWKS cache.
func WithContext(ctx context.Context) AuthOption {
	return func(opts *AuthOptions) {
		opts.Context = ctx
	}
}

// WithValidator allows the user to specify an alternative validator to the auth
// middleware. This is particularly useful for testing authentication.
func WithValidator(validator tokens.Validator) AuthOption {
	return func(opts *AuthOptions) {
		opts.validator = validator
	}
}
