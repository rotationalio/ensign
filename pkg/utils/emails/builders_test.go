package emails_test

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
)

func TestEmailBuilders(t *testing.T) {
	setupMIMEDir(t)

	emailData := emails.EmailData{
		SenderName:     "Lewis Hudson",
		SenderEmail:    "lewis@example.com",
		RecipientName:  "Rachel Lendt",
		RecipientEmail: "rachel@example.com",
	}

	welcomeData := emails.WelcomeData{
		EmailData:    emailData,
		FirstName:    "Rachel",
		LastName:     "Lendt",
		Email:        "rachel@example.com",
		Organization: "Events R Us",
		Domain:       "eventsrus.com",
	}
	mail, err := emails.WelcomeEmail(welcomeData)
	require.NoError(t, err, "expected no error when building welcome email")
	require.Equal(t, emails.WelcomeRE, mail.Subject, "expected welcome email subject to match")
	generateMIME(t, mail, "welcome.mime")
}

func TestLoadAttachment(t *testing.T) {
	// Add an attachment to a new email
	msg := mail.NewV3Mail()
	err := emails.LoadAttachment(msg, filepath.Join("testdata", "foo.zip"))
	require.NoError(t, err, "expected no error when adding attachment")

	// Ensure the attachment was added
	require.Len(t, msg.Attachments, 1, "expected one attachment")
	require.Equal(t, "foo.zip", msg.Attachments[0].Filename, "expected attachment to have correct filename")
	require.Equal(t, "application/zip", msg.Attachments[0].Type, "expected attachment to have correct type")
	require.Equal(t, "attachment", msg.Attachments[0].Disposition, "expected attachment to have correct disposition")

	// Ensure that we can decode the attachment data
	var data []byte
	data, err = base64.StdEncoding.DecodeString(msg.Attachments[0].Content)
	require.NoError(t, err, "expected no error when decoding attachment data")
	require.NotEmpty(t, data, "attachment has no data")
}

func TestAttachJSON(t *testing.T) {
	foo := map[string]string{"foo": "bar"}
	data, err := json.Marshal(foo)
	require.NoError(t, err, "expected no error when marshaling JSON")

	// Add an attachment to a new email
	msg := mail.NewV3Mail()
	err = emails.AttachJSON(msg, data, "foo.json")
	require.NoError(t, err, "expected no error when adding attachment")

	// Ensure the attachment was added
	require.Len(t, msg.Attachments, 1, "expected one attachment")
	require.Equal(t, "foo.json", msg.Attachments[0].Filename, "expected attachment to have correct filename")
	require.Equal(t, "application/json", msg.Attachments[0].Type, "expected attachment to have correct type")
	require.Equal(t, "attachment", msg.Attachments[0].Disposition, "expected attachment to have correct disposition")

	// Ensure that we can decode the attachment data
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(msg.Attachments[0].Content)
	require.NoError(t, err, "expected no error when decoding attachment data")
	actual := make(map[string]string)
	err = json.Unmarshal(decoded, &actual)
	require.NoError(t, err, "expected no error when unmarshaling JSON attachment")
	require.Equal(t, foo, actual, "expected JSON to match")
}
