package emails_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/stretchr/testify/require"
)

func TestSendGrid(t *testing.T) {
	conf := &emails.Config{}
	require.False(t, conf.Enabled(), "sendgrid should be disabled when there is no API key")
	require.NoError(t, conf.Validate(), "no validation error should be returned when sendgrid is disabled")

	conf.APIKey = "SG.testing123"
	require.True(t, conf.Enabled(), "sendgrid should be enabled when there is an API key")

	// FromEmail is required when enabled
	conf.FromEmail = ""
	conf.AdminEmail = "test@example.com"
	require.Error(t, conf.Validate(), "expected from email to be required")

	// AdminEmail is required when enabled
	conf.FromEmail = "test@example.com"
	conf.AdminEmail = ""
	require.Error(t, conf.Validate(), "expected admin email to be required")

	// Should be valid when enabled and emails are specified
	conf = &emails.Config{
		APIKey:     "testing123",
		FromEmail:  "test@example.com",
		AdminEmail: "admin@example.com",
	}
	require.NoError(t, conf.Validate(), "expected configuration to be valid")

	// Archive is only supported in testing mode
	conf.Archive = "fixtures/emails"
	require.Error(t, conf.Validate(), "expected error when archive is set in non-testing mode")
	conf.Testing = true
	require.NoError(t, conf.Validate(), "expected configuration to be valid with archive in testing mode")
}
