package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *tenantTestSuite) TestTenantMemberList() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	member := &db.Member{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
		Created:  time.Unix(1670424445, 0),
		Modified: time.Unix(1670424445, 0),
	}

	prefix := member.TenantID[:]
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
		for i := 0; i < 7; i++ {
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     []byte(fmt.Sprintf("value %d", i)),
				Namespace: in.Namespace,
			})
		}
		return nil
	}

	// Should return an error if the tenant does not exist.
	_, err := suite.client.TenantMemberList(ctx, "invalid", &api.PageQuery{})
	suite.requireError(err, http.StatusBadRequest, "could not parse tenant ulid", "expected error when tenant does not exist")

	members, err := suite.client.TenantMemberList(ctx, member.TenantID.String(), &api.PageQuery{})
	require.NoError(err, "could not list tenant members")
	require.Len(members.TenantMembers, 7, "expected 7 members")
}

func (suite *tenantTestSuite) TestTenantMemberCreate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	tenantID := ulid.Make().String()

	defer cancel()

	// Connect to a mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Call the OnPut method and return a PutReply
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if tenant id is not a valid ULID.
	_, err := suite.client.TenantMemberCreate(ctx, "tenantID", &api.Member{ID: "", Name: "member-example"})
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
		ID: ulid.Make().String(),
	}

	// Create a member test fixture
	req := &api.Member{
		Name: "member-example",
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

	// TODO: Test length of values assigned to *api.MemberPage
	_, err := suite.client.MemberList(ctx, &api.PageQuery{})
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

	// Should return an error if member id exists.
	_, err := suite.client.MemberCreate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member-example", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member id cannot be specified on create", "expected error when member id exists")

	// Should return an error if the member name does not exist
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member name is required", "expected error when member name does not exist")

	// Should return an error if the member role does not exist.
	_, err = suite.client.MemberCreate(ctx, &api.Member{ID: "", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	// Create a member test fixture
	req := &api.Member{
		Name: "member-example",
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
	member := &db.Member{
		ID:   ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name: "member-example",
		Role: "Admin",
	}
	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// Unmarshal the data with msgpack
	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the member")

	// Call the OnGet method and return test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

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
	member := &db.Member{
		ID:   ulid.MustParse("01ARZ3NDEKTSV4RRFFQ69G5FAV"),
		Name: "member-example",
		Role: "Admin",
	}

	defer cancel()

	// Connect to mock trtl database
	trtl := db.GetMock()
	defer trtl.Reset()

	// Marshal the data with msgpack
	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member")

	// Unmarshal the data with msgpack
	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(err, "could not unmarshal the member")

	// Call the OnGet method and return the test data.
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return an error if the member does not exist
	_, err = suite.client.MemberDetail(ctx, "invalid")
	suite.requireError(err, http.StatusBadRequest, "could not parse member id", "expected error when member does not exist")

	// Call the OnPut method and return a PutReply.
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Should return an error if the member name does not exist.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Role: "Admin"})
	suite.requireError(err, http.StatusBadRequest, "member name is required", "expected error when member name does not exist")

	// Should return an error if the member role does not exist.
	_, err = suite.client.MemberUpdate(ctx, &api.Member{ID: "01ARZ3NDEKTSV4RRFFQ69G5FAV", Name: "member-example"})
	suite.requireError(err, http.StatusBadRequest, "member role is required", "expected error when member role does not exist")

	req := &api.Member{
		ID:   "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		Name: "member-example",
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

	// Should return an error if the member does not exist.
	err := suite.client.MemberDelete(ctx, "invalid")
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
