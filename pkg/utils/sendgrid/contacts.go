package sendgrid

import (
	"strings"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Contact contains contact information for a user, for adding users to marketing lists
// and sending emails.
type Contact struct {
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Email        string        `json:"email"`
	Country      string        `json:"country"`
	CustomFields *CustomFields `json:"custom_fields"`
}

// ParseName attempts to parse a full name into first and last names on the contact.
func (c *Contact) ParseName(name string) {
	parts := strings.Split(name, " ")
	c.FirstName = parts[0]
	if len(parts) > 1 {
		c.LastName = strings.Join(parts[1:], " ")
	} else {
		c.LastName = ""
	}
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
