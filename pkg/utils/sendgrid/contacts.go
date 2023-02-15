package sendgrid

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Contact struct {
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Email        string        `json:"email"`
	Country      string        `json:"country"`
	CustomFields *CustomFields `json:"custom_fields"`
}

// FullName attempts to construct the contact's full name from existing name fields.
func (c Contact) FullName() string {
	switch {
	case c.FirstName == "" && c.LastName == "":
		return ""
	case c.FirstName != "" && c.LastName == "":
		return c.FirstName
	case c.FirstName == "" && c.LastName != "":
		return c.LastName
	default:
		return c.FirstName + " " + c.LastName
	}
}

// NewEmail returns the sendgrid email object for constructing emails.
func (c Contact) NewEmail() *mail.Email {
	return mail.NewEmail(c.FullName(), c.Email)
}

// TODO: make custom fields request to get field IDs rather than hardcoding.
type CustomFields struct {
	Title                string `json:"e1_T"`
	Organization         string `json:"e2_T"`
	CloudServiceProvider string `json:"e3_T"`
}
