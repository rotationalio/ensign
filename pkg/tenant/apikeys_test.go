package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (s *tenantTestSuite) TestAPIKeyCreate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	var err error
	id := "01GQ38J5YWH4DCYJ6CZ2P5FA2G"
	key := &qd.APIKey{
		ID:           ulid.MustParse(id),
		ClientID:     "ABCDEFGHIJKLMNOP",
		ClientSecret: "A1B2C3D4E5F6G7H8I9J0",
		Name:         "Leopold's API Key",
		OrgID:        ulid.MustParse("01GQ38QWNR7MYQXSQ682PJQM7T"),
		ProjectID:    ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5FA2G"),
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
	_, err = s.client.APIKeyCreate(ctx, &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.APIKeyCreate(ctx, &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Name is required
	claims.Permissions = []string{tenant.WriteAPIKey}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.APIKeyCreate(ctx, &api.APIKey{})
	s.requireError(err, http.StatusBadRequest, "API key name is required", "expected error when name is missing")

	// Permissions are required
	req := &api.APIKey{
		Name: key.Name,
	}
	_, err = s.client.APIKeyCreate(ctx, req)
	s.requireError(err, http.StatusBadRequest, "API key permissions are required", "expected error when permissions are missing")

	// ProjectID is required
	req.Permissions = key.Permissions
	_, err = s.client.APIKeyCreate(ctx, &api.APIKey{})
	s.requireError(err, http.StatusBadRequest, "API key name is required", "expected error when name is missing")

	// Successfully creating an API key
	req.ProjectID = key.ProjectID.String()
	expected := &api.APIKey{
		ID:           id,
		ClientID:     key.ClientID,
		ClientSecret: key.ClientSecret,
		Name:         req.Name,
		ProjectID:    req.ProjectID,
		Owner:        claims.Name,
		Permissions:  req.Permissions,
		Created:      key.Created.Format(time.RFC3339Nano),
		Modified:     key.Modified.Format(time.RFC3339Nano),
	}
	out, err := s.client.APIKeyCreate(ctx, req)
	require.NoError(err, "expected no error when creating API key")
	require.Equal(expected, out, "expected API key to be created")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys("", mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	_, err = s.client.APIKeyCreate(ctx, req)
	s.requireError(err, http.StatusInternalServerError, "could not create API key", "expected error when quarterdeck returns an error")
}
