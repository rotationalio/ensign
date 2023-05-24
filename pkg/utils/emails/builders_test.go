package emails_test

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

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

	verifyEmail := emails.VerifyEmailData{
		EmailData: data,
		FullName:  "Rachel Lendt",
		VerifyURL: "https://rotational.app/verify?token=1234567890",
	}
	mail, err = emails.VerifyEmail(verifyEmail)
	require.NoError(t, err, "expected no error when building verify email")
	require.Equal(t, emails.VerifyEmailRE, mail.Subject, "expected verify email subject to match")
	generateMIME(t, mail, "verify_email.mime")

	inviteData := emails.InviteData{
		EmailData:   data,
		Email:       "rachel@example.com",
		InviterName: "Lewis Hudson",
		OrgName:     "Events R Us",
		Role:        "Member",
		InviteURL:   "https://rotational.app/invite?token=1234567890",
	}
	mail, err = emails.InviteEmail(inviteData)
	require.NoError(t, err, "expected no error when building invite email")
	require.Equal(t, fmt.Sprintf(emails.InviteRE, "Lewis Hudson"), mail.Subject, "expected invite email subject to match")
	generateMIME(t, mail, "invite.mime")

	dailyUsersData := emails.DailyUsersData{
		EmailData:           data,
		Date:                time.Date(2023, 4, 7, 0, 0, 0, 0, time.UTC),
		InactiveDate:        time.Date(2023, 3, 8, 0, 0, 0, 0, time.UTC),
		Domain:              "ensign.local",
		EnsignDashboardLink: "http://grafana.ensign.local/dashboards/ensign",
		NewUsers:            2,
		DailyUsers:          8,
		ActiveUsers:         102,
		InactiveUsers:       3,
		APIKeys:             58,
		ActiveKeys:          52,
		InactiveKeys:        6,
		RevokedKeys:         12,
		Organizations:       87,
		NewOrganizations:    1,
		Projects:            87,
		NewProjects:         1,
	}
	mail, err = emails.DailyUsersEmail(dailyUsersData)
	require.NoError(t, err, "expected no error when building daily users email")
	require.Equal(t, fmt.Sprintf(emails.DailyUsersRE, "ensign.local", "April 7, 2023"), mail.Subject, "expected daily users email subject to be dynamic")
	generateMIME(t, mail, "dailyusers.mime")
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

func ExampleDailyUsersData_TabTable() {
	dailyUsersData := emails.DailyUsersData{
		Domain:           "ensign.local",
		NewUsers:         2,
		DailyUsers:       8,
		ActiveUsers:      102,
		InactiveUsers:    3,
		APIKeys:          58,
		ActiveKeys:       52,
		InactiveKeys:     6,
		RevokedKeys:      12,
		Organizations:    87,
		NewOrganizations: 1,
		Projects:         87,
		NewProjects:      1,
	}
	fmt.Println(dailyUsersData.TabTable())
	// Output:
	// New Users:          2    Daily Users:        8
	// Active Users:       102  Inactive Users:     3
	// API Keys:           58   Revoked API Keys:   12
	// Active API Keys:    52   Inactive API Keys:  6
	// New Organizations:  1    Organizations:      87
	// New Projects:       1    Projects:           87
}

func ExampleDailyUsersData_NewAccountsCSV() {
	data := emails.DailyUsersData{
		NewAccounts: []*emails.NewAccountData{
			{
				Name:          "Wiley E. Coyote",
				Email:         "wiley@acme.co",
				EmailVerified: true,
				Role:          "owner",
				LastLogin:     time.Date(2023, 7, 8, 19, 21, 39, 0, time.UTC),
				Created:       time.Date(2023, 7, 8, 12, 2, 52, 0, time.UTC),
				Organization:  "Acme, Inc.",
				Domain:        "acme.co",
				Projects:      3,
				APIKeys:       7,
				Invitations:   3,
				Users:         2,
			},
			{
				Name:          "Rod P. Runner",
				Email:         "rod@acme.co",
				EmailVerified: false,
				Role:          "member",
				LastLogin:     time.Date(2023, 7, 8, 13, 12, 42, 0, time.UTC),
				Created:       time.Date(2023, 7, 8, 12, 2, 52, 0, time.UTC),
				Organization:  "Acme, Inc.",
				Domain:        "acme.co",
				Projects:      3,
				APIKeys:       7,
				Invitations:   3,
				Users:         2,
			},
			{
				Name:          "Julie Smith Lee",
				Email:         "jlee@foundations.io",
				EmailVerified: true,
				Role:          "owner",
				LastLogin:     time.Date(2023, 7, 8, 8, 22, 27, 0, time.UTC),
				Created:       time.Date(2023, 7, 8, 8, 21, 1, 0, time.UTC),
				Organization:  "Foundations",
				Domain:        "foundations.io",
				Projects:      1,
				APIKeys:       1,
				Invitations:   8,
				Users:         1,
			},
		},
	}

	csv, _ := data.NewAccountsCSV()
	fmt.Println(string(csv))
	// Output:
	// name,email,email_verified,role,last_login,created,organization,domain,projects,apikeys,users,invitations
	// Wiley E. Coyote,wiley@acme.co,true,owner,2023-07-08T19:21:39Z,2023-07-08T12:02:52Z,"Acme, Inc.",acme.co,3,7,2,3
	// Rod P. Runner,rod@acme.co,false,member,2023-07-08T13:12:42Z,2023-07-08T12:02:52Z,"Acme, Inc.",acme.co,3,7,2,3
	// Julie Smith Lee,jlee@foundations.io,true,owner,2023-07-08T08:22:27Z,2023-07-08T08:21:01Z,Foundations,foundations.io,1,1,1,8
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

func TestAttachCSV(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 10))
	w := csv.NewWriter(buf)
	w.Write([]string{"foo", "bar"})
	w.Flush()
	data := buf.Bytes()

	// Add an attachment to a new email
	msg := mail.NewV3Mail()
	err := emails.AttachCSV(msg, data, "foo.csv")
	require.NoError(t, err, "expected no error when adding attachment")

	// Ensure the attachment was added
	require.Len(t, msg.Attachments, 1, "expected one attachment")
	require.Equal(t, "foo.csv", msg.Attachments[0].Filename, "expected attachment to have correct filename")
	require.Equal(t, "text/csv", msg.Attachments[0].Type, "expected attachment to have correct type")
	require.Equal(t, "attachment", msg.Attachments[0].Disposition, "expected attachment to have correct disposition")

	// Ensure that we can decode the attachment data
	var decoded []byte
	decoded, err = base64.StdEncoding.DecodeString(msg.Attachments[0].Content)
	require.NoError(t, err, "expected no error when decoding attachment data")

	r := csv.NewReader(bytes.NewReader(decoded))
	actual, err := r.ReadAll()
	require.NoError(t, err, "expected no error when reading CSV attachment")
	require.Len(t, actual, 1)
	require.Equal(t, actual[0], []string{"foo", "bar"})
}
