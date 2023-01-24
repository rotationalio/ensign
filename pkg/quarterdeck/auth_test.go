package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (s *quarterdeckTestSuite) TestRegister() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: only happy path test is implemented; implement error paths as well.
	req := &api.RegisterRequest{
		Name:     "Rachel Johnson",
		Email:    "rachel@example.com",
		Password: "supers3cretSquirrel?",
		PwCheck:  "supers3cretSquirrel?",
	}

	rep, err := s.client.Register(ctx, req)
	require.NoError(err, "unable to create user from valid request")

	require.NotEmpty(rep.ID, "did not get a user ID back from the database")
	require.Equal(req.Email, rep.Email)
	require.Equal("Welcome to Ensign!", rep.Message)
	require.NotEmpty(rep.Created, "did not get a created timestamp back")

	// TODO: test that the user actually made it into the database
}

func (s *quarterdeckTestSuite) TestLogin() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: actually implement the login test!
	req := &api.LoginRequest{}
	_, err := s.client.Login(ctx, req)
	require.Error(err, "expected bad request")
}

func (s *quarterdeckTestSuite) TestAuthenticate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: actually implement the authenticate test!
	req := &api.APIAuthentication{}
	_, err := s.client.Authenticate(ctx, req)
	require.Error(err, "expected bad request")
}

func (s *quarterdeckTestSuite) TestRefresh() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test empty RefreshRequest returns error
	req := &api.RefreshRequest{}
	_, err := s.client.Refresh(ctx, req)
	require.Error(err, "missing credentials")

	// Test invalid refersh token returns error
	req = &api.RefreshRequest{RefreshToken: "refresh"}
	_, err = s.client.Refresh(ctx, req)
	require.Error(err, "could not verify refresh token")

	// Happy path test
	registerReq := &api.RegisterRequest{
		Name:     "Rachel Johnson",
		Email:    "rachel@example.com",
		Password: "supers3cretSquirrel?",
		PwCheck:  "supers3cretSquirrel?",
	}
	registerRep, err := s.client.Register(ctx, registerReq)
	require.NoError(err)
	loginReq := &api.LoginRequest{
		Email:    registerRep.Email,
		Password: "supers3cretSquirrel?",
	}
	loginRep, err := s.client.Login(ctx, loginReq)
	require.NoError(err)
	refreshReq := &api.RefreshRequest{
		RefreshToken: loginRep.RefreshToken,
	}
	refreshRep, err := s.client.Refresh(ctx, refreshReq)
	require.NoError(err, "could not create credentials")
	require.NotNil(refreshRep)
	require.NotEqual(loginRep.AccessToken, refreshRep.AccessToken)
	require.NotEqual(loginRep.RefreshToken, refreshRep.RefreshToken)
}
