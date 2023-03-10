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
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *tenantTestSuite) TestProjectAPIKeyList() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	projectID := "01GQ38QWNR7MYQXSQ682PJQM7T"
	orgID := "02ABC8QWNR7MYQXSQ682PJQM7T"
	tenantID := "03ABC8QWNR7MYQXSQ682PJQM7Y"
	project := &db.Project{
		TenantID: ulid.MustParse(tenantID),
		ID:       ulid.MustParse(projectID),
		OrgID:    ulid.MustParse(orgID),
	}

	key, err := project.Key()
	require.NoError(err, "could not create project key")

	var data []byte
	data, err = project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	// Trtl Get should return project key or project data
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{Value: key}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "namespace %s not found", gr.Namespace)
		}
	}

	// Create initial fixtures
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
	_, err = s.client.ProjectAPIKeyList(ctx, "invalid", &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyList(ctx, "invalid", &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Should fail if OrgID is not in the claims
	req := &api.PageQuery{
		PageSize: 10,
	}
	claims.Permissions = []string{perms.ReadAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyList(ctx, projectID, req)
	s.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when user does not have an OrgID")

	// TODO: Add test for wrong orgID in claims

	// Successfully listing API keys
	claims.OrgID = orgID
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
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

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	projectID := "01GQ38J5YWH4DCYJ6CZ2P5BA2G"
	orgID := "02ABC8QWNR7MYQXSQ682PJQM7T"
	tenantID := "03DEF8QWNR7MYQXSQ682PJQM7T"
	project := &db.Project{
		ID:       ulid.MustParse(projectID),
		OrgID:    ulid.MustParse(orgID),
		TenantID: ulid.MustParse(tenantID),
	}
	keyData, err := project.Key()
	require.NoError(err, "could not generate project key")

	var projectData []byte
	projectData, err = project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	// OnGet should return success for project retrieval
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: keyData,
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: projectData,
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", gr.Namespace)
		}
	}

	// Create initial fixtures
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
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyCreate(ctx, "invalid", &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Should fail if the OrgID is not in the claims
	claims.Permissions = []string{perms.EditAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.ProjectAPIKeyCreate(ctx, projectID, &api.APIKey{})
	s.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when OrgID is not in claims")

	// Name is required
	claims.OrgID = orgID
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

	// TODO: Add test for wrong orgID in claims

	// Successfully creating an API key
	claims.OrgID = orgID
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
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

func (s *tenantTestSuite) TestAPIKeyUpdate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	id := "01GQ38J5YWH4DCYJ6CZ2P5DA2G"
	orgID := "01GQ38QWNR7MYQXSQ682PJQM7T"
	key := &qd.APIKey{
		ID:          ulid.MustParse(id),
		ClientID:    "ABCDEFGHIJKLMNOP",
		Name:        "Leopold's Renamed API Key",
		OrgID:       ulid.MustParse(orgID),
		ProjectID:   ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5BA2G"),
		CreatedBy:   ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		LastUsed:    time.Now(),
		Permissions: []string{"publish", "subscribe"},
		Created:     time.Now(),
		Modified:    time.Now(),
	}

	// Initial mock checks for an auth token and returns 200 with the key fixture
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusOK), mock.UseJSONFixture(key), mock.RequireAuth())

	// Create initial user claims
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
		OrgID:       orgID,
	}

	// Endpoint must be authenticated
	req := &api.APIKey{
		ID: "invalid",
	}
	require.NoError(s.SetClientCSRFProtection(), "could not set client CSRF protection")
	_, err := s.client.APIKeyUpdate(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.APIKeyUpdate(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Should return an error when the API key is not parseable
	claims.Permissions = []string{perms.EditAPIKeys}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.APIKeyUpdate(ctx, req)
	s.requireError(err, http.StatusBadRequest, "could not parse API key ID from URL", "expected error when API key ID is not parseable")

	// Should return an error when the name is not provided
	req.ID = id
	_, err = s.client.APIKeyUpdate(ctx, req)
	s.requireError(err, http.StatusBadRequest, "API key name is required for update", "expected error when name is not provided")

	// Sucessfully update an API key
	expected := &api.APIKey{
		ID:          id,
		ClientID:    "ABCDEFGHIJKLMNOP",
		Name:        "Leopold's Renamed API Key",
		Owner:       key.CreatedBy.String(),
		Permissions: key.Permissions,
		Created:     key.Created.Format(time.RFC3339Nano),
		Modified:    key.Modified.Format(time.RFC3339Nano),
	}
	req.Name = "Leoopold's Renamed API Key"
	reply, err := s.client.APIKeyUpdate(ctx, req)
	require.NoError(err, "expected no error when updating API key")
	require.Equal(expected, reply, "expected updated API key to be returned")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys(id, mock.UseStatus(http.StatusInternalServerError), mock.RequireAuth())
	_, err = s.client.APIKeyUpdate(ctx, req)
	s.requireError(err, http.StatusInternalServerError, "could not update API key", "expected error when quarterdeck returns an error")
}

func (s *tenantTestSuite) TestAPIKeyPermissions() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	perms := []string{perms.Publisher, perms.Subscriber, perms.ReadTopics, perms.ReadMetrics}

	// Initial mock returns 200 with the permissions fixture
	s.quarterdeck.OnAPIKeys("permissions", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(perms))

	// Endpoint must be authenticated
	_, err := s.client.APIKeyPermissions(ctx)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Create valid claims for the user
	claims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Test successful response
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	reply, err := s.client.APIKeyPermissions(ctx)
	require.NoError(err, "expected no error when getting API key permissions")
	require.Equal(perms, reply, "expected API key permissions to be returned")

	// Ensure an error is returned when quarterdeck returns an error
	s.quarterdeck.OnAPIKeys("permissions", mock.UseStatus(http.StatusUnauthorized))
	_, err = s.client.APIKeyPermissions(ctx)
	s.requireError(err, http.StatusUnauthorized, "could not retrieve API key permissions for user", "expected error when quarterdeck returns an error")
}
