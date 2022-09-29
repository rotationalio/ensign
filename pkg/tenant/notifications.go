package tenant

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
)

const (
	sgHost     = "https://api.sendgrid.com"
	sgContacts = "/v3/marketing/contacts"
)

type sgAddContact struct {
	ListIDs  []string     `json:"list_ids"`
	Contacts []*sgContact `json:"contacts"`
}

type sgContact struct {
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	Email        string          `json:"email"`
	Country      string          `json:"country"`
	CustomFields *sgCustomFields `json:"custom_fields"`
}

// TODO: make custom fields request to get field IDs rather than hardcoding.
type sgCustomFields struct {
	Title        string `json:"e1_T"`
	Organization string `json:"e2_T"`
}

// Signs up a contact to receive notifications from SendGrid by making a request to the
// SendGrid add contacts marketing API. The SendGrid API is asynchronous, which means
// that it doesn't return success if the user is registered. Instead a job ID is
// returned and the endpoint has to check if the registration was actually successful
// or not. To not block this endpoint, sign up doesn't check success but returns ok if
// the registration was correctly submitted.
//
// TODO: check for when the user is successfully signed up then send a welcome email.
// TODO: move all sendgrid-specific functionality into its own helper package.
// TODO: search for ensign list ID rather than configuring it externally.
func (s *Server) SignUp(c *gin.Context) {
	// Ensure SendGrid is enabled before making the request.
	if !s.conf.SendGrid.Enabled() {
		log.Error().Msg("sendgrid is not enabled: cannot register user for ensign notifications")
		c.JSON(http.StatusInternalServerError, "could not register user for notifications")
		return
	}

	// Parse the POST request from the user
	contact := &api.ContactInfo{}
	if err := c.BindJSON(contact); err != nil {
		log.Warn().Err(err).Msg("could not parse contact info")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Check the required fields are set
	if contact.FirstName == "" || contact.LastName == "" || contact.Email == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("first and last name and email are required"))
		return
	}

	// Create the SendGrid request
	var buf bytes.Buffer
	sgdata := &sgAddContact{
		Contacts: []*sgContact{
			{
				FirstName: contact.FirstName,
				LastName:  contact.LastName,
				Email:     contact.Email,
				Country:   contact.Country,
				CustomFields: &sgCustomFields{
					Title:        contact.Title,
					Organization: contact.Organization,
				},
			},
		},
	}

	if s.conf.SendGrid.EnsignListID != "" {
		sgdata.ListIDs = []string{s.conf.SendGrid.EnsignListID}
	}

	if err := json.NewEncoder(&buf).Encode(sgdata); err != nil {
		log.Error().Err(err).Msg("could not json encode sendgrid contact data")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not encode sendgrid contact data"))
		return
	}

	// Execute the SendGrid request
	req := sendgrid.GetRequest(s.conf.SendGrid.APIKey, sgContacts, sgHost)
	req.Method = http.MethodPut
	req.Body = buf.Bytes()

	rep, err := sendgrid.API(req)
	if err != nil {
		log.Error().Err(err).Msg("could not execute sendgrid api request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user notifications signup"))
		return
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		log.Error().Int("status_code", rep.StatusCode).Str("body", rep.Body).Msg("non-200 status returned from sendgrid api request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user notifications signup"))
		return
	}

	// Return 204 if the signup was successful.
	log.Info().Str("jobid", rep.Body).Msg("contact signed up for ensign private beta access")
	c.JSON(http.StatusNoContent, nil)
}
