package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMemberModel(t *testing.T) {
	member := &db.Member{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
		Created:  time.Unix(1670424445, 0).In(time.UTC),
		Modified: time.Unix(1670424445, 0).In(time.UTC),
	}

	err := member.Validate()
	require.NoError(t, err, "could not validate member data")

	key, err := member.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, member.TenantID[:], key[0:16], "unexpected marshaling of the tenant id half of the key")
	require.Equal(t, member.ID[:], key[16:], "unexpected marshaling of the member id half of the key")

	require.Equal(t, db.MembersNamespace, member.Namespace(), "unexpected member namespace")

	// Test marshal and unmarshal
	data, err := member.MarshalValue()
	require.NoError(t, err, "could not marshal the member")

	other := &db.Member{}
	err = other.UnmarshalValue(data)
	require.NoError(t, err, "could not unmarshal the member")

	MembersEqual(t, member, other)
}

func (s *dbTestSuite) TestCreateTenantMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		TenantID: ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
	}

	err := member.Validate()
	require.NoError(err, "could not validate member data")

	// Call OnPut method from mock trtl database
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.CreateTenantMember(ctx, member)
	require.NoError(err, "could not create member")

	require.NotEmpty(member.ID, "expected non-zero ulid to be populated")
	require.NotZero(member.Created, "expected member to have a created timestamp")
	require.Equal(member.Created, member.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestCreateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		TenantID: ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
	}

	// Call OnPut method from mock trtl database
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

	require.NotEmpty(member.ID, "expected non-zero ulid to be populated")
	require.NotZero(member.Created, "expected member to have a created timestamp")
	require.Equal(member.Created, member.Modified, "expected the same created and modified timestamp")
}

func (s *dbTestSuite) TestRetrieveMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		TenantID: ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
	}

	// Call OnGet method from mock trtl database
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}
		if !bytes.Equal(in.Key[16:], member.ID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		// TODO: Add msgpack fixture helpers

		// Marshal the data with msgpack
		data, err := member.MarshalValue()
		require.NoError(err, "could not marshal the member")

		// Unmarshal the data with msgpack
		other := &db.Member{}
		err = other.UnmarshalValue(data)
		require.NoError(err, "could not unmarshal the member")

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	member, err := db.RetrieveMember(ctx, member.ID)
	require.NoError(err, "could not retrieve member")

	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), member.ID, "expected member id to match")
	require.Equal("member001", member.Name, "expected member name to match")
	require.Equal("role-example", member.Role, "expected member role to match")

	// TODO: Use crypto rand and monotonic entropy with ulid.New
	_, err = db.RetrieveMember(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListMembers() {
	require := s.Require()
	ctx := context.Background()

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

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
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

	values, err := db.List(ctx, prefix, namespace)
	require.NoError(err, "could not get member values")
	require.Len(values, 7)

	members := make([]*db.Member, 0, len(values))
	members = append(members, member)
	require.Len(members, 1)

	_, err = db.ListMembers(ctx, member.TenantID)
	require.Error(err, "could not list members")
}

func (s *dbTestSuite) TestUpdateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		TenantID: ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Name:     "member001",
		Role:     "role-example",
		Created:  time.Unix(1670424445, 0),
		Modified: time.Unix(1670424467, 0),
	}

	err := member.Validate()
	require.NoError(err, "could not validate member data")

	// Call OnPut method from mock trtl database
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if len(in.Key) == 0 || len(in.Value) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Put request")
		}

		if !bytes.Equal(in.Key[0:16], member.TenantID[:]) {
			return nil, status.Error(codes.NotFound, "tenant not found")
		}

		if !bytes.Equal(in.Key[16:], member.ID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	err = db.UpdateMember(ctx, member)
	require.NoError(err, "could not update member")

	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), member.ID, "member ID should not have changed")
	require.Equal(time.Unix(1670424445, 0), member.Created, "expected created timestamp to not change")
	require.True(time.Unix(1670424467, 0).Before(member.Modified), "expected modified timestamp to be updated")

	// Test NotFound path
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	err = db.UpdateMember(ctx, &db.Member{TenantID: ulid.Make(), ID: ulid.Make(), Name: "member002"})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteMember() {
	require := s.Require()
	ctx := context.Background()
	memberID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	// Call OnDelete method from mock trtl database
	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}

		if !bytes.Equal(in.Key[16:], memberID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}
	err := db.DeleteMember(ctx, memberID)
	require.NoError(err, "could not delete member")

	// Test NotFound path
	// TODO: Use crypto rand and monotonic entropy with ulid.New
	err = db.DeleteMember(ctx, ulid.Make())
	require.ErrorIs(err, db.ErrNotFound)
}

// MembersEqual tests assertions in the MemberModel.
// Note: require.True compares the actual.Created and actual.Modified
// timestamps because MsgPack does not preserve time zone information.
func MembersEqual(t *testing.T, expected, actual *db.Member, msgAndArgs ...interface{}) {
	require.Equal(t, expected.ID, actual.ID, msgAndArgs...)
	require.Equal(t, expected.Name, actual.Name, msgAndArgs...)
	require.Equal(t, expected.Role, actual.Role, msgAndArgs...)
	require.True(t, expected.Created.Equal(actual.Created), msgAndArgs...)
	require.True(t, expected.Modified.Equal(actual.Modified), msgAndArgs...)
}
