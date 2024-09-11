package ensign_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type serverTestSuite struct {
	suite.Suite
	conf        config.Config
	quarterdeck *authtest.Server
	store       *mock.Store
	srv         *ensign.Server
	client      api.EnsignClient
	conn        *bufconn.Listener
	dataDir     string
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}

func (s *serverTestSuite) SetupSuite() {
	var err error
	assert := s.Assert()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a temporary data directory
	s.dataDir, err = os.MkdirTemp("", "ensign-data-*")
	assert.NoError(err)

	// Initialize authtest server
	s.quarterdeck, err = authtest.NewServer()
	assert.NoError(err, "could not initialize authtest server")

	// This configuration will run the ensign server as a fully functional gRPC service
	// on an in-memory socket allowing the testing of RPCs from the client perspective.
	s.conf, err = config.Config{
		Maintenance: false,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
		BindAddr:    "127.0.0.1:0",
		Monitoring: config.MonitoringConfig{
			Enabled: false,
			NodeID:  "localtest",
		},
		Storage: config.StorageConfig{
			Testing:  true,
			ReadOnly: false,
			DataPath: s.dataDir,
		},
		Auth: config.AuthConfig{
			KeysURL:            s.quarterdeck.KeysURL(),
			Audience:           authtest.Audience,
			Issuer:             authtest.Issuer,
			MinRefreshInterval: 5 * time.Minute,
		},
	}.Mark()
	assert.NoError(err, "could not mark test configuration as valid")

	// Create the server and run it on a bufconn.
	s.srv, err = ensign.New(s.conf)
	assert.NoError(err, "could not create server with a test configuration")

	s.conn = bufconn.New()
	go s.srv.Run(s.conn.Sock())
	time.Sleep(750 * time.Millisecond)

	// Create a client for testing purposes
	cc, err := s.conn.Connect(grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(err, "could not connect to bufconn")
	s.client = api.NewEnsignClient(cc)

	// Keep a reference to the mock store to make testing easier
	s.store = s.srv.StoreMock()

	// Register metrics without starting server
	err = o11y.Serve(s.conf.Monitoring)
	assert.NoError(err, "could not register o11y collectors")

	// Run the broker for handling events
	s.srv.RunBroker()
}

func (s *serverTestSuite) TearDownSuite() {
	assert := s.Assert()
	assert.NoError(s.srv.Shutdown(), "could not shutdown the ensign server")

	// Close the authtest server
	s.quarterdeck.Close()

	// Shutdown o11y server (which should not be running) but as a safeguard
	err := o11y.Shutdown(context.Background())
	assert.NoError(err, "could not shutdown o11y collectors")

	assert.NoError(os.RemoveAll(s.dataDir), "could not clean up temporary data directory")
	logger.ResetLogger()
}

func (s *serverTestSuite) AfterTest(_, _ string) {
	s.store.Reset()
}

// Check an error response from the gRPC Ensign client, ensuring that it is a) a status
// error, b) has the code specified, and c) (if supplied) that the message matches.
func (s *serverTestSuite) GRPCErrorIs(err error, code codes.Code, msg string) {
	require := s.Require()
	serr, ok := status.FromError(err)
	require.True(ok, "err is not a grpc status error")
	require.Equal(code, serr.Code(), "status code %s did not match expected %s", serr.Code(), code)
	if msg != "" {
		require.Equal(msg, serr.Message(), "status message did not match the expected message")
	}
}

// Check an error response from the gRPC Ensign client, ensuring that it is a) a status
// error, b) has the code specified, and c) (if supplied) that the message matches.
func GRPCErrorIs(t *testing.T, err error, code codes.Code, msg string) {
	serr, ok := status.FromError(err)
	require.True(t, ok, "err is not a grpc status error")
	require.Equal(t, code, serr.Code(), "status code %s did not match expected %s", serr.Code(), code)

	if msg != "" {
		require.Equal(t, msg, serr.Message(), "status message did not match the expected message")
	}
}

func (s *serverTestSuite) CheckTopicMap(projectID string, topics map[string][]byte) {
	require := s.Require()
	expected, ok := projectTopics[projectID]
	require.True(ok, "could not find project %q in fixtures", projectID)

	require.Len(topics, len(expected), "expected actual length to match fixtures")
	for tid, name := range expected {
		require.Contains(topics, name, "topic map should contain the topic name")

		expectedTopicID := ulid.MustParse(tid)
		actualTopicID, err := ulids.Parse(topics[name])
		require.NoError(err, "could not parse bytes returned in topic map")
		require.Equal(0, expectedTopicID.Compare(actualTopicID))
	}
}
