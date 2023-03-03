package tenant_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
)

func (s *tenantTestSuite) TestRegister() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
		Email:   "leopold.wentzel@gmail.com",
		Message: "Welcome to Ensign!",
		Role:    "member",
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

	// Register method should handle errors from Quarterdeck
	s.quarterdeck.OnRegister(mock.UseStatus(http.StatusBadRequest))
	err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "could not complete registration")
}

func (s *tenantTestSuite) TestLogin() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Configure the initial mock to return a 200 response with the reply
	s.quarterdeck.OnLogin(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

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

	// TODO: Verify that CSRF cookies are set on the HTTP response

	// Login method should handle errors from Quarterdeck
	s.quarterdeck.OnLogin(mock.UseStatus(http.StatusInternalServerError))
	_, err = s.client.Login(ctx, req)
	s.requireError(err, http.StatusInternalServerError, "could not complete login")
}

func (s *tenantTestSuite) TestRefresh() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	reply := &qd.LoginReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
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
	s.quarterdeck.OnRefresh(mock.UseStatus(http.StatusUnauthorized))
	_, err = s.client.Refresh(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "could not complete refresh")
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
	s.quarterdeck.OnVerify(mock.UseStatus(http.StatusBadRequest))
	err = s.client.VerifyEmail(ctx, req)
	s.requireError(err, http.StatusBadRequest, "could not complete email verification")
}
