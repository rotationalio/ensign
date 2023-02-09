package emails

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// New email manager with the specified configuration.
func New(conf Config) (m *EmailManager, err error) {
	m = &EmailManager{conf: conf}
	if conf.Testing {
		log.Warn().Str("storage", conf.Archive).Msg("using mock email client")
		m.client = &mock.SendGridClient{
			Storage: conf.Archive,
		}
	} else {
		if conf.APIKey == "" {
			return nil, errors.New("cannot create email client without API key")
		}
		m.client = sendgrid.NewSendClient(conf.APIKey)
	}

	// Parse the from and admin emails from the configuration
	if m.fromEmail, err = mail.ParseAddress(conf.FromEmail); err != nil {
		return nil, fmt.Errorf("could not parse 'from' email %q: %s", conf.FromEmail, err)
	}

	if m.adminsEmail, err = mail.ParseAddress(conf.AdminEmail); err != nil {
		return nil, fmt.Errorf("could not parse admin email %q: %s", conf.AdminEmail, err)
	}

	return m, nil
}

// EmailManager allows a server to send rich emails using the SendGrid service.
type EmailManager struct {
	conf        Config
	client      SendGridClient
	fromEmail   *mail.Address
	adminsEmail *mail.Address
}

func (m *EmailManager) Send(message *sgmail.SGMailV3) (err error) {
	var rep *rest.Response
	if rep, err = m.client.Send(message); err != nil {
		return err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return errors.New(rep.Body)
	}

	return nil
}
