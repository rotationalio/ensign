package middleware

import "errors"

var (
	ErrUnauthenticated  = errors.New("request is unauthenticated")
	ErrNoClaims         = errors.New("no claims found on the request context")
	ErrNoUserInfo       = errors.New("no user info found on the request context")
	ErrInvalidAuthToken = errors.New("invalid authorization token")
	ErrAuthRequired     = errors.New("this endpoint requires authentication")
	ErrNoPermission     = errors.New("user does not have permission to perform this operation")
	ErrNoAuthUser       = errors.New("could not identify authenticated user in request")
	ErrNoAuthUserData   = errors.New("could not retrieve user data")
	ErrIncompleteUser   = errors.New("user is missing required fields")
	ErrUnverifiedUser   = errors.New("user is not verified")
	ErrCSRFVerification = errors.New("csrf verification failed for request")
	ErrParseBearer      = errors.New("could not parse Bearer token from Authorization header")
	ErrNoAuthorization  = errors.New("no authorization header in request")
	ErrNoRequest        = errors.New("no request found on the context")
	ErrRateLimit        = errors.New("rate limit reached: too many requests")
	ErrNoRefreshToken   = errors.New("no refresh token available on request")
	ErrRefreshDisabled  = errors.New("reauthentication with refresh tokens disabled")
)
