package quarterdeck

import (
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
			Sender: sendgrid.Contact{
				Email: s.conf.SendGrid.FromEmail,
			},
			Recipient: sendgrid.Contact{
				Email: user.Email,
			},
		},
		FullName: user.Name,
	}

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
			Sender: sendgrid.Contact{
				Email: s.conf.SendGrid.FromEmail,
			},
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
