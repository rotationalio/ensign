package db_test

import (
	"bytes"
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	mrpc "github.com/trisacrypto/directory/pkg/trtl/mock"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *dbTestSuite) TestVerifyOrg() {
	// Setup test fixtures and variables
	require := s.Require()
	ctx := context.Background()

	orgID := ulid.MustParse("02ABC8QWNR7MYQXSQ682PJQM7T")
	resourceID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	// Setup the trtl database mock
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.OrganizationNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		// If the resource ID is incorrect return not found
		if !bytes.Equal(resourceID[:], in.Key) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}

		// Return the org ID for the specified resourceID
		return &pb.GetReply{
			Value: orgID[:],
		}, nil
	}

	// Should not be able to verify an incorrect orgID
	claimsOrgID := ulid.MustParse("01GWT0E850YBSDQH0EQFXRCMGB")
	err := db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.ErrorIs(err, db.ErrOrgNotVerified, "expected error when claims orgID and resourceID do not match")

	// Should be able to verify a correct orgID
	claimsOrgID = ulid.MustParse("02ABC8QWNR7MYQXSQ682PJQM7T")
	err = db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.NoError(err, "expected no verification error with correct claims orgID")

	// Errors should be returned if claimsOrgID or resourceID is zero
	err = db.VerifyOrg(ctx, ulid.ULID{}, resourceID)
	require.ErrorIs(err, db.ErrMissingOrgID)

	err = db.VerifyOrg(ctx, claimsOrgID, ulid.ULID{})
	require.ErrorIs(err, db.ErrMissingID)

	// If the database returns an error then verify org should return the error
	s.mock.UseError(mrpc.GetRPC, codes.NotFound, "resource not found")
	err = db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.ErrorIs(err, db.ErrNotFound)

	// If there is a more significant error, that should also be returned
	s.mock.UseError(mrpc.GetRPC, codes.Internal, "something bad happened")
	err = db.VerifyOrg(ctx, claimsOrgID, resourceID)
	require.EqualError(err, "rpc error: code = Internal desc = something bad happened")
}

func (s *dbTestSuite) TestGetOrgIndex() {
	// Setup variables and fixtures for tests
	require := s.Require()
	ctx := context.Background()

	orgID := ulid.MustParse("02ABC8QWNR7MYQXSQ682PJQM7T")
	resourceID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	// Mock the trtl database functionality
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.OrganizationNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		// If the resource ID is incorrect return not found
		if !bytes.Equal(resourceID[:], in.Key) {
			return nil, status.Error(codes.NotFound, "resource not found")
		}

		// Return the org ID for the specified resourceID
		return &pb.GetReply{
			Value: orgID[:],
		}, nil
	}

	// Should be able to retreive the organization from the database
	actual, err := db.GetOrgIndex(ctx, resourceID)
	require.NoError(err, "could not get orgID from the database")
	require.Equal(orgID, actual, "expected the returned result to equal the fixture")

	// The resource ID is required
	_, err = db.GetOrgIndex(ctx, ulid.ULID{})
	require.ErrorIs(err, db.ErrMissingID)

	// Not found error is returned when it is not found
	_, err = db.GetOrgIndex(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)

	// If there is a more significant error, that should also be returned
	s.mock.UseError(mrpc.GetRPC, codes.Internal, "something bad happened")
	_, err = db.GetOrgIndex(ctx, resourceID)
	require.EqualError(err, "rpc error: code = Internal desc = something bad happened")

	// The data returned from the database must be a parseable ULID
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: []byte("notaulid"),
		}, nil
	}

	_, err = db.GetOrgIndex(ctx, resourceID)
	require.EqualError(err, "ulid: bad data size when unmarshaling")
}

func (s *dbTestSuite) TestPutOrgIndex() {
	// Setup variables and fixtures for tests
	require := s.Require()
	ctx := context.Background()

	orgID := ulid.MustParse("02ABC8QWNR7MYQXSQ682PJQM7T")
	resourceID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	// Mock the trtl database functionality
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.OrganizationNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}
		return &pb.PutReply{}, nil
	}

	// Test valid storage of resourceID in organization index.
	err := db.PutOrgIndex(ctx, resourceID, orgID)
	require.NoError(err, "could not store resourceID and orgID in the database")

	// Ensure resourceID is required
	err = db.PutOrgIndex(ctx, ulid.ULID{}, orgID)
	require.ErrorIs(err, db.ErrMissingID, "resource id should be required")

	// Ensure orgID is required
	err = db.PutOrgIndex(ctx, resourceID, ulid.ULID{})
	require.ErrorIs(err, db.ErrMissingOrgID, "org id should be required")

	// If there is a more significant error, that should also be returned
	s.mock.UseError(mrpc.PutRPC, codes.Internal, "something bad happened")
	err = db.PutOrgIndex(ctx, resourceID, orgID)
	require.EqualError(err, "rpc error: code = Internal desc = something bad happened")
}
