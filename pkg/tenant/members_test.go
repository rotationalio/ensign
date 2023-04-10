package tenant_test

import (
	"bytes"
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
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestMemberList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")

	members := []*db.Member{
		{
			OrgID:        orgID,
			ID:           ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7"),
			Email:        "test@testing.com",
			Name:         "member001",
			Role:         "Admin",
			Created:      time.Unix(1670424445, 0),
			Modified:     time.Unix(1670424445, 0),
			LastActivity: time.Unix(1670424445, 0),
			DateAdded:    time.Unix(1670424445, 0),
		},

		{
			OrgID:        orgID,
			ID:           ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Email:        "wilder@testing.com",
			Name:         "member002",
			Role:         "Member",
			Created:      time.Unix(1673659941, 0),
			Modified:     time.Unix(1673659941, 0),
			LastActivity: time.Unix(1673659941, 0),
			DateAdded:    time.Unix(1673659941, 0),
		},

		{
			OrgID:        orgID,
			ID:           ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Email:        "moore@testing.com",
			Name:         "member003",
			Role:         "Admin",
			Created:      time.Unix(1674073941, 0),
			Modified:     time.Unix(1674073941, 0),
			LastActivity: time.Unix(1674073941, 0),
			DateAdded:    time.Unix(1674073941, 0),
		},
	}

	prefix := orgID[:]
	namespace := "members"

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
		for _, member := range members {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, member.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := member.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       member.ID[:],
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

	// Endpoint must be authenticated
	_, err := suite.client.MemberList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberList(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Retrieve all members.
	rep, err := suite.client.MemberList(ctx, req)
	require.NoError(err, "could not list members")
	require.Len(rep.Members, 3, "expected 3 members")
	require.Empty(rep.NextPageToken, "expected next page token")

	// Verify member data has been populated.
	for i := range members {
		require.Equal(members[i].ID.String(), rep.Members[i].ID, "expected member id to match")
		require.Equal(members[i].Email, rep.Members[i].Email, "expected member email to match")
		require.Equal(members[i].Name, rep.Members[i].Name, "expected member name to match")
		require.Equal(members[i].Role, rep.Members[i].Role, "expected member role to match")
		require.Equal(members[i].Created.Format(time.RFC3339Nano), rep.Members[i].Created, "expected member created time to match")
		require.Equal(members[i].Modified.Format(time.RFC3339Nano), rep.Members[i].Modified, "expected member modified time to match")
		require.Equal(members[i].LastActivity.Format(time.RFC3339), rep.Members[i].LastActivity, "expected last activity to match")
		require.Equal(members[i].DateAdded.Format(time.RFC3339), rep.Members[i].DateAdded, "expected date added to match")
	}

	// Set page size to test pagination.
	req.PageSize = 2
	rep, err = suite.client.MemberList(ctx, req)
	require.NoError(err, "could not list members")
	require.Len(rep.Members, 2, "expected 2 members")
	require.NotEmpty(rep.NextPageToken, "next page token expected")

	// Test next page token.
	req.NextPageToken = rep.NextPageToken
	rep2, err := suite.client.MemberList(ctx, req)
	require.NoError(err, "could not list members")
	require.Len(rep2.Members, 1, "expected 1 member")
	require.NotEqual(rep.Members[0].ID, rep2.Members[0].ID, "should not have same member ID")
	require.Empty(rep2.NextPageToken, "should be empty when a next page does not exist")

	// Limit maximum number of requests to 3, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for i := 0; i < 3; i++ {
		page, err := suite.client.MemberList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Members)

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
		OrgID:       "0000000000000000",
		Permissions: []string{perms.ReadCollaborators},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.MemberList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestMemberCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	userOrg := "01GMBVR86186E0EKCHQK4ESJB1"
	orgID := ulid.MustParse(userOrg)
	email := "newuser@example.com"
	role := perms.RoleMember

	members := []*db.Member{
		{
			ID:    ulids.New(),
			Email: "leopold.wentzel@gmail.com",
		},
	}

	// Configure initial Trtl Cursor to return the requesting user's member data
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, orgID[:]) || in.Namespace != db.MembersNamespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, member := range members {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, member.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := member.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       member.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Configure quarterdeck mock to return the invite token
	invite := &qd.UserInviteReply{
		UserID:    ulids.New(),
		OrgID:     orgID,
		Email:     email,
		Role:      role,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339Nano),
		CreatedBy: members[0].ID,
		Created:   time.Now().Format(time.RFC3339Nano),
	}
	suite.quarterdeck.OnInvites("", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(invite))

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMBVR86186E0EKCHQK4ESJB1",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.MemberCreate(ctx, &api.Member{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberCreate(ctx, &api.Member{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.AddCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if member id exists.
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member-example", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member id cannot be specified on create", "expected error when member id exists")

	// Should return an error if member email does not exist.
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "", Name: "member-example", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member email is required", "expected error when member email does not exist")

	// Should return an error if the member role does not exist.
	_, err = suite.client.MemberCreate(ctx, &api.Member{Email: "test@testing.com", ID: "", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	// Should return an error if the member role is invalid.
	_, err = suite.client.MemberCreate(ctx, &api.Member{Email: "test@testing.com", Role: "invalid"})
	suite.requireError(err, http.StatusBadRequest, "unknown member role", "expected error when member role is invalid")

	// Create a member test fixture
	req := &api.Member{
		Role:  role,
		Email: email,
	}

	rep, err := suite.client.MemberCreate(ctx, req)
	require.NoError(err, "could not add member")
	require.NotEmpty(rep.ID, "expected non-zero ulid to be populated")
	require.Equal(req.Email, rep.Email, "expected member email to match")
	require.Empty(rep.Name, "expected member name to be empty")
	require.Equal(req.Role, rep.Role, "expected member role to match")
	require.Equal(rep.Status, string(db.MemberStatusPending), "expected member status to be pending")
	require.NotEmpty(rep.Created, "expected created time to be populated")
	require.NotEmpty(rep.Modified, "expected modified time to be populated")
	require.NotEmpty(rep.LastActivity, "expected last activity time to be populated")
	require.NotEmpty(rep.DateAdded, "expected date added timem to be populated")

	// Should not be able to create a member with the same email.
	members = append(members, &db.Member{
		ID:    invite.UserID,
		Email: rep.Email,
	})

	_, err = suite.client.MemberCreate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "team member already exists with this email address", "expected error when member email already exists")

	// Test that the endpoint returns an error if quarterdeck returns an error.
	suite.quarterdeck.OnInvites("", mock.UseError(http.StatusUnauthorized, "invalid user claims"))
	req.Email = "other@example.com"
	_, err = suite.client.MemberCreate(ctx, req)
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when quarterdeck returns an error")

	// Create a test fixture.
	test := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "0000000000000000",
		Permissions: []string{perms.AddCollaborators},
	}

	// User org id is required.
	require.NoError(suite.SetClientCredentials(test))
	_, err = suite.client.MemberCreate(ctx, &api.Member{})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")
}

func (suite *tenantTestSuite) TestMemberDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		ID:       ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:     "member-example",
		Role:     "Admin",
		Created:  time.Now().Add(-time.Hour),
		Modified: time.Now(),
	}

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// OnGet should return member data or member ID.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.GetReply{Value: data}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: member.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err = suite.client.MemberDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing in the claims
	_, err = suite.client.MemberDetail(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

	// Should return an error if org verification fails.
	claims.OrgID = "01GMBVR86186E0EKCHQK4ESJB1"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberDetail(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the member does not exist.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberDetail(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	// Create a member test fixture.
	req := &api.Member{
		ID:   "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name: "member-example",
		Role: "Admin",
	}
	rep, err := suite.client.MemberDetail(ctx, req.ID)
	require.NoError(err, "could not retrieve member")
	require.Equal(req.ID, rep.ID, "expected member id to match")
	require.Equal(req.Name, rep.Name, "expected member name to match")
	require.Equal(req.Role, rep.Role, "expected member role to match")
	require.NotEmpty(rep.Created, "expected created time to be populated")
	require.NotEmpty(rep.Modified, "expected modified time to be populated")

	// Test the not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: member.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}

	_, err = suite.client.MemberDetail(ctx, req.ID)
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")
}

func (suite *tenantTestSuite) TestMemberUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		OrgID:  ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:     ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Email:  "test@testing.com",
		Name:   "member001",
		Role:   "Admin",
		Status: "Confirmed",
	}

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// OnGet should return member data or member ID.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.GetReply{Value: data}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: member.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
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
		OrgID:       "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the member ID is not parseable.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "invalid", Email: "test@testing.com", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	// Should return an error if the member email is not provided.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member email is required", "expected error when member email does not exist")

	// Should return an error if the member name is not provided.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member name is required", "expected error when member name does not exist")

	// Should return an error if the member role is not provided.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Name: "member001"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	// Should return an error if the member role provided is not valid.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Email: "test@testing.com", Name: "member001", Role: "Guest"})
	suite.requireError(err, http.StatusBadRequest, "unknown member role", "expected error when member role is not valid")

	req := &api.Member{
		ID:    "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Email: "test@testing.com",
		Name:  "member001",
		Role:  "Admin",
	}

	rep, err := suite.client.MemberUpdate(ctx, req)
	require.NoError(err, "could not update member")
	require.NotEqual(req.ID, "01GM8MEZ097ZC7RQRCWMPRPS0T", "member id should not match")
	require.Equal(rep.Email, req.Email, "expected member email to match")
	require.Equal(rep.Name, req.Name, "expected member name to match")
	require.Equal(rep.Role, req.Role, "expected member role to match")
	require.NotEmpty(rep.Created, "expected created time to be populated")
	require.NotEmpty(rep.Modified, "expected modified time to be populated")

	// Test the not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace == db.OrganizationNamespace {
			return &pb.GetReply{
				Value: member.ID[:],
			}, nil
		}
		return nil, status.Error(codes.NotFound, "not found")
	}
	_, err = suite.client.MemberUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")
}

func (suite *tenantTestSuite) TestMemberRoleUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	orgID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	prefix := orgID[:]
	namespace := "members"

	member := &db.Member{
		OrgID:  orgID,
		ID:     ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Email:  "test@testing.com",
		Name:   "member001",
		Role:   perms.RoleOwner,
		Status: "Confirmed",
	}

	// Marshal the member data with msgpack.
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// OnGet should return member data or member ID.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.GetReply{Value: data}, nil
		case db.OrganizationNamespace:
			return &pb.GetReply{Value: member.ID[:]}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Create members in the database.
	members := []*db.Member{
		{
			OrgID:  orgID,
			ID:     ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
			Email:  "test@testing.com",
			Name:   "member001",
			Role:   perms.RoleOwner,
			Status: db.MemberStatusConfirmed,
		},
		{
			OrgID:  orgID,
			ID:     ulid.MustParse("01GX1FCEYW8NFYRBHAFFHWD45C"),
			Email:  "ryan@testing.com",
			Name:   "member002",
			Role:   perms.RoleOwner,
			Status: db.MemberStatusConfirmed,
		},

		{
			OrgID:  orgID,
			ID:     ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Email:  "wilder@testing.com",
			Name:   "member003",
			Role:   perms.RoleAdmin,
			Status: db.MemberStatusConfirmed,
		},

		{
			OrgID:  orgID,
			ID:     ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Email:  "moore@testing.com",
			Name:   "member004",
			Role:   perms.RoleMember,
			Status: db.MemberStatusConfirmed,
		},
	}

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some members and terminate.
		for _, dbMember := range members {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, dbMember.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				dbMemData, err := dbMember.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       dbMember.ID[:],
					Value:     dbMemData,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated.
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permission")

	// Set valid permissions for the rest of the tests.
	claims.Permissions = []string{perms.EditCollaborators, perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the orgID is missing in the claims.
	_, err = suite.client.MemberRoleUpdate(ctx, "invalid", &api.UpdateMemberParams{Role: perms.RoleAdmin})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

	// Should return an error if org verification fails.
	claims.OrgID = "01GWT0E850YBSDQH0EQFXRCMGB"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: perms.RoleAdmin})
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the member does not exist.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberRoleUpdate(ctx, "invalid", &api.UpdateMemberParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	// Should return an error if the member role is not provided.
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	// Should return an errror if the member role provided is not valid.
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: "Viewer"})
	suite.requireError(err, http.StatusBadRequest, "unknown member role", "expected error when member role is not valid")

	// Should return an error if the member id in the database does not match the id in the URL.
	_, err = suite.client.MemberRoleUpdate(ctx, "01GQ2XB2SCGY5RZJ1ZGYSEMNDE", &api.UpdateMemberParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusInternalServerError, "member id does not match id in URL", "expected error when member id does not match")

	// Set database to have one owner. Should return an error if org does not have an owner.
	members[1].Role = perms.RoleAdmin
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusBadRequest, "organization must have at least one owner", "expected error when org does not have an owner")

	// Set more than one member role to owner for test.
	members[1].Role = perms.RoleOwner
	rep, err := suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateMemberParams{Role: perms.RoleObserver})
	require.NoError(err, "could not update member role")
	require.Equal(rep.Role, perms.RoleObserver, "expected member role to update")
}

func (suite *tenantTestSuite) TestMemberDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	memberID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// OnGet returns the memberID.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: memberID[:],
		}, nil
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
	err := suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.RemoveCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing
	err = suite.client.MemberDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when member ID is not parseable")

	// Should return an error if org verification fails.
	claims.OrgID = "01GMBVR86186E0EKCHQK4ESJB1"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.MemberDelete(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusUnauthorized, "could not verify organization", "expected error when org verification fails")

	// Should return an error if the member does not exist.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.MemberDelete(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	err = suite.client.MemberDelete(ctx, memberID.String())
	require.NoError(err, "could not delete member")

	// Should return an error if the member ID is parsed but not found.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}

	err = suite.client.MemberDelete(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member ID is not found")
}

func (suite *tenantTestSuite) TestInvitePreview() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initial Quarterdeck mock returns a valid preview
	preview := &qd.UserInvitePreview{
		Email:       "leopold.wentzel@gmail.com",
		OrgName:     "Events R Us",
		InviterName: "Geoffrey",
		Role:        "Member",
		UserExists:  true,
	}
	suite.quarterdeck.OnInvites("token1234", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(preview))

	// Test successful preview request
	rep, err := suite.client.InvitePreview(ctx, "token1234")
	require.NoError(err, "could not get preview invite")
	require.Equal(preview.Email, rep.Email, "expected email to match")
	require.Equal(preview.OrgName, rep.OrgName, "expected org name to match")
	require.Equal(preview.InviterName, rep.InviterName, "expected inviter name to match")
	require.Equal(preview.Role, rep.Role, "expected role to match")
	require.True(rep.HasAccount, "expected user to exist")

	// Test invalid invitation response is correctly forwarded by Tenant
	suite.quarterdeck.OnInvites("token1234", mock.UseError(http.StatusBadRequest, "invalid invitation"))
	_, err = suite.client.InvitePreview(ctx, "token1234")
	suite.requireError(err, http.StatusBadRequest, "invalid invitation", "expected error when token is invalid")
}
