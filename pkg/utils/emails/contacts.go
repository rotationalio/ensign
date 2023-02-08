package emails

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	sg "github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/sendgrid/sendgrid-go"
)

// AddContact adds a contact to the SendGrid marketing contacts list.
// TODO: Allow the user to specify list IDs using variadic arguments.
func (m *EmailManager) AddContact(contact *sg.Contact) error {
	if !m.conf.Enabled() {
		return errors.New("sendgrid is not enabled, cannot add contact")
	}

	// Create the SendGrid request
	var buf bytes.Buffer
	sgdata := &sg.AddContact{
		Contacts: []*sg.Contact{contact},
	}

	// TODO: What happens if no list IDs are specified?
	if m.conf.EnsignListID != "" {
		sgdata.ListIDs = []string{m.conf.EnsignListID}
	}

	if err := json.NewEncoder(&buf).Encode(sgdata); err != nil {
		return fmt.Errorf("could not encode json sendgrid contact data: %w", err)
	}

	// Execute the SendGrid request
	req := sendgrid.GetRequest(m.conf.APIKey, sg.Contacts, sg.Host)
	req.Method = http.MethodPut
	req.Body = buf.Bytes()

	rep, err := sendgrid.API(req)
	if err != nil {
		return fmt.Errorf("could not execute sendgrid api request: %w", err)
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return fmt.Errorf("received non-200 status code: %d", rep.StatusCode)
	}
	return nil
}
