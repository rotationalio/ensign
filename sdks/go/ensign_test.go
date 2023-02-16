package ensign_test

import (
	"testing"

	ensign "github.com/rotationalio/ensign/sdks/go"
	"github.com/stretchr/testify/require"
)

func TestNilWilNilOpts(t *testing.T) {
	_, err := ensign.New(nil)
	require.NoError(t, err, "could not pass nil into ensign")
}

func TestOptions(t *testing.T) {
	opts := &ensign.Options{
		Endpoint:     "localhost:443",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}

	// Test Endpoint is required
	opts.Endpoint = ""
	require.EqualError(t, opts.Validate(), ensign.ErrMissingEndpoint.Error(), "opts should be invalid with missing endpoint")

	// Test ClientID is required
	opts.Endpoint = "localhost:443"
	opts.ClientID = ""
	require.EqualError(t, opts.Validate(), ensign.ErrMissingClientID.Error(), "opts should be invalid with missing client ID")

	// Test ClientSecret is required
	opts.ClientID = "client-id"
	opts.ClientSecret = ""
	require.EqualError(t, opts.Validate(), ensign.ErrMissingClientSecret.Error(), "opts should be invalid with missing client secret")

	// Test valid options
	opts.ClientSecret = "client-secret"
	require.NoError(t, opts.Validate(), "opts should be valid")
}
