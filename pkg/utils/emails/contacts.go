package emails

import (
	"errors"
	"fmt"

	sg "github.com/rotationalio/ensign/pkg/utils/sendgrid"
)

// AddContact adds a contact to SendGrid, adding them to the Ensign marketing list if
// it is configured. This is an upsert operation so existing contacts will be updated.
// The caller can optionally specify additional lists that the contact should be added
// to. If no lists are configured or specified, then the contact is added or updated in
// SendGrid but is not added to any marketing lists.
func (m *EmailManager) AddContact(contact *sg.Contact, listIDs ...string) (err error) {
	if !m.conf.Enabled() {
		return errors.New("sendgrid is not enabled, cannot add contact")
	}

	// Setup the request data
	sgdata := &sg.AddContactData{
		Contacts: []*sg.Contact{contact},
	}

	// Add the contact to the specified marketing lists
	if m.conf.EnsignListID != "" {
		sgdata.ListIDs = append(sgdata.ListIDs, m.conf.EnsignListID)
	}

	for _, listID := range listIDs {
		if listID != "" {
			sgdata.ListIDs = append(sgdata.ListIDs, listID)
		}
	}

	// Invoke the SendGrid API to add the contact
	if err = sg.AddContacts(m.conf.APIKey, sgdata); err != nil {
		return fmt.Errorf("could not add contact to sendgrid: %w", err)
	}
	return nil
}
