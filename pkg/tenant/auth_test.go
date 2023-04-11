package tenant_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
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
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Create initial fixtures
	reply := &qd.RegisterReply{
		ID:      ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5DA2G"),
		OrgID:   ulid.MustParse("02GQ38J5YWH4DCYJ6CZ2P5DA35"),
		Email:   "leopold.wentzel@gmail.com",
		Message: "Welcome to Ensign!",
		Role:    perms.RoleAdmin,
		Created: time.Now().Format(time.RFC3339Nano),
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
		Name:         "Leopold Wentzel",
		Email:        "leopold.wentzel@gmail.com",
		Password:     "ajdfsd943%^&xbs",
		PwCheck:      "ajdfsd943%^&xbs",
		Organization: "Rotational Labs",
		Domain:       "rotational.io",
		AgreeToS:     true,
		AgreePrivacy: true,
	}
	testCases := []struct {
		missing string
		err     string
	}{
		{"name", "name is required"},
		{"email", "email is required"},
		{"password", "password is required"},
		{"pwcheck", "passwords do not match"},
		{"organization", "organization is required"},
		{"domain", "domain is required"},
		{"agreetos", "you must agree to the terms of service"},
		{"agreeprivacy", "you must agree to the privacy policy"},
	}
	for _, tc := range testCases {
		s.Run("missing_"+tc.missing, func() {
			// Create local copy for this test
			req := *req

			// Set the field to the default value
			switch tc.missing {
			case "name":
				req.Name = ""
			case "email":
				req.Email = ""
			case "password":
				req.Password = ""
			case "pwcheck":
				req.PwCheck = ""
			case "organization":
				req.Organization = ""
			case "domain":
				req.Domain = ""
			case "agreetos":
				req.AgreeToS = false
			case "agreeprivacy":
				req.AgreePrivacy = false
			default:
				require.Fail("invalid test case")
			}

			// Should return a validation error
			err := s.client.Register(ctx, &req)
			s.requireError(err, http.StatusBadRequest, tc.err)
		})
	}

	// Test mismatched passwords
	req.PwCheck = "hunter3"
	err := s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "passwords do not match")

	// Successful registration
	req.PwCheck = req.Password
	err = s.client.Register(ctx, req)
	require.NoError(err, "could not complete registration")

	// Test registration with an invite token.
	req.InviteToken = "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA"
	err = s.client.Register(ctx, req)
	require.NoError(err, "could not complete registration with invite token")

	// Test that a tenant, member, and project were created without error
	s.StopTasks()
	require.Equal(7, trtl.Calls[trtlmock.PutRPC], "expected 7 Put calls to trtl for two puts for each tenant, member, and project (store and org index) and one for object_keys.")
	require.Equal(0, trtl.Calls[trtlmock.GetRPC], "expected no gets on register")
	require.Equal(0, trtl.Calls[trtlmock.DeleteRPC], "expected no deletes on register")
	require.Equal(0, trtl.Calls[trtlmock.CursorRPC], "expected no cursors on register")

	// Register method should handle errors from Quarterdeck
	s.quarterdeck.OnRegister(mock.UseError(http.StatusBadRequest, "password too weak"))
	err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "password too weak")
}

func (s *tenantTestSuite) TestLogin() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := ulid.MustParse("01GX647S8PCVBCPJHXGJSPM87P")
	memberID := ulid.MustParse("01GQ2XA3ZFR8FYG6W6ZZM1FFS7")

	// Connect to mock trtl database.
	trtl := db.GetMock()
	defer trtl.Reset()

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwiZXhwIjoxNjgwNjE1MzMwLCJuYmYiOjE2ODA2MTE3MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AiLCJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6Impkb2VAZXhhbXBsZS5jb20iLCJvcmciOiIxMjMiLCJwcm9qZWN0IjoiYWJjIiwicGVybWlzc2lvbnMiOlsicmVhZDpkYXRhIiwid3JpdGU6ZGF0YSJdfQ.LLb6c2RdACJmoT3IFgJEwfu2_YJMcKgM2bF3ISF41A37gKTOkBaOe-UuTmjgZ7WEcuQ-cVkht0KI_4zqYYctB_WB9481XoNwff5VgFf3xrPdOYxS00YXQnl09RRqt6Fmca8nvd4mXfdO7uvpyNVuCIqNxBPXdSnRhreSoFB1GtFm42sBPAD7vF-MQUmU0c4PTsbiCfhR1_buH0NYEE1QFp3vYcgoiXOJHh9VStmRscqvLB12AQrcs26G9opdTCCORmvR2W3JLJ_hliHyp-d9lhXmCDFyiGkDEhTAUglqwBjqz5SO1UfAThWJO18PvZl4QPhb724oNT82VPh0DMDwfw",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiLCJodHRwOi8vMTI3LjAuMC4xL3YxL3JlZnJlc2giXSwiZXhwIjoxNjgwNjE4OTMwLCJuYmYiOjE2ODA2MTQ0MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AifQ.CLHmtZwSPFCPoMBX06D_C3h3WuEonUbvbfWLvtmrMmIwnTwQ4hxsaRJo_a4qI-emp1HNg-yu_7c3VNwjkti-d0c7CAGApTaf5eRdGJ5HGUkI8RDHbbMFaOK86nAFnzdPJ2JLmGtLzvpF9eFXFllDhRiAB-2t0uKcOdN7cFghdwyWXIVJIJNjngF_WUFklmLKnqORtj_tA6UJ6NJnZln34eMGftAHbuH8x-xUiRePHnro4ydS43CKNOgRP8biMHiRR2broBz0apIt30TeQShaBSbmGx__LYdm7RKPJNVHAn_3h_PwwKQG567-Aqabg6TSmpwhXCk_RfUyQVGv2b997w",
	}

	// Configure the initial mock to return a 200 response with the reply
	s.quarterdeck.OnLogin(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	members := []*db.Member{
		{
			OrgID:  orgID,
			ID:     memberID,
			Email:  "leopold.wentzel@gmail.com",
			Name:   "Leopold Wentzel",
			Role:   perms.RoleAdmin,
			Status: db.MemberStatusPending,
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

	// Connect to trtl mock and call OnPut to update the member status.
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	// Email is required
	req := &api.LoginRequest{
		Password: "hunter2",
	}

	_, err := s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing email/password for login")

	// Password is required
	req.Email = "leopold.wentzel@gmail.com"
	req.Password = ""
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing email/password for login")

	// Successful login
	expected := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}
	req.Password = "hunter2"
	rep, err := s.client.Login(ctx, req)
	require.NoError(err, "could not complete login")
	require.Equal(expected, rep, "unexpected login reply")

	// Set invite token and test login.
	req.InviteToken = "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA"
	rep, err = s.client.Login(ctx, req)
	require.NoError(err, "could not complete login")
	require.Equal(expected, rep, "unexpected login reply")

	// TODO: Verify that CSRF cookies are set on the HTTP response

	// Login method should handle errors from Quarterdeck
	s.quarterdeck.OnLogin(mock.UseError(http.StatusForbidden, "invalid login credentials"))
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusForbidden, "invalid login credentials")
}

func (s *tenantTestSuite) TestRefresh() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwiZXhwIjoxNjgwNjE1MzMwLCJuYmYiOjE2ODA2MTE3MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AiLCJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6Impkb2VAZXhhbXBsZS5jb20iLCJvcmciOiIxMjMiLCJwcm9qZWN0IjoiYWJjIiwicGVybWlzc2lvbnMiOlsicmVhZDpkYXRhIiwid3JpdGU6ZGF0YSJdfQ.LLb6c2RdACJmoT3IFgJEwfu2_YJMcKgM2bF3ISF41A37gKTOkBaOe-UuTmjgZ7WEcuQ-cVkht0KI_4zqYYctB_WB9481XoNwff5VgFf3xrPdOYxS00YXQnl09RRqt6Fmca8nvd4mXfdO7uvpyNVuCIqNxBPXdSnRhreSoFB1GtFm42sBPAD7vF-MQUmU0c4PTsbiCfhR1_buH0NYEE1QFp3vYcgoiXOJHh9VStmRscqvLB12AQrcs26G9opdTCCORmvR2W3JLJ_hliHyp-d9lhXmCDFyiGkDEhTAUglqwBjqz5SO1UfAThWJO18PvZl4QPhb724oNT82VPh0DMDwfw",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiLCJodHRwOi8vMTI3LjAuMC4xL3YxL3JlZnJlc2giXSwiZXhwIjoxNjgwNjE4OTMwLCJuYmYiOjE2ODA2MTQ0MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AifQ.CLHmtZwSPFCPoMBX06D_C3h3WuEonUbvbfWLvtmrMmIwnTwQ4hxsaRJo_a4qI-emp1HNg-yu_7c3VNwjkti-d0c7CAGApTaf5eRdGJ5HGUkI8RDHbbMFaOK86nAFnzdPJ2JLmGtLzvpF9eFXFllDhRiAB-2t0uKcOdN7cFghdwyWXIVJIJNjngF_WUFklmLKnqORtj_tA6UJ6NJnZln34eMGftAHbuH8x-xUiRePHnro4ydS43CKNOgRP8biMHiRR2broBz0apIt30TeQShaBSbmGx__LYdm7RKPJNVHAn_3h_PwwKQG567-Aqabg6TSmpwhXCk_RfUyQVGv2b997w",
	}

	// Configure the initial mock to return a 200 response with the reply
	s.quarterdeck.OnRefresh(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	// Refresh token is required
	req := &api.RefreshRequest{}
	_, err := s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing refresh token")

	// Successful refresh
	expected := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}
	req.RefreshToken = "refresh"
	rep, err := s.client.Refresh(ctx, req)
	require.NoError(err, "could not complete refresh")
	require.Equal(expected, rep, "unexpected refresh reply")

	// Refresh method should handle errors from Quarterdeck
	s.quarterdeck.OnRefresh(mock.UseError(http.StatusUnauthorized, "expired token"))
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "expired token")
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
	s.requireError(err, http.StatusBadRequest, "missing token in request")

	// Successful verification
	req.Token = "token"
	err = s.client.VerifyEmail(ctx, req)
	require.NoError(err, "expected successful verification")

	// VerifyEmail method should handle errors from Quarterdeck
	s.quarterdeck.OnVerify(mock.UseError(http.StatusBadRequest, "invalid token"))
	err = s.client.VerifyEmail(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid token")
}
