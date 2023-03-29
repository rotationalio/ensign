package db_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *dbTestSuite) TestVerifyOrg() {
	require := s.Require()
	ctx := context.Background()

	claimsOrgID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	resourceID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	namespace := "organizations"

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != namespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		return &pb.GetReply{
			Value: []byte("test"),
		}, nil
	}

	orgID, err := db.GetOrgIndex(ctx, resourceID)
	require.NoError(err, "could not get orgID")

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != namespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}
		return &pb.PutReply{}, nil
	}

	err = db.PutOrgIndex(ctx, resourceID, orgID)
	require.NoError(err, "could not store resourceID and orgID")

	ok, err := db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.NoError(err, "could not verify org")
	require.True(ok)
}
