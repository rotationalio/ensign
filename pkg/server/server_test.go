package server_test

import (
	"context"
	"testing"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/config"
	"github.com/rotationalio/ensign/pkg/server"
	"github.com/rotationalio/ensign/pkg/utils/bufconn"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serverTestSuite struct {
	suite.Suite
	conf   config.Config
	srv    *server.Server
	client api.EnsignClient
	conn   *bufconn.Listener
}

func (s *serverTestSuite) SetupSuite() {
	var err error
	require := s.Require()

	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog must be false otherwise this will be overridden
	logger.Discard()

	// This configuration will run the ensign server as a fully functional gRPC service
	// on an in-memory socket allowing the testing of RPCs from the client perspective.
	s.conf, err = config.Config{
		Maintenance: false,
		LogLevel:    logger.LevelDecoder(zerolog.DebugLevel),
		ConsoleLog:  false,
		BindAddr:    "127.0.0.1:0",
	}.Mark()
	require.NoError(err, "could not mark test configuration as valid")

	// Create the server and run it on a bufconn.
	s.srv, err = server.New(s.conf)
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

	logger.ResetLogger()
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
