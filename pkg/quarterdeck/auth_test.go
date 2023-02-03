package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (s *quarterdeckTestSuite) TestRegister() {
	defer s.ResetDatabase()
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

	// Test error paths
	// Test password mismatch
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
}

func (s *quarterdeckTestSuite) TestLogin() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Happy Path: user and password expected to be in database fixtures.
	req := &api.LoginRequest{
		Email:    "jannel@example.com",
		Password: "theeaglefliesatmidnight",
	}
	tokens, err := s.client.Login(ctx, req)
	require.NoError(err, "was unable to login with valid credentials, have fixtures changed?")
	require.NotEmpty(tokens.AccessToken, "missing access token in response")
	require.NotEmpty(tokens.RefreshToken, "missing refresh token in response")

	// Validate claims are as expected
	claims, err := s.srv.VerifyToken(tokens.AccessToken)
	require.NoError(err, "could not verify token")
	require.Equal("01GKHJSK7CZW0W282ZN3E9W86Z", claims.Subject)
	require.Equal("Jannel P. Hudson", claims.Name)
	require.Equal("jannel@example.com", claims.Email)
	require.NotEmpty(claims.Picture)
	require.Equal("01GKHJRF01YXHZ51YMMKV3RCMK", claims.OrgID)
	require.Len(claims.Permissions, 18)

	// Test password incorrect
	req.Password = "this is not the right password"
	_, err = s.client.Login(ctx, req)
	s.CheckError(err, http.StatusForbidden, "invalid login credentials")

	// Test email and password are required
	_, err = s.client.Login(ctx, &api.LoginRequest{Email: "jannel@example.com"})
	s.CheckError(err, http.StatusBadRequest, "missing credentials")

	_, err = s.client.Login(ctx, &api.LoginRequest{Password: "theeaglefliesatmidnight"})
	s.CheckError(err, http.StatusBadRequest, "missing credentials")

	// Test user not found
	_, err = s.client.Login(ctx, &api.LoginRequest{Email: "jonsey@example.com", Password: "logmeinplease"})
	s.CheckError(err, http.StatusForbidden, "invalid login credentials")
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
	require.Len(claims.Permissions, 6)

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
		Email:    "jannel@example.com",
		Password: "theeaglefliesatmidnight",
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

	// Test empty RefreshRequest returns error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{})
	s.CheckError(err, http.StatusBadRequest, "missing credentials")

	// Test invalid refresh token returns error
	_, err = s.client.Refresh(ctx, &api.RefreshRequest{RefreshToken: "refresh"})
	s.CheckError(err, http.StatusForbidden, "could not verify refresh token")
}
