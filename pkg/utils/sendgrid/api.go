package sendgrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const (
	Host       = "https://api.sendgrid.com"
	ContactsEP = "/v3/marketing/contacts"
	ListsEP    = "/v3/marketing/lists"
	FieldsEP   = "/v3/marketing/field_definitions"
)

type AddContactData struct {
	ListIDs  []string   `json:"list_ids"`
	Contacts []*Contact `json:"contacts"`
}

// Add contacts to SendGrid marketing lists.
func AddContacts(apiKey string, data *AddContactData) (err error) {
	// Include the contacts and list IDs in the request body.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return fmt.Errorf("could not encode json sendgrid contact data: %w", err)
	}

	// Create the PUT request
	req := sendgrid.GetRequest(apiKey, ContactsEP, Host)
	req.Method = http.MethodPut
	req.Body = buf.Bytes()

	// Do the request
	_, err = doRequest(req)
	return err
}

// Fetch lists of contacts from SendGrid.
func MarketingLists(apiKey, pageToken string) (string, error) {
	params := map[string]string{
		"page_size": "100",
	}

	if pageToken != "" {
		params["page_token"] = pageToken
	}

	// Create the GET request
	req := sendgrid.GetRequest(apiKey, ListsEP, Host)
	req.Method = http.MethodGet
	req.QueryParams = params

	return doRequest(req)
}

// Fetch field definitions from SendGrid.
func FieldDefinitions(apiKey string) (string, error) {
	req := sendgrid.GetRequest(apiKey, FieldsEP, Host)
	req.Method = http.MethodGet
	return doRequest(req)
}

// Helper to perform a SendGrid request, handling errors and returning the response
// body.
func doRequest(req rest.Request) (_ string, err error) {
	rep, err := sendgrid.MakeRequest(req)
	if err != nil {
		return "", fmt.Errorf("could not make sendgrid request: %w", err)
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return "", fmt.Errorf("received status code %d: could not make sendgrid request", rep.StatusCode)
	}

	return rep.Body, nil
}
