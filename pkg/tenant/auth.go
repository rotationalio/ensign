package tenant

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/radish"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog/log"
)

// Lifetimes for CSRF token cookies that the server generates
const (
	protectLoginCSRFLifetime = time.Minute * 10
	authCSRFLifetime         = time.Hour * 12 // Should be longer than the access token lifetime
)

// Register is a publicly accessible endpoint that allows new users to create an
// account via Quarterdeck by providing an email address and password.
//
// Route: POST /v1/register
func (s *Server) Register(c *gin.Context) {
	var err error
	ctx := c.Request.Context()

	// Parse the request body
	params := &api.RegisterRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse register request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// Filter bad requests before they reach Quarterdeck
	// Note: This is a simple check to ensure that all required fields are present.
	if err = params.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// Make the register request to Quarterdeck
	req := &qd.RegisterRequest{
		Email:        params.Email,
		Password:     params.Password,
		PwCheck:      params.PwCheck,
		AgreeToS:     params.AgreeToS,
		AgreePrivacy: params.AgreePrivacy,
	}

	var reply *qd.RegisterReply
	if reply, err = s.quarterdeck.Register(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create member model for the new user
	member := &db.Member{
		ID:           reply.ID,
		OrgID:        reply.OrgID,
		Email:        reply.Email,
		Organization: reply.OrgName,
		Workspace:    reply.OrgDomain,
		Role:         reply.Role,
		JoinedAt:     time.Now(),
	}

	// Create a default tenant and member for the new user
	// Note: This task will error if the member model is invalid
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return db.CreateUserResources(ctx, member)
	}), radish.WithRetries(3),
		radish.WithBackoff(backoff.NewExponentialBackOff()),
		radish.WithErrorf("could not create default tenant and member for new user %s", reply.ID.String()),
		radish.WithContext(sentry.CloneContext(c)),
	)

	// Add to SendGrid Ensign Marketing list in go routine
	// TODO: use worker queue to limit number of go routines for tasks like this
	// TODO: test in live integration tests to make sure this works
	hub := sentrygin.GetHubFromContext(c).Clone()
	go func() {
		// TODO: We don't have the name of the user here, so we would have to add it
		// elsewhere.
		contact := &sendgrid.Contact{
			Email: params.Email,
		}

		if err := s.sendgrid.AddContact(contact); err != nil {
			log.Error().Err(err).Msg("could not add newly registered user to sendgrid ensign marketing list")
			if hub != nil {
				hub.CaptureException(err)
			}
		}
	}()

	// Return the response from Quarterdeck
	c.Status(http.StatusNoContent)
}

// Login is a publically accessible endpoint that allows users to login into their
// account via Quarterdeck and receive access and refresh tokens for future requests.
// Access and refresh tokens are set in the cookies for the convenience of frontends.
//
// Route: POST /v1/login
func (s *Server) Login(c *gin.Context) {
	var err error

	// Parse the request body
	params := &api.LoginRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse login request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// Validate that required fields were provided
	if params.Email == "" || params.Password == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
		return
	}

	// TODO: Add validation method for login request.
	if params.InviteToken != "" && params.OrgID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot provide both invite token and org id"))
		return
	}

	// Make the login request to Quarterdeck
	req := &qd.LoginRequest{
		Email:       params.Email,
		Password:    params.Password,
		InviteToken: params.InviteToken,
	}

	if params.OrgID != "" {
		if req.OrgID, err = ulid.Parse(params.OrgID); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse orgID from the request")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid org id"))
			return
		}
	}

	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Login(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Update the status of a member with an invite token to Confirmed.
	if params.InviteToken != "" {
		// Parse access token to get the orgID.
		var claims *tokens.Claims
		if claims, err = tokens.ParseUnverifiedTokenClaims(reply.AccessToken); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse access token from the claims")
			c.JSON(http.StatusUnauthorized, api.ErrorResponse(responses.ErrTryLoginAgain))
			return
		}

		var orgID ulid.ULID
		if orgID, err = ulid.Parse(claims.OrgID); err != nil {
			sentry.Error(c).Err(err).Msg("could not parse orgID from access token")
			c.JSON(http.StatusUnauthorized, api.ErrorResponse(responses.ErrTryLoginAgain))
			return
		}

		// Get member from the database by their email.
		var member *db.Member
		if member, err = db.GetMemberByEmail(c, orgID, params.Email); err != nil {
			sentry.Error(c).Str("orgID", orgID.String()).Err(err).Msg("could not get member from the database")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryLoginAgain))
			return
		}

		// Verify the member ID matches the ID in the claims. If they do not match delete the
		// member from the database and create a new member record in with the ID from the claims.
		if claims.Subject != member.ID.String() {
			if err = db.DeleteMember(c, orgID, member.ID); err != nil {
				sentry.Error(c).Err(err).Msg("could not delete member from the database")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
				return
			}

			if member.ID, err = ulid.Parse(claims.Subject); err != nil {
				sentry.Error(c).Err(err).Msg("could not claims subject")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
				return
			}

			if err = db.CreateMember(c, member); err != nil {
				sentry.Error(c).Err(err).Msg("could not recreate member record for invited user")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
				return
			}
		}

		// User has accepted the invite
		member.JoinedAt = time.Now()
		member.Name = claims.Name
		if err = db.UpdateMember(c, member); err != nil {
			sentry.Error(c).Err(err).Msg("could not update member in the database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
			return
		}
	}

	// Set the access and refresh tokens as cookies for the front-end
	if err = middleware.SetAuthCookies(c, reply.AccessToken, reply.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
		sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err = middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Update last login time for the member record in a background task
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return db.UpdateLastLogin(ctx, reply.AccessToken, time.Now())
	}), radish.WithErrorf("could not update last login for user after login"),
		radish.WithContext(sentry.CloneContext(c)),
	)

	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		LastLogin:    reply.LastLogin,
	}
	c.JSON(http.StatusOK, out)
}

// Refresh is a publicly accessible endpoint that allows users to refresh their
// access token using their refresh token. This enables frontend clients to provide a
// seamless login experience for the user.
//
// Route: POST /v1/refresh
func (s *Server) Refresh(c *gin.Context) {
	var err error

	// Parse the request body
	params := &api.RefreshRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse refresh request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// Validate that required fields were provided
	if params.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrLogBackIn))
		return
	}

	// Construct the refresh request to Quarterdeck
	req := &qd.RefreshRequest{
		RefreshToken: params.RefreshToken,
	}

	// Add the orgID if provided
	if params.OrgID != "" {
		if req.OrgID, err = ulid.Parse(params.OrgID); err != nil {
			sentry.Warn(c).Err(err).Msg("could not parse orgID from the request")
			c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid org_id"))
			return
		}
	}

	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Refresh(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Set the access and refresh tokens as cookies for the front-end
	if err = middleware.SetAuthCookies(c, reply.AccessToken, reply.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
		sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Protect the frontend from CSRF attacks by setting the double cookie tokens
	expiresAt := time.Now().Add(authCSRFLifetime)
	if err = middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Set the access and refresh tokens as cookies for the frontend
	if err = middleware.SetAuthCookies(c, reply.AccessToken, reply.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
		sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Update last login time for the member record in a background task
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return db.UpdateLastLogin(ctx, reply.AccessToken, time.Now())
	}), radish.WithErrorf("could not update last login for user after refresh"),
		radish.WithContext(sentry.CloneContext(c)),
	)

	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		LastLogin:    reply.LastLogin,
	}
	c.JSON(http.StatusOK, out)
}

// Switch is an authenticated endpoint that allows human users to switch between
// organizations that they are a member of. This exists to allow users to fetch new
// access and refresh tokens without having to re-enter their credentials. This
// endpoint is not available to machine users with API key credentials, since API keys
// can only exist in one project in one organization. If the user is already
// authenticated with the requested organization, this endpoint returns an error. The
// refresh endpoint should be used if the access token simply needs to be refreshed.
func (s *Server) Switch(c *gin.Context) {
	var (
		ctx context.Context
		err error
	)

	// Context with the credentials is required for the Quarterdeck request
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve credentials from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("missing user credentials"))
		return
	}

	// Fetch the user claims
	var claims *tokens.Claims
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Warn(c).Err(err).Msg("could not retrieve claims from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("missing user credentials"))
		return
	}

	// Parse the request body
	params := &api.SwitchRequest{}
	if err = c.BindJSON(params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse switch request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Validate that required fields were provided
	if params.OrgID == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("missing org_id in request"))
		return
	}

	// Parse the orgID from the request
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(params.OrgID); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse orgID from the request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid org_id in request"))
		return
	}

	// Prevent switching to the same organization
	if orgID.String() == claims.OrgID {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("already logged in to this organization"))
		return
	}

	// Construct the switch request to Quarterdeck
	req := &qd.SwitchRequest{
		OrgID: orgID,
	}

	var reply *qd.LoginReply
	if reply, err = s.quarterdeck.Switch(ctx, req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Set the access and refresh tokens as cookies for the front-end
	if err = middleware.SetAuthCookies(c, reply.AccessToken, reply.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
		sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set auth cookies"))
		return
	}

	// Update last login time for the member record in a background task
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return db.UpdateLastLogin(ctx, reply.AccessToken, time.Now())
	}), radish.WithErrorf("could not update last login for user after switch"),
		radish.WithContext(sentry.CloneContext(c)),
	)

	out := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
		LastLogin:    reply.LastLogin,
	}
	c.JSON(http.StatusOK, out)
}

// ProtectLogin prepares the front-end for login by setting the double cookie
// tokens for CSRF protection.
func (s *Server) ProtectLogin(c *gin.Context) {
	expiresAt := time.Now().Add(protectLoginCSRFLifetime)
	if err := middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set csrf login protection cookies")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}
	c.JSON(http.StatusOK, &api.Reply{Success: true})
}

// VerifyEmail is a publicly accessible endpoint that allows users to verify their
// email address by supplying a token that was sent to their email address. If the
// token has already been verified, this endpoint returns a 202 Accepted response.
//
// Route: POST /v1/verify
func (s *Server) VerifyEmail(c *gin.Context) {
	var (
		err    error
		params *api.VerifyRequest
	)

	// Parse the request body
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse verify email request")
		// TODO: What action can the user take if their attempt to verify their email fails?
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrVerificationFailed))
		return
	}

	// Validate that required fields were provided
	if params.Token == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrVerificationFailed))
		return
	}

	req := &qd.VerifyRequest{
		Token: params.Token,
	}

	// Parse the orgID if provided
	if params.OrgID != "" {
		if req.OrgID, err = ulid.Parse(params.OrgID); err != nil {
			sentry.Warn(c).Str("org_id", params.OrgID).Err(err).Msg("could not parse orgID from the request")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrVerificationFailed))
			return
		}
	}

	// Make the verify request to Quarterdeck
	var rep *qd.LoginReply
	if rep, err = s.quarterdeck.VerifyEmail(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// If we got credentials from Quarterdeck, set them in the cookies
	if rep != nil && rep.AccessToken != "" && rep.RefreshToken != "" {
		if err = middleware.SetAuthCookies(c, rep.AccessToken, rep.RefreshToken, s.conf.Auth.CookieDomain); err != nil {
			sentry.Error(c).Err(err).Msg("could not set access and refresh token cookies")
			c.Status(http.StatusNoContent)
			return
		}

		if err = middleware.SetDoubleCookieToken(c, s.conf.Auth.CookieDomain, time.Now().Add(authCSRFLifetime)); err != nil {
			sentry.Error(c).Err(err).Msg("could not set csrf protection cookies")
			c.Status(http.StatusNoContent)
			return
		}

		// Return the credentials in the reply
		out := &api.AuthReply{
			AccessToken:  rep.AccessToken,
			RefreshToken: rep.RefreshToken,
		}
		c.JSON(http.StatusOK, out)
		return
	}

	// Return 204 if already verified
	c.Status(http.StatusNoContent)
}

// ResendEmail is a publicly accessible endpoint that allows users to resend emails to
// the email address in the POST request by forwarding the request to Quarterdeck. If
// the email address belongs to a user who has not been verified then this endpoint
// will send a new verification email by forwarding the request to Quarterdeck. If
// there is an orgID in the request and the user is invited to that organization but
// has not accepted the invite then the invitation email is resent. Because this is an
// unauthenticated endpoint, it always returns a 204 No Content response to prevent
// revealing information about registered email addresses and users.
//
// Route: POST /v1/resend
func (s *Server) ResendEmail(c *gin.Context) {
	var (
		err    error
		params *api.ResendRequest
	)

	// Parse the request body
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse resend email request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrTryResendAgain))
		return
	}

	// Email is always required for this endpoint
	req := &qd.ResendRequest{}
	req.Email = strings.TrimSpace(params.Email)
	if req.Email == "" || len(req.Email) > 254 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrBadResendRequest))
		return
	}

	// Parse the orgID if provided
	if params.OrgID != "" {
		if req.OrgID, err = ulid.Parse(params.OrgID); err != nil {
			sentry.Warn(c).Str("org_id", params.OrgID).Err(err).Msg("could not parse orgID from the request")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrBadResendRequest))
			return
		}
	}

	// Make the resend request to Quarterdeck
	// Note: We are relying on Quarterdeck to adhere to best security practices and
	// only return an error if the request is not parseable. Otherwise, it should
	// return 204.
	if err = s.quarterdeck.ResendEmail(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ForgotPassword is a publicly accessible endpoint that allows users to request a
// password reset by forwarding a POST request with an email address to Quarterdeck. If
// the email exists in the database then an email is sent to the user with a password
// reset link. This endpoint always returns a 204 No Content response to prevent
// revealing information about registered email addresses and users.
//
// Route: POST /v1/forgot-password
func (s *Server) ForgotPassword(c *gin.Context) {
	var (
		err    error
		params *api.ForgotPasswordRequest
	)

	// Parse the request body
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse forgot password request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrSendPasswordResetFailed))
		return
	}

	// Email is always required for this endpoint
	params.Email = strings.TrimSpace(params.Email)
	if params.Email == "" || len(params.Email) > 254 {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrInvalidEmail))
		return
	}

	// Make the forgot password request to Quarterdeck
	// Note: We are relying on Quarterdeck to adhere to best security practices and
	// only return an error if the request is not parseable. Otherwise, it should
	// return 204.
	req := &qd.ForgotPasswordRequest{
		Email: params.Email,
	}
	if err = s.quarterdeck.ForgotPassword(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ResetPassword is a publicly accessible endpoint that allows users to reset their
// password by forwarding a POST request with a reset token and a new password to
// Quarterdeck. If the password reset was successful then this endpoint returns a
// confirmation email to the user and a 204 No Content response.
//
// Route: POST /v1/reset-password
func (s *Server) ResetPassword(c *gin.Context) {
	var (
		err    error
		params *api.ResetPasswordRequest
	)

	// Parse the request body
	if err = c.BindJSON(&params); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse reset password request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(responses.ErrPasswordResetFailed))
		return
	}

	// Validate that required fields were provided
	if err = params.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make the reset password request to Quarterdeck
	req := &qd.ResetPasswordRequest{
		Token:    params.Token,
		Password: params.Password,
		PwCheck:  params.PwCheck,
	}
	if err = s.quarterdeck.ResetPassword(c.Request.Context(), req); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Clear the authentication cookies to log out the user
	middleware.ClearAuthCookies(c, s.conf.Auth.CookieDomain)

	// Return 204 for the response
	c.Status(http.StatusNoContent)
}
