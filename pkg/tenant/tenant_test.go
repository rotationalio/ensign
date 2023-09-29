package tenant_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/service"
	"github.com/rotationalio/ensign/pkg/utils/tlstest"
	emock "github.com/rotationalio/go-ensign/mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type tenantTestSuite struct {
	suite.Suite
	srv         *tenant.Server
	auth        *authtest.Server
	client      api.TenantClient
	quarterdeck *mock.Server
	ensign      *emock.Ensign
	metatopic   *emock.Ensign
	subscriber  *tenant.TopicSubscriber
	stop        chan bool
}

// Runs once before all tests are executed
func (suite *tenantTestSuite) SetupSuite() {
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	var err error
	assert := suite.Assert()
	suite.stop = make(chan bool, 1)

	// Discards logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Start the authtest server for authentication verification
	suite.auth, err = authtest.NewServer()
	assert.NoError(err, "could not start the authtest server")

	// Start an httptest server to handle mock requests to Quarterdeck
	suite.quarterdeck, err = mock.NewServer()
	assert.NoError(err, "could not start the quarterdeck mock server")

	// Ensure Quarterdeck returns a 200 on status so Tenant knows it's ready
	suite.quarterdeck.OnStatus(mock.UseStatus(http.StatusOK))

	// Creates a test configuration to run the Tenant API server as a fully
	// functional server on an open port using the local-loopback for networking.
	conf, err := config.Config{
		Maintenance:  false,
		BindAddr:     "127.0.0.1:0",
		Mode:         gin.TestMode,
		LogLevel:     logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:   false,
		AllowOrigins: []string{"http://localhost:3000"},
		Auth: config.AuthConfig{
			Audience:     authtest.Audience,
			Issuer:       authtest.Issuer,
			KeysURL:      suite.auth.KeysURL(),
			CookieDomain: "127.0.0.1",
		},
		SendGrid: emails.Config{
			FromEmail:  "ensign@rotational.io",
			AdminEmail: "admins@rotational.io",
			Testing:    true,
		},
		Quarterdeck: config.QuarterdeckConfig{
			URL:          suite.quarterdeck.URL(),
			WaitForReady: 1 * time.Second,
		},
		Database: config.DatabaseConfig{
			Testing: true,
		},
		Ensign: config.SDKConfig{
			Enabled:          true,
			Endpoint:         "bufconn",
			Insecure:         true,
			NoAuthentication: true,
			WaitForReady:     1 * time.Second,
			Testing:          true,
		},
		MetaTopic: config.MetaTopicConfig{
			TopicName: "meta",
			SDKConfig: config.SDKConfig{
				Enabled:          true,
				Endpoint:         "bufconn",
				Insecure:         true,
				NoAuthentication: true,
				WaitForReady:     1 * time.Second,
				Testing:          true,
			},
		},
	}.Mark()
	assert.NoError(err, "test configuration is invalid")

	suite.srv, err = tenant.New(conf)
	assert.NoError(err, "could not create the tenant api server from the test configuration")

	suite.srv.Server = *service.New(conf.BindAddr, service.WithMode(conf.Mode), service.WithTLS(tlstest.Config()))
	suite.srv.Server.Register(suite.srv)

	// Starts the Tenant server. Server will run for the duration of all tests.
	// Implements reset methods to ensure the server state doesn't change
	// between tests in Before/After.
	go func() {
		if err := suite.srv.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// This is a bad enough error that we should panic, otherwise tests will
			// fail and it will be hard to debug why the tests failed.
			panic(err)
		}
		suite.stop <- true
	}()

	// Wait for 500ms to ensure the API server starts.
	time.Sleep(500 * time.Millisecond)

	// Creates a Tenant client to make requests to the server.
	assert.NotEmpty(suite.srv.URL(), "no url to connect the client on")
	suite.client, err = api.New(suite.srv.URL(), api.WithClient(tlstest.Client()))
	assert.NoError(err, "could not initialize the Tenant client")

	// Fetch the ensign mock for the tests
	client := suite.srv.GetEnsignClient()
	suite.ensign = client.GetMockServer()

	// Create the topic subscriber and fetch the mock for the tests
	suite.subscriber, err = tenant.NewTopicSubscriber(conf.MetaTopic)
	assert.NoError(err, "could not create the meta topic subscriber")
	suite.metatopic = suite.subscriber.GetEnsignClient().GetMockServer()
}

func (suite *tenantTestSuite) TearDownSuite() {
	assert := suite.Assert()

	// Shutdown the quarterdeck mock server
	suite.quarterdeck.Close()

	// Shutdown the ensign mock server
	suite.ensign.Shutdown()

	// Shutdown the authtest server
	suite.auth.Close()

	// Shuts down the tenant API server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := suite.srv.Shutdown(ctx)
	assert.NoError(err, "could not gracefully shut down the tenant test server")

	// Waits for server to stop in order to prevent race conditions.
	<-suite.stop

	// Cleanup logger
	logger.ResetLogger()
}

func (suite *tenantTestSuite) AfterTest(suiteName, testName string) {
	// Ensure any credentials set on the client are reset
	suite.client.(*api.APIv1).SetCredentials("")
	suite.client.(*api.APIv1).SetCSRFProtect(false)

	// Reset the quarterdeck mock server
	suite.quarterdeck.Reset()
}

// Stop the task manager, waiting for all the tasks to finish. Tests should defer
// ResetTasks() to ensure that the task manager is available to the other tests.
func (suite *tenantTestSuite) StopTasks() {
	tasks := suite.srv.GetTaskManager()
	tasks.Stop()
}

// Reset the task manager to ensure that other tests have access to it.
func (suite *tenantTestSuite) ResetTasks() {
	suite.srv.ResetTaskManager()
}

// Helper function to set cookies for CSRF protection on the tenant client
func (s *tenantTestSuite) SetClientCSRFProtection() error {
	s.client.(*api.APIv1).SetCSRFProtect(true)
	return nil
}

// Helper function to clear cookies for CSRF protection on the tenant client
func (s *tenantTestSuite) ClearClientCSRFProtection() error {
	s.client.(*api.APIv1).SetCSRFProtect(false)
	return nil
}

// Helper function to set the credentials on the test client from claims, reducing 3 or
// 4 lines of code into a single helper function call to make tests more readable.
func (s *tenantTestSuite) SetClientCredentials(claims *tokens.Claims) error {
	token, err := s.auth.CreateAccessToken(claims)
	if err != nil {
		return err
	}

	s.client.(*api.APIv1).SetCredentials(token)
	return nil
}

// Helper function to add the user claims to the current context.
func (s *tenantTestSuite) ContextWithClaims(ctx context.Context, claims *tokens.Claims) (c context.Context, err error) {
	token, err := s.auth.CreateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	return qd.ContextWithToken(ctx, token), nil
}

// Helper function to get the access token from the cookies.
func (s *tenantTestSuite) GetClientAccessToken() (string, error) {
	return s.client.(*api.APIv1).AccessToken()
}

// Helper function to get the refresh token from the cookies.
func (s *tenantTestSuite) GetClientRefreshToken() (string, error) {
	return s.client.(*api.APIv1).RefreshToken()
}

// Helper function to set the access and refresh tokens in the cookies.
func (s *tenantTestSuite) SetAuthTokens(access, refresh string) {
	s.client.(*api.APIv1).SetAuthTokens(access, refresh)
}

// Helper function to clear the access and refresh tokens from the cookies.
func (s *tenantTestSuite) ClearAuthTokens() {
	s.client.(*api.APIv1).ClearAuthTokens()
}

func TestTenant(t *testing.T) {
	suite.Run(t, &tenantTestSuite{})
}

func statusMessage(status int, message string) string {
	return fmt.Sprintf("[%d] %s", status, message)
}

var httpErrorRE = regexp.MustCompile(`^\[(?P<status>\d+)\] (?P<message>.*)$`)

// requireHTTPError asserts that an HTTP error has the matching status code and that
// the message matches one of the standard error responses.
func (s *tenantTestSuite) requireHTTPError(err error, status int) {
	require := s.Require()
	require.Error(err, "expected an error but didn't get one")
	matches := httpErrorRE.FindStringSubmatch(err.Error())
	require.NotNil(matches, "expected error message to be in the format '[status] message'")

	// Status code must match
	statusIndex := httpErrorRE.SubexpIndex("status")
	require.GreaterOrEqual(statusIndex, 0, "could not parse status code from error message")
	code, err := strconv.Atoi(matches[statusIndex])
	require.NoError(err, "could not parse status code as integer")
	require.Equal(status, code, "expected error status code to match")

	// Message must be one of the standard error responses
	msgIndex := httpErrorRE.SubexpIndex("message")
	require.GreaterOrEqual(msgIndex, 0, "could not parse message from error message")
	require.Contains(responses.AllResponses, strings.TrimSpace(matches[msgIndex]), "expected error message to be one of the standard error responses")
}

func (s *tenantTestSuite) requireError(err error, status int, message string, msgAndArgs ...interface{}) {
	require := s.Require()
	require.EqualError(err, statusMessage(status, message), msgAndArgs...)
}

func (s *tenantTestSuite) requireMultiError(err error, messages ...string) {
	require := s.Require()
	require.IsType(&multierror.Error{}, err)

	var actual []string
	for _, e := range err.(*multierror.Error).Errors {
		actual = append(actual, e.Error())
	}

	require.ElementsMatch(messages, actual)
}

// Asserts that the access and refresh tokens were set in the cookies by checking the
// client's cookie jar.
func (s *tenantTestSuite) requireAuthCookies(access, refresh string) {
	require := s.Require()
	token, err := s.GetClientAccessToken()
	require.NoError(err, "could not get access token from client")
	require.Equal(access, token, "wrong access token in cookies")

	token, err = s.GetClientRefreshToken()
	require.NoError(err, "could not get refresh token from client")
	require.Equal(refresh, token, "wrong refresh token in cookies")
}

// Asserts that the access and refresh tokens are not set in the cookies by checking
// the client's cookie jar.
func (s *tenantTestSuite) requireNoAuthCookies() {
	require := s.Require()
	_, err := s.GetClientAccessToken()
	require.ErrorIs(err, api.ErrNoAccessToken, "expected no access token in cookies")
	_, err = s.GetClientRefreshToken()
	require.ErrorIs(err, api.ErrNoRefreshToken, "expected no refresh token in cookies")
}

func (s *tenantTestSuite) TestRefreshCookies() {
	// This test asserts that the Authenticate middleware is properly configured to
	// automatically refresh the access token when it expires.
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup claims with an invalid access token but refresh token is in the cookies
	claims := &tokens.Claims{}
	refresh, err := s.auth.CreateToken(claims)
	require.NoError(err, "could not create refresh token")
	s.SetAuthTokens("", refresh)

	// Setup the Quarterdeck mock to return a valid token pair
	reply := &qd.LoginReply{}
	reply.AccessToken, reply.RefreshToken, err = s.auth.CreateTokenPair(claims)
	require.NoError(err, "could not create token pair")
	s.quarterdeck.OnRefresh(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(reply))

	// Should be no authentication issues on request
	require.NoError(s.SetClientCSRFProtection())
	_, err = s.client.InviteAccept(ctx, &api.MemberInviteToken{})
	s.requireHTTPError(err, http.StatusBadRequest)

	// New tokens should be available in the cookies
	s.requireAuthCookies(reply.AccessToken, reply.RefreshToken)
	s.ClearAuthTokens()
}
