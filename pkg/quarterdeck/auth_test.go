package quarterdeck_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *quarterdeckTestSuite) TestRegister() {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		s.T().Skip("skipping long running test in short mode")
	}

	defer s.ResetDatabase()
	defer mock.Reset()
	defer s.ResetTasks()
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	project := "01GKHJRF01YXHZ51YMMKV3RCMK"
	projectID := ulid.MustParse(project)
	req := &api.RegisterRequest{
		Name:         "Rachel Johnson",
		Email:        "rachel@example.com",
		Password:     "supers3cretSquirrel?",
		PwCheck:      "supers3cretSquirrel?",
		Organization: "Financial Services Ltd",
		Domain:       "financial-services",
		AgreeToS:     true,
		AgreePrivacy: true,
	}

	sent := time.Now()
	rep, err := s.client.Register(ctx, req)
	require.NoError(err, "unable to create user from valid request")

	require.False(ulids.IsZero(rep.ID), "did not get a user ID back from the database")
	require.False(ulids.IsZero(rep.OrgID), "did not get back an orgID from the database")
	require.Equal(req.Email, rep.Email)
	require.Equal("Welcome to Ensign!", rep.Message)
	require.Equal(rep.Role, permissions.RoleOwner)
	require.NotEmpty(rep.Created, "did not get a created timestamp back")

	// Test that the user actually made it into the database
	user, err := models.GetUser(context.Background(), rep.ID, rep.OrgID)
	require.NoError(err, "could not get user from database")
	require.Equal(rep.Email, user.Email, "user creation check failed")

	// Test that the verification fields were set on the user
	require.False(user.EmailVerified, "user should not be verified")
	require.NotEmpty(user.GetVerificationToken(), "user should have a verification token")
	require.True(user.EmailVerificationExpires.Valid, "user should have an email verification expiration")
	expiresAt, err := time.Parse(time.RFC3339Nano, user.EmailVerificationExpires.String)
	require.NoError(err, "could not parse email verification expiration")
	require.True(expiresAt.After(sent), "email verification expiration should be after the email was sent")
	require.NotEmpty(user.EmailVerificationSecret, "user should have an email verification secret")

	// Test with a project ID provided
	req.Email = "jane@example.com"
	req.Domain = "it-services"
	req.ProjectID = project
	rep, err = s.client.Register(ctx, req)
	require.NoError(err, "unable to create user from valid request")

	// Test that the user made it into the database
	user, err = models.GetUser(context.Background(), rep.ID, rep.OrgID)
	require.NoError(err, "could not get user from database")
	require.Equal(rep.Email, user.Email, "user creation check failed")

	// Test that the organization project link was created in the database
	op := &models.OrganizationProject{
		OrgID:     rep.OrgID,
		ProjectID: projectID,
	}
	ok, err := op.Exists(context.Background())
	require.NoError(err, "could not check if organization project link exists")
	require.True(ok, "organization project link was not created")

	// Test when email in token does not match request
	req.ProjectID = ""
	req.Organization = ""
	req.Domain = ""
	token := "s6jsNBizyGh_C_ZsUSuJsquONYa-KH_2cmoJZd-jnIk"
	req.InviteToken = token
	req.Email = "wrong@example.com"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "invalid invitation")

	// Test invite token exists but is expired
	req.InviteToken = "s6jsNBizyGh_C_ZsUSuJsquONYa--gpcfzorN8DsdjIA"
	req.Email = "eefrank@checkers.io"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "invalid invitation")

	// Test with a valid invite token provided
	req.InviteToken = token
	req.Email = "joe@checkers.io"
	rep, err = s.client.Register(ctx, req)
	require.NoError(err, "unable to create invited user from valid request")

	// Test that the user made it into the database
	user, err = models.GetUser(context.Background(), rep.ID, rep.OrgID)
	require.NoError(err, "could not get user from database")
	require.Equal(rep.Email, user.Email, "user creation check failed")

	// Test error paths
	// Test password mismatch
	req.Organization = "Financial Services Ltd"
	req.Domain = "financial-services"
	req.PwCheck = "notthe same"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "passwords do not match")

	// Test no agreement
	req.PwCheck = req.Password
	req.AgreeToS = false
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: terms_agreement")

	// Test no email address
	req.AgreeToS = true
	req.Email = ""
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: email")

	// Test invalid project ID
	req.InviteToken = ""
	req.Email = "jannel@example.com"
	req.ProjectID = "invalid"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "could not parse project ID in request")

	// Test cannot create existing user
	req.ProjectID = ""
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusConflict, "user or organization already exists")

	// Test cannot create existing organization
	req.Email = "freddy@example.com"
	req.Domain = "example.com"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusConflict, "user or organization already exists")

	// Test with unknown invite token
	req.InviteToken = "notatoken"
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "invalid invitation")

	// Test error is returned when both invite token and project ID provided
	req.InviteToken = token
	req.ProjectID = project
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "only one field can be set: invite_token, project_id")

	// Test error is returned when organization/domain is missing but project ID is provided
	req.InviteToken = ""
	req.Organization = ""
	req.Domain = ""
	_, err = s.client.Register(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: organization")

	// Wait for all async tasks to finish
	s.StopTasks()

	// Check that the invite token was deleted
	_, err = models.GetUserInvite(context.Background(), token)
	require.ErrorIs(err, models.ErrNotFound, "invite token should have been deleted")

	// Test that one verify email was sent to each user
	messages := []*mock.EmailMeta{
		{
			To:        "rachel@example.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.VerifyEmailRE,
			Timestamp: sent,
		},
		{
			To:        "jane@example.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.VerifyEmailRE,
			Timestamp: sent,
		},
		{
			To:        "joe@checkers.io",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.VerifyEmailRE,
			Timestamp: sent,
		},
	}
	mock.CheckEmails(s.T(), messages)
}

func (s *quarterdeckTestSuite) TestLogin() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetDatabase()
	defer s.ResetTasks()

	// Test Happy Path: user and password expected to be in database fixtures.
	req := &api.LoginRequest{
		Email:    "zendaya@testing.io",
		Password: "iseeallthings",
	}
	tokens, err := s.client.Login(ctx, req)
	require.NoError(err, "was unable to login with valid credentials, have fixtures changed?")
	require.NotEmpty(tokens.AccessToken, "missing access token in response")
	require.NotEmpty(tokens.RefreshToken, "missing refresh token in response")

	// Validate claims are as expected
	claims, err := s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")
	require.Equal("01GQYYKY0ECGWT5VJRVR32MFHM", claims.Subject)
	require.Equal("Zendaya Longeye", claims.Name)
	require.Equal("zendaya@testing.io", claims.Email)
	require.NotEmpty(claims.Picture)
	require.Equal("01GQFQ14HXF2VC7C1HJECS60XX", claims.OrgID, "expected most recent login org to be set in the claims (Checkers)")
	require.Len(claims.Permissions, 13)

	// Test login fails when email in request does not match email in token
	token := "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA"
	req.InviteToken = token
	req.Email = "wrong@example.com"
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusBadRequest, responses.ErrRequestNewInvite)

	// Test invite token exists but is expired
	req.InviteToken = "s6jsNBizyGh_C_ZsUSuJsquONYa--gpcfzorN8DsdjIA"
	req.Email = "eefrank@checkers.io"
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusBadRequest, responses.ErrRequestNewInvite)

	// Test valid login with invite token
	req.Email = "eefrank@checkers.io"
	req.Password = "supersecretssquirrel"
	req.InviteToken = token
	tokens, err = s.client.Login(ctx, req)
	require.NoError(err, "was unable to login with valid credentials, have fixtures changed?")
	require.NotEmpty(tokens.AccessToken, "missing access token in response")
	require.NotEmpty(tokens.RefreshToken, "missing refresh token in response")

	// Validate claims are as expected
	claims, err = s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")
	require.Equal("01GQFQ4475V3BZDMSXFV5DK6XX", claims.Subject)
	require.Equal("eefrank@checkers.io", claims.Email)
	require.NotEmpty(claims.Picture)
	require.Equal("01GKHJRF01YXHZ51YMMKV3RCMK", claims.OrgID)
	require.Len(claims.Permissions, 16)

	// Test login fails with invalid invite token
	req.InviteToken = "notatoken"
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusBadRequest, responses.ErrRequestNewInvite)

	// Test orgID and invite token cannot be used together
	req.OrgID = ulids.New()
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	// Test password incorrect
	req.InviteToken = ""
	req.OrgID = ulid.ULID{}
	req.Password = "this is not the right password"
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusForbidden, responses.ErrTryLoginAgain)

	// Test email and password are required
	_, err = s.client.Login(ctx, &api.LoginRequest{Email: "jannel@example.com"})
	s.CheckError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	_, err = s.client.Login(ctx, &api.LoginRequest{Password: "theeaglefliesatmidnight"})
	s.CheckError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	// Test user not found
	_, err = s.client.Login(ctx, &api.LoginRequest{Email: "jonsey@example.com", Password: "logmeinplease"})
	s.CheckError(err, http.StatusForbidden, responses.ErrTryLoginAgain)

	// Test user not verified
	req = &api.LoginRequest{
		Email:    "jannel@example.com",
		Password: "theeaglefliesatmidnight",
	}
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusForbidden, responses.ErrVerifyEmail)

	// Test that the invite token was deleted after use
	s.StopTasks()
	_, err = models.GetUserInvite(context.Background(), token)
	require.ErrorIs(err, models.ErrNotFound, "invite token should have been deleted")
}

func (s *quarterdeckTestSuite) TestLoginMultiOrg() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Happy Path: user and password expected to be in database fixtures.
	req := &api.LoginRequest{
		Email:    "zendaya@testing.io",
		Password: "iseeallthings",
		OrgID:    ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
	}

	tokens, err := s.client.Login(ctx, req)
	require.NoError(err, "was unable to login with valid credentials, have fixtures changed?")

	claims, err := s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")

	require.Equal("01GKHJRF01YXHZ51YMMKV3RCMK", claims.OrgID)
	require.Len(claims.Permissions, 18, "expected 18 permissions for the user, have the fixtures changed?")

	// Should be able to log into a different organization now
	req.OrgID = ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")

	tokens, err = s.client.Login(ctx, req)
	require.NoError(err, "was unable to login with valid credentials, have fixtures changed?")

	claims, err = s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")

	require.Equal("01GQFQ14HXF2VC7C1HJECS60XX", claims.OrgID)
	require.Len(claims.Permissions, 13)
}

func (s *quarterdeckTestSuite) TestAuthenticate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Happy Path: user and password expected to be in database fixtures.
	req := &api.APIAuthentication{
		ClientID:     "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa",
		ClientSecret: "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS",
	}
	tokens, err := s.client.Authenticate(ctx, req)
	require.NoError(err, "was unable to authenticate with valid api credentials, have fixtures changed?")
	require.NotEmpty(tokens.AccessToken, "missing access token in response")
	require.NotEmpty(tokens.RefreshToken, "missing refresh token in response")

	// Validate claims are as expected
	claims, err := s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")
	require.Equal("01GME02TJP2RRP39MKR525YDQ6", claims.Subject)
	require.Empty(claims.Name)
	require.Empty(claims.Email)
	require.Empty(claims.Picture)
	require.Equal("01GKHJRF01YXHZ51YMMKV3RCMK", claims.OrgID)
	require.Equal("01GQ7P8DNR9MR64RJR9D64FFNT", claims.ProjectID)
	require.Len(claims.Permissions, 5)

	// Test client secret incorrect
	req.ClientSecret = "this is not the right secret"
	_, err = s.client.Authenticate(ctx, req)
	s.CheckError(err, http.StatusForbidden, "invalid credentials")

	// Test email and password are required
	_, err = s.client.Authenticate(ctx, &api.APIAuthentication{ClientID: "DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa"})
	s.CheckError(err, http.StatusBadRequest, "missing credentials")

	_, err = s.client.Authenticate(ctx, &api.APIAuthentication{ClientSecret: "wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS"})
	s.CheckError(err, http.StatusBadRequest, "missing credentials")

	// Test user not found
	_, err = s.client.Authenticate(ctx, &api.APIAuthentication{ClientID: "PBWNdzLwHpcgVBEhocVtRcCWShAYVefe", ClientSecret: "hvXZZcouqH9SKnT6meloCYn2IvkOhYfXuxJzb8Wy9w690BGOKBP0VjQ9vrdv0spI"})
	s.CheckError(err, http.StatusForbidden, "invalid credentials")
}

func (s *quarterdeckTestSuite) TestRefresh() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Happy Path: user and password expected to be in database fixtures.
	req := &api.LoginRequest{
		Email:    "zendaya@testing.io",
		Password: "iseeallthings",
	}
	tokens, err := s.client.Login(ctx, req)
	require.NoError(err, "could not login user to begin authenticate tests, have fixtures changed?")
	require.NotEmpty(tokens.RefreshToken, "no refresh token returned")

	// Get the claims from the refresh token
	origClaims, err := s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify refresh token")

	// Refresh the access and refresh tokens
	newTokens, err := s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: tokens.RefreshToken})
	require.NoError(err, "could not refresh credentials with refresh token")
	require.NotEmpty(newTokens.AccessToken)

	require.NotEqual(tokens.AccessToken, newTokens.AccessToken)
	require.NotEqual(tokens.RefreshToken, newTokens.RefreshToken)

	claims, err := s.srv.VerifyToken(newTokens.AccessToken)
	require.NoError(err, "could not verify new access token")

	require.Equal(origClaims.Subject, claims.Subject)
	require.Equal(origClaims.Name, claims.Name)
	require.Equal(origClaims.Email, claims.Email)
	require.Equal(origClaims.Picture, claims.Picture)
	require.Equal(origClaims.OrgID, claims.OrgID)
	require.Equal(origClaims.ProjectID, claims.ProjectID)

	// Refresh with a specified orgID rather than the one in the token
	orgID := ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")
	newTokens, err = s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: tokens.RefreshToken, OrgID: orgID})
	require.NoError(err, "could not refresh credentials with refresh token")
	require.NotEmpty(newTokens.AccessToken)
	require.NotEqual(tokens.AccessToken, newTokens.AccessToken)
	require.NotEqual(tokens.RefreshToken, newTokens.RefreshToken)

	// Verify the new claims are for the specified org
	claims, err = s.srv.VerifyToken(newTokens.AccessToken)
	require.NoError(err, "could not verify new access token")
	require.Equal(orgID.String(), claims.OrgID)
	require.Equal(origClaims.Subject, claims.Subject)
	require.Equal(origClaims.Name, claims.Name)
	require.Equal(origClaims.Email, claims.Email)
	require.Equal(origClaims.Picture, claims.Picture)
	require.Equal(origClaims.ProjectID, claims.ProjectID)

	// Test passing in an orgID the user is not associated with returns an error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: tokens.RefreshToken, OrgID: ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XY")})
	s.CheckError(err, http.StatusForbidden, responses.ErrLogBackIn)

	// Test empty RefreshRequest returns error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{})
	s.CheckError(err, http.StatusBadRequest, responses.ErrLogBackIn)

	// Test invalid refresh token returns error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: "refresh"})
	s.CheckError(err, http.StatusForbidden, responses.ErrLogBackIn)

	// Test validating with an access token returns an error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: newTokens.AccessToken})
	s.CheckError(err, http.StatusForbidden, responses.ErrLogBackIn)
}

func (s *quarterdeckTestSuite) TestSwitch() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer s.ResetDatabase()
	defer s.ResetTasks()

	// Switching organizations requires authentication
	req := &api.SwitchRequest{}
	_, err := s.client.Switch(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Create valid claims for accessing the API
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01GQYYKY0ECGWT5VJRVR32MFHM",
		},
		Name:        "Zendaya Longeye",
		Email:       "zendaya@testing.io",
		OrgID:       "01GKHJRF01YXHZ51YMMKV3RCMK",
		Permissions: []string{permissions.ReadAPIKeys},
	}

	ctx = s.AuthContext(ctx, claims)

	// An orgID is required in the request
	_, err = s.client.Switch(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing organization id")

	// The orgID cannot be the same as the orgID in the claims
	req.OrgID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	_, err = s.client.Switch(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "cannot switch into the organization you are currently logged into")

	// Cannot switch into an organization that does not exist
	req.OrgID = ulid.Make()
	_, err = s.client.Switch(ctx, req)
	s.CheckError(err, http.StatusForbidden, "invalid credentials")

	// Cannot switch into an organization the user doesn't belong to
	req.OrgID = ulid.MustParse("01GYAVA5ARPRC5Y5CHRJDV34CT")
	_, err = s.client.Switch(ctx, req)
	s.CheckError(err, http.StatusForbidden, "invalid credentials")

	// Happy path: new credentials should be issued
	require := s.Require()
	req.OrgID = ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX")
	rep, err := s.client.Switch(ctx, req)
	require.NoError(err, "could not switch organizations")
	require.NotEmpty(rep.AccessToken, "missing access token")
	require.NotEmpty(rep.RefreshToken, "missing refresh token")

	// Validate claims
	newClaims, err := s.srv.VerifyToken(rep.AccessToken)
	require.NoError(err, "could not verify access token")
	require.Equal(claims.Subject, newClaims.Subject)
	require.Equal(claims.Name, newClaims.Name)
	require.Equal(claims.Email, newClaims.Email)
	require.NotEqual(claims.OrgID, newClaims.OrgID)
	require.Equal(req.OrgID.String(), newClaims.OrgID)
	require.NotEmpty(newClaims.Permissions)
	require.NotEqual(claims.Permissions, newClaims.Permissions)
}

func (s *quarterdeckTestSuite) TestVerify() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetTasks()
	defer s.ResetDatabase()
	defer mock.Reset()

	// Test that an empty token is rejected
	err := s.client.VerifyEmail(ctx, &api.VerifyRequest{})
	s.CheckError(err, http.StatusBadRequest, "missing token in request")

	// Test that an error is returned if it doesn't exist in the database
	req := &api.VerifyRequest{
		Token: "wrongtoken",
	}
	err = s.client.VerifyEmail(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "invalid token")

	// Test that 410 is returned if the token is expired
	// jannel@example.com
	req.Token = "EpiLbYGb58xsOsjk2CWaNMOS0s-LCyW1VVvKrZNg7dI"
	sent := time.Now()
	err = s.client.VerifyEmail(ctx, req)
	s.CheckError(err, http.StatusGone, "token expired, a new verification token has been sent to the email associated with the account")

	// User should be issued a new token
	user, err := models.GetUser(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z", "01GKHJRF01YXHZ51YMMKV3RCMK")
	require.NoError(err, "could not get user from database")
	require.False(user.EmailVerified, "user should not be verified")
	token := user.GetVerificationToken()
	require.NotEmpty(token, "user should have a verification token")
	require.NotEqual(req.Token, token, "user should have a new verification token")
	expiresAt, err := user.GetVerificationExpires()
	require.NoError(err, "could not parse email verification expiration")
	require.True(expiresAt.After(sent), "new token should not be expired")
	require.NotEmpty(user.EmailVerificationSecret, "user should have an email verification secret")

	// Happy path - verifying the user
	req.Token = token
	err = s.client.VerifyEmail(ctx, req)
	require.NoError(err, "could not verify user")

	// User should be verified
	user, err = models.GetUser(ctx, "01GKHJSK7CZW0W282ZN3E9W86Z", "01GKHJRF01YXHZ51YMMKV3RCMK")
	require.NoError(err, "could not get user from database")
	require.True(user.EmailVerified, "user should be verified")

	// Test that 202 is returned if the user is already verified
	err = s.client.VerifyEmail(ctx, req)
	require.NoError(err, "expected no error when user is already verified")

	// Test that the verification email was sent for the expired case
	s.StopTasks()
	messages := []*mock.EmailMeta{
		{
			To:        "jannel@example.com",
			From:      s.conf.SendGrid.FromEmail,
			Subject:   emails.VerifyEmailRE,
			Timestamp: sent,
		},
	}
	mock.CheckEmails(s.T(), messages)
}
