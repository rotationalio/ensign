package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (s *tenantTestSuite) TestProjectAPIKeyList() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	projectID := "01GQ38QWNR7MYQXSQ682PJQM7T"
	page := &qd.APIKeyList{
		APIKeys: []*qd.APIKey{
			{
				ID:       ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5BA2G"),
				ClientID: "ABCDEFGHIJKLMNOP",
				Name:     "Leopold's Publish Key",
			},
			{
				ID:       ulid.MustParse("02GQ38J5YWH4DCYJ6CZ2P5BA2G"),
				ClientID: "QRSTUVWXYZABCDEF",
				Name:     "Leopold's Subscribe Key",
			},
			{
				ID:       ulid.MustParse("03GQ38J5YWH4DCYJ6CZ2P5BA2G"),
				ClientID: "GHIJKLMNOPQRSTUV",
				Name:     "Leopold's PubSub Key",
			},
		},
		NextPageToken: "next_page_token",
	}

	// Initial mock checks for an auth token and returns 200 with the page fixture
	s.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(page), mock.RequireAuth())

	// Create initial user claims
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.ProjectAPIKeyList(ctx, "invalid", &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyList(ctx, "invalid", &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Successfully listing API keys
	claims.Permissions = []string{perms.ReadAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	req := &api.PageQuery{
		PageSize: 10,
	}
	reply, err := s.client.ProjectAPIKeyList(ctx, projectID, req)
	require.NoError(err, "expected no error when listing API keys")
	require.Equal(projectID, reply.ProjectID, "expected project ID to match")
	require.Equal(page.NextPageToken, reply.NextPageToken, "expected next page token to match")
	require.Equal(len(page.APIKeys), len(reply.APIKeys), "expected API key count to match")
	for i, key := range reply.APIKeys {
		require.Equal(page.APIKeys[i].ID.String(), key.ID, "expected API key ID to match")
		require.Equal(page.APIKeys[i].ClientID, key.ClientID, "expected API key Client ID to match")
		require.Equal(page.APIKeys[i].Name, key.Name, "expected API key name to match")
	}

	// Error should be returned when Quarterdeck returns an error
	s.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	_, err = s.client.ProjectAPIKeyList(ctx, projectID, req)
	s.requireError(err, http.StatusInternalServerError, "could not list API keys", "expected error when Quarterdeck returns an error")
}

func (s *tenantTestSuite) TestProjectAPIKeyCreate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	projectID := "01GQ38J5YWH4DCYJ6CZ2P5BA2G"
	key := &qd.APIKey{
		ID:           ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5DA2G"),
		ClientID:     "ABCDEFGHIJKLMNOP",
		ClientSecret: "A1B2C3D4E5F6G7H8I9J0",
		Name:         "Leopold's API Key",
		OrgID:        ulid.MustParse("01GQ38QWNR7MYQXSQ682PJQM7T"),
		ProjectID:    ulid.MustParse(projectID),
		CreatedBy:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		LastUsed:     time.Now(),
		Permissions:  []string{"publish", "subscribe"},
		Created:      time.Now(),
		Modified:     time.Now(),
	}

	// Initial mock checks for an auth token and returns 201 with the key fixture
	s.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusCreated), mock.UseJSONFixture(key), mock.RequireAuth())

	// Create initial user claims
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"edit:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set CSRF protection on client")
	_, err := s.client.ProjectAPIKeyCreate(ctx, "invalid", &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Name is required
	claims.Permissions = []string{perms.EditAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", &api.APIKey{})
	s.requireError(err, http.StatusBadRequest, "API key name is required", "expected error when name is missing")

	// Permissions are required
	req := &api.APIKey{
		Name: key.Name,
	}
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", req)
	s.requireError(err, http.StatusBadRequest, "API key permissions are required", "expected error when permissions are missing")

	// ProjectID is required
	req.Permissions = key.Permissions
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", req)
	s.requireError(err, http.StatusBadRequest, "invalid project ID", "expected error when project id is missing")

	// Successfully creating an API key
	expected := &api.APIKey{
		ID:           key.ID.String(),
		ClientID:     key.ClientID,
		ClientSecret: key.ClientSecret,
		Name:         req.Name,
		Owner:        key.CreatedBy.String(),
		Permissions:  req.Permissions,
		Created:      key.Created.Format(time.RFC3339Nano),
		Modified:     key.Modified.Format(time.RFC3339Nano),
	}
	out, err := s.client.ProjectAPIKeyCreate(ctx, projectID, req)
	require.NoError(err, "expected no error when creating API key")
	require.Equal(expected, out, "expected API key to be created")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	_, err = s.client.ProjectAPIKeyCreate(ctx, projectID, req)
	s.requireError(err, http.StatusInternalServerError, "could not create API key", "expected error when quarterdeck returns an error")
}

func (s *tenantTestSuite) TestAPIKeyDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	id := "01GQ38J5YWH4DCYJ6CZ2P5DA2G"
	orgID := "01GQ38QWNR7MYQXSQ682PJQM7T"
	key := &qd.APIKey{
		ID:           ulid.MustParse(id),
		ClientID:     "ABCDEFGHIJKLMNOP",
		ClientSecret: "A1B2C3D4E5F6G7H8I9J0",
		Name:         "Leopold's API Key",
		OrgID:        ulid.MustParse(orgID),
		ProjectID:    ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5BA2G"),
		CreatedBy:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		LastUsed:     time.Now(),
		Permissions:  []string{"publish", "subscribe"},
		Created:      time.Now(),
		Modified:     time.Now(),
	}

	// Initial mock checks for an auth token and returns 200 with the key fixture
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusOK), mock.UseJSONFixture(key), mock.RequireAuth())

	// Create initial user claims
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		OrgID:       orgID,
	}

	// Endpoint must be authenticated
	_, err := s.client.APIKeyDetail(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.APIKeyDetail(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Successfully retrieving an API key
	claims.Permissions = []string{perms.ReadAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	expected := &api.APIKey{
		ID:          id,
		ClientID:    key.ClientID,
		Name:        key.Name,
		Owner:       key.CreatedBy.String(),
		Permissions: key.Permissions,
		Created:     key.Created.Format(time.RFC3339Nano),
		Modified:    key.Modified.Format(time.RFC3339Nano),
	}
	out, err := s.client.APIKeyDetail(ctx, id)
	require.NoError(err, "expected no error when retrieving API key")
	require.Equal(expected, out, "expected API key to be retrieved")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	_, err = s.client.APIKeyDetail(ctx, id)
	s.requireError(err, http.StatusInternalServerError, "could not retrieve API key", "expected error when quarterdeck returns an error")
}

func (s *tenantTestSuite) TestAPIKeyDelete() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := "01GQ38J5YWH4DCYJ6CZ2P5DA2G"
	orgID := "01GQ38QWNR7MYQXSQ682PJQM7T"

	// Initial mock checks for an auth token and returns 204
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusNoContent), mock.RequireAuth())

	// Create initial user claims
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
		OrgID:       orgID,
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set client CSRF protection")
	err := s.client.APIKeyDelete(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	err = s.client.APIKeyDelete(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Successfully deleting an API key
	claims.Permissions = []string{perms.DeleteAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	err = s.client.APIKeyDelete(ctx, id)
	require.NoError(err, "expected no error when deleting API key")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	err = s.client.APIKeyDelete(ctx, id)
	s.requireError(err, http.StatusInternalServerError, "could not delete API key", "expected error when quarterdeck returns an error")
}
