package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/gravatar"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	en "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestTenantProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	orgID := ulid.MustParse("02GMTWFK4XZY597Y128KXQ4ABC")
	ownerID := ulids.New()

	projects := []*db.Project{
		{
			OrgID:       orgID,
			TenantID:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:          ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5FA2G"),
			OwnerID:     ownerID,
			Name:        "project001",
			Description: "This is an archived project.",
			Archived:    true,
			APIKeys:     2,
			Created:     time.Unix(1670424445, 0),
			Modified:    time.Unix(1670424445, 0),
		},
		{
			OrgID:       orgID,
			TenantID:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:          ulid.MustParse("01GQ38JP6CCWPNDS6KG5WDA59T"),
			OwnerID:     ownerID,
			Name:        "project002",
			Description: "This is a project with no topics.",
			APIKeys:     2,
			Created:     time.Unix(1673659941, 0),
			Modified:    time.Unix(1673659941, 0),
		},
		{
			OrgID:       orgID,
			TenantID:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:          ulid.MustParse("01GQ38K6YPE0ZA9ADC2BGSVWRM"),
			OwnerID:     ownerID,
			Name:        "project003",
			Description: "This is an active project.",
			APIKeys:     2,
			Topics:      3,
			Created:     time.Unix(1674073941, 0),
			Modified:    time.Unix(1674073941, 0),
		},
	}

	states := []string{
		db.ProjectStatusArchived,
		db.ProjectStatusIncomplete,
		db.ProjectStatusActive,
	}

	prefix := tenantID[:]
	namespace := "projects"

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	key, err := db.CreateKey(orgID, tenantID)
	require.NoError(err, "could not create tenant key")

	keyData, err := key.MarshalValue()
	require.NoError(err, "could not marshal data")

	member := &db.Member{
		OrgID: orgID,
		ID:    ownerID,
		Email: "leopold.wentzel@gmail.com",
		Name:  "Leopold Wentzel",
	}

	memberKey, err := member.Key()
	require.NoError(err, "could not create member key from struct fixture")

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal data")

	// Trtl should return the Tenant key or member info on get
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.KeysNamespace:
			if !bytes.Equal(in.Key, tenantID[:]) {
				return nil, status.Error(codes.FailedPrecondition, "unexpected get request")
			}
			return &pb.GetReply{
				Value: keyData,
			}, nil
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, memberKey) {
				return nil, status.Error(codes.FailedPrecondition, "unexpected get request")
			}
			return &pb.GetReply{
				Value: memberData,
			}, nil
		default:
			return nil, status.Error(codes.FailedPrecondition, "unexpected namespace in get request")
		}
	}

	// Call the OnCursor method
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, project := range projects {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, project.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := project.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       project.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	req := &api.PageQuery{}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.TenantProjectList(ctx, "invalid", req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectList(ctx, "invalid", req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// TODO: Add test for wrong orgID in claims

	// Should return an error if the tenant does not exist.
	claims.OrgID = orgID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectList(ctx, "invalid", req)
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant does not exist")

	rep, err := suite.client.TenantProjectList(ctx, tenantID.String(), req)
	require.NoError(err, "could not list tenant projects")
	require.Len(rep.TenantProjects, 3, "expected 3 projects")
	require.Empty(rep.NextPageToken, "next page token should not be set when there is only 1 page")

	// Verify project data has been populated.
	for i := range projects {
		require.Equal(projects[i].ID.String(), rep.TenantProjects[i].ID, "expected project id to match")
		require.Equal(projects[i].Name, rep.TenantProjects[i].Name, "expected project name to match")
		require.Equal(member.Name, rep.TenantProjects[i].Owner.Name, "expected project owner name to match")
		require.Equal(projects[i].Description, rep.TenantProjects[i].Description, "expected project description to match")
		require.Equal(member.Picture(), rep.TenantProjects[i].Owner.Picture, "expected project owner picture to match")
		require.Equal(states[i], rep.TenantProjects[i].Status, "expected project state to match")
		require.Equal(projects[i].Topics, rep.TenantProjects[i].ActiveTopics, "expected project topics to match")
		require.Zero(rep.TenantProjects[i].DataStorage.Value, "expected project data storage to be zero")
		require.Equal("GB", rep.TenantProjects[i].DataStorage.Units, "expected project data storage units to be GB")
		require.Equal(projects[i].Created.Format(time.RFC3339Nano), rep.TenantProjects[i].Created, "expected project created time to match")
	}

	// Set page size and test pagination.
	req.PageSize = 2
	rep, err = suite.client.TenantProjectList(ctx, tenantID.String(), req)
	require.NoError(err, "could not list projects")
	require.Len(rep.TenantProjects, 2, "expected 2 projects")
	require.NotEmpty(rep.NextPageToken, "next page token should bet set")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.TenantProjectList(ctx, tenantID.String(), req)
	require.NoError(err, "could not list projects")
	require.Len(rep2.TenantProjects, 1, "expected 1 project")
	require.NotEqual(rep.TenantProjects[0].ID, rep2.TenantProjects[0].ID, "should not have same project ID")
	require.Empty(rep2.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 3, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 3; i++ {
		page, err := suite.client.TenantProjectList(ctx, tenantID.String(), req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.TenantProjects)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 2, "expected 3 results in 2 pages")
	require.Equal(nResults, 3, "expected 3 results in 2 pages")
}

func (suite *tenantTestSuite) TestTenantProjectCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// OnGet returns the tenantID.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: tenantID[:],
		}, nil
	}

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Quarterdeck server mock expects authentication and returns 200 OK
	suite.quarterdeck.OnProjectsCreate(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(&qd.Project{}), mock.RequireAuth())

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"write:nothing"},
	}
	claims.Subject = ulids.New().String()

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectCreate(ctx, tenantID.String(), &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if tenant id is not a valid ULID.
	claims.OrgID = "01GMBVR86186E0EKCHQK4ESJB1"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectCreate(ctx, "tenantID", &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusNotFound, "tenant not found", "expected error when tenant id does not exist")

	// Should return an error if the project ID exists.
	_, err = suite.client.TenantProjectCreate(ctx, tenantID.String(), &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "project id cannot be specified on create", "expected error when project id exists")

	// Should return an error if the project name does not exist.
	_, err = suite.client.TenantProjectCreate(ctx, tenantID.String(), &api.Project{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	// Should return an error if the project description is too long.
	_, err = suite.client.TenantProjectCreate(ctx, tenantID.String(), &api.Project{ID: "", Name: "project001", Description: strings.Repeat("a", 2001)})
	suite.requireError(err, http.StatusBadRequest, "project description is too long", "expected error when project description is too long")

	// Create a project test fixture.
	req := &api.Project{
		Name: "project001",
	}

	project, err := suite.client.TenantProjectCreate(ctx, tenantID.String(), req)
	require.NoError(err, "could not add project")
	require.NotEmpty(project.ID, "expected non-zero ulid to be populated")
	require.Equal(req.Name, project.Name, "project name should match")
	require.Equal(claims.Name, project.Owner.Name, "project owner name should match")
	require.Equal(gravatar.New(claims.Email, nil), project.Owner.Picture, "project owner picture should match")
	require.NotEmpty(project.Created, "expected non-zero created time to be populated")
	require.NotEmpty(project.Modified, "expected non-zero modified time to be populated")

	// Should return an error if the Quarterdeck returns an error
	suite.quarterdeck.OnProjectsCreate(mock.UseError(http.StatusInternalServerError, "could not create project"), mock.RequireAuth())
	_, err = suite.client.TenantProjectCreate(ctx, tenantID.String(), req)
	suite.requireError(err, http.StatusInternalServerError, "could not create project", "expected error when quarterdeck returns an error")

	// Quarterdeck mock should have been called
	require.Equal(2, suite.quarterdeck.ProjectsCreateCount(), "expected quarterdeck mock to be called")
}

func (suite *tenantTestSuite) TestTenantProjectPatch() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	projectID := ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5GA2G")
	ownerID := ulid.MustParse("02ABC8J5YWH4DCYJ6CZ2P5GA3H")

	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		OrgID:       orgID,
		TenantID:    tenantID,
		ID:          projectID,
		OwnerID:     ownerID,
		Name:        "Bouldering Project",
		Description: "Track your epic bouldering sends",
	}

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	projectKey, err := project.Key()
	require.NoError(err, "could not create project key from fixture")

	member := &db.Member{
		OrgID: orgID,
		ID:    ownerID,
		Name:  "Magnus Midtbø",
	}

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	memberKey, err := member.Key()
	require.NoError(err, "could not create member key from fixture")

	newOwner := &db.Member{
		OrgID: orgID,
		ID:    ulids.New(),
		Name:  "Adam Ondra",
	}

	newOwnerData, err := newOwner.MarshalValue()
	require.NoError(err, "could not marshal new owner data")

	newOwnerKey, err := newOwner.Key()
	require.NoError(err, "could not create new owner key from fixture")

	// Trtl Get should return the requested resources by default
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: projectKey,
			}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: orgID[:],
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: projectData,
			}, nil
		case db.MembersNamespace:
			switch {
			case bytes.Equal(gr.Key, memberKey):
				return &pb.GetReply{
					Value: memberData,
				}, nil
			case bytes.Equal(gr.Key, newOwnerKey):
				return &pb.GetReply{
					Value: newOwnerData,
				}, nil
			default:
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, errors.New("unexpected namespace")
		}
	}

	// Trtl Put should succeed by default
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       ulids.New().String(),
		Permissions: []string{"edit:nothing"},
	}

	// Endpoint must be authenticated
	req := &api.Project{}
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), projectID.String(), req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the right permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), projectID.String(), req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have the right permissions")

	// Error should be returned if the tenantID is not parseable
	claims.Permissions = []string{perms.EditProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantProjectPatch(ctx, "invalid", projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrTenantNotFound, "expected error when tenant id is not parseable")

	// Error should be returned if the projectID is not parseable
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), "invalid", req)
	suite.requireError(err, http.StatusNotFound, responses.ErrProjectNotFound, "expected error when project id is not parseable")

	// User must be in the same org as the tenant
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrTenantNotFound, "expected error when user is not in the same org as the tenant")

	// Set claims to the correct orgID
	claims.OrgID = orgID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Test field validation
	testCases := []struct {
		req *api.Project
		err string
	}{
		{&api.Project{Description: strings.Repeat("a", db.MaxDescriptionLength+1)}, "validation error: project description is too long"},
		{&api.Project{Status: "Archived"}, api.InvalidFieldError("status").Error()},
		{&api.Project{Name: "Moonboard Sends", Description: "Crushin' v7s on the Moonboard"}, ""},
		{&api.Project{Owner: api.Member{ID: newOwner.ID.String()}}, ""},
		{&api.Project{Owner: api.Member{ID: newOwner.ID.String(), Name: "Adam Ondra"}}, ""},
		{&api.Project{Owner: api.Member{}, Name: "Moonboard Sends"}, ""},
		{&api.Project{}, ""},
	}

	for i, tc := range testCases {
		_, err = suite.client.ProjectPatch(ctx, projectID.String(), tc.req)
		if tc.err == "" {
			require.NoError(err, "expected no error with valid fields")
		} else {
			suite.requireError(err, http.StatusBadRequest, tc.err, fmt.Sprintf("expected error when field validation fails for test case %d", i))
		}
	}

	// Should return an error if the owner is not parseable as a ULID
	req.Owner = api.Member{ID: "invalid"}
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when owner is not parseable as a ULID")

	// Should return an error if the owner does not exist
	req.Owner = api.Member{ID: ulids.New().String()}
	_, err = suite.client.TenantProjectPatch(ctx, tenantID.String(), projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when owner is not found")
}

func (suite *tenantTestSuite) TestProjectList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")

	defer cancel()

	projects := []*db.Project{
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5FA2G"),
			Name:     "project001",
			Created:  time.Unix(1670424445, 0),
			Modified: time.Unix(1670424445, 0),
		},
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38JP6CCWPNDS6KG5WDA59T"),
			Name:     "project002",
			Created:  time.Unix(1673659941, 0),
			Modified: time.Unix(1673659941, 0),
		},
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ38K6YPE0ZA9ADC2BGSVWRM"),
			Name:     "project003",
			Created:  time.Unix(1674073941, 0),
			Modified: time.Unix(1674073941, 0),
		},
	}

	prefix := tenantID[:]
	namespace := "projects"

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, project := range projects {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, project.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := project.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       project.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	req := &api.PageQuery{}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMTWFK4XZY597Y128KXQ4WHP",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated.
	_, err := suite.client.ProjectList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permission")

	// Set valid permissions for the rest of the tests.
	claims.Permissions = []string{perms.ReadProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Retrieve all projects.
	rep, err := suite.client.ProjectList(ctx, req)
	require.NoError(err, "could not list projects")
	require.Len(rep.Projects, 3, "expected 3 projects")
	require.Empty(rep.NextPageToken, "did not expect next page token when there is only 1 page")

	// Verify project data has been populated.
	for i := range projects {
		require.Equal(projects[i].ID.String(), rep.Projects[i].ID, "project id should match")
		require.Equal(projects[i].Name, rep.Projects[i].Name, "project name should match")
		require.Equal(projects[i].Created.Format(time.RFC3339Nano), rep.Projects[i].Created, "project created should match")
		require.Equal(projects[i].Modified.Format(time.RFC3339Nano), rep.Projects[i].Modified, "project modified should match")
	}

	// Set page size and test pagination.
	req.PageSize = 2
	rep, err = suite.client.ProjectList(ctx, req)
	require.NoError(err, "could not list projects")
	require.Len(rep.Projects, 2, "expected 2 projects")
	require.NotEmpty(rep.NextPageToken, "next page token should be set")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.ProjectList(ctx, req)
	require.NoError(err, "could not list projects")
	require.Len(rep2.Projects, 1, "expected 1 project")
	require.NotEqual(rep.Projects[0].ID, rep2.Projects[0].ID, "should not have same project ID")
	require.Empty(rep2.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 3, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 3; i++ {
		page, err := suite.client.ProjectList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Projects)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 2, "expected 3 results in 2 pages")
	require.Equal(nResults, 3, "expected 3 results in 2 pages")

	// Set test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "",
		Permissions: []string{perms.ReadProjects},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.ProjectList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestProjectCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// OnGet returns the tenantID.
	tenantID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")

	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: tenantID[:],
		}, nil
	}

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Quarterdeck server mock expects authentication and returns 200 OK
	suite.quarterdeck.OnProjectsCreate(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(&qd.Project{}), mock.RequireAuth())

	// Set the initial claims fixture.
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"write:nothing"},
	}
	claims.Subject = ulids.New().String()

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectCreate(ctx, &api.Project{TenantID: "01GMBVR86186E0EKCHQK4ESJB1", Name: "project001"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if a project ID exists.
	claims.OrgID = "01GMBVR86186E0EKCHQK4ESJB1"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "project id cannot be specified on create", "expected error when project id exists")

	// Should return an error if a project name does not exist.
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: ""})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	// Should return an error if the tenant ID is missing from the request
	_, err = suite.client.ProjectCreate(ctx, &api.Project{ID: "", Name: "project001"})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant id is missing")

	// Create a project test fixture.
	req := &api.Project{
		Name:     "project001",
		TenantID: "01GMBVR86186E0EKCHQK4ESJB1",
	}

	project, err := suite.client.ProjectCreate(ctx, req)
	require.NoError(err, "could not add project")
	require.Equal(req.Name, project.Name)
	require.Equal(claims.Name, project.Owner.Name, "expected owner name to be set")
	require.Equal(gravatar.New(claims.Email, nil), project.Owner.Picture, "expected owner gravatar to be set")
	require.NotEmpty(project.Created, "project created should not be empty")
	require.NotEmpty(project.Modified, "project modified should not be empty")

	// Should return an error if the Quarterdeck returns an error
	suite.quarterdeck.OnProjectsCreate(mock.UseError(http.StatusInternalServerError, "could not create project"), mock.RequireAuth())
	_, err = suite.client.ProjectCreate(ctx, req)
	suite.requireError(err, http.StatusInternalServerError, "could not create project", "expected error when quarterdeck returns an error")

	// Quarterdeck mock should have been called
	require.Equal(2, suite.quarterdeck.ProjectsCreateCount(), "expected quarterdeck mock to be called")
}

func (suite *tenantTestSuite) TestProjectDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		OrgID:       ulids.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		TenantID:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:          ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		OwnerID:     ulids.New(),
		Name:        "project001",
		Description: "My first project",
		Created:     time.Now().Add(-time.Hour),
		Modified:    time.Now(),
	}
	key, err := project.Key()
	require.NoError(err, "could not create project key")

	// Marshal the project data with msgpack.
	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	member := &db.Member{
		OrgID: project.OrgID,
		ID:    ulids.New(),
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member fixture")

	// Call the OnGet method and return test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: key,
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: projectData,
			}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		case db.MembersNamespace:
			return &pb.GetReply{
				Value: memberData,
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", gr.Namespace)
		}
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectDetail(ctx, project.ID.String())
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the project id is not parseable
	claims.OrgID = "01GKKYAWC4PA72YC53RVXAEC67"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectDetail(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project does not exist")

	rep, err := suite.client.ProjectDetail(ctx, project.ID.String())
	require.NoError(err, "could not retrieve project")
	require.Equal(project.ID.String(), rep.ID, "expected project id to match")
	require.Equal(member.Name, rep.Owner.Name, "expected project owner name to match")
	require.Equal(member.Picture(), rep.Owner.Picture, "expected project owner picture to match")
	require.Equal(project.Name, rep.Name, "expected project name to match")
	require.Equal(project.Description, rep.Description, "expected project description to match")
	require.Equal(project.Created.Format(time.RFC3339Nano), rep.Created, "expected project created to match")
	require.Equal(project.Modified.Format(time.RFC3339Nano), rep.Modified, "expected project modified to match")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.ProjectDetail(ctx, project.ID.String())
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project ID is not found")
}

func (suite *tenantTestSuite) TestProjectUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		OrgID:       ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		TenantID:    ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:          ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		OwnerID:     ulids.New(),
		Name:        "project001",
		Description: "My first project",
	}

	key, err := project.Key()
	require.NoError(err, "could not create project key")

	// Marshal the project data with msgpack.
	data, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	// Trtl Get should return the project key or the project data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: key,
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: data,
			}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", gr.Namespace)
		}
	}

	// OnPut method should return a success response.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}
	claims.Subject = ulids.New().String()

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditProjects}

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the project ID is not parseable.
	claims.OrgID = "01GKKYAWC4PA72YC53RVXAEC67"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "invalid"})
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project does not exist")

	// Should return an error if the project name is missing.
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67"})
	suite.requireError(err, http.StatusBadRequest, "project name is required", "expected error when project name does not exist")

	// Should return an error if the project description is too long.
	_, err = suite.client.ProjectUpdate(ctx, &api.Project{ID: "01GKKYAWC4PA72YC53RVXAEC67", Name: "project001", Description: strings.Repeat("a", 2001)})
	suite.requireError(err, http.StatusBadRequest, "project description is too long", "expected error when project description is too long")

	req := &api.Project{
		ID:          "01GKKYAWC4PA72YC53RVXAEC67",
		TenantID:    "01GMTWFK4XZY597Y128KXQ4WHP",
		Name:        "project001",
		Description: "New project",
	}

	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	rep, err := suite.client.ProjectUpdate(ctx, req)
	require.NoError(err, "could not update project")
	require.NotEqual(req.ID, "01GMTWFK4XZY597Y128KXQ4WHP", "project id should not match")
	require.Equal(rep.Name, req.Name, "expected project name to match")
	require.Equal(rep.Description, req.Description, "expected project description to match")
	require.NotEmpty(rep.Created, "expected project created to be set")
	require.NotEmpty(rep.Modified, "expected project modified to be set")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.ProjectUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project ID is not found")
}

func (suite *tenantTestSuite) TestProjectPatch() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	projectID := ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5GA2G")
	ownerID := ulid.MustParse("02ABC8J5YWH4DCYJ6CZ2P5GA3H")

	trtl := db.GetMock()
	defer trtl.Reset()

	project := &db.Project{
		OrgID:       orgID,
		TenantID:    tenantID,
		ID:          projectID,
		OwnerID:     ownerID,
		Name:        "Bouldering Project",
		Description: "Track your epic bouldering sends",
	}

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	projectKey, err := project.Key()
	require.NoError(err, "could not create project key from fixture")

	member := &db.Member{
		OrgID: orgID,
		ID:    ownerID,
		Name:  "Magnus Midtbø",
	}

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	memberKey, err := member.Key()
	require.NoError(err, "could not create member key from fixture")

	newOwner := &db.Member{
		OrgID: orgID,
		ID:    ulids.New(),
		Name:  "Adam Ondra",
	}

	newOwnerData, err := newOwner.MarshalValue()
	require.NoError(err, "could not marshal new owner data")

	newOwnerKey, err := newOwner.Key()
	require.NoError(err, "could not create new owner key from fixture")

	// Trtl Get should return the requested resources by default
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: projectKey,
			}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: orgID[:],
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: projectData,
			}, nil
		case db.MembersNamespace:
			switch {
			case bytes.Equal(gr.Key, memberKey):
				return &pb.GetReply{
					Value: memberData,
				}, nil
			case bytes.Equal(gr.Key, newOwnerKey):
				return &pb.GetReply{
					Value: newOwnerData,
				}, nil
			default:
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, errors.New("unexpected namespace")
		}
	}

	// Trtl Put should succeed by default
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       ulids.New().String(),
		Permissions: []string{"edit:nothing"},
	}

	// Endpoint must be authenticated
	req := &api.Project{}
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.ProjectPatch(ctx, projectID.String(), req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when not authenticated")

	// User must have the right permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectPatch(ctx, projectID.String(), req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have the right permissions")

	// Error should be returned if the projectID is not parseable
	claims.Permissions = []string{perms.EditProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.ProjectPatch(ctx, "invalid", req)
	suite.requireError(err, http.StatusNotFound, responses.ErrProjectNotFound, "expected error when project id is not parseable")

	// User must be in the same org as the project
	req.Name = "New Project Name"
	_, err = suite.client.ProjectPatch(ctx, projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrProjectNotFound, "expected error when user is not in the same org as the project")

	// Set the user claims to the correct org for the project
	claims.OrgID = orgID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Test field validation
	testCases := []struct {
		req *api.Project
		err string
	}{
		{&api.Project{Description: strings.Repeat("a", db.MaxDescriptionLength+1)}, "validation error: project description is too long"},
		{&api.Project{Status: "Archived"}, api.InvalidFieldError("status").Error()},
		{&api.Project{Name: "Moonboard Sends", Description: "Crushin' v7s on the Moonboard"}, ""},
		{&api.Project{Owner: api.Member{ID: newOwner.ID.String()}}, ""},
		{&api.Project{Owner: api.Member{ID: newOwner.ID.String(), Name: "Adam Ondra"}}, ""},
		{&api.Project{Owner: api.Member{}, Name: "Moonboard Sends"}, ""},
	}

	for i, tc := range testCases {
		_, err = suite.client.ProjectPatch(ctx, projectID.String(), tc.req)
		if tc.err == "" {
			require.NoError(err, "expected no error with valid fields")
		} else {
			suite.requireError(err, http.StatusBadRequest, tc.err, fmt.Sprintf("expected error when field validation fails for test case %d", i))
		}
	}

	// Should return an error if the owner is not parseable as a ULID
	req.Owner = api.Member{ID: "invalid"}
	_, err = suite.client.ProjectPatch(ctx, projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when owner is not parseable as a ULID")

	// Should return an error if the owner does not exist
	req.Owner = api.Member{ID: ulids.New().String()}
	_, err = suite.client.ProjectPatch(ctx, projectID.String(), req)
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when owner is not found")
}

func (suite *tenantTestSuite) TestProjectDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := "01GMTWFK4XZY597Y128KXQ4WHP"
	projectID := "01GKKYAWC4PA72YC53RVXAEC67"
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	key, err := db.CreateKey(ulid.MustParse(tenantID), ulid.MustParse(projectID))
	require.NoError(err, "could not create project key")

	keyData, err := key.MarshalValue()
	require.NoError(err, "could not marshal the project key")

	project := &db.Project{
		OrgID:    ulids.New(),
		TenantID: ulid.MustParse(tenantID),
		ID:       ulid.MustParse(projectID),
	}

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal the project")

	// OnGet method should return the project key or the project data.
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
		case db.OrganizationNamespace:
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", gr.Namespace)
		}
	}

	// Call the OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	err = suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client claims")
	err = suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.DeleteProjects}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.ProjectDelete(ctx, project.ID.String())
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the project id is not parseable.
	claims.OrgID = project.ID.String()
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.ProjectDelete(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project does not exist")

	err = suite.client.ProjectDelete(ctx, projectID)
	require.NoError(err, "could not delete project")

	// Should return an error if the project ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: project.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	err = suite.client.ProjectDelete(ctx, projectID)
	suite.requireError(err, http.StatusNotFound, "project not found", "expected error when project ID is not found")
}

func (suite *tenantTestSuite) TestUpdateProjectStats() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Project info to return on the Quarterdeck call
	orgID := ulids.New()
	userID := ulids.New()
	projectID := ulids.New()
	qdProject := &qd.Project{
		OrgID:        orgID,
		ProjectID:    projectID,
		APIKeysCount: 2,
	}

	// Project info to return on the Ensign call
	enProject := &en.ProjectInfo{
		Topics:         7,
		ReadonlyTopics: 4,
	}

	expectedAPIKeys := uint64(2)
	expectedTopics := uint64(3)

	// Project that trtl should be updating
	project := &db.Project{
		OrgID:    orgID,
		TenantID: ulids.New(),
		ID:       projectID,
		OwnerID:  ulids.New(),
		Name:     "project-1",
	}

	projectData, err := project.MarshalValue()
	require.NoError(err, "could not marshal project data")

	objectKey, err := project.Key()
	require.NoError(err, "could not create project key")

	// Init the trtl mock
	trtl := db.GetMock()
	defer trtl.Reset()

	// Initial trtl get should return the project
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.KeysNamespace:
			return &pb.GetReply{
				Value: objectKey[:],
			}, nil
		case db.ProjectNamespace:
			return &pb.GetReply{
				Value: projectData,
			}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Initial trtl put should verify that api keys and topics were counted correctly
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if !bytes.Equal(in.Key, objectKey) || in.Namespace != db.ProjectNamespace {
			return nil, status.Error(codes.FailedPrecondition, "unexpected key or namespace")
		}

		p := &db.Project{}
		if err := p.UnmarshalValue(in.Value); err != nil {
			return nil, err
		}

		require.Equal(expectedAPIKeys, p.APIKeys, "api keys were not counted correctly")
		require.Equal(expectedTopics, p.Topics, "topics were not counted correctly")
		return &pb.PutReply{}, nil
	}

	// Initial quarterdeck mock should return the project info
	suite.quarterdeck.OnProjectsDetail(projectID.String(), mock.UseStatus(http.StatusOK), mock.UseJSONFixture(qdProject), mock.RequireAuth())

	// Project access should return the access token
	login := &qd.LoginReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	suite.quarterdeck.OnProjectsAccess(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(login), mock.RequireAuth())

	// Initial ensign mock should return the project info
	suite.ensign.OnInfo = func(ctx context.Context, in *en.InfoRequest) (*en.ProjectInfo, error) {
		return enProject, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Should return an error if credentials are not in the context.
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	expected := statusMessage(http.StatusUnauthorized, "missing authorization header")
	suite.requireMultiError(err, expected, expected)

	// Successfully updating the project
	ctx, err = suite.ContextWithClaims(ctx, claims)
	require.NoError(err, "could not set claims on the context")
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	require.NoError(err, "could not update project stats")

	// Test that the topic count is 0 if ensign returns inconsistent values
	enProject.ReadonlyTopics = 10
	expectedTopics = 0
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	require.NoError(err, "could not update project stats")

	// Test that no topics are counted if the ensign call fails
	suite.ensign.OnInfo = func(ctx context.Context, in *en.InfoRequest) (*en.ProjectInfo, error) {
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	require.ErrorIs(err, status.Error(codes.Unauthenticated, "missing credentials"), "expected an error if only the ensign rpc fails")

	// Test that no API keys are counted if the quarterdeck call fails
	enProject.ReadonlyTopics = 4
	suite.ensign.OnInfo = func(ctx context.Context, in *en.InfoRequest) (*en.ProjectInfo, error) {
		return enProject, nil
	}
	expectedTopics = 3
	suite.quarterdeck.OnProjectsDetail(projectID.String(), mock.UseError(http.StatusUnauthorized, "invalid claims"), mock.RequireAuth())
	expectedAPIKeys = 0
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	require.ErrorContains(err, statusMessage(http.StatusUnauthorized, "invalid claims"), "expected an error if only the quarterdeck rpc fails")

	// Test that the method returns an error if both rpcs fail
	suite.ensign.OnInfo = func(ctx context.Context, in *en.InfoRequest) (*en.ProjectInfo, error) {
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}
	expectedTopics = 0
	err = suite.srv.UpdateProjectStats(ctx, userID, projectID)
	suite.requireMultiError(err, statusMessage(http.StatusUnauthorized, "invalid claims"), status.Error(codes.Unauthenticated, "missing credentials").Error())
}
