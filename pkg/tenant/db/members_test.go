package db_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMemberModel(t *testing.T) {
	member := &db.Member{
		ID:   ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name: "member-example",
		Role: "role-example",
	}

	key, err := member.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, member.ID[:], key, "unexpected marshaling of the key")

	require.Equal(t, db.MembersNamespace, member.Namespace(), "unexpected member namespace")

	// Test marshal and unmarshal
	data, err := member.MarshalValue()
	require.NoError(t, err, "could not marshal the member")

	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the member")

	require.Equal(t, member, other, "unmarshaled member does not match marshaled tenant")
}

func (s *dbTestSuite) TestCreateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{Name: "member-example"}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.CreateMember(ctx, member)
	require.NoError(err, "could not create member")

	require.NotEqual("", member.ID, "expected non-zero ulid to be populated")
}

func (s *dbTestSuite) TestRetrieveMember() {
	require := s.Require()
	ctx := context.Background()
	memberID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}

		if !bytes.Equal(in.Key, memberID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		// Load fixture from disk
		data, err := os.ReadFile("testdata/member.json")
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	member, err := db.RetrieveMember(ctx, memberID)
	require.NoError(err, "could not retrieve member")

	require.Equal(memberID, member.ID)
	require.Equal("member-example", member.Name)

	_, err = db.RetrieveMember(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestUpdateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		ID:   ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name: "member-example",
		Role: "role-example",
	}

	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key, member.ID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err := db.UpdateMember(ctx, member)
	require.NoError(err, "could not update member")

	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), member.ID, "member ID should not have changed")

	// Test NotFound path
	err = db.UpdateMember(ctx, &db.Member{ID: ulid.Make()})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteMember() {
	require := s.Require()
	ctx := context.Background()
	memberID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}
		if !bytes.Equal(in.Key, memberID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}
	err := db.DeleteMember(ctx, memberID)
	require.NoError(err, "could not delete member")

	// Test NotFound path
	err = db.DeleteMember(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}
