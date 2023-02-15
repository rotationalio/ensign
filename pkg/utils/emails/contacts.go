package emails

import (
	"errors"
	"fmt"

	sg "github.com/rotationalio/ensign/pkg/utils/sendgrid"
)

// AddContact adds a contact to the SendGrid marketing contacts list.
// TODO: Allow the user to specify list IDs using variadic arguments.
func (m *EmailManager) AddContact(contact *sg.Contact) (err error) {
	if !m.conf.Enabled() {
		return errors.New("sendgrid is not enabled, cannot add contact")
	}

	// Setup the request data
	sgdata := &sg.AddContactData{
		Contacts: []*sg.Contact{contact},
	}

	// TODO: What happens if no list IDs are specified?
	if m.conf.EnsignListID != "" {
		sgdata.ListIDs = []string{m.conf.EnsignListID}
	}

	// Invoke the SendGrid API to add the contact
	if err = sg.AddContacts(m.conf.APIKey, sgdata); err != nil {
		return fmt.Errorf("could not add contact to sendgrid: %w", err)
	}
	return nil
}
