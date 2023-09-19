package tenant_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	trtlmock "github.com/trisacrypto/directory/pkg/trtl/mock"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *tenantTestSuite) TestRegister() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetTasks()

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Set up the mock to return success for put requests
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (_ *pb.PutReply, err error) {
		switch pr.Namespace {
		case db.MembersNamespace:
			// Verify that the correct member record is being created
			member := &db.Member{}
			if err = member.UnmarshalValue(pr.Value); err != nil {
				return nil, status.Errorf(codes.Internal, "could not unmarshal member data in put request: %v", err)
			}

			if member.Organization == "" {
				return nil, status.Errorf(codes.FailedPrecondition, "missing organization in member record being created")
			}
		}

		return &pb.PutReply{}, nil
	}

	orgID := ulid.MustParse("02GQ38J5YWH4DCYJ6CZ2P5DA35")
	id := ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5DA2G")

	members := []*db.Member{
		{
			OrgID: orgID,
			ID:    id,
			Email: "leopold.wentzel@gmail.com",
			Name:  "Leopold Wentzel",
			Role:  perms.RoleAdmin,
		},
	}

	// Connect to trtl mock and call OnCursor to loop through members in the database.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
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

	// Create initial fixtures
	reply := &qd.RegisterReply{
		ID:        id,
		OrgID:     orgID,
		Email:     "leopold.wentzel@gmail.com",
		OrgName:   "Rotational Labs",
		OrgDomain: "rotational-io",
		Message:   "Welcome to Ensign!",
		Role:      perms.RoleAdmin,
		Created:   time.Now().Format(time.RFC3339Nano),
	}

	// Make sure that we are passing all required fields to Quarterdeck
	s.quarterdeck.OnRegister(mock.UseHandler(func(w http.ResponseWriter, r *http.Request) {
		var err error
		req := &qd.RegisterRequest{}
		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err = req.Validate(); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(reply)
		require.NoError(err, "could not encode quarterdeck reply from mock")
	}))

	// Test missing fields
	req := &api.RegisterRequest{
		Email:        "leopold.wentzel@gmail.com",
		Password:     "ajdfsd943%^&xbs",
		PwCheck:      "ajdfsd943%^&xbs",
		AgreeToS:     true,
		AgreePrivacy: true,
	}
	testCases := []struct {
		missing string
	}{
		{"email"},
		{"password"},
		{"pwcheck"},
		{"agreetos"},
		{"agreeprivacy"},
	}
	for _, tc := range testCases {
		s.Run("missing_"+tc.missing, func() {
			// Create local copy for this test
			req := *req

			// Set the field to the default value
			switch tc.missing {
			case "email":
				req.Email = ""
			case "password":
				req.Password = ""
			case "pwcheck":
				req.PwCheck = ""
			case "agreetos":
				req.AgreeToS = false
			case "agreeprivacy":
				req.AgreePrivacy = false
			default:
				require.Fail("invalid test case")
			}

			// Should return a validation error
			err := s.client.Register(ctx, &req)
			s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)
		})
	}

	// Test mismatched passwords
	req.PwCheck = "hunter3"
	err := s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	// Successful registration
	req.PwCheck = req.Password
	err = s.client.Register(ctx, req)
	require.NoError(err, "could not complete registration")

	// Test that trtl was called the correct number of times across all register calls
	s.StopTasks()
	require.Equal(3, trtl.Calls[trtlmock.PutRPC], "expected 3 Put calls to trtl, 2 tenant puts (store, org_index), one member put for the user")
	require.Equal(0, trtl.Calls[trtlmock.GetRPC], "expected no gets on register")
	require.Equal(0, trtl.Calls[trtlmock.DeleteRPC], "expected no deletes on register")

	// Register method should handle errors from Quarterdeck
	s.quarterdeck.OnRegister(mock.UseError(http.StatusBadRequest, "password too weak"))
	err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "password too weak")
}

func (s *tenantTestSuite) TestLogin() {
	require := s.Require()

	// Connect to mock trtl database
	// Since the mock is shared between routines, the very last thing the test should
	// do is reset the mock to avoid interfering with background tasks.
	trtl := db.GetMock()
	defer trtl.Reset()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetTasks()

	orgID := ulid.MustParse("01GX647S8PCVBCPJHXGJSPM87P")
	memberID := ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7")

	// Create member fixture
	members := []*db.Member{
		{
			OrgID: orgID,
			ID:    memberID,
			Email: "leopold.wentzel@gmail.com",
			Name:  "Leopold Wentzel",
			Role:  perms.RoleAdmin,
		},
	}

	memberData, err := members[0].MarshalValue()
	require.NoError(err, "could not marshal member data")

	// Connect to trtl mock and call OnCursor to loop through members in the database.
	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
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

	// Trtl Get should return a valid member record for update.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		return &pb.GetReply{
			Value: memberData,
		}, nil
	}

	// Connect to trtl mock and call OnPut to update the member status.
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Quarterdeck mock should return auth tokens on login
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: memberID.String(),
		},
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
		OrgID: orgID.String(),
	}
	creds := &qd.LoginReply{}
	creds.AccessToken, creds.RefreshToken, err = s.quarterdeck.CreateTokenPair(claims)
	require.NoError(err, "could not create token pair from claims fixture")
	s.quarterdeck.OnLogin(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(creds))

	// Email is required
	req := &api.LoginRequest{
		Password: "hunter2",
	}

	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	// Password is required
	req.Email = "leopold.wentzel@gmail.com"
	req.Password = ""
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain)

	// Successful login
	req.Password = "hunter2"
	rep, err := s.client.Login(ctx, req)
	require.NoError(err, "could not complete login")
	require.Equal(creds.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(creds.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	s.requireAuthCookies(creds.AccessToken, creds.RefreshToken)
	s.ClearAuthTokens()
	s.ResetTasks()

	// Set invite token and test login.
	req.InviteToken = "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA"
	rep, err = s.client.Login(ctx, req)
	require.NoError(err, "could not complete login")
	require.Equal(creds.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(creds.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	s.requireAuthCookies(creds.AccessToken, creds.RefreshToken)
	s.ClearAuthTokens()
	s.ResetTasks()

	// Set orgID and return an error if invite token is set.
	req.OrgID = orgID.String()
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, "cannot provide both invite token and org id")

	// Should return an error if org ID is not valid.
	req.InviteToken = ""
	req.OrgID = "invalid"
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid org id")

	// Test login with orgID and no invite token.
	req.OrgID = orgID.String()
	rep, err = s.client.Login(ctx, req)
	require.NoError(err, "could not complete login")
	require.Equal(creds.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(creds.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	s.requireAuthCookies(creds.AccessToken, creds.RefreshToken)
	s.ClearAuthTokens()
	s.ResetTasks()

	// TODO: Test case where return user ID is different from the existing member ID.

	// Login method should handle errors from Quarterdeck
	s.quarterdeck.OnLogin(mock.UseError(http.StatusForbidden, "invalid login credentials"))
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusForbidden, "invalid login credentials")

	// Test returning an error with valid org ID when Quarterdeck returns an error.
	s.quarterdeck.OnLogin(mock.UseError(http.StatusInternalServerError, "could not create valid credentials"))
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusInternalServerError, "could not create valid credentials")
}

func (s *tenantTestSuite) TestRefresh() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetTasks()

	// Setup the trtl mock
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		OrgID: ulids.New(),
		ID:    ulids.New(),
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	key, err := member.Key()
	require.NoError(err, "could not get member key")

	// Trtl Get should return a valid member record for update.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, status.Errorf(codes.NotFound, "namespace not found")
		}

		return &pb.GetReply{
			Value: memberData,
		}, nil
	}

	// Trtl Put should update the last activity timestamp for the member.
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}

			// Verify that last activity time is updated
			member := &db.Member{}
			if err = member.UnmarshalValue(in.Value); err != nil {
				return nil, status.Errorf(codes.Internal, "could not unmarshal member data in put request: %v", err)
			}

			if member.LastActivity.IsZero() {
				return nil, status.Errorf(codes.FailedPrecondition, "expected last activity to be set in updated member model")
			}
		default:
			return nil, status.Errorf(codes.NotFound, "namespace not found")
		}

		return &pb.PutReply{}, nil
	}

	// Quarterdeck mock should return auth tokens on refresh
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: ulids.New().String(),
		},
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
		OrgID: ulids.New().String(),
	}
	creds := &qd.LoginReply{}
	creds.AccessToken, creds.RefreshToken, err = s.quarterdeck.CreateTokenPair(claims)
	require.NoError(err, "could not create token pair from claims fixture")
	s.quarterdeck.OnRefresh(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(creds))

	// Refresh token is required
	req := &api.RefreshRequest{}
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrLogBackIn)

	// Should return an error if the orgID is not parseable
	req.RefreshToken = "refresh"
	req.OrgID = "not-a-ulid"
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid org_id")

	// Successful refresh
	req.OrgID = ulids.New().String()
	rep, err := s.client.Refresh(ctx, req)
	require.NoError(err, "could not complete refresh")
	require.Equal(creds.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(creds.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	s.requireAuthCookies(creds.AccessToken, creds.RefreshToken)
	s.ClearAuthTokens()
	s.ResetTasks()

	// Refresh method should handle errors from Quarterdeck
	s.quarterdeck.OnRefresh(mock.UseError(http.StatusUnauthorized, "expired token"))
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "expired token")
}

func (s *tenantTestSuite) TestSwitch() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.ResetTasks()

	// Setup the trtl mock
	trtl := db.GetMock()
	defer trtl.Reset()

	member := &db.Member{
		OrgID: ulids.New(),
		ID:    ulids.New(),
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	memberData, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	key, err := member.Key()
	require.NoError(err, "could not get member key")

	// Trtl Get should return a valid member record for update.
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, status.Errorf(codes.NotFound, "namespace not found")
		}

		return &pb.GetReply{
			Value: memberData,
		}, nil
	}

	// Trtl Put should update the last activity timestamp for the member.
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}

			// Verify that last activity time is updated
			member := &db.Member{}
			if err = member.UnmarshalValue(in.Value); err != nil {
				return nil, status.Errorf(codes.Internal, "could not unmarshal member data in put request: %v", err)
			}

			if member.LastActivity.IsZero() {
				return nil, status.Errorf(codes.FailedPrecondition, "expected last activity to be set in updated member model")
			}
		default:
			return nil, status.Errorf(codes.NotFound, "namespace not found")
		}

		return &pb.PutReply{}, nil
	}

	// Quarterdeck mock should return auth tokens on switch
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: ulids.New().String(),
		},
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
		OrgID: ulids.New().String(),
	}
	creds := &qd.LoginReply{}
	creds.AccessToken, creds.RefreshToken, err = s.quarterdeck.CreateTokenPair(claims)
	require.NoError(err, "could not create token pair from claims fixture")
	s.quarterdeck.OnSwitch(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(creds))

	// Endpoint must be authenticated
	req := &api.SwitchRequest{}
	_, err = s.client.Switch(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Should return an error if the orgID is not provided
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.Switch(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing org_id in request")

	// Should return an error if the orgID is invalid
	req.OrgID = "not-a-ulid"
	_, err = s.client.Switch(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid org_id in request")

	// Should return an error if the user is already authenticated with the org
	req.OrgID = claims.OrgID
	_, err = s.client.Switch(ctx, req)
	s.requireError(err, http.StatusBadRequest, "already logged in to this organization")

	// Successfully switching to a new organization
	req.OrgID = "02GMTWFK4XZY597Y128KXQ4ABC"
	rep, err := s.client.Switch(ctx, req)
	require.NoError(err, "expected successful switch")
	require.Equal(creds.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(creds.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	s.requireAuthCookies(creds.AccessToken, creds.RefreshToken)
	s.ClearAuthTokens()
	s.ResetTasks()

	// Switch method should handle errors from Quarterdeck
	s.quarterdeck.OnSwitch(mock.UseError(http.StatusForbidden, "invalid credentials"), mock.RequireAuth())
	_, err = s.client.Switch(ctx, req)
	s.requireError(err, http.StatusForbidden, "invalid credentials")
}

func (s *tenantTestSuite) TestVerifyEmail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Configure the initial mock to return a 200 response
	s.quarterdeck.OnVerify(mock.UseStatus(http.StatusOK))

	// Token is required
	req := &api.VerifyRequest{}
	err := s.client.VerifyEmail(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrVerificationFailed)

	// Successful verification
	req.Token = "token"
	err = s.client.VerifyEmail(ctx, req)
	require.NoError(err, "expected successful verification")

	// VerifyEmail method should handle errors from Quarterdeck
	s.quarterdeck.OnVerify(mock.UseError(http.StatusBadRequest, "invalid token"))
	err = s.client.VerifyEmail(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid token")
}

func (s *tenantTestSuite) TestResendEmail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Run("Happy Path", func() {
		// Quarterdeck returns 204 response
		s.quarterdeck.OnResendEmail(mock.UseStatus(http.StatusNoContent))

		// Should return success if email is provided
		req := &api.ResendRequest{Email: "leopold.wentzel@gmail.com"}
		err := s.client.ResendEmail(ctx, req)
		require.NoError(err, "expected successful resend if email is provided")

		// Should return success if orgID is provided
		req = &api.ResendRequest{
			Email: "leopold.wentzel@gmail.com",
			OrgID: ulids.New().String(),
		}
		err = s.client.ResendEmail(ctx, req)
		require.NoError(err, "expected successful resend if email and orgID are provided")
	})

	s.Run("Bad Email", func() {
		// Should return 400 if email is not provided
		testCases := []struct {
			email string
		}{
			{""}, {"\t\t"}, {"\n\n"}, {strings.Repeat("a", 256)},
		}

		for _, tc := range testCases {
			err := s.client.ResendEmail(ctx, &api.ResendRequest{Email: tc.email})
			s.requireHTTPError(err, http.StatusBadRequest)
		}
	})

	s.Run("Bad orgID", func() {
		// Should return 400 if orgID is not parseable
		err := s.client.ResendEmail(ctx, &api.ResendRequest{
			Email: "leopold.wentzel@gmail.com",
			OrgID: "not-a-ulid",
		})
		s.requireHTTPError(err, http.StatusBadRequest)
	})

	s.Run("Quarterdeck Error", func() {
		// Should forward errors from Quarterdeck
		s.quarterdeck.OnResendEmail(mock.UseError(http.StatusBadRequest, responses.ErrBadResendRequest))
		err := s.client.ResendEmail(ctx, &api.ResendRequest{Email: "leopold.wentzel@gmail.com"})
		s.requireHTTPError(err, http.StatusBadRequest)
	})
}
