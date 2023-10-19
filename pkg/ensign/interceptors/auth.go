package interceptors

import (
	"context"
	"net"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	health "github.com/rotationalio/ensign/pkg/utils/probez/grpc/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	header = "authorization" // MUST BE LOWER CASE!
	bearer = "Bearer "       // MUST INCLUDE TRAILING SPACE!
)

// Authenticator ensures that the RPC request has a valid Quarterdeck-issued JWT token
// in the credentials metadata of the request, otherwise it stops processing and returns
// an Unauthenticated error. A valid JWT token means that the token is supplied in the
// credentials, is unexpired, was signed by Quarterdeck private keys, and has the
// correct audience and issuer.
//
// This interceptor extracts the claims from the JWT token and adds them to the context
// of the request, ensuring that downstream interceptors and the handlers can access the
// claims without having to parse the JWT token in the credentials.
//
// In order to perform authentication, this middleware fetches public JSON Web Key Sets
// (JWKS) from the authorizing Quarterdeck server and caches them according to the
// Cache-Control or Expires headers in the response. As Quarterdeck keys are rotated,
// the cache must refresh the public keys in a background routine to correctly
// authenticate incoming credentials. Users can control how the JWKS are fetched and
// cached using AuthOptions from the Quarterdeck middleware package.
//
// Both Unary and Streaming interceptors can be returned from this middleware handler.
type Authenticator struct {
	conf         middleware.AuthOptions
	validator    tokens.Validator
	publicRoutes map[string]struct{}
}

// Create an authenticator to handle both unary and streaming RPC calls, modifying the
// behavior of the authenticator using auth options from Quarterdeck middleware.
func NewAuthenticator(opts ...middleware.AuthOption) (auth *Authenticator, err error) {
	auth = &Authenticator{
		conf: middleware.NewAuthOptions(opts...),
	}

	if auth.validator, err = auth.conf.Validator(); err != nil {
		return nil, err
	}

	// Setup public routes (e.g. routes that do not require authentication)
	auth.publicRoutes = map[string]struct{}{
		statusEndpoint:       {},
		health.CheckEndpoint: {},
		health.WatchEndpoint: {},
	}

	return auth, nil
}

// Authenticate a request using the access token credentials provided in the metadata.
func (a *Authenticator) authenticate(ctx context.Context) (_ context.Context, err error) {
	var (
		claims *tokens.Claims
		md     metadata.MD
		ok     bool
	)

	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Extract the authorization credentials (we expect [at least] 1 JWT token)
	values := md[header]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Loop through credentials to find the first valid claims
	// NOTE: we only expect one token but are trying to future-proof the interceptor
	for _, value := range values {
		if !strings.HasPrefix(value, bearer) {
			continue
		}

		token := strings.TrimPrefix(value, bearer)
		if claims, err = a.validator.Verify(token); err == nil {
			break
		}
	}

	// Check to see if we found any valid claims in the request
	if claims == nil {
		log.Debug().Err(err).Int("tokens", len(values)).Msg("could not find a valid access token in request")
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Specify user for Sentry if Sentry is configured
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		var remoteIP string
		if remote, ok := peer.FromContext(ctx); ok {
			remoteIP = UserIP(remote.Addr)

		}

		hub.Scope().SetUser(sentry.User{
			ID:        claims.Subject,
			Email:     claims.Email,
			Name:      claims.Name,
			IPAddress: remoteIP,
		})
	}

	// Add the claims to the context so that downstream handlers can access it
	return contexts.WithClaims(ctx, claims), nil
}

// Return the Unary interceptor that uses the Authenticator handler.
func (a *Authenticator) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if !a.isPublic(info.FullMethod) {
			if ctx, err = a.authenticate(ctx); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

// Return the Stream interceptor that uses the Authenticator handler.
func (a *Authenticator) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		if !a.isPublic(info.FullMethod) {
			var ctx context.Context
			if ctx, err = a.authenticate(stream.Context()); err != nil {
				return err
			}
			stream = contexts.Stream(stream, ctx)
		}

		return handler(srv, stream)
	}
}

func (a *Authenticator) isPublic(route string) bool {
	_, ok := a.publicRoutes[route]
	return ok
}

func UserIP(addr net.Addr) string {
	// If this is a TCP ip address then handle it directly
	if tcpaddr, ok := addr.(*net.TCPAddr); ok {
		return tcpaddr.IP.String()
	}

	// Try parsing the IP address the hard way
	ipaddr := addr.String()
	if host, _, err := net.SplitHostPort(ipaddr); err == nil {
		return host
	}
	return ipaddr
}
