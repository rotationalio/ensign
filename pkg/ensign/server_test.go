package ensign_test

import (
	"context"
	"os"
	"testing"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serverTestSuite struct {
	suite.Suite
	conf    config.Config
	srv     *ensign.Server
	client  api.EnsignClient
	conn    *bufconn.Listener
	dataDir string
}

func (s *serverTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// Create a temporary data directory
	s.dataDir, err = os.MkdirTemp("", "ensign-data-*")
	require.NoError(err)

	// This configuration will run the ensign server as a fully functional gRPC service
	// on an in-memory socket allowing the testing of RPCs from the client perspective.
	s.conf, err = config.Config{
		Maintenance: false,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
		BindAddr:    "127.0.0.1:0",
		Storage: config.StorageConfig{
			ReadOnly: false,
			DataPath: s.dataDir,
		},
	}.Mark()
	require.NoError(err, "could not mark test configuration as valid")

	// Create the server and run it on a bufconn.
	s.srv, err = ensign.New(s.conf)
	require.NoError(err, "could not create server with a test configuration")

	s.conn = bufconn.New()
	go s.srv.Run(s.conn.Sock())

	// Create a client for testing purposes
	cc, err := s.conn.Connect(context.Background(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(err, "could not connect to bufconn")
	s.client = api.NewEnsignClient(cc)
}

func (s *serverTestSuite) TearDownSuite() {
	require := s.Require()
	require.NoError(s.srv.Shutdown(), "could not shutdown the ensign server")

	require.NoError(os.RemoveAll(s.dataDir), "could not clean up temporary data directory")
	logger.ResetLogger()
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
