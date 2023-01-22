package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

func (s *quarterdeckTestSuite) TestAPIKeyList() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: implement actual tests
	req := &api.PageQuery{}
	_, err := s.client.APIKeyList(ctx, req)
	require.Error(err, "unauthorized requests should not return a response")

	// require.NoError(err, "should return an empty list")
	// require.Empty(rep.APIKeys)
	// require.Empty(rep.NextPageToken)
}

func (s *quarterdeckTestSuite) TestAPIKeyCreate() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Creating an API Key requires an authenticated endpoint
	req := &api.APIKey{}
	_, err := s.client.APIKeyCreate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Creating an API Key requires the apikeys:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.APIKeyCreate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	claims.Permissions = []string{"apikeys:edit"}
	ctx = s.AuthContext(ctx, claims)

	// TODO: test invalid requests

	// Test Happy Path
	req = &api.APIKey{
		Name:        "Testing Keys",
		Source:      "Test Client",
		ProjectID:   ulids.New(),
		Permissions: []string{"publisher", "subscriber"},
	}

	rep, err := s.client.APIKeyCreate(ctx, req)
	require.NoError(err, "could not execute happy path request")
	require.NotEmpty(s, rep, "expected an API key response from the server")

	// Validate the response returned by the server
	require.False(ulids.IsZero(rep.ID), "no id returned in response")
	require.NotEmpty(rep.ClientID, "no client_id returned in response")
	require.NotEmpty(rep.ClientSecret, "no client_secret returned in response")
	require.NotEmpty(rep.Name, "no name returned in response")
	require.False(ulids.IsZero(rep.OrgID), "no org_id returned in response")
	require.False(ulids.IsZero(rep.ProjectID), "no project_id returned in response")
	require.False(ulids.IsZero(rep.CreatedBy), "no created_by returned in response")
	require.NotEmpty(rep.Source, "no source returned in response")
	require.NotEmpty(rep.UserAgent, "no user agent returned in response")
	require.True(rep.LastUsed.IsZero(), "expected an empty last_used after creating a key")
	require.NotEmpty(rep.Permissions, "no permissions returned in response")
	require.False(rep.Created.IsZero(), "no created returned in response")
	require.False(rep.Modified.IsZero(), "no modified returned in response")

	// Specific assertions about API Key creation
	require.Len(rep.ClientID, 32, "expected len 32 client id")
	require.Len(rep.ClientSecret, 64, "expected len 64 client secret")
	require.Equal(req.Name, rep.Name, "expected name to match request")
	require.Equal(claims.OrgID, rep.OrgID.String(), "expected orgID to match claims")
	require.Equal(req.ProjectID, rep.ProjectID, "expected projectID to match request")
	require.Equal(claims.Subject, rep.CreatedBy.String(), "expected created_by to match claims")
	require.Equal(req.Source, rep.Source, "expected source to match request")
	require.Equal("Quarterdeck API Client/v1", rep.UserAgent, "expected user agent to match client")
	require.Equal(req.Permissions, rep.Permissions, "expected permissions to match request")

	// Assert that the key has been created in the database
	model, err := models.GetAPIKey(ctx, rep.ClientID)
	require.NoError(err, "apikey could not be fetched or was not created")
	require.Equal(rep.ID, model.ID, "apikey fetched from database does not match response")
}

func (s *quarterdeckTestSuite) TestAPIKeyDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieving an API Key requires an authenticated enpdoint
	_, err := s.client.APIKeyDetail(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Retrieving an API Key requires the apikeys:read permission
	claims := &tokens.Claims{
		Name:  "Tom Riddle",
		Email: "voldy@example.com",
		OrgID: ulids.New().String(),
	}

	ctx = s.AuthContext(ctx, claims)
	apiKey, err := s.client.APIKeyDetail(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	require.Nil(apiKey, "no reply should be returned")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Cannot retrieve a key that is not in the same organization
	claims.Permissions = []string{"apikeys:read"}
	ctx = s.AuthContext(ctx, claims)

	apiKey, err = s.client.APIKeyDetail(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	require.Nil(apiKey, "no reply should be returned")
	s.CheckError(err, http.StatusNotFound, "api key not found")

	// Test happy path and fetch the key
	claims = &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01GKHJSK7CZW0W282ZN3E9W86Z",
		},
		Name:        "Jannel P. Hudson",
		Email:       "jannel@example.com",
		OrgID:       "01GKHJRF01YXHZ51YMMKV3RCMK",
		ProjectID:   "01GQ7P8DNR9MR64RJR9D64FFNT",
		Permissions: []string{"apikeys:read"},
	}
	ctx = s.AuthContext(ctx, claims)

	apiKey, err = s.client.APIKeyDetail(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	require.NoError(err, "could not fetch valid API Key detail")
	require.NotNil(apiKey, "expected API Key to be retrieved")

	// IMPORTANT: require that the key secret is empty!
	require.Empty(apiKey.ClientSecret, "the client secret should not be populated in a default response")

	// Check that the model is populated with all expected fields.
	require.False(ulids.IsZero(apiKey.ID), "no id was returned on the response")
	require.NotEmpty(apiKey.ClientID, "no client_id was returned on the response")
	require.NotEmpty(apiKey.Name, "no name was returned on the response")
	require.False(ulids.IsZero(apiKey.OrgID), "no org_id was returned on the response")
	require.False(ulids.IsZero(apiKey.ProjectID), "no project_id was returned on the response")
	require.False(ulids.IsZero(apiKey.CreatedBy), "no created_by was returned on the response")
	require.NotEmpty(apiKey.Source, "no source was returned on the response")
	require.NotEmpty(apiKey.UserAgent, "no user_agent was returned on the response")
	require.False(apiKey.LastUsed.IsZero(), "no last_used was returned on the response")
	require.False(apiKey.Created.IsZero(), "no created was returned on the response")
	require.False(apiKey.Modified.IsZero(), "no modified was returned on the response")

	// Test cannot parse ULID returns not found
	apiKey, err = s.client.APIKeyDetail(ctx, "notaulid")
	require.Nil(apiKey, "no reply should be returned on not found")
	s.CheckError(err, http.StatusNotFound, "api key not found")

	// Test database not found
	apiKey, err = s.client.APIKeyDetail(ctx, ulids.New().String())
	require.Nil(apiKey, "no reply should be returned on not found")
	s.CheckError(err, http.StatusNotFound, "api key not found")
}

func (s *quarterdeckTestSuite) TestAPIKeyUpdate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.APIKey{ID: ulids.New()}
	_, err := s.client.APIKeyUpdate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (s *quarterdeckTestSuite) TestAPIKeyDelete() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.client.APIKeyDelete(ctx, "42")
	require.Error(err, "expected unimplemented error")
}
