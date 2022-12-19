package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

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

	// Should return an error if the member does not exist.
	err = suite.client.MemberDelete(ctx, "invalid")
	require.Error(err, "member does not exist")

	err = suite.client.MemberDelete(ctx, memberID)
	require.NoError(err, "could not delete member")
}
