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
	"github.com/rotationalio/ensign/pkg/quarterdeck/responses"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
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
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (*pb.PutReply, error) {
		return &pb.PutReply{}, nil
	}

	orgID := ulid.MustParse("02GQ38J5YWH4DCYJ6CZ2P5DA35")
	id := ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5DA2G")

	members := []*db.Member{
		{
			OrgID:  orgID,
			ID:     id,
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

	// Create initial fixtures
	reply := &qd.RegisterReply{
		ID:      id,
		OrgID:   orgID,
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
	}{
		{"name"},
		{"email"},
		{"password"},
		{"pwcheck"},
		{"organization"},
		{"domain"},
		{"agreetos"},
		{"agreeprivacy"},
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

	// Test registration with an invite token.
	req.InviteToken = "pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA"
	req.Organization = ""
	req.Domain = ""
	err = s.client.Register(ctx, req)
	require.NoError(err, "could not complete registration with invite token")

	// Test that a tenant, member, and project were created without error
	s.StopTasks()
	require.Equal(7, trtl.Calls[trtlmock.PutRPC], "expected 7 Put calls to trtl, 2 tenant puts (store, org_index), 3 project puts (store, org_index, object_keys), one member put for each user")
	require.Equal(0, trtl.Calls[trtlmock.GetRPC], "expected no gets on register")
	require.Equal(0, trtl.Calls[trtlmock.DeleteRPC], "expected no deletes on register")
	require.Equal(1, trtl.Calls[trtlmock.CursorRPC], "expected 1 cursor call on register")

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

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR0U2MkVYWFIwWDA1NjFYRDUzUkRGQlFKIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDEiLCJzdWIiOiIwMUdRMlhBM1pGUjhGWUc2VzZaWk0xRkZTNyIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjgxMjM4ODI0LCJuYmYiOjE2ODEyMzUyMjQsImlhdCI6MTY4MTIzNTIyNCwianRpIjoiMDFneHJwdmE4eDk1dDJuMTlkcDZ2eG52dG0iLCJuYW1lIjoiTGVvcG9sZCBXZW50emVsIiwiZW1haWwiOiJsZW9wb2xkLndlbnR6ZWxAZ21haWwuY29tIiwib3JnIjoiMDFHWDY0N1M4UENWQkNQSkhYR0pTUE04N1AifQ.XY5-E1MJftIaTO--eusGGdUpjz-s6XIIynKQG7GlZ4VAe1HI5aVWKUhXmwKWSpQk0QElRcjTO_nNuOMB0jpmoYYKccJh8-dBFD1zltbSxAjhKqqKmEiY828ZoO6b66_B8jT0l1FmmYS8KafTPBTmP_t-u-CwJPMjEgkonbuXTIg7lIZ7F1CPrx5j3Ga5xq7asuxU4YqOPPXfdX3oSTsKojWRBL3kw7HkxeQBXzZ1say7xHu8iDYAbQw1L6JW_XDaEFYptQvLysEGwPG-uEp21gw_RSmNmLq0ANlgrcdBBcAaqk2_1L8lYIjPpcv3l7uFWN82T46iybP9XJLv9bOGq0g7eoversMx12D2IUpDFn32V_Gqp1lPUoikqqrM_hwnAXkH0qnwGbVcF7yttsjGKz72qUAiwaY2RH0QMAVaq1ElDrgqsvzx160ivzpvvN-7mKJ8WjZ4ZAPq8fyj1WcziNqfGqPgNer-PDav_Q59JOjLZG6A54FPoxHxNAPJofRES60K4XM06JnWOlfwI8tsUhwP5CKbEbEm3Ol6RZSlf1nUbJubHkGcgnem1DAiXoW0igB1aKMeCyzP1JfM3x93YWFzSLdLEah-y38UPei6sCDr1o_qeacSOErfJDHcYIElgmehkr4N4TlfcVjpoay0muREUYEKovgCBImzwv8YSVY",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR0U2MkVYWFIwWDA1NjFYRDUzUkRGQlFKIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDEiLCJzdWIiOiIwMUdRMlhBM1pGUjhGWUc2VzZaWk0xRkZTNyIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiLCJodHRwOi8vbG9jYWxob3N0OjMwMDEvdjEvcmVmcmVzaCJdLCJleHAiOjE2ODEyNDI0MjQsIm5iZiI6MTY4MTIzNzkyNCwiaWF0IjoxNjgxMjM1MjI0LCJqdGkiOiIwMWd4cnB2YTh4OTV0Mm4xOWRwNnZ4bnZ0bSJ9.IX2rrVOhV9JAQMf0RECpLuf5szQ7AeXT8SI-1G49-7U7-UAkGeukbgdDRCyZ7Ai5BI6MvDC1LycQJVpx0DX5L1xTekz30T09Nu8kpU4paWUe3PMA4b7XttU2dkmesPpq4fywGOmOu0tQhsyh_sWWLAqk4qBSwxQEGjq4b0UI2N_egX2FquhDrGWomPlGEoBNQBPdWK9h1zTXOUsrxtn0K7C5jRzrj3GpL0wHwa_sGgAs5kCq7d5QLYwIHeE_MvWPmpXz1geDykFOxHFPZvfgBKUxDxxiKfYIlAxmfGiBws34Y2nv23rv__nE0mv0_a2IezGv_emDVUY01YYHUBamIqlD0GBFHCabcYwN1suAuvNcJCnJozcENlLOf-RfD2HHH7WrkrBWkuABgjiePwsBExhfsE-HbGFDzRGe-bEpypeS8ubjXUpdF_waRbOdogenzVfOiXQFzF8fyD09OQ26cBlk2LidSjJsm0P_6qmJF2T83DaSSOASozGyj376p5vHBKJxXbyb425vLJkvOciidtpai-UXe1WtX_arqCcO_03s_FiJj_ctp2J1DFVcTTkM0GXECnnTkYmBHZ0muS476KYOfmkJnQdrSujYIUtKvQyiPyD3afFGXCrbwzNmSMDDcryZm11tIKq5uNEa22235P4rqcdGu05VZuBZMaolpzU",
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
	require.Equal(expected, rep, "unexpected login reply")

	// TODO: Verify that CSRF cookies are set on the HTTP response

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
	s.requireError(err, http.StatusBadRequest, responses.ErrLogBackIn)

	// Should return an error if the orgID is not parseable
	req.RefreshToken = "refresh"
	req.OrgID = "not-a-ulid"
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusBadRequest, "invalid org_id")

	// Successful refresh
	expected := &api.AuthReply{
		AccessToken:  reply.AccessToken,
		RefreshToken: reply.RefreshToken,
	}
	req.OrgID = ulids.New().String()
	rep, err := s.client.Refresh(ctx, req)
	require.NoError(err, "could not complete refresh")
	require.Equal(expected, rep, "unexpected refresh reply")

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

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwiZXhwIjoxNjgwNjE1MzMwLCJuYmYiOjE2ODA2MTE3MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AiLCJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6Impkb2VAZXhhbXBsZS5jb20iLCJvcmciOiIxMjMiLCJwcm9qZWN0IjoiYWJjIiwicGVybWlzc2lvbnMiOlsicmVhZDpkYXRhIiwid3JpdGU6ZGF0YSJdfQ.LLb6c2RdACJmoT3IFgJEwfu2_YJMcKgM2bF3ISF41A37gKTOkBaOe-UuTmjgZ7WEcuQ-cVkht0KI_4zqYYctB_WB9481XoNwff5VgFf3xrPdOYxS00YXQnl09RRqt6Fmca8nvd4mXfdO7uvpyNVuCIqNxBPXdSnRhreSoFB1GtFm42sBPAD7vF-MQUmU0c4PTsbiCfhR1_buH0NYEE1QFp3vYcgoiXOJHh9VStmRscqvLB12AQrcs26G9opdTCCORmvR2W3JLJ_hliHyp-d9lhXmCDFyiGkDEhTAUglqwBjqz5SO1UfAThWJO18PvZl4QPhb724oNT82VPh0DMDwfw",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiLCJodHRwOi8vMTI3LjAuMC4xL3YxL3JlZnJlc2giXSwiZXhwIjoxNjgwNjE4OTMwLCJuYmYiOjE2ODA2MTQ0MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AifQ.CLHmtZwSPFCPoMBX06D_C3h3WuEonUbvbfWLvtmrMmIwnTwQ4hxsaRJo_a4qI-emp1HNg-yu_7c3VNwjkti-d0c7CAGApTaf5eRdGJ5HGUkI8RDHbbMFaOK86nAFnzdPJ2JLmGtLzvpF9eFXFllDhRiAB-2t0uKcOdN7cFghdwyWXIVJIJNjngF_WUFklmLKnqORtj_tA6UJ6NJnZln34eMGftAHbuH8x-xUiRePHnro4ydS43CKNOgRP8biMHiRR2broBz0apIt30TeQShaBSbmGx__LYdm7RKPJNVHAn_3h_PwwKQG567-Aqabg6TSmpwhXCk_RfUyQVGv2b997w",
		LastLogin:    time.Now().Format(time.RFC3339Nano),
	}

	// Configure the initial mock to return the tokens
	s.quarterdeck.OnSwitch(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply), mock.RequireAuth())

	// Create some claims for the user
	claims := &tokens.Claims{
		Name:        "Leopold Wentzel",
		Email:       "leopold.wentzel@gmail.com",
		OrgID:       "01GMTWFK4XZY597Y128KXQ4WHP",
		Permissions: []string{"read:organizations"},
	}

	// Endpoint must be authenticated
	req := &api.SwitchRequest{}
	_, err := s.client.Switch(ctx, req)
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
	require.Equal(reply.AccessToken, rep.AccessToken, "expected access token to match")
	require.Equal(reply.RefreshToken, rep.RefreshToken, "expected refresh token to match")
	require.Equal(reply.LastLogin, rep.LastLogin, "expected last login to match")

	// TODO: Verify that token cookies were set

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
