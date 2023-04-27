package emails

import (
	"errors"
	"net/mail"

	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
)

// Configures SendGrid for sending emails and managing marketing contacts.
type Config struct {
	APIKey       string `split_words:"true" required:"false"`
	FromEmail    string `split_words:"true" default:"ensign@rotational.io"`
	AdminEmail   string `split_words:"true" default:"admins@rotational.io"`
	EnsignListID string `split_words:"true" required:"false"`
	Testing      bool   `split_words:"true" default:"false"`
	Archive      string `split_words:"true"`
}

// From and admin emails are required if the SendGrid API is enabled.
func (c Config) Validate() (err error) {
	if c.Enabled() {
		if c.AdminEmail == "" || c.FromEmail == "" {
			return errors.New("invalid configuration: admin and from emails are required if sendgrid is enabled")
		}

		if _, err = c.AdminContact(); err != nil {
			return errors.New("invalid configuration: admin email is unparsable")
		}

		if _, err = c.FromContact(); err != nil {
			return errors.New("invalid configuration: from email is unparsable")
		}

		if !c.Testing && c.Archive != "" {
			return errors.New("invalid configuration: email archiving is only supported in testing mode")
		}
	}

	return nil
}

// Returns true if there is a SendGrid API key available
func (c Config) Enabled() bool {
	return c.APIKey != ""
}

// Parses the FromEmail and returns a sendgrid contact for ease of mailing.
func (c Config) FromContact() (sendgrid.Contact, error) {
	return parseEmail(c.FromEmail)
}

// Parses the AdminEmail and returns a sendgrid contact for ease of mailing.
func (c Config) AdminContact() (sendgrid.Contact, error) {
	return parseEmail(c.AdminEmail)
}

func (c Config) MustFromContact() sendgrid.Contact {
	contact, err := c.FromContact()
	if err != nil {
		panic(err)
	}
	return contact
}

func (c Config) MustAdminContact() sendgrid.Contact {
	contact, err := c.AdminContact()
	if err != nil {
		panic(err)
	}
	return contact
}

func parseEmail(email string) (contact sendgrid.Contact, err error) {
	if email == "" {
		return contact, ErrUnparsable
	}

	var addr *mail.Address
	if addr, err = mail.ParseAddress(email); err != nil {
		return contact, ErrUnparsable
	}

	contact = sendgrid.Contact{
		Email: addr.Address,
	}
	contact.ParseName(addr.Name)

	return contact, nil
}
