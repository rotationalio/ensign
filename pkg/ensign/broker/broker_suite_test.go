package broker_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/stretchr/testify/suite"
)

type brokerTestSuite struct {
	suite.Suite
	broker *broker.Broker
	events *mock.Store
	echan  chan error
}

func TestBroker(t *testing.T) {
	suite.Run(t, new(brokerTestSuite))
}

func (s *brokerTestSuite) SetupSuite() {
	// Create a mock store that returns no error when events are inserted
	s.events = &mock.Store{}
	s.events.UseError(mock.Insert, nil)

	// Discard all logging to prevent verbose test output
	logger.Discard()
}

func (s *brokerTestSuite) TearDownSuite() {
	logger.ResetLogger()
}

func (s *brokerTestSuite) BeforeTest(suiteName, testName string) {
	// Create a new broker that isn't running before each test
	// NOTE: tests must run the broker if they need it running.
	s.broker = broker.New(s.events)
	s.echan = make(chan error, 1)
}

func (s *brokerTestSuite) AfterTest(suiteName, testName string) {
	// Ensure the broker is shutdown after each test and the echan is closed
	s.broker.Shutdown()
	close(s.echan)
}
