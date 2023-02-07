package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rs/zerolog/log"
)

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
	params := &api.ContactInfo{}
	if err := c.BindJSON(params); err != nil {
		log.Warn().Err(err).Msg("could not parse contact info")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Check the required fields are set
	if params.FirstName == "" || params.LastName == "" || params.Email == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("first and last name and email are required"))
		return
	}

	// Add the contact to SendGrid
	contact := &emails.Contact{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Country:   params.Country,
		CustomFields: &emails.CustomFields{
			Title:                params.Title,
			Organization:         params.Organization,
			CloudServiceProvider: params.CloudServiceProvider,
		},
	}
	if err := s.sendgrid.AddContact(contact); err != nil {
		log.Error().Err(err).Msg("could not add contact to sendgrid")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add contact to sendgrid"))
		return
	}

	// Return 204 if the signup was successful.
	log.Info().Msg("contact signed up for ensign private beta access")
	c.JSON(http.StatusNoContent, nil)
}
