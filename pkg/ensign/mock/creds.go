package mock

import (
	"context"

	"google.golang.org/grpc"
)

type Credentials struct {
	token string
}

func (t *Credentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.token,
	}, nil
}

func (t *Credentials) RequireTransportSecurity() bool {
	return false
}

func PerRPCToken(token string) grpc.CallOption {
	return grpc.PerRPCCredentials(&Credentials{token: token})
}

func WithPerRPCToken(token string) grpc.DialOption {
	return grpc.WithPerRPCCredentials(&Credentials{token: token})
}
