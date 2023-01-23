package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestTenantMemberList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")

	members := []*db.Member{
		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7"),
			Name:     "member001",
			Role:     "Admin",
			Created:  time.Unix(1670424445, 0),
			Modified: time.Unix(1670424445, 0),
		},

		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Name:     "member002",
			Role:     "Member",
			Created:  time.Unix(1673659941, 0),
			Modified: time.Unix(1673659941, 0),
		},

		{
			TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
			ID:       ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Name:     "member003",
			Role:     "Admin",
			Created:  time.Unix(1674073941, 0),
			Modified: time.Unix(1674073941, 0),
		},
	}

	prefix := tenantID[:]
	namespace := "members"

	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i, member := range members {
			data, err := member.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     data,
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := suite.client.TenantMemberList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantMemberList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the tenant ID is not parseable
	_, err = suite.client.TenantMemberList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant ulid", "expected error when tenant ID is missing")

	rep, err := suite.client.TenantMemberList(ctx, tenantID.String(), &api.PageQuery{})
	require.NoError(err, "could not list tenant members")
	require.Len(rep.TenantMembers, 3, "expected 3 members")

	// Test first member data has been populated.
	require.Equal(members[0].ID.String(), rep.TenantMembers[0].ID, "expected member id to match")
	require.Equal(members[0].Name, rep.TenantMembers[0].Name, "expected member name to match")
	require.Equal(members[0].Role, rep.TenantMembers[0].Role, "expected member role to match")

	// Test second member data has been populated.
	require.Equal(members[1].ID.String(), rep.TenantMembers[1].ID, "expected member id to match")
	require.Equal(members[1].Name, rep.TenantMembers[1].Name, "expected member name to match")
	require.Equal(members[1].Role, rep.TenantMembers[1].Role, "expected member role to match")

	// Test third member data has been populated.
	require.Equal(members[2].ID.String(), rep.TenantMembers[2].ID, "expected member id to match")
	require.Equal(members[2].Name, rep.TenantMembers[2].Name, "expected member name to match")
	require.Equal(members[2].Role, rep.TenantMembers[2].Role, "expected member role to match")
}

func (suite *tenantTestSuite) TestTenantMemberCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulids.New().String()
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"write:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err := suite.client.TenantMemberCreate(ctx, tenantID, &api.Member{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.TenantMemberCreate(ctx, tenantID, &api.Member{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.AddCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if tenant id is not a valid ULID.
	_, err = suite.client.TenantMemberCreate(ctx, "tenantID", &api.Member{ID: "", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant id", "expected error when tenant id does not exist")

	// Should return an error if the member id exists.
	_, err = suite.client.TenantMemberCreate(ctx, tenantID, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member-example", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member id cannot be specified on create", "expected error when member id exists")

	// Should return an error if the member name does not exist
	_, err = suite.client.TenantMemberCreate(ctx, tenantID, &api.Member{ID: "", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "tenant member name is required", "expected error when tenant member name does not exist")

	// Should return an error if the member role does not exist.
	_, err = suite.client.TenantMemberCreate(ctx, tenantID, &api.Member{ID: "", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "tenant member role is required", "expected error when tenant member role does not exist")

	tenant := &api.Tenant{
		ID: ulids.New().String(),
	}

	// Create a member test fixture
	req := &api.Member{
		Name: "member001",
		Role: "Admin",
	}

	member, err := suite.client.TenantMemberCreate(ctx, tenant.ID, req)
	require.NoError(err, "could not add member")
	require.Equal(req.Name, member.Name, "member name should match")
	require.Equal(req.Role, member.Role, "member role should match")
}

func (suite *tenantTestSuite) TestMemberList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnCursor method.
	trtl.OnCursor = func(cr *pb.CursorRequest, t pb.Trtl_CursorServer) error {
		return nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := suite.client.MemberList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberList(ctx, &api.PageQuery{})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.ReadCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// TODO: Test length of values assigned to *api.MemberPage
	_, err = suite.client.MemberList(ctx, &api.PageQuery{})
	require.NoError(err, "could not list members")
}

func (suite *tenantTestSuite) TestMemberCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
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

	// Should return an error if the member name does not exist
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member name is required", "expected error when member name does not exist")

	// Should return an error if the member role does not exist.
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	// Create a member test fixture
	req := &api.Member{
		Name: "member001",
		Role: "Admin",
	}

	member, err := suite.client.MemberCreate(ctx, req)
	require.NoError(err, "could not add member")
	require.Equal(req.Name, member.Name, "expected memeber name to match")
	require.Equal(req.Role, member.Role, "expected member role to match")
}

func (suite *tenantTestSuite) TestMemberDetail() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		ID:   ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name: "member-example",
		Role: "Admin",
	}

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// Unmarshal the data with msgpack
	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the member")

	// OnGet method should return test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
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

	// Should return an error if the member does not exist.
	_, err = suite.client.MemberDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse member id", "expected error when member does not exist")

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
}

func (suite *tenantTestSuite) TestMemberUpdate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name:     "member001",
		Role:     "Admin",
	}

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// Unmarshal the data with msgpack
	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the member")

	// OnGet method should return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
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

	// Endpoint must be authenticated
	require.NoError(suite.SetClientCSRFProtection(), "could not set csrf protection")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.EditCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the member ID is not parseable.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "invalid", Name: "member001", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "could not parse member id", "expected error when member does not exist")

	// Should return an error if the member name is not provided.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member name is required", "expected error when member name does not exist")

	// Should return an error if the member role is not provided.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member001"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	req := &api.Member{
		ID:   "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name: "member001",
		Role: "Admin",
	}

	rep, err := suite.client.MemberUpdate(ctx, req)
	require.NoError(err, "could not update member")
	require.NotEqual(req.ID, "01GM8MEZ097ZC7RQRCWMPRPS0T", "member id should not match")
	require.Equal(rep.Name, req.Name, "expected member name to match")
	require.Equal(rep.Role, req.Role, "expected member role to match")
}

func (suite *tenantTestSuite) TestMemberDelete() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	memberID := "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

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
	err := suite.client.MemberDelete(ctx, memberID)
	suite.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// User must have the correct permissions
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")
	err = suite.client.MemberDelete(ctx, memberID)
	suite.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user does not have permissions")

	// Set valid permissions for the rest of the tests
	claims.Permissions = []string{perms.RemoveCollaborators}
	require.NoError(suite.SetClientCredentials(claims), "could not set client credentials")

	// Should return an error if the member does not exist.
	err = suite.client.MemberDelete(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse member id", "expected error when member does not exist")

	err = suite.client.MemberDelete(ctx, memberID)
	require.NoError(err, "could not delete member")

	// Should return an error if the member ID is parsed but not found.
	trtl.OnDelete = func(ctx context.Context, dr *pb.DeleteRequest) (*pb.DeleteReply, error) {
		return nil, errors.New("key not found")
	}

	err = suite.client.MemberDelete(ctx, "01ARZ3NDEKTSV4RRFFQ69G5FAV")
	suite.requireError(err, http.StatusNotFound, "could not delete member", "expected error when member ID is not found")
}
