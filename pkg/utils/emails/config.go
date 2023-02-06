package emails

import "errors"

// Configures SendGrid for sending emails and managing marketing contacts.
type Config struct {
	APIKey       string `split_words:"true" required:"false"`
	FromEmail    string `split_words:"true" default:"ensign@rotational.io"`
	AdminEmail   string `split_words:"true" default:"admins@rotational.io"`
	EnsignListID string `split_words:"true" required:"false"`
}

// From and admin emails are required if the SendGrid API is enabled.
func (c Config) Validate() (err error) {
	if c.Enabled() {
		if c.AdminEmail == "" || c.FromEmail == "" {
			return errors.New("invalid configuration: admin and from emails are required if sendgrid is enabled")
		}
	}
	return nil
}

// Returns true if there is a SendGrid API key available
func (c Config) Enabled() bool {
	return c.APIKey != ""
}
