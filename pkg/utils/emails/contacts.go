package emails

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
)

const (
	sgHost     = "https://api.sendgrid.com"
	sgContacts = "/v3/marketing/contacts"
)

type AddContact struct {
	ListIDs  []string   `json:"list_ids"`
	Contacts []*Contact `json:"contacts"`
}

type Contact struct {
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Email        string        `json:"email"`
	Country      string        `json:"country"`
	CustomFields *CustomFields `json:"custom_fields"`
}

// TODO: make custom fields request to get field IDs rather than hardcoding.
type CustomFields struct {
	Title                string `json:"e1_T"`
	Organization         string `json:"e2_T"`
	CloudServiceProvider string `json:"e3_T"`
}

// AddContactToSendGrid adds a contact to the SendGrid marketing contacts list.
func AddContactToSendGrid(conf Config, contact *Contact) error {
	if !conf.Enabled() {
		return errors.New("sendgrid is not enabled, cannot add contact")
	}

	// Create the SendGrid request
	var buf bytes.Buffer
	sgdata := &AddContact{
		Contacts: []*Contact{contact},
	}

	if conf.EnsignListID != "" {
		sgdata.ListIDs = []string{conf.EnsignListID}
	}

	if err := json.NewEncoder(&buf).Encode(sgdata); err != nil {
		return fmt.Errorf("could not encode json sendgrid contact data: %w", err)
	}

	// Execute the SendGrid request
	req := sendgrid.GetRequest(conf.APIKey, sgContacts, sgHost)
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
