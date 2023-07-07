package broker_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	. "github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroker(t *testing.T) {
	logger.Discard()
	defer logger.ResetLogger()

	var wg, pubwg sync.WaitGroup

	broker := New()
	broker.Run(nil)

	var sent, recv, acks uint32

	// Create publisher and subscriber go routines
	nevents := 12
	npubs, nsubs := 10, 10
	topics := make([]ulid.ULID, nsubs)

	for i := 0; i < nsubs; i++ {
		wg.Add(1)
		topics[i] = ulid.Make()

		go func(i int) {
			defer wg.Done()
			_, C, err := broker.Subscribe(topics[i])
			assert.NoError(t, err, "could not register subscriber")

			for range C {
				atomic.AddUint32(&recv, 1)
			}
		}(i)
	}

	for i := 0; i < npubs; i++ {
		wg.Add(1)
		pubwg.Add(1)

		pubID, C, err := broker.Register()
		require.NoError(t, err, "could not registered publisher")

		go func(i int, C <-chan PublishResult) {
			defer wg.Done()
			for range C {
				atomic.AddUint32(&acks, 1)
			}
		}(i, C)

		go func(i int, pubID rlid.RLID) {
			defer pubwg.Done()
			topic := topics[i]

			for n := 0; n < nevents; n++ {
				broker.Publish(pubID, &api.EventWrapper{TopicId: topic.Bytes()})
				atomic.AddUint32(&sent, 1)
			}
		}(i, pubID)
	}

	// Wait for all publishers to finish sending their events, then shutdown.
	pubwg.Wait()

	err := broker.Shutdown()
	require.NoError(t, err, "could not shutdown broker")

	// Wait for all go routines to stop to start checking results
	wg.Wait()

	nacks := atomic.LoadUint32(&acks)
	nsent := atomic.LoadUint32(&sent)
	nrecv := atomic.LoadUint32(&recv)
	require.Equal(t, nsent, nacks, "the expected number of events were not published with acks")
	require.Equal(t, nsent, nrecv, "the expected number of events was not received by subs")

}

const runRoutines = 2

func TestBrokerStartupShutdown(t *testing.T) {
	logger.Discard()
	defer logger.ResetLogger()

	broker := New()
	nroutines := runtime.NumGoroutine()

	// Test shutdown with no pubs/subs
	broker.Run(nil)
	require.Equal(t, nroutines+runRoutines, runtime.NumGoroutine())

	err := broker.Shutdown()
	require.NoError(t, err, "could not shutdown broker")
	require.Equal(t, nroutines, runtime.NumGoroutine())
	require.NoError(t, broker.Shutdown(), "should be able to call shutdown when broker is not running")
	time.Sleep(50 * time.Millisecond)

	// Test shutdown with pubs/subs
	broker.Run(nil)
	time.Sleep(50 * time.Millisecond)
	require.Equal(t, nroutines+runRoutines, runtime.NumGoroutine(), "unable to start broker after shutdown")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		pubID, ch, err := broker.Register()
		require.NoError(t, err, "could not register publisher")
		require.NotZero(t, pubID, "no publisher ID returned")
		go func(C <-chan PublishResult) {
			<-C
			wg.Done()
		}(ch)

		subID, evts, err := broker.Subscribe(ulid.Make())
		require.NoError(t, err, "could not register subscriber")
		require.NotZero(t, subID, "no subscriber ID returned")

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

	time.Sleep(50 * time.Millisecond)
	require.Equal(t, 0, broker.NumPublishers())
	require.Equal(t, 0, broker.NumSubscribers())
}

func TestRegisterClose(t *testing.T) {
	broker := New()
	broker.Run(nil)
	defer broker.Shutdown()

	require.Equal(t, 0, broker.NumPublishers(), "expected 0 publishers in initialized broker")
	require.Equal(t, 0, broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Closing an unknown publisher should return an error
	err := broker.Close(rlid.Make(42))
	require.ErrorIs(t, err, ErrUnknownID)

	// Register a publisher
	pubID, cb, err := broker.Register()
	require.NoError(t, err, "could not register publisher")
	require.Equal(t, 1, broker.NumPublishers(), "expected publisher after register")
	require.Equal(t, 0, broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Register a subscriber
	subID, evts, err := broker.Subscribe(ulid.Make())
	require.NoError(t, err, "could not register subscriber")
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

func TestNoPublishNotRunning(t *testing.T) {
	// Should not be able to register, subscribe, or publish if broker is not running.
	broker := New()

	pubID, cb, err := broker.Register()
	require.ErrorIs(t, err, ErrBrokerNotRunning)
	require.Nil(t, cb)
	require.Zero(t, pubID)

	subID, evts, err := broker.Subscribe(ulid.Make())
	require.ErrorIs(t, err, ErrBrokerNotRunning)
	require.Nil(t, evts)
	require.Zero(t, subID)

	err = broker.Publish(rlid.Make(42), &api.EventWrapper{})
	require.ErrorIs(t, err, ErrBrokerNotRunning)

	broker.Run(nil)
	broker.Shutdown()

	pubID, cb, err = broker.Register()
	require.ErrorIs(t, err, ErrBrokerNotRunning)
	require.Nil(t, cb)
	require.Zero(t, pubID)

	subID, evts, err = broker.Subscribe(ulid.Make())
	require.ErrorIs(t, err, ErrBrokerNotRunning)
	require.Nil(t, evts)
	require.Zero(t, subID)

	err = broker.Publish(rlid.Make(24), &api.EventWrapper{})
	require.ErrorIs(t, err, ErrBrokerNotRunning)
}
