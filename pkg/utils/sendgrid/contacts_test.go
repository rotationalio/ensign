package sendgrid_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/stretchr/testify/require"
)

func TestContact(t *testing.T) {
	// Test parsing the full name
	contact := sendgrid.Contact{}
	require.Equal(t, "", contact.FullName())

	contact.FirstName = "John"
	require.Equal(t, "John", contact.FullName())

	contact.LastName = "Doe"
	require.Equal(t, "John Doe", contact.FullName())

	contact.FirstName = ""
	require.Equal(t, "Doe", contact.FullName())

	// Test creating an email object from the contact
	contact = sendgrid.Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	email := contact.NewEmail()
	require.Equal(t, "John Doe", email.Name)
	require.Equal(t, "john@example.com", email.Address)
}
