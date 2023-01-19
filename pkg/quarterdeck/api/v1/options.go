package api

import (
	"net/http"
)

// ClientOption allows us to configure the APIv1 client when it is created.
type ClientOption func(c *APIv1) error

func WithClient(client *http.Client) ClientOption {
	return func(c *APIv1) error {
		c.client = client
		return nil
	}
}

func WithCredentials(creds Credentials) ClientOption {
	return func(c *APIv1) error {
		c.creds = creds
		return nil
	}
}

// RequestOption allows us to configure individual APIv1 client requests
// TODO: this is just a hack that modifies the request to get us started, but we will
// likely want to consider what design pattern we want to use in SC-12797.
type RequestOption func(req *http.Request) error

// WithRPCCredentials overwrites any existing Authorize header with the specified creds.
func WithRPCCredentials(creds Credentials) RequestOption {
	return func(req *http.Request) (err error) {
		var token string
		if token, err = creds.AccessToken(); err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
