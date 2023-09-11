package broker_test

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	. "github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *brokerTestSuite) TestBroker() {
	var wg, pubwg, readywg sync.WaitGroup

	assert := s.Assert()
	require := s.Require()
	s.broker.Run(s.echan)

	var sent, recv, acks uint32

	// Create publisher and subscriber go routines
	nevents := 12
	npubs, nsubs := 10, 10
	topics := make([]ulid.ULID, nsubs)

	for i := 0; i < nsubs; i++ {
		wg.Add(1)
		readywg.Add(1)
		topics[i] = ulid.Make()

		go func(i int) {
			defer wg.Done()

			s.T().Logf("subscriber %d started", i)
			_, C, err := s.broker.Subscribe(topics[i])
			assert.NoError(err, "could not register subscriber")
			readywg.Done()

			for range C {
				atomic.AddUint32(&recv, 1)
			}
			s.T().Logf("subscriber %d finished", i)
		}(i)
	}

	// Wait for all subscribers to come online before starting publishers.
	readywg.Wait()

	for i := 0; i < npubs; i++ {
		wg.Add(1)
		pubwg.Add(1)

		pubID, C, err := s.broker.Register()
		require.NoError(err, "could not registered publisher")

		go func(i int, C <-chan PublishResult) {
			defer wg.Done()

			s.T().Logf("publisher ack recv %d started", i)
			for range C {
				atomic.AddUint32(&acks, 1)
			}
			s.T().Logf("publisher ack recv %d finished", i)
		}(i, C)

		go func(i int, pubID rlid.RLID) {
			defer pubwg.Done()
			topic := topics[i]

			s.T().Logf("publisher %d (%s) publishing %d events to topic %s", i, pubID, nevents, topic)
			for n := 0; n < nevents; n++ {
				s.broker.Publish(pubID, &api.EventWrapper{TopicId: topic.Bytes()})
				atomic.AddUint32(&sent, 1)
			}
			s.T().Logf("publisher %d (%s) finished", i, pubID)
		}(i, pubID)
	}

	// Wait for all publishers to finish sending their events, then shutdown.
	s.T().Log("waiting for publishers to finish sending events")
	pubwg.Wait()

	s.T().Log("shutting down the broker")
	err := s.broker.Shutdown()
	require.NoError(err, "could not shutdown broker")

	// Wait for all go routines to stop to start checking results
	s.T().Log("waiting for all go routines to stop")
	wg.Wait()

	nacks := atomic.LoadUint32(&acks)
	nsent := atomic.LoadUint32(&sent)
	nrecv := atomic.LoadUint32(&recv)
	s.T().Logf("%d sent %d acks %d recv", nsent, nacks, nrecv)
	require.Equal(nsent, nacks, "the expected number of events were not published with acks")
	require.Equal(nsent, nrecv, "the expected number of events was not received by subs")
}

func (s *brokerTestSuite) TestWriteErrors() {
	// If the events cannot be written to disk, nacks should be sent to publisher and
	// the subscribers should not receive any events.
	assert := s.Assert()
	require := s.Require()

	s.events.UseError(mock.Insert, errors.New("unable to write event to disk"))
	s.broker.Run(s.echan)

	topic := ulid.Make()
	nEvents := 21
	var wg, pubs, ready sync.WaitGroup
	var sent, recv, acks, nacks uint32

	// Create publisher and subscriber go routines (only one of each for this test)
	wg.Add(1)
	ready.Add(1)
	go func() {
		defer wg.Done()

		_, C, err := s.broker.Subscribe(topic)
		assert.NoError(err, "could not register subscriber")
		ready.Done()

		for range C {
			atomic.AddUint32(&recv, 1)
		}
	}()

	// Wait for subscribers to come online before starting publishers
	ready.Wait()

	wg.Add(1)
	pubs.Add(1)

	pubID, C, err := s.broker.Register()
	require.NoError(err, "could not register publisher")

	go func(C <-chan PublishResult) {
		defer wg.Done()
		for result := range C {
			if result.IsNack() {
				atomic.AddUint32(&nacks, 1)
			} else {
				atomic.AddUint32(&acks, 1)
			}
		}
	}(C)

	go func(pubID rlid.RLID) {
		defer pubs.Done()
		for n := 0; n < nEvents; n++ {
			s.broker.Publish(pubID, &api.EventWrapper{TopicId: topic.Bytes()})
			atomic.AddUint32(&sent, 1)
		}
	}(pubID)

	// Wait for publishers to finish then shutdown
	pubs.Wait()
	err = s.broker.Shutdown()
	require.NoError(err, "could not shutdown broker")

	// Wait for subscribers to finish before checking result
	wg.Wait()

	nNacks := atomic.LoadUint32(&nacks)
	nAcks := atomic.LoadUint32(&acks)
	nSent := atomic.LoadUint32(&sent)
	nRecv := atomic.LoadUint32(&recv)
	require.Equal(uint32(nEvents), nSent, "expected a fixed number of events to be sent")
	require.Zero(nRecv, "expected no messages to be recv by subscribers when erroring")
	require.Zero(nAcks, "expected no acks back to the publisher, only nacks")
	require.Equal(nSent, nNacks, "expected a nack for every message sent")
}

func (s *brokerTestSuite) TestBrokerStartupShutdown() {
	require := s.Require()
	nroutines := runtime.NumGoroutine()

	// Test shutdown with no pubs/subs
	s.broker.Run(s.echan)
	require.Greater(runtime.NumGoroutine(), nroutines, "expected more go routines to be running than before")

	err := s.broker.Shutdown()
	require.NoError(err, "could not shutdown broker")
	require.Less(runtime.NumGoroutine(), nroutines+2, "expected fewer go routines afer shutdown")
	require.NoError(s.broker.Shutdown(), "should be able to call shutdown when broker is not running")
	time.Sleep(50 * time.Millisecond)

	// Test shutdown with pubs/subs
	s.broker.Run(s.echan)
	time.Sleep(50 * time.Millisecond)
	require.Greater(runtime.NumGoroutine(), nroutines, "unable to start broker after shutdown")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		pubID, ch, err := s.broker.Register()
		require.NoError(err, "could not register publisher")
		require.NotZero(pubID, "no publisher ID returned")
		go func(C <-chan PublishResult) {
			<-C
			wg.Done()
		}(ch)

		subID, evts, err := s.broker.Subscribe(ulid.Make())
		require.NoError(err, "could not register subscriber")
		require.NotZero(subID, "no subscriber ID returned")

		go func(C <-chan *api.EventWrapper) {
			<-C
			wg.Done()
		}(evts)
	}

	require.Equal(10, s.broker.NumPublishers())
	require.Equal(10, s.broker.NumSubscribers())

	err = s.broker.Shutdown()
	require.NoError(err, "could not shutdown broker with pubs/subs")

	// If the tests times out, it is because the broker didn't correctly close channels
	wg.Wait()

	time.Sleep(50 * time.Millisecond)
	require.Equal(0, s.broker.NumPublishers())
	require.Equal(0, s.broker.NumSubscribers())
}

func (s *brokerTestSuite) TestRegisterClose() {
	s.broker.Run(s.echan)
	require := s.Require()

	require.Equal(0, s.broker.NumPublishers(), "expected 0 publishers in initialized broker")
	require.Equal(0, s.broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Closing an unknown publisher should return an error
	err := s.broker.Close(rlid.Make(42))
	require.ErrorIs(err, ErrUnknownID)

	// Register a publisher
	pubID, cb, err := s.broker.Register()
	require.NoError(err, "could not register publisher")
	require.Equal(1, s.broker.NumPublishers(), "expected publisher after register")
	require.Equal(0, s.broker.NumSubscribers(), "expected 0 subscribers in initialized broker")

	// Register a subscriber
	subID, evts, err := s.broker.Subscribe(ulid.Make())
	require.NoError(err, "could not register subscriber")
	require.Equal(1, s.broker.NumPublishers(), "expected publisher after register")
	require.Equal(1, s.broker.NumSubscribers(), "expected subscriber after subscribe")

	// Close the publisher
	require.NoError(s.broker.Close(pubID))

	// Perform a non-blocking read of the channel so tests don't timeout
	open := true
	select {
	case _, open = <-cb:
	default:
	}

	require.False(open, "expected publisher channel to be closed")
	require.Equal(0, s.broker.NumPublishers(), "expected no publishers after close")
	require.Equal(1, s.broker.NumSubscribers(), "expected subscriber after subscribe")

	// Close the subscriber
	require.NoError(s.broker.Close(subID))

	// Perform non blocking read of the channel so tests don't timeout
	open = true
	select {
	case _, open = <-evts:
	default:
	}

	require.False(open, "expected subscriber channel to be closed")
	require.Equal(0, s.broker.NumPublishers(), "expected no publishers after close")
	require.Equal(0, s.broker.NumSubscribers(), "expected no subscribers after close")
}

func (s *brokerTestSuite) TestNoPublishNotRunning() {
	require := s.Require()

	// Should not be able to register, subscribe, or publish if broker is not running.
	pubID, cb, err := s.broker.Register()
	require.ErrorIs(err, ErrBrokerNotRunning)
	require.Nil(cb)
	require.Zero(pubID)

	subID, evts, err := s.broker.Subscribe(ulid.Make())
	require.ErrorIs(err, ErrBrokerNotRunning)
	require.Nil(evts)
	require.Zero(subID)

	err = s.broker.Publish(rlid.Make(42), &api.EventWrapper{})
	require.ErrorIs(err, ErrBrokerNotRunning)

	s.broker.Run(nil)
	s.broker.Shutdown()

	pubID, cb, err = s.broker.Register()
	require.ErrorIs(err, ErrBrokerNotRunning)
	require.Nil(cb)
	require.Zero(pubID)

	subID, evts, err = s.broker.Subscribe(ulid.Make())
	require.ErrorIs(err, ErrBrokerNotRunning)
	require.Nil(evts)
	require.Zero(subID)

	err = s.broker.Publish(rlid.Make(24), &api.EventWrapper{})
	require.ErrorIs(err, ErrBrokerNotRunning)
}

func (s *brokerTestSuite) TestAckCommittedTimestamp() {
	require := s.Require()
	s.broker.Run(s.echan)

	var wg sync.WaitGroup
	pubID, C, err := s.broker.Register()
	require.NoError(err, "could not register publisher")

	acks := make([]PublishResult, 0, 10)
	wg.Add(1)
	go func(C <-chan PublishResult) {
		defer wg.Done()
		for result := range C {
			acks = append(acks, result)
		}
	}(C)

	wg.Add(1)
	go func(pubID rlid.RLID) {
		defer wg.Done()
		topicID := ulids.New()
		for n := 0; n < 10; n++ {
			s.broker.Publish(pubID, &api.EventWrapper{TopicId: topicID[:]})
		}
		s.broker.Shutdown()
	}(pubID)

	wg.Wait()

	require.Len(acks, 10, "unexpected number of acks recieved")
	for _, result := range acks {
		require.True(result.Committed.IsValid(), "committed timestamp is not valid")
		require.False(result.Committed.AsTime().IsZero(), "committed timestamp is zero valued")
	}
}
