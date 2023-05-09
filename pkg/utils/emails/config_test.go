package emails_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
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

	// Require parsable emails when enabled
	conf.FromEmail = "foo"
	conf.AdminEmail = "test@example.com"
	require.Error(t, conf.Validate())

	conf.FromEmail = "test@example.com"
	conf.AdminEmail = "foo"
	require.Error(t, conf.Validate())

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

func TestContactParsing(t *testing.T) {
	testCases := []struct {
		FromEmail    string
		AdminEmail   string
		FromContact  *sendgrid.Contact
		AdminContact *sendgrid.Contact
		FromErr      error
		AdminErr     error
	}{
		{"", "", nil, nil, emails.ErrUnparsable, emails.ErrUnparsable},
		{"foo", "bar", nil, nil, emails.ErrUnparsable, emails.ErrUnparsable},
		{"enson@rotational.io", "ensign@rotational.io", &sendgrid.Contact{Email: "enson@rotational.io"}, &sendgrid.Contact{Email: "ensign@rotational.io"}, nil, nil},
		{"The Ensign Team at Rotational <enson@rotational.io>", "Rotational Admins <ensign@rotational.io>", &sendgrid.Contact{FirstName: "The", LastName: "Ensign Team at Rotational", Email: "enson@rotational.io"}, &sendgrid.Contact{FirstName: "Rotational", LastName: "Admins", Email: "ensign@rotational.io"}, nil, nil},
	}

	for i, tc := range testCases {
		conf := emails.Config{
			FromEmail:  tc.FromEmail,
			AdminEmail: tc.AdminEmail,
		}

		actualFromContact, err := conf.FromContact()
		if tc.FromErr != nil {
			require.Error(t, err, "expected from contact error for test case %d", i)
			require.ErrorIs(t, err, tc.FromErr, "test case %d failed", i)
		} else {
			require.Equal(t, tc.FromContact, &actualFromContact, "test case %d failed", i)
		}

		actualAdminContact, err := conf.AdminContact()
		if tc.AdminErr != nil {
			require.Error(t, err, "expected from contact error for test case %d", i)
			require.ErrorIs(t, err, tc.AdminErr, "test case %d failed", i)
		} else {
			require.Equal(t, tc.AdminContact, &actualAdminContact, "test case %d failed", i)
		}
	}
}
