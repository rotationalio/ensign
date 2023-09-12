package tenant_test

import (
	"bytes"
	"context"
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
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtlmock "github.com/trisacrypto/directory/pkg/trtl/mock"
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
			JoinedAt:     time.Unix(1670424445, 0),
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
			JoinedAt:     time.Unix(1673659941, 0),
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
			JoinedAt:     time.Unix(1674073941, 0),
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
		require.Equal(members[i].LastActivity.Format(time.RFC3339), rep.Members[i].LastActivity, "expected last activity to match")
		require.Equal(members[i].JoinedAt.Format(time.RFC3339), rep.Members[i].DateAdded, "expected date added to match")
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
	organization := "Cloud Services"
	workspace := "cloud-services"

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
		UserID:       ulids.New(),
		OrgID:        orgID,
		Email:        email,
		Role:         role,
		Organization: organization,
		Workspace:    workspace,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339Nano),
		CreatedBy:    members[0].ID,
		Created:      time.Now().Format(time.RFC3339Nano),
	}
	suite.quarterdeck.OnInvitesCreate(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(invite), mock.RequireAuth())

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
	require.NotEmpty(rep.Organization, "expected organization to be populated")
	require.NotEmpty(rep.Workspace, "expected workspace to be populated")
	require.True(rep.Invited, "expected member to have the invited flag set")
	require.Equal(rep.OnboardingStatus, db.MemberStatusPending.String(), "expected member status to be pending")
	require.NotEmpty(rep.Created, "expected created time to be populated")

	// Should not be able to create a member with the same email.
	members = append(members, &db.Member{
		ID:    invite.UserID,
		Email: rep.Email,
	})

	_, err = suite.client.MemberCreate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, "team member already exists with this email address", "expected error when member email already exists")

	// Test that the endpoint returns an error if quarterdeck returns an error.
	suite.quarterdeck.OnInvitesCreate(mock.UseError(http.StatusUnauthorized, "invalid user claims"))
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

	// Test the not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
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
		OrgID:        ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:           ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Email:        "test@testing.com",
		Name:         "member001",
		Organization: "Cloud Services",
		Workspace:    "cloud-services",
		Role:         "Admin",
		JoinedAt:     time.Now(),
	}

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// OnGet should return member data or member ID.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.GetReply{Value: data}, nil
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

	// Should return an error if the member ID is not parseable.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "invalid", Email: "test@testing.com", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when member does not exist")

	// Should return validation errors
	req := &api.Member{
		ID:           "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Email:        "test@testing.com",
		Name:         strings.Repeat("a", 4096),
		Organization: "Rotational Labs",
		Workspace:    "rotational-io",
		Role:         "Admin",
	}
	expected := &api.FieldValidationErrors{
		{
			Field: "name",
			Err:   db.ErrNameTooLong.Error(),
			Index: -1,
		},
	}
	_, err = suite.client.MemberUpdate(ctx, req)
	suite.requireError(err, http.StatusBadRequest, expected.Error(), "expected validation errors")

	// Test updating the member record with a valid name
	req.Name = "member002"
	rep, err := suite.client.MemberUpdate(ctx, req)
	require.NoError(err, "could not update member")
	require.Equal(member.ID.String(), rep.ID, "expected member id to match")
	require.Equal(req.Name, rep.Name, "expected member name to be updated")
	require.Equal(member.Email, rep.Email, "expected member email to be unchanged")
	require.Equal(member.Organization, rep.Organization, "expected organization to be unchanged")
	require.Equal(member.Workspace, rep.Workspace, "expected workspace to be unchanged")

	// Test the not found path
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	_, err = suite.client.MemberUpdate(ctx, req)
	suite.requireError(err, http.StatusNotFound, responses.ErrMemberNotFound, "expected error when member does not exist")
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

	member := &db.Member{
		OrgID: orgID,
		ID:    ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Email: "test@testing.com",
		Name:  "member001",
		Role:  perms.RoleOwner,
	}

	userID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	userReply := &qd.User{
		UserID: userID,
		Role:   perms.RoleObserver,
	}

	suite.quarterdeck.OnUsersRoleUpdate(userID.String(), mock.UseStatus(http.StatusOK), mock.UseJSONFixture(userReply), mock.RequireAuth())

	// Marshal the member data with msgpack.
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// OnGet should return member data or member ID.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.GetReply{Value: data}, nil
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
			OrgID: orgID,
			ID:    ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
			Email: "test@testing.com",
			Name:  "member001",
			Role:  perms.RoleOwner,
		},
		{
			OrgID: orgID,
			ID:    ulid.MustParse("01GX1FCEYW8NFYRBHAFFHWD45C"),
			Email: "ryan@testing.com",
			Name:  "member002",
			Role:  perms.RoleOwner,
		},

		{
			OrgID: orgID,
			ID:    ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Email: "wilder@testing.com",
			Name:  "member003",
			Role:  perms.RoleAdmin,
		},

		{
			OrgID: orgID,
			ID:    ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Email: "moore@testing.com",
			Name:  "member004",
			Role:  perms.RoleMember,
		},
	}

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != db.MembersNamespace {
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
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions.
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have correct permission")

	// Set valid permissions for the rest of the tests.
	claims.Permissions = []string{perms.EditCollaborators, perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the orgID is missing in the claims.
	_, err = suite.client.MemberRoleUpdate(ctx, "invalid", &api.UpdateRoleParams{Role: perms.RoleAdmin})
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when org id is missing or not a valid ulid")

	// Should return an error if the member role is not provided.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{})
	suite.requireError(err, http.StatusBadRequest, "team member role is required", "expected error when member role does not exist")

	// Should return an errror if the member role provided is not valid.
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: "Viewer"})
	suite.requireError(err, http.StatusBadRequest, "unknown team member role", "expected error when member role is not valid")

	// Should return an error if the member does not exist.
	_, err = suite.client.MemberRoleUpdate(ctx, "invalid", &api.UpdateRoleParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	// Should return an error if org does not have an owner.
	members[0].Role = perms.RoleMember
	members[1].Role = perms.RoleAdmin
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusInternalServerError, "could not update member role", "expected error when org does not have an owner")

	// Should return an error if the member is not confirmed.
	members[0].Role = perms.RoleOwner
	members[1].Role = perms.RoleOwner
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleObserver})
	suite.requireError(err, http.StatusBadRequest, "cannot update role for pending team member", "expected error when member is not confirmed")

	// Should return an error if the member already has the specified role.
	member.Organization = "testorg"
	member.Workspace = "testorg"
	member.ProfessionSegment = "Personal"
	member.DeveloperSegment = []string{"Application Development"}
	data, err = member.MarshalValue()
	require.NoError(err, "could not marshal the member")
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleOwner})
	suite.requireError(err, http.StatusBadRequest, "team member already has the requested role", "expected error when member already has the specified role")

	// Successfully updating the member role.
	reply, err := suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleObserver})
	require.NoError(err, "expected no error when updating the member role")
	require.Equal(perms.RoleObserver, reply.Role, "expected the member role to be updated")

	// Test Tenant returns an error if Quarterdeck returns an error
	suite.quarterdeck.OnUsersRoleUpdate(userID.String(), mock.UseError(http.StatusBadRequest, "organization must have at least one owner"), mock.RequireAuth())
	_, err = suite.client.MemberRoleUpdate(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV", &api.UpdateRoleParams{Role: perms.RoleAdmin})
	suite.requireError(err, http.StatusBadRequest, "organization must have at least one owner", "expected error when quarterdeck returns an error")
}

func (suite *tenantTestSuite) TestMemberDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	orgID := ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	memberID := ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE")
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	prefix := orgID[:]

	member := &db.Member{
		OrgID: orgID,
		ID:    memberID,
		Email: "cmoon@test.com",
		Name:  "Cindy Moon",
		Role:  perms.RoleOwner,
	}

	members := []*db.Member{
		{
			OrgID: orgID,
			ID:    memberID,
			Email: "cmoon@test.com",
			Name:  "Cindy Moon",
			Role:  perms.RoleOwner,
		},
		{
			OrgID: orgID,
			ID:    ulid.MustParse("01GX1FCEYW8NFYRBHAFFHWD45C"),
			Email: "leopold.wentzel@gmail.com",
			Name:  "Leopold Wentzel",
			Role:  perms.RoleAdmin,
		},
	}

	// Marshal member data.
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	// OnGet returns the member data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Call the OnCursor method and add some members to the database.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) (err error) {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != db.MembersNamespace {
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

	// Call the OnDelete method and return a DeleteReply.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return &pb.DeleteReply{}, nil
	}

	// Initial Quarterdeck mock returns a confirmation token
	qdReply := &qd.UserRemoveReply{
		APIKeys: []string{"key1", "key2"},
		Token:   "token",
	}
	suite.quarterdeck.OnUsersRemove(memberID.String(), mock.UseStatus(http.StatusOK), mock.UseJSONFixture(qdReply), mock.RequireAuth())

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"delete:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.RemoveCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should error if the orgID is missing
	_, err = suite.client.MemberDelete(ctx, "invalid")
	suite.requireError(err, http.StatusUnauthorized, "invalid user claims", "expected error when member ID is not parseable")

	// Should return an error if the member does not exist.
	claims.OrgID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberDelete(ctx, "invalid")
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member does not exist")

	// Should return an error if org does not have an owner.
	members[0].Role = perms.RoleAdmin
	_, err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusInternalServerError, "could not delete member", "expected error when org does not have an owner")

	// Should return an error if request is made to delete last org owner.
	members[0].Role = perms.RoleOwner
	_, err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusBadRequest, "organization must have at least one owner", "expected error when deleting last org owner")

	// Should return the token to the client when Quarterdeck returns a token.
	expected := &api.MemberDeleteReply{
		APIKeys: qdReply.APIKeys,
		Token:   qdReply.Token,
	}
	members[1].Role = perms.RoleOwner
	rep, err := suite.client.MemberDelete(ctx, memberID.String())
	require.NoError(err, "could not delete member")
	require.Equal(expected, rep, "expected response to match")

	// Should delete the member from the database if Quarterdeck did not return a token.
	qdReply = &qd.UserRemoveReply{
		Deleted: true,
	}
	suite.quarterdeck.OnUsersRemove(memberID.String(), mock.UseStatus(http.StatusOK), mock.UseJSONFixture(qdReply), mock.RequireAuth())
	rep, err = suite.client.MemberDelete(ctx, memberID.String())
	require.NoError(err, "could not delete member")
	require.True(rep.Deleted, "expected deleted to be returned in response")
	require.Equal(1, trtl.Calls[trtlmock.DeleteRPC], "expected delete to be called once")

	// Should return an error if Quarterdeck returns an error.
	suite.quarterdeck.OnUsersRemove(memberID.String(), mock.UseError(http.StatusInternalServerError, "could not delete user"))
	_, err = suite.client.MemberDelete(ctx, memberID.String())
	suite.requireError(err, http.StatusInternalServerError, "could not delete user", "expected error when Quarterdeck returns an error")

	// Should return an error if the member ID is parsed but not found.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Error(codes.NotFound, "member not found")
	}

	_, err = suite.client.MemberDelete(ctx, "01GQ2XB2SCGY5RZJ1ZGYSEMNDE")
	suite.requireError(err, http.StatusNotFound, "member not found", "expected error when member ID is not found")
}
