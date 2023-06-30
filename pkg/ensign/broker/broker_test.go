package broker_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	. "github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
)

func TestBrokerStartupShutdown(t *testing.T) {
	broker := New()
	nroutines := runtime.NumGoroutine()

	// Test shutdown with no pubs/subs
	broker.Run(nil)
	require.Equal(t, nroutines+1, runtime.NumGoroutine())

	err := broker.Shutdown()
	require.NoError(t, err, "could not shutdown broker")
	require.Equal(t, nroutines, runtime.NumGoroutine())
	require.NoError(t, broker.Shutdown(), "should be able to call shutdown when broker is not running")

	// Test shutdown with pubs/subs
	broker.Run(nil)
	require.Equal(t, nroutines+1, runtime.NumGoroutine(), "unable to start broker after shutdown")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		_, ch := broker.Register()
		go func(C <-chan bool) {
			<-C
			wg.Done()
		}(ch)

		_, evts := broker.Subscribe(ulid.Make())
		go func(C <-chan *api.EventWrapper) {
			<-C
			wg.Done()
		}(evts)
	}

	require.Equal(t, 10, broker.NumPublishers())
	require.Equal(t, 10, broker.NumSubscribers())

	err = broker.Shutdown()
	require.NoError(t, err, "could not shutdown broker with pubs/subs")

	// If the tests times out, it is because the broker didn't correctly close channels
	wg.Wait()

	require.Equal(t, 0, broker.NumPublishers())
	require.Equal(t, 0, broker.NumSubscribers())
}

func TestRegisterClose(t *testing.T) {
	broker := New()
	require.Equal(t, 0, broker.NumPublishers(), "expected 0 publishers in initialized broker")
	require.Equal(t, 0, broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Closing an unknown publisher should return an error
	err := broker.Close(rlid.Make(42))
	require.ErrorIs(t, err, ErrUnknownID)

	// Register a publisher
	pubID, cb := broker.Register()
	require.Equal(t, 1, broker.NumPublishers(), "expected publisher after register")
	require.Equal(t, 0, broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Register a subscriber
	subID, evts := broker.Subscribe(ulid.Make())
	require.Equal(t, 1, broker.NumPublishers(), "expected publisher after register")
	require.Equal(t, 1, broker.NumSubscribers(), "expected subscriber after subscribe")

	// Close the publisher
	require.NoError(t, broker.Close(pubID))

	// Perform a non-blocking read of the channel so tests don't timeout
	open := true
	select {
	case _, open = <-cb:
	default:
	}

	require.False(t, open, "expected publisher channel to be closed")
	require.Equal(t, 0, broker.NumPublishers(), "expected no publishers after close")
	require.Equal(t, 1, broker.NumSubscribers(), "expected subscriber after subscribe")

	// Close the subscriber
	require.NoError(t, broker.Close(subID))

	// Perform non blocking read of the channel so tests don't timeout
	open = true
	select {
	case _, open = <-evts:
	default:
	}

	require.False(t, open, "expected subscriber channel to be closed")
	require.Equal(t, 0, broker.NumPublishers(), "expected no publishers after close")
	require.Equal(t, 0, broker.NumSubscribers(), "expected no subscribers after close")
}
