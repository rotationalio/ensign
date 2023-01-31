package db_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *dbTestSuite) TestCreateUser() {
	require := s.Require()
	ctx := context.Background()

	// Configure trtl to return success for all requests
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	// Should return an error if organization is missing
	member := &db.Member{
		Name: "Leopold Wentzel",
	}
	require.ErrorIs(db.CreateUser(ctx, member), db.ErrMissingOrgID, "expected error when orgID is missing")

	// Should return an error if user name is missing
	member.Name = ""
	member.OrgID = ulid.MustParse("02ABCYAWC4PA72YC53RVXAEC67")
	require.ErrorIs(db.CreateUser(ctx, member), db.ErrMissingMemberName, "expected error when member name is missing")

	// TODO: Is the user role required?

	// Succesfully creating all the required resources
	member.Name = "Leopold Wentzel"
	require.NoError(db.CreateUser(ctx, member), "expected no error when creating user resources")
	require.NotEmpty(member.ID, "expected member ID to be set")
	require.NotEmpty(member.TenantID, "expected tenant ID to be set")
	require.NotEmpty(member.Created, "expected created time to be set")
	require.NotEmpty(member.Modified, "expected modified time to be set")

	// Test that the method returns an error if trtl returns an error
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.Internal, "trtl error")
	}
	require.Error(db.CreateUser(ctx, member), "expected error when trtl returns an error")
}
