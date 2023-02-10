package emails_test

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
)

func TestEmailBuilders(t *testing.T) {
	setupMIMEDir(t)

	sender := sendgrid.Contact{
		FirstName: "Lewis",
		LastName:  "Hudson",
		Email:     "lewis@example.com",
	}
	recipient := sendgrid.Contact{
		FirstName: "Rachel",
		LastName:  "Lendt",
		Email:     "rachel@example.com",
	}
	data := emails.EmailData{
		Sender:    sender,
		Recipient: recipient,
	}

	welcomeData := emails.WelcomeData{
		EmailData:    data,
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

func TestEmailData(t *testing.T) {
	sender := sendgrid.Contact{
		FirstName: "Lewis",
		LastName:  "Hudson",
		Email:     "lewis@example.com",
	}
	recipient := sendgrid.Contact{
		FirstName: "Rachel",
		LastName:  "Lendt",
		Email:     "rachel@example.com",
	}
	data := emails.EmailData{
		Sender:    sender,
		Recipient: recipient,
	}

	// Email is not valid without a subject
	require.EqualError(t, data.Validate(), emails.ErrMissingSubject.Error(), "email subject should be required")

	// Email is not valid without a sender
	data.Subject = "Subject Line"
	data.Sender.Email = ""
	require.EqualError(t, data.Validate(), emails.ErrMissingSender.Error(), "email sender should be required")

	// Email is not valid without a recipient
	data.Sender.Email = sender.Email
	data.Recipient.Email = ""
	require.EqualError(t, data.Validate(), emails.ErrMissingRecipient.Error(), "email recipient should be required")

	// Successful validation
	data.Recipient.Email = recipient.Email
	require.NoError(t, data.Validate(), "expected no error when validating email data")
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
