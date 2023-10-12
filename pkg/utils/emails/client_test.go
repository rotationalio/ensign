package emails_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rotationalio/confire"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/stretchr/testify/require"
)

// TestLiveSend enables you to send a test email to an actual email address so long as
// the $SENDGRID_TEST_SENDING_LIVE_EMAILS environment variable is set to "1". The
// $SENDGRID_API_KEY, and $SENDGRID_ADMINS_EMAIL environment variables are also required
// the emails will be sent to the $SENDGRID_ADMINS_EMAIL address.
func TestLiveSend(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	// Only run the test if an environment variable is set
	if os.Getenv("SENDGRID_TEST_SENDING_LIVE_EMAILS") != "1" {
		t.Skip("not testing live email sends without environment configuration")
	}

	var conf emails.Config
	err := confire.Process("SENDGRID", &conf)
	require.NoError(t, err, "could not load email configuration")

	client, err := emails.New(conf)
	require.NoError(t, err, "could not create email client for live sending")

	emailData := emails.EmailData{
		Sender:    conf.MustFromContact(),
		Recipient: conf.MustAdminContact(),
	}

	t.Run("WelcomeEmail", func(t *testing.T) {
		data := emails.WelcomeData{
			EmailData:    emailData,
			FirstName:    "Rico",
			LastName:     "Hernandez",
			Email:        "rico@example.com",
			Organization: "Checkers Labs",
			Domain:       "checkers",
		}

		message, err := emails.WelcomeEmail(data)
		require.NoError(t, err, "could not create welcome email")

		err = client.Send(message)
		require.NoError(t, err, "could not send welcome email")
	})

	t.Run("VerifyEmail", func(t *testing.T) {
		data := emails.VerifyEmailData{
			EmailData: emailData,
			FullName:  "Rodrigo Balentine",
			VerifyURL: "https://bbengfort.github.io",
		}

		message, err := emails.VerifyEmail(data)
		require.NoError(t, err, "could not create verify email")

		err = client.Send(message)
		require.NoError(t, err, "could not send verify email")
	})

	t.Run("InviteEmail", func(t *testing.T) {
		data := emails.InviteData{
			EmailData:   emailData,
			Email:       "rico@example.com",
			InviterName: "Bella Washington",
			OrgName:     "Chess Strategies, Inc.",
			Role:        "Observer",
			InviteURL:   "https://rotational.io/blog/year-one-lessons/",
		}

		message, err := emails.InviteEmail(data)
		require.NoError(t, err, "could not create invite email")

		err = client.Send(message)
		require.NoError(t, err, "could not send invite email")
	})

	t.Run("PasswordResetRequest", func(t *testing.T) {
		data := emails.ResetRequestData{
			EmailData: emailData,
			ResetURL:  "https://bbengfort.github.io",
		}

		message, err := emails.PasswordResetRequestEmail(data)
		require.NoError(t, err, "could not create password reset request email")

		err = client.Send(message)
		require.NoError(t, err, "could not send password reset request email")
	})

	t.Run("PasswordResetSuccess", func(t *testing.T) {
		message, err := emails.PasswordResetSuccessEmail(emailData)
		require.NoError(t, err, "could not create password reset success email")

		err = client.Send(message)
		require.NoError(t, err, "could not send password reset success email")
	})

	t.Run("DailyUsers", func(t *testing.T) {
		// TODO: populate data!
		data := emails.DailyUsersData{
			EmailData: emailData,
		}

		message, err := emails.DailyUsersEmail(data)
		require.NoError(t, err, "could not create daily users email")

		err = client.Send(message)
		require.NoError(t, err, "could not send daily users email")
	})
}
