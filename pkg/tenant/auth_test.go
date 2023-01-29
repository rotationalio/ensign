package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (s *tenantTestSuite) TestRegister() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create initial fixtures
	reply := &qd.RegisterReply{
		ID:      ulid.MustParse("01GQ38J5YWH4DCYJ6CZ2P5DA2G"),
		Email:   "leopold.wentzel@gmail.com",
		Message: "Welcome to Ensign!",
		Role:    "member",
		Created: time.Now().Format(time.RFC3339Nano),
	}

	// Configure the initial mock to return a 200 response with the reply
	s.quarterdeck.OnRegister(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	// Name is required
	req := &api.RegisterRequest{
		Email:    "leopold.wentzel@gmail.com",
		Password: "hunter2",
		PwCheck:  "hunter2",
	}
	_, err := s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing required fields for registration")

	// Email is required
	req.Name = "Leopold Wentzel"
	req.Email = ""
	_, err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing required fields for registration")

	// Password is required
	req.Email = "leopold.wentzel@gmail.com"
	req.Password = ""
	_, err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing required fields for registration")

	// Password check is required
	req.Password = "hunter2"
	req.PwCheck = ""
	_, err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "missing required fields for registration")

	// Password and password check must match
	req.PwCheck = "hunter3"
	_, err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusBadRequest, "passwords do not match")

	// Successful registration
	expected := &api.RegisterReply{
		Email:   reply.Email,
		Message: reply.Message,
		Role:    reply.Role,
	}
	req.PwCheck = "hunter2"
	rep, err := s.client.Register(ctx, req)
	require.NoError(err, "could not complete registration")
	require.Equal(expected, rep, "unexpected registration reply")

	// Register method should handle errors from Quarterdeck
	s.quarterdeck.OnRegister(mock.UseStatus(http.StatusInternalServerError))
	_, err = s.client.Register(ctx, req)
	s.requireError(err, http.StatusInternalServerError, "could not complete registration")
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
