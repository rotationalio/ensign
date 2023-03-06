package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *quarterdeckTestSuite) TestAPIKeyList() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Listing API Keys requires authentication
	req := &api.APIPageQuery{}
	_, err := s.client.APIKeyList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Listing API Keys requires the apikeys:read permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)

	_, err = s.client.APIKeyList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GKHJSK7CZW0W282ZN3E9W86Z"
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	claims.Permissions = []string{perms.ReadAPIKeys}
	ctx = s.AuthContext(ctx, claims)

	// Should be able to list all keys for the specified organization
	page, err := s.client.APIKeyList(ctx, req)
	require.NoError(err, "could not fetch api keys")
	require.Len(page.APIKeys, 11, "expected 11 results back from the fixtures")
	require.Empty(page.NextPageToken, "expected no next page token in response")

	// Should be able to pagination the request for the specified organization
	req.PageSize = 3
	page, err = s.client.APIKeyList(ctx, req)
	require.NoError(err, "could not fetch paginated api keys")
	require.Len(page.APIKeys, 3, "expected 3 results back from the fixtures")
	require.NotEmpty(page.NextPageToken, "expected next page token in response")

	// Test fetching the next page with the next page token
	req.NextPageToken = page.NextPageToken
	page2, err := s.client.APIKeyList(ctx, req)
	require.NoError(err, "could not fetch paginated api keys")
	require.Len(page2.APIKeys, 3, "expected 3 results back from the fixtures")
	require.NotEmpty(page2.NextPageToken, "expected next page token in response")
	require.NotEqual(page.APIKeys[2].ID, page2.APIKeys[0].ID, "expected a new page of results")

	// Test filtering with ProjectID and complete pagination with multiple requests
	req = &api.APIPageQuery{
		ProjectID: "01GQFR0KM5S2SSJ8G5E086VQ9K",
		PageSize:  3,
	}

	// Limit maximum number of request to 10, break when pagination is complete.
	nPages, nResults := 0, 0
	for i := 0; i < 10; i++ {
		page, err = s.client.APIKeyList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.APIKeys)

		for _, key := range page.APIKeys {
			// Ensure the project filter is working properly
			require.Equal(req.ProjectID, key.ProjectID.String())
			require.Equal(claims.OrgID, key.OrgID.String())
		}

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 3, "expected 9 results in 3 pages")
	require.Equal(nResults, 9, "expected 9 results in 3 pages")

	// TODO: test edge cases and bad requests
}

func (s *quarterdeckTestSuite) TestAPIKeyCreate() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Creating an API Key requires authentication
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
	claims.Permissions = []string{perms.EditAPIKeys}
	ctx = s.AuthContext(ctx, claims)

	// TODO: test invalid requests

	// Test Happy Path
	req = &api.APIKey{
		Name:        "Testing Keys",
		Source:      "Test Client",
		ProjectID:   ulid.MustParse("01GQ7P8DNR9MR64RJR9D64FFNT"),
		Permissions: []string{"publisher", "subscriber"},
	}

	rep, err := s.client.APIKeyCreate(ctx, req)
	require.NoError(err, "could not execute happy path request")
	require.NotEmpty(rep, "expected an API key response from the server")

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

func (s *quarterdeckTestSuite) TestCannotCreateAPIKeyInUnownedProject() {
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Creating an API Key requires the apikeys:edit permission
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01GKHJSK7CZW0W282ZN3E9W86Z",
		},
		Name:        "Jannel P. Hudson",
		Email:       "jannel@example.com",
		OrgID:       "01GKHJRF01YXHZ51YMMKV3RCMK",
		Permissions: []string{perms.EditAPIKeys},
	}
	ctx = s.AuthContext(ctx, claims)

	// User should not be able to create an APIKey in a project not owned by that org.
	req := &api.APIKey{
		Name:        "Sneaky Key",
		Source:      "Hacker",
		ProjectID:   ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"),
		Permissions: []string{"publisher", "subscriber"},
	}

	_, err := s.client.APIKeyCreate(ctx, req)
	s.CheckError(err, 400, "validation error: invalid project id for apikey")

	// Ensure that OrgID comes from claims and not user input
	req.OrgID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	_, err = s.client.APIKeyCreate(ctx, req)
	s.CheckError(err, 400, "field restricted for request: org_id")
}

func (s *quarterdeckTestSuite) TestAPIKeyDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieving an API Key requires authentication
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
	claims.Permissions = []string{perms.ReadAPIKeys}
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
		Permissions: []string{perms.ReadAPIKeys},
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
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Updating an API Key requires authentication
	in := &api.APIKey{ID: ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6"), Name: "changed"}
	out, err := s.client.APIKeyUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	require.Nil(out, "expected no data returned after an error")

	// Updating an API Key requires the apikeys:edit permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)

	out, err = s.client.APIKeyUpdate(ctx, in)
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")
	require.Nil(out, "expected no data returned after an error")

	// Cannot update a key that is not in the same organization
	claims.Permissions = []string{perms.EditAPIKeys}
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.APIKeyUpdate(ctx, in)
	s.CheckError(err, http.StatusNotFound, "api key not found")
	require.Nil(out, "expected no data returned after an error")

	// Test happy path and delete the key
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	out, err = s.client.APIKeyUpdate(ctx, in)
	require.NoError(err, "should have been able to update the key")
	require.NotSame(in, out, "expected a different object to be returned")

	// IMPORTANT: ClientSecret should not be returned in this endpoint!
	require.Empty(out.ClientSecret, "client secret should not be returned in update")

	// Check that the model is populated with all expected fields.
	require.False(ulids.IsZero(out.ID), "no id was returned on the response")
	require.NotEmpty(out.ClientID, "no client_id was returned on the response")
	require.NotEmpty(out.Name, "no name was returned on the response")
	require.False(ulids.IsZero(out.OrgID), "no org_id was returned on the response")
	require.False(ulids.IsZero(out.ProjectID), "no project_id was returned on the response")
	require.False(ulids.IsZero(out.CreatedBy), "no created_by was returned on the response")
	require.NotEmpty(out.Source, "no source was returned on the response")
	require.NotEmpty(out.UserAgent, "no user_agent was returned on the response")
	require.False(out.LastUsed.IsZero(), "no last_used was returned on the response")
	require.False(out.Created.IsZero(), "no created was returned on the response")
	require.False(out.Modified.IsZero(), "no modified was returned on the response")

	// Verify key was updated
	key, err := models.RetrieveAPIKey(ctx, in.ID)
	require.NoError(err, "could not retrieve key from database")
	require.Equal("changed", key.Name, "key was not updated in the database")

	// Test database not found
	_, err = s.client.APIKeyUpdate(ctx, &api.APIKey{ID: ulids.New(), Name: "changed"})
	s.CheckError(err, http.StatusNotFound, "api key not found")

	// TODO: Test cannot parse ULID returns not found (needs direct http client)
	// TODO: test other validation cases and bad requests
}

func (s *quarterdeckTestSuite) TestAPIKeyDelete() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Deleting an API Key requires authentication
	err := s.client.APIKeyDelete(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Deleting an API Key requires the apikeys:delete permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
		OrgID: ulids.New().String(),
	}
	ctx = s.AuthContext(ctx, claims)

	err = s.client.APIKeyDelete(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	s.CheckError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Cannot delete a key that is not in the same organization
	claims.Permissions = []string{perms.DeleteAPIKeys}
	ctx = s.AuthContext(ctx, claims)
	err = s.client.APIKeyDelete(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	s.CheckError(err, http.StatusNotFound, "api key not found")

	// Test happy path and delete the key
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	err = s.client.APIKeyDelete(ctx, "01GME02TJP2RRP39MKR525YDQ6")
	require.NoError(err, "should have been able to delete the key")

	// Verify key was deleted
	_, err = models.RetrieveAPIKey(ctx, ulid.MustParse("01GME02TJP2RRP39MKR525YDQ6"))
	require.ErrorIs(err, models.ErrNotFound)

	// Test cannot parse ULID returns not found
	err = s.client.APIKeyDelete(ctx, "notaulid")
	s.CheckError(err, http.StatusNotFound, "api key not found")

	// Test database not found
	err = s.client.APIKeyDelete(ctx, ulids.New().String())
	s.CheckError(err, http.StatusNotFound, "api key not found")
}
