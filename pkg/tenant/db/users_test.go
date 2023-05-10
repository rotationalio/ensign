package db_test

import (
	"bytes"
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *dbTestSuite) TestCreateUserResources() {
	require := s.Require()
	ctx := context.Background()

	orgName := "Rotational Labs"

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
		Email:  "lwentzel@email.com",
		Name:   "Leopold Wentzel",
		Role:   "Member",
		Status: db.MemberStatusConfirmed,
	}
	require.ErrorIs(db.CreateUserResources(ctx, orgName, member), db.ErrMissingOrgID, "expected error when orgID is missing")

	// Should return an error if user email is missing
	member.Email = ""
	member.OrgID = ulid.MustParse("02ABCYAWC4PA72YC53RVXAEC67")
	require.ErrorIs(db.CreateUserResources(ctx, orgName, member), db.ErrMissingMemberEmail, "expected error when member email is missing")

	// Should return an error if user role is missing
	member.Name = "Leopold Wentzel"
	member.Email = "lwentzel@email.com"
	member.Role = ""
	require.ErrorIs(db.CreateUserResources(ctx, orgName, member), db.ErrMissingMemberRole, "expected error when member role is missing")

	// Should return an error if the org name is empty
	member.Role = "Member"
	require.ErrorIs(db.CreateUserResources(ctx, "", member), db.ErrMissingTenantName, "expected error when org name is not provided")

	// Succesfully creating all the required resources
	require.NoError(db.CreateUserResources(ctx, orgName, member), "expected no error when creating user resources")
	require.NotEmpty(member.ID, "expected member ID to be set")
	require.NotEmpty(member.Created, "expected created time to be set")
	require.NotEmpty(member.Modified, "expected modified time to be set")

	// Test that the method returns an error if trtl returns an error
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return nil, status.Error(codes.Internal, "trtl error")
	}
	require.Error(db.CreateUserResources(ctx, orgName, member), "expected error when trtl returns an error")
}

func (s *dbTestSuite) TestUpdateLastLogin() {
	require := s.Require()
	ctx := context.Background()

	trtl := db.GetMock()
	defer trtl.Reset()

	orgID := ulids.MustParse("01GX647S8PCVBCPJHXGJSPM87P")
	memberID := ulids.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7")
	accessToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR0U2MkVYWFIwWDA1NjFYRDUzUkRGQlFKIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDEiLCJzdWIiOiIwMUdRMlhBM1pGUjhGWUc2VzZaWk0xRkZTNyIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjgxMjM4ODI0LCJuYmYiOjE2ODEyMzUyMjQsImlhdCI6MTY4MTIzNTIyNCwianRpIjoiMDFneHJwdmE4eDk1dDJuMTlkcDZ2eG52dG0iLCJuYW1lIjoiTGVvcG9sZCBXZW50emVsIiwiZW1haWwiOiJsZW9wb2xkLndlbnR6ZWxAZ21haWwuY29tIiwib3JnIjoiMDFHWDY0N1M4UENWQkNQSkhYR0pTUE04N1AifQ.XY5-E1MJftIaTO--eusGGdUpjz-s6XIIynKQG7GlZ4VAe1HI5aVWKUhXmwKWSpQk0QElRcjTO_nNuOMB0jpmoYYKccJh8-dBFD1zltbSxAjhKqqKmEiY828ZoO6b66_B8jT0l1FmmYS8KafTPBTmP_t-u-CwJPMjEgkonbuXTIg7lIZ7F1CPrx5j3Ga5xq7asuxU4YqOPPXfdX3oSTsKojWRBL3kw7HkxeQBXzZ1say7xHu8iDYAbQw1L6JW_XDaEFYptQvLysEGwPG-uEp21gw_RSmNmLq0ANlgrcdBBcAaqk2_1L8lYIjPpcv3l7uFWN82T46iybP9XJLv9bOGq0g7eoversMx12D2IUpDFn32V_Gqp1lPUoikqqrM_hwnAXkH0qnwGbVcF7yttsjGKz72qUAiwaY2RH0QMAVaq1ElDrgqsvzx160ivzpvvN-7mKJ8WjZ4ZAPq8fyj1WcziNqfGqPgNer-PDav_Q59JOjLZG6A54FPoxHxNAPJofRES60K4XM06JnWOlfwI8tsUhwP5CKbEbEm3Ol6RZSlf1nUbJubHkGcgnem1DAiXoW0igB1aKMeCyzP1JfM3x93YWFzSLdLEah-y38UPei6sCDr1o_qeacSOErfJDHcYIElgmehkr4N4TlfcVjpoay0muREUYEKovgCBImzwv8YSVY"

	// Configure a fixture to return on Trtl Get
	member := &db.Member{
		OrgID: orgID,
		ID:    memberID,
		Email: "leopold.wentzel@gmail.com",
		Name:  "Leopold Wentzel",
		Role:  "Member",
	}

	key, err := member.Key()
	require.NoError(err, "expected the member fixture to have enough data to create a key")

	data, err := member.MarshalValue()
	require.NoError(err, "expected no error when marshaling member fixture")

	// Should return success for all get requests
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		if in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "unexpected namespace in Get request")
		}

		if !bytes.Equal(in.Key, key) {
			return nil, status.Error(codes.FailedPrecondition, "unexpected key in Get request")
		}

		return &pb.GetReply{
			Value: data,
		}, nil
	}

	// Should return success on put if the timestamp was updated
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		if in.Namespace != db.MembersNamespace {
			return nil, status.Error(codes.FailedPrecondition, "unexpected namespace in Put request")
		}

		if !bytes.Equal(in.Key, key) {
			return nil, status.Error(codes.FailedPrecondition, "unexpected key in Put request")
		}

		if len(in.Value) == 0 {
			return nil, status.Error(codes.FailedPrecondition, "unexpected value in Put request")
		}

		updated := &db.Member{}
		if err := updated.UnmarshalValue(in.Value); err != nil {
			return nil, status.Error(codes.FailedPrecondition, "could not unmarshal value sent in Put request")
		}

		if updated.LastActivity.IsZero() {
			return nil, status.Error(codes.FailedPrecondition, "expected last activity to be set in updated model")
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	// Should return an error if the access token is not parseable
	err = db.UpdateLastLogin(ctx, "bad-token", time.Now())
	require.Error(err, "expected error when access token is not parseable")

	// Succesfully updating the last login
	err = db.UpdateLastLogin(ctx, accessToken, time.Now())
	require.NoError(err, "expected no error when updating last login")

	// Should return an error if trtl returns an error
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Error(codes.Internal, "trtl error")
	}

	err = db.UpdateLastLogin(ctx, accessToken, time.Now())
	require.ErrorIs(err, status.Error(codes.Internal, "trtl error"), "expected error when trtl returns an error")
}
