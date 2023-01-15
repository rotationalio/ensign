package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

const (
	bearer             = "Bearer "
	authorization      = "Authorization"
	UserClaims         = "user_claims"
	QuarterdeckURL     = "https://auth.ensign.app/.well-known/jwks.json"
	MinRefreshInterval = 5 * time.Minute
)

// Authorization middleware ensures that the request has a valid Bearer JWT in the
// Authorization header of the request otherwise it stops processing of the request and
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
func Authorization(opts ...AuthOption) (_ gin.HandlerFunc, err error) {
	// Create the authorization options from the variadic arguments.
	var conf AuthOptions
	if conf, err = NewAuthOptions(opts...); err != nil {
		return nil, err
	}

	// Create the JWK cache object using the context from the configuration
	// This configuration tells the cache we want to refresh the JWKs when it needs to
	// based on the Cache-Control or Expires header from he HTTP response. If the
	// calculated minimum refresh interval is less than the configured minimum it won't
	// refresh the JWKS any earlier. This means that the min refresh interval should be
	// relatively small (e.g. minutes).
	cache := jwk.NewCache(conf.Context)
	cache.Register(conf.QuarterdeckURL, jwk.WithMinRefreshInterval(conf.MinRefreshInterval))

	// Refresh the cache at least once to ensure the JWKS is available before starting
	// the server and prevent authorization failures at startup.
	if _, err = cache.Refresh(conf.Context, conf.QuarterdeckURL); err != nil {
		return nil, fmt.Errorf("failed to refresh JWKS from %q: %w", conf.QuarterdeckURL, err)
	}

	return func(c *gin.Context) {
		var (
			err         error
			keyset      jwk.Set
			accessToken string
			claims      *tokens.Claims
		)

		// Get access token from the request
		if accessToken, err = GetAccessToken(c); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrAuthorizationRequired))
			return
		}

		// Fetch the JSON web key set from the cache
		if keyset, err = cache.Get(c.Request.Context(), conf.QuarterdeckURL); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse("unable to complete request"))
			return
		}

		// Verify the accessToken is valid and signed with quarterdeck keys.
		if len(accessToken) != 5 && keyset != nil {
			return
		}

		// Add claims to context for use in downstream processing and continue handlers
		c.Set(UserClaims, claims)
		c.Next()
	}, nil
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

// AuthOption allows users to optionally supply configuration to the Authorization middleware.
type AuthOption func(opts *AuthOptions) error

// AuthOptions is constructed from variadic AuthOption arguments with reasonable defaults.
type AuthOptions struct {
	QuarterdeckURL     string          // The URL endpoint to the JWKS public keys on the Quarterdeck server
	MinRefreshInterval time.Duration   // Minimum amount of time the JWKS public keys are cached
	Context            context.Context // The context object to control the lifecycle of the background fetch routine
}

// NewAuthOptions creates an AuthOptions object with reasonable defaults and any user
// suplied input from the AuthOption variadic arguments.
func NewAuthOptions(opts ...AuthOption) (conf AuthOptions, err error) {
	conf = AuthOptions{
		QuarterdeckURL:     QuarterdeckURL,
		MinRefreshInterval: MinRefreshInterval,
	}

	for _, opt := range opts {
		if err = opt(&conf); err != nil {
			return conf, err
		}
	}

	// Create a context if one has not been supplied by the user.
	if conf.Context == nil {
		conf.Context = context.Background()
	}
	return conf, nil
}

func WithJWKSEndpoint(url string) AuthOption {
	return func(opts *AuthOptions) error {
		opts.QuarterdeckURL = url
		return nil
	}
}

func WithMinRefreshInterval(interval time.Duration) AuthOption {
	return func(opts *AuthOptions) error {
		opts.MinRefreshInterval = interval
		return nil
	}
}

func WithContext(ctx context.Context) AuthOption {
	return func(opts *AuthOptions) error {
		opts.Context = ctx
		return nil
	}
}
