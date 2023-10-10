package db_test

import (
	"bytes"
	"context"
	"strings"
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

func TestProfessionSegment(t *testing.T) {
	// Test that all the profession segments can be parsed into the enum
	// TODO: This does not validate that all enum values have a string representation
	for enum, segment := range db.ProfessionSegmentStrings {
		val, err := db.ParseProfessionSegment(segment)
		require.NoError(t, err, "could not parse profession segment %s", segment)
		require.Equal(t, enum, val, "wrong enum value for %s", segment)
	}

	testCases := []struct {
		segment string
		enum    db.ProfessionSegment
		err     error
	}{
		{"", db.ProfessionSegmentUnspecified, nil},
		{"NotARealSegment", db.ProfessionSegmentUnspecified, db.ErrProfessionUnknown},
		{" work ", db.ProfessionSegmentWork, nil},
		{"Work", db.ProfessionSegmentWork, nil},
		{"WORK", db.ProfessionSegmentWork, nil},
		{"EduCatIon", db.ProfessionSegmentEducation, nil},
	}

	for i, tc := range testCases {
		enum, err := db.ParseProfessionSegment(tc.segment)
		if tc.err == nil {
			require.NoError(t, err, "expected no error for test case: %d", i)
			require.Equal(t, tc.enum, enum, "wrong enum value for test case: %d", i)
		} else {
			require.ErrorIs(t, err, tc.err, "expected error for test case: %d", i)
		}
	}
}

func TestDeveloperSegment(t *testing.T) {
	// Test that all the developer segments can be parsed into the enum
	// TODO: This does not validate that all enum values have a string representation
	for enum, segment := range db.DeveloperSegmentStrings {
		val, err := db.ParseDeveloperSegment(segment)
		require.NoError(t, err, "could not parse developer segment %s", segment)
		require.Equal(t, enum, val, "wrong enum value for %s", segment)
	}

	testCases := []struct {
		segment string
		enum    db.DeveloperSegment
		err     error
	}{
		{"", db.DeveloperSegmentUnspecified, nil},
		{"NotARealSegment", db.DeveloperSegmentUnspecified, db.ErrDeveloperUnknown},
		{"Application Development", db.DeveloperSegmentApplicationDevelopment, nil},
		{"application development", db.DeveloperSegmentApplicationDevelopment, nil},
		{" data science ", db.DeveloperSegmentDataScience, nil},
	}

	for i, tc := range testCases {
		enum, err := db.ParseDeveloperSegment(tc.segment)
		if tc.err == nil {
			require.NoError(t, err, "expected no error for test case: %d", i)
			require.Equal(t, tc.enum, enum, "wrong enum value for test case: %d", i)
		} else {
			require.ErrorIs(t, err, tc.err, "expected error for test case: %d", i)
		}
	}
}

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
		JoinedAt:     time.Unix(1670424445, 0).In(time.UTC),
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

	// Role is required
	member.Email = "test@testing.com"
	member.Name = "Leopold Wentzel"
	member.Role = ""
	require.ErrorIs(t, member.Validate(), db.ErrMissingMemberRole, "expected validate to fail with missing role")

	// Unknown roles are rejected
	member.Role = "NotARealRole"
	require.ErrorIs(t, member.Validate(), db.ErrUnknownMemberRole, "expected validate to fail with invalid role")

	// Correct validation
	member.Role = perms.RoleAdmin
	require.NoError(t, member.Validate(), "expected validate to succeed with required org id")

	// Test the onboarding validation errors
	testCases := []struct {
		name              string
		organization      string
		workspace         string
		professionSegment db.ProfessionSegment
		developerSegment  []db.DeveloperSegment
		errs              db.ValidationErrors
	}{
		{name: strings.Repeat("a", 1025), errs: db.ValidationErrors{{Field: "name", Err: db.ErrNameTooLong, Index: -1}}},
		{organization: strings.Repeat("a", 1025), errs: db.ValidationErrors{{Field: "organization", Err: db.ErrOrganizationTooLong, Index: -1}}},
		{workspace: strings.Repeat("a", 1025), errs: db.ValidationErrors{{Field: "workspace", Err: db.ErrWorkspaceTooLong, Index: -1}}},
		{workspace: "rotational io", errs: db.ValidationErrors{{Field: "workspace", Err: db.ErrInvalidWorkspace, Index: -1}}},
		{workspace: "2bornot2b", errs: db.ValidationErrors{{Field: "workspace", Err: db.ErrInvalidWorkspace, Index: -1}}},
		{workspace: "hi", errs: db.ValidationErrors{{Field: "workspace", Err: db.ErrInvalidWorkspace, Index: -1}}},
		{developerSegment: []db.DeveloperSegment{db.DeveloperSegmentApplicationDevelopment, db.DeveloperSegmentUnspecified}, errs: db.ValidationErrors{{Field: "developer_segment", Err: db.ErrDeveloperUnspecified, Index: 1}}},
		{name: strings.Repeat("a", 1025), workspace: "not a valid workspace", errs: db.ValidationErrors{{Field: "name", Err: db.ErrNameTooLong, Index: -1}, {Field: "workspace", Err: db.ErrInvalidWorkspace, Index: -1}}},
		{name: "Leopold Wentzel", organization: "Rotational Labs", workspace: "rotational-io", professionSegment: db.ProfessionSegmentEducation, developerSegment: []db.DeveloperSegment{db.DeveloperSegmentApplicationDevelopment}},
	}

	for i, tc := range testCases {
		member := &db.Member{
			OrgID:             orgID,
			Email:             "test@testing.com",
			Role:              perms.RoleAdmin,
			Name:              tc.name,
			Organization:      tc.organization,
			Workspace:         tc.workspace,
			ProfessionSegment: tc.professionSegment,
			DeveloperSegment:  tc.developerSegment,
		}
		err := member.Validate()
		if tc.errs == nil {
			require.NoError(t, err, "expected no validation errors for test case: %d", i)
		} else {
			var verrs db.ValidationErrors
			require.ErrorAs(t, err, &verrs, "expected error to be a ValidationErrors for test case: %d", i)
			require.Equal(t, tc.errs, verrs, "wrong validation errors for test case: %d", i)
		}
	}
}

func TestMemberStatus(t *testing.T) {
	// Default member should have status onboarding (new users without an invite)
	member := &db.Member{}
	require.Equal(t, db.MemberStatusOnboarding, member.OnboardingStatus(), "expected default member status to be onboarding")

	// Member who has only completed some steps should have status onboarding
	member.Name = "Leopold Wentzel"
	member.ProfessionSegment = db.ProfessionSegmentPersonal
	require.Equal(t, db.MemberStatusOnboarding, member.OnboardingStatus(), "expected partial member record to be onboarding")

	// Member who has not accepted an invite should have status pending
	member.Invited = true
	member.JoinedAt = time.Time{}
	require.Equal(t, db.MemberStatusPending, member.OnboardingStatus(), "expected member status to be pending")

	// Member who has accepted an invite but not completed onboarding should have status onboarding
	member.JoinedAt = time.Now()
	require.Equal(t, db.MemberStatusOnboarding, member.OnboardingStatus(), "expected member status to be onboarding")

	// Member who has completed onboarding should have status active
	member.Organization = "Rotational"
	member.Workspace = "rotational-io"
	member.DeveloperSegment = []db.DeveloperSegment{db.DeveloperSegmentApplicationDevelopment}
	require.Equal(t, db.MemberStatusActive, member.OnboardingStatus(), "expected member status to be active")
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
		OrgID: ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1"),
		Email: "test@testing.com",
		Name:  "member001",
		Role:  "Admin",
	}

	// Call OnPut method from mock trtl database
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
		switch in.Namespace {
		case db.MembersNamespace:
			return &pb.PutReply{Success: true}, nil
		case db.OrganizationNamespace:
			return &pb.PutReply{}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	err := db.CreateMember(ctx, member)
	require.NoError(err, "could not create member")

	require.NotEmpty(member.ID, "expected non-zero ulid to be populated")
	require.Equal(db.MemberStatusOnboarding, member.OnboardingStatus(), "expected member to have onboarding status")
	require.NotZero(member.Created, "expected member to have a created timestamp")
	require.Equal(member.Created, member.Modified, "expected the same created and modified timestamp")
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

	// Configure trtl to return the member records on cursor
	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back the member data
		for _, member := range members {
			key, err := member.Key()
			if err != nil {
				return status.Error(codes.FailedPrecondition, "could not marshal member key for trtl response")
			}
			data, err := member.MarshalValue()
			if err != nil {
				return status.Error(codes.FailedPrecondition, "could not marshal member data for trtl response")
			}
			stream.Send(&pb.KVPair{
				Key:       key,
				Value:     data,
				Namespace: in.Namespace,
			})
		}

		return nil
	}

	s.Run("Single Page", func() {
		// If all the results are on a single page then the next cursor is nil
		cursor := &pg.Cursor{
			PageSize: 100,
		}

		rep, cursor, err := db.ListMembers(ctx, orgID, cursor)
		require.NoError(err, "could not list members")
		require.Len(rep, 3, "expected 3 members")
		require.Nil(cursor, "next page cursor should not be set since there isn't a next page")

		for i := range members {
			require.Equal(members[i].ID, rep[i].ID, "expected member id to match")
			require.Equal(members[i].Email, rep[i].Email, "expected member name to match")
			require.Equal(members[i].Name, rep[i].Name, "expected member name to match")
			require.Equal(members[i].Role, rep[i].Role, "expected member role to match")
		}
	})

	s.Run("Multiple Pages", func() {
		// If results are on multiple pages then the next cursor is not nil
		cursor := &pg.Cursor{
			PageSize: 2,
		}
		rep, cursor, err := db.ListMembers(ctx, orgID, cursor)
		require.NoError(err, "could not list members")
		require.Len(rep, 2, "expected 2 members on the first page")
		require.NotNil(cursor, "expected cursor to be not nil because there is a next page")

		// Ensure the new start index is correct
		startBytes, err := members[2].Key()
		require.NoError(err, "could not marshal member key")
		startKey := &db.Key{}
		require.NoError(startKey.UnmarshalValue(startBytes), "could not unmarshal member key")
		startString, err := startKey.String()
		require.NoError(err, "could not convert member key to string")
		require.Equal(startString, cursor.StartIndex, "expected cursor start index to match")
		require.Empty(cursor.EndIndex, "expected cursor end index to be empty")

		// Configure trtl to return the rest of the members
		members = members[2:]
		rep, cursor, err = db.ListMembers(ctx, orgID, cursor)
		require.NoError(err, "could not list members")
		require.Len(rep, 1, "expected 1 member on the second page")
		require.Nil(cursor, "expected cursor to be nil because there is no next page")
	})
}

func (s *dbTestSuite) TestGetMemberByEmail() {
	require := s.Require()
	ctx := context.Background()

	orgID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	email := "test3@testing.com"

	members := []*db.Member{
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7"),
			Email:    "test@testing.com",
			Name:     "member001",
			Role:     perms.RoleAdmin,
			Created:  time.Unix(1670424445, 0),
			Modified: time.Unix(1670424445, 0),
		},
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XAMGG9N7DF7KSRDQVFZ2A"),
			Email:    "test2@testing.com",
			Name:     "member002",
			Role:     perms.RoleMember,
			Created:  time.Unix(1673659941, 0),
			Modified: time.Unix(1673659941, 0),
		},
		{
			OrgID:    orgID,
			ID:       ulid.MustParse("01GQ2XB2SCGY5RZJ1ZGYSEMNDE"),
			Email:    "test3@testing.com",
			Name:     "member003",
			Role:     perms.RoleAdmin,
			Created:  time.Unix(1674073941, 0),
			Modified: time.Unix(1674073941, 0),
		},
	}

	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, orgID[:]) {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		for _, member := range members {
			data, err := member.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(member.Email),
				Value:     data,
				Namespace: db.MembersNamespace,
			})
		}
		return nil
	}

	// Should return an error if email is not provided.
	_, err := db.GetMemberByEmail(ctx, orgID, "")
	require.ErrorIs(err, db.ErrMissingMemberEmail, "expected error when email is not provided")

	// Should return an error if email does not exist.
	_, err = db.GetMemberByEmail(ctx, orgID, "test4@testing.com")
	require.ErrorIs(err, db.ErrMemberEmailNotFound, "expected error when email does not exist")

	rep, err := db.GetMemberByEmail(ctx, orgID, email)
	require.NoError(err, "could not get member by email")
	require.Equal(rep.Email, email, "expected member email to match")
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

	// If created timestamp is missing then it should be updated
	member.Created = time.Time{}
	require.NoError(db.UpdateMember(ctx, member), "could not update member")
	require.Equal(member.Modified, member.Created, "expected created timestamp to be updated")

	// Should fail if member ID is missing
	member.ID = ulid.ULID{}
	require.ErrorIs(db.UpdateMember(ctx, member), db.ErrMissingID, "expected error for missing member ID")

	// Should fail if member model is invalid
	member.ID = ulid.MustParse("01GKKYAWC4PA72YC53RVXAEC67")
	member.Role = ""
	require.ErrorIs(db.UpdateMember(ctx, member), db.ErrMissingMemberRole, "expected error for invalid member model")

	// Test NotFound path
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.NotFound, "not found")
	}
	req := &db.Member{OrgID: ulids.New(), ID: ulids.New(), Email: "test@testing.com", Name: "member002", Role: "Admin"}
	err = db.UpdateMember(ctx, req)
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
	require.True(t, expected.JoinedAt.Equal(actual.JoinedAt), msgAndArgs...)
}
