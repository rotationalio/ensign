package tenant

import (
	sdk "github.com/rotationalio/go-ensign"
	"github.com/rotationalio/go-ensign/auth"
)

// EnsignClient wraps an Ensign SDK client for specific usage. This is not strictly
// necessary but it allows us to specify how Tenant is interacting with Ensign. For
// example, in request handlers we want to make Ensign requests with the user's
// permissions which requires per-rpc authentication.
type EnsignClient struct {
	client *sdk.Client
	opts   *sdk.Options
}

// NewEnsignClient creates an Ensign client from the configuration
func NewEnsignClient(opts *sdk.Options) (*EnsignClient, error) {
	if opts == nil {
		opts = sdk.NewOptions()
	}

	client, err := sdk.New(opts)
	if err != nil {
		return nil, err
	}

	return &EnsignClient{client: client, opts: opts}, nil
}

// InvokeOnce exposes a clone of the SDK client for a single call using the provided
// token for per-rpc authentication. This should be used in request handlers where
// Ensign requests are made on behalf of the user.
func (c *EnsignClient) InvokeOnce(token string) *sdk.Client {
	return c.client.WithCallOptions(auth.PerRPCToken(token, c.opts.Insecure))
}

// Set an SDK client on the client for testing purposes
func (c *EnsignClient) SetClient(client *sdk.Client) {
	c.client = client
}

// Set SDK options on the client for testing purposes
func (c *EnsignClient) SetOpts(opts *sdk.Options) {
	c.opts = opts
}
