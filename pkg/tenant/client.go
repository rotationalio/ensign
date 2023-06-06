package tenant

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	sdk "github.com/rotationalio/go-ensign"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/rotationalio/go-ensign/auth"
)

// EnsignClient wraps an Ensign SDK client for specific usage. This is not strictly
// necessary but it allows us to specify how Tenant is interacting with Ensign. For
// example, in request handlers we want to make Ensign requests with the user's
// permissions which requires per-rpc authentication.
type EnsignClient struct {
	client *sdk.Client
	conf   config.SDKConfig
}

// NewEnsignClient creates an Ensign client from the configuration
func NewEnsignClient(conf config.SDKConfig) (_ *EnsignClient, err error) {
	var client *sdk.Client
	if client, err = sdk.New(conf.ClientOptions()...); err != nil {
		return nil, err
	}

	return &EnsignClient{client: client, conf: conf}, nil
}

// InvokeOnce exposes a clone of the SDK client for a single call using the provided
// token for per-rpc authentication. This should be used in request handlers where
// Ensign requests are made on behalf of the user.
func (c *EnsignClient) InvokeOnce(token string) *sdk.Client {
	return c.client.WithCallOptions(auth.PerRPCToken(token, c.conf.Insecure))
}

// WaitForReady is a client-side wait that blocks until the Ensign server is ready to
// accept requests or the timeout is exceeded.
func (c *EnsignClient) WaitForReady() (attempts int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.conf.WaitForReady)
	defer cancel()

	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
	defer ticker.Stop()

	for {
		attempts++
		select {
		case <-ctx.Done():
			return attempts, ctx.Err()
		case <-ticker.C:
			var rep *api.ServiceState
			if rep, err = c.client.Status(ctx); err != nil {
				continue
			}

			if rep.Status == api.ServiceState_HEALTHY {
				return attempts, nil
			}
		}
	}
}

// Subscribe uses the credentials in the client to subscribe to the configured topic
// and returns the subscriber channel.
func (c *EnsignClient) Subscribe() (sub *sdk.Subscription, err error) {
	return c.client.Subscribe(c.conf.TopicName)
}

// Set an SDK client on the client for testing purposes
func (c *EnsignClient) SetClient(client *sdk.Client) {
	c.client = client
}

// Set SDK options on the client for testing purposes
func (c *EnsignClient) SetOpts(conf config.SDKConfig) {
	c.conf = conf
}
