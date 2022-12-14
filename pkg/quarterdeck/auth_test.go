package quarterdeck_test

import (
	"context"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
)

func (suite *quarterdeckTestSuite) TestRegister() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: only happy path test is implemented; implement error paths as well.
	req := &api.RegisterRequest{
		Name:     "Rachel Johnson",
		Email:    "rachel@example.com",
		Password: "supers3cretSquirrel?",
		PwCheck:  "supers3cretSquirrel?",
	}

	rep, err := suite.client.Register(ctx, req)
	require.NoError(err, "unable to create user from valid request")

	require.NotEmpty(rep.ID, "did not get a user ID back from the database")
	require.Equal(req.Email, rep.Email)
	require.Equal("Welcome to Ensign!", rep.Message)
	require.NotEmpty(rep.Created, "did not get a created timestamp back")

	// TODO: test that the user actually made it into the database
}

func (suite *quarterdeckTestSuite) TestLogin() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: actually implement the login test!
	req := &api.LoginRequest{}
	_, err := suite.client.Login(ctx, req)
	require.Error(err, "expected bad request")
}

func (suite *quarterdeckTestSuite) TestAuthenticate() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &api.APIAuthentication{}
	_, err := suite.client.Authenticate(ctx, req)
	require.Error(err, "expected unimplemented error")
}

func (suite *quarterdeckTestSuite) TestRefresh() {
	require := suite.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := suite.client.Refresh(ctx)
	require.Error(err, "expected unimplemented error")
}
