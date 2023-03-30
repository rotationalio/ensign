package db_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMemberModel(t *testing.T) {
	member := &db.Member{
		OrgID:        ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:           ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Email:        "test@testing.com",
		Name:         "member001",
		Role:         "Admin",
		Created:      time.Unix(1670424445, 0).In(time.UTC),
		Modified:     time.Unix(1670424445, 0).In(time.UTC),
		LastActivity: time.Unix(1670424445, 0).In(time.UTC),
		DateAdded:    time.Unix(1670424445, 0).In(time.UTC),
	}

	// Successful validation
	err := member.Validate()
	require.NoError(t, err, "could not validate member data")

	key, err := member.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, member.OrgID[:], key[0:16], "unexpected marshaling of the org id half of the key")
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

func TestMemberValidation(t *testing.T) {
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	member := &db.Member{
		OrgID: orgID,
		Email: "test@testing.com",
		Name:  "Leopold Wentzel",
		Role:  perms.RoleAdmin,
	}

	// OrgID is required
	member.OrgID = ulids.Null
	require.ErrorIs(t, member.Validate(), db.ErrMissingOrgID, "expected validate to fail with missing org id")

	// Email is required
	member.OrgID = orgID
	member.Email = ""
	require.ErrorIs(t, member.Validate(), db.ErrMissingMemberEmail, "expected validate to fail with missing email")

	// Name is required
	member.Email = "test@testing.com"
	member.Name = ""
	require.ErrorIs(t, member.Validate(), db.ErrMissingMemberName, "expected validate to fail with missing name")

	// Name must have non-whitespace characters
	member.Name = " "
	require.ErrorIs(t, member.Validate(), db.ErrMissingMemberName, "expected validate to fail with missing name")

	// Role is required
	member.Name = "Leopold Wentzel"
	member.Role = ""
	require.ErrorIs(t, member.Validate(), db.ErrMissingMemberRole, "expected validate to fail with missing role")

	// Unknown roles are rejected
	member.Role = "NotARealRole"
	require.ErrorIs(t, member.Validate(), db.ErrUnknownMemberRole, "expected validate to fail with invalid role")

	// Correct validation
	member.Role = perms.RoleAdmin
	require.NoError(t, member.Validate(), "expected validate to succeed with required org id")
}

func TestMemberKey(t *testing.T) {
	// Test that the key can't be created when ID is missing
	id := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	member := &db.Member{
		OrgID: orgID,
	}
	_, err := member.Key()
	require.ErrorIs(t, err, db.ErrMissingID, "expected error when missing member id")

	// Test that the key can't be created when OrgID is missing
	member.ID = id
	member.OrgID = ulids.Null
	_, err = member.Key()
	require.ErrorIs(t, err, db.ErrMissingOrgID, "expected error when missing org id")

	// Test that the key is composed correctly
	member.OrgID = orgID
	key, err := member.Key()
	require.NoError(t, err, "could not marshal the key")
	require.Equal(t, member.OrgID[:], key[0:16], "unexpected marshaling of the org id half of the key")
	require.Equal(t, member.ID[:], key[16:], "unexpected marshaling of the member id half of the key")
}

func (s *dbTestSuite) TestCreateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		OrgID:  ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		Email:  "test@testing.com",
		Name:   "member001",
		Role:   "Admin",
		Status: "Confirmed",
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
	require.Equal(member.Status, "Confirmed", "expected member to have confirmed status")
	require.NotZero(member.Created, "expected member to have a created timestamp")
	require.Equal(member.Created, member.Modified, "expected the same created and modified timestamp")
	require.Equal(member.Created, member.LastActivity, "expected the same created and last activity timestamp")
	require.Equal(member.Created, member.DateAdded, "expected the same created and date added timestamp")
}

func (s *dbTestSuite) TestRetrieveMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		OrgID: ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:    ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Email: "test@testing.com",
		Name:  "member001",
		Role:  "Admin",
	}

	// Call OnGet method from mock trtl database
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Get request")
		}
		if !bytes.Equal(in.Key[0:16], member.OrgID[:]) || !bytes.Equal(in.Key[16:], member.ID[:]) {
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

	member, err := db.RetrieveMember(ctx, member.OrgID, member.ID)
	require.NoError(err, "could not retrieve member")

	require.Equal(ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"), member.ID, "expected member id to match")
	require.Equal("test@testing.com", member.Email, "expected member email to match")
	require.Equal("member001", member.Name, "expected member name to match")
	require.Equal("Admin", member.Role, "expected member role to match")

	_, err = db.RetrieveMember(ctx, member.OrgID, ulids.New())
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestListMembers() {
	require := s.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")

	members := []*db.Member{
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7"),
			Email:    "test@testing.com",
			Name:     "member001",
			Role:     "Admin",
			Created:  time.Unix(1670424445, 0),
			Modified: time.Unix(1670424445, 0),
		},
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Email:    "test2@testing.com",
			Name:     "member002",
			Role:     "Member",
			Created:  time.Unix(1673659941, 0),
			Modified: time.Unix(1673659941, 0),
		},
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Email:    "test3@testing.com",
			Name:     "member003",
			Role:     "Admin",
			Created:  time.Unix(1674073941, 0),
			Modified: time.Unix(1674073941, 0),
		},
	}

	prefix := orgID[:]
	namespace := "members"

	// Call the OnCursor method
	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		var start bool
		// Send back some data and terminate
		for _, member := range members {
			if in.SeekKey != nil && bytes.Equal(in.SeekKey, member.ID[:]) {
				start = true
			}
			if in.SeekKey == nil || start {
				data, err := member.MarshalValue()
				require.NoError(err, "could not marshal data")
				stream.Send(&pb.KVPair{
					Key:       member.ID[:],
					Value:     data,
					Namespace: in.Namespace,
				})
			}
		}
		return nil
	}

	prev := &pg.Cursor{
		StartIndex: "",
		EndIndex:   "",
		PageSize:   100,
	}

	// Return all members and verify next page token is not set.
	rep, next, err := db.ListMembers(ctx, orgID, prev)
	require.NoError(err, "could not list members")
	require.Len(rep, 3, "expected 3 members")
	require.Nil(next, "next page cursor should not be set since there isn't a next page")

	for i := range members {
		require.Equal(members[i].ID, rep[i].ID, "expected member id to match")
		require.Equal(members[i].Email, rep[i].Email, "expected member name to match")
		require.Equal(members[i].Name, rep[i].Name, "expected member name to match")
		require.Equal(members[i].Role, rep[i].Role, "expected member role to match")
	}

	// Test pagination by setting a page size.
	prev.PageSize = 2
	rep, next, err = db.ListMembers(ctx, orgID, prev)
	require.NoError(err, "could not list members")
	require.Len(rep, 2, "expected 2 members")
	require.NotEqual(prev.StartIndex, next.StartIndex, "starting index should not be the same")
	require.NotEqual(prev.EndIndex, next.EndIndex, "ending index should not be the same")
	require.Equal(prev.PageSize, next.PageSize, "page size should be the same")
	require.NotEmpty(next.Expires, "expires timestamp should not be empty")
}

func (s *dbTestSuite) TestUpdateMember() {
	require := s.Require()
	ctx := context.Background()
	member := &db.Member{
		OrgID:    ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		ID:       ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67"),
		Email:    "test@testing.com",
		Name:     "member001",
		Role:     "Admin",
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

		if !bytes.Equal(in.Key[0:16], member.OrgID[:]) {
			return nil, status.Error(codes.NotFound, "organization not found")
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
	require.Equal(member.Modified, member.LastActivity, "expected modified and last activity timestamp to match")

	// If created timestamp is missing then it should be updated
	member.Created = time.Time{}
	require.NoError(db.UpdateMember(ctx, member), "could not update member")
	require.Equal(member.Modified, member.Created, "expected created timestamp to be updated")

	// Should fail if member ID is missing
	member.ID = ulid.ULID{}
	require.ErrorIs(db.UpdateMember(ctx, member), db.ErrMissingID, "expected error for missing member ID")

	// Should fail if member model is invalid
	member.ID = ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	member.Name = ""
	require.ErrorIs(db.UpdateMember(ctx, member), db.ErrMissingMemberName, "expected error for invalid member model")

	// Test NotFound path
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	err = db.UpdateMember(ctx, &db.Member{OrgID: ulids.New(), ID: ulids.New(), Email: "test@testing.com", Name: "member002", Role: "Admin"})
	require.ErrorIs(err, db.ErrNotFound)
}

func (s *dbTestSuite) TestDeleteMember() {
	require := s.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	memberID := ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")

	// Call OnDelete method from mock trtl database
	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
		if len(in.Key) == 0 || in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "bad Delete request")
		}

		if !bytes.Equal(in.Key[0:16], orgID[:]) || !bytes.Equal(in.Key[16:], memberID[:]) {
			return nil, status.Error(codes.NotFound, "member not found")
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}
	err := db.DeleteMember(ctx, orgID, memberID)
	require.NoError(err, "could not delete member")

	// Test NotFound path
	err = db.DeleteMember(ctx, orgID, ulids.New())
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
	require.True(t, expected.LastActivity.Equal(actual.LastActivity), msgAndArgs...)
	require.True(t, expected.DateAdded.Equal(actual.DateAdded), msgAndArgs...)
}
