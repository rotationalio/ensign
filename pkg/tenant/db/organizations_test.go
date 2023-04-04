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

	resourceID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.OrganizationNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		return &pb.GetReply{
			Value: resourceID[:],
		}, nil
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.OrganizationNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}
		return &pb.PutReply{}, nil
	}

	orgID, err := db.GetOrgIndex(ctx, resourceID)
	require.NoError(err, "could not get orgID from the database")

	err = db.PutOrgIndex(ctx, resourceID, orgID)
	require.NoError(err, "could not store resourceID and orgID in the database")

	claimsOrgID := ulid.MustParse("01GWT0E850YBSDQH0EQFXRCMGB")
	err = db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.ErrorIs(err, db.ErrOrgNotVerified, "expected error when claims orgID and resourceID do not match")

	claimsOrgID = ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	err = db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.NoError(err, "could not verify org")
}