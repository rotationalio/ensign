package db_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *dbTestSuite) TestCreateUserResources() {
	require := s.Require()
	ctx := context.Background()

	projectID := ulids.New()

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
		Role: "Member",
	}
	require.ErrorIs(db.CreateUserResources(ctx, projectID, member), db.ErrMissingOrgID, "expected error when orgID is missing")

	// Should return an error if user name is missing
	member.Name = ""
	member.OrgID = ulid.MustParse("02ABCYAWC4PA72YC53RVXAEC67")
	require.ErrorIs(db.CreateUserResources(ctx, projectID, member), db.ErrMissingMemberName, "expected error when member name is missing")

	// Should return an error if user role is missing
	member.Name = "Leopold Wentzel"
	member.Role = ""
	require.ErrorIs(db.CreateUserResources(ctx, projectID, member), db.ErrMissingMemberRole, "expected error when member role is missing")

	// Succesfully creating all the required resources
	member.Role = "Member"
	require.NoError(db.CreateUserResources(ctx, projectID, member), "expected no error when creating user resources")
	require.NotEmpty(member.ID, "expected member ID to be set")
	require.NotEmpty(member.Created, "expected created time to be set")
	require.NotEmpty(member.Modified, "expected modified time to be set")

	// Test that the method returns an error if trtl returns an error
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.Internal, "trtl error")
	}
	require.Error(db.CreateUserResources(ctx, projectID, member), "expected error when trtl returns an error")
}
