package ensign_test

import (
	"context"
	"os"
	"testing"
	"time"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/authtest"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
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

	// Create a client for testing purposes
	cc, err := s.conn.Connect(context.Background(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(err, "could not connect to bufconn")
	s.client = api.NewEnsignClient(cc)

	// Keep a reference to the mock store to make testing easier
	s.store = s.srv.StoreMock()
}

func (s *serverTestSuite) TearDownSuite() {
	assert := s.Assert()
	assert.NoError(s.srv.Shutdown(), "could not shutdown the ensign server")

	// Close the authtest server
	s.quarterdeck.Close()

	assert.NoError(os.RemoveAll(s.dataDir), "could not clean up temporary data directory")
	logger.ResetLogger()
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

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
