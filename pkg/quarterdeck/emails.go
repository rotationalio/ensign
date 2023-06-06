package quarterdeck

import (
	"encoding/json"
	"fmt"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send an email to a user to verify their email address. This method requires the user
// object to have an existing verification token.
func (s *Server) SendVerificationEmail(user *models.User) (err error) {
	data := emails.VerifyEmailData{
		EmailData: emails.EmailData{
			Sender: s.conf.SendGrid.MustFromContact(),
			Recipient: sendgrid.Contact{
				Email: user.Email,
			},
		},
		FullName: user.Name,
	}
	data.Recipient.ParseName(user.Name)

	if data.VerifyURL, err = s.conf.EmailURL.VerifyURL(user.GetVerificationToken()); err != nil {
		return err
	}

	var msg *mail.SGMailV3
	if msg, err = emails.VerifyEmail(data); err != nil {
		return err
	}

	// Send the email
	return s.sendgrid.Send(msg)
}

// Send an email to a user to invite them to join an organization.
func (s *Server) SendInviteEmail(inviter *models.User, org *models.Organization, invite *models.UserInvitation) (err error) {
	data := emails.InviteData{
		EmailData: emails.EmailData{
			Sender: s.conf.SendGrid.MustFromContact(),
			Recipient: sendgrid.Contact{
				Email: invite.Email,
			},
		},
		Email:       invite.Email,
		InviterName: inviter.Name,
		OrgName:     org.Name,
		Role:        invite.Role,
	}

	if data.InviteURL, err = s.conf.EmailURL.InviteURL(invite.Token); err != nil {
		return err
	}

	var msg *mail.SGMailV3
	if msg, err = emails.InviteEmail(data); err != nil {
		return err
	}

	// Send the email
	return s.sendgrid.Send(msg)
}

// Send the daily users report to the Rotational admins.
// This method overwrites the email data on the report with the configured sender and
// recipient of the server so it should not be specified by the user (e.g. the user
// should only supply the report data for the email template).
func (s *Server) SendDailyUsers(data *emails.DailyUsersData) (err error) {
	data.EmailData = emails.EmailData{
		Sender:    s.conf.SendGrid.MustFromContact(),
		Recipient: s.conf.SendGrid.MustAdminContact(),
	}

	data.Domain = s.conf.Reporting.Domain
	data.EnsignDashboardLink = s.conf.Reporting.DashboardURL

	var msg *mail.SGMailV3
	if msg, err = emails.DailyUsersEmail(*data); err != nil {
		return err
	}

	// Attach the report as json
	var attachment []byte
	if attachment, err = json.MarshalIndent(data, "", " "); err != nil {
		return err
	}

	if err = emails.AttachJSON(msg, attachment, fmt.Sprintf("daily_users_%s.json", data.Date.Format("20060102"))); err != nil {
		return err
	}

	// Attach the new accounts CSV if there are any available
	if len(data.NewAccounts) > 0 {
		var accounts []byte
		if accounts, err = data.NewAccountsCSV(); err != nil {
			return err
		}

		if err = emails.AttachCSV(msg, accounts, fmt.Sprintf("new_accounts_%s.json", data.Date.Format("20060102"))); err != nil {
			return err
		}
	}

	// Send the email
	return s.sendgrid.Send(msg)
}
