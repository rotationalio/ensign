package ensign_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverTestSuite) TestPublishEvents() {
	// Should be able to publish a series of valid events
	require := s.Require()
	stream := s.setupValidPublisher()
	s.store.UseError(store.Insert, nil)
	defer s.store.Reset()

	events := make([]*api.EventWrapper, 0, 10)
	for i := 0; i < 10; i++ {
		events = append(events, MakeEmpty("01H6XTAVNM21F6JXNGAJF1SJ4S"))
	}

	results := stream.WithEventResults(&api.OpenStream{ClientId: "tester"}, events...)
	err := s.srv.Publish(stream)
	require.NoError(err, "was not able to publish events")

	var acks, nacks int
	for _, event := range events {
		if ack := results.Ack(event); ack != nil {
			acks++
		}

		if nack := results.Nack(event); nack != nil {
			nacks++
		}
	}

	require.Equal(12, stream.Calls(mock.StreamRecv), "expected 1 open stream, 10 events, and one io.EOF")
	require.Equal(12, stream.Calls(mock.StreamSend), "expected 1 ready, 10 events, and 1 close stream")

	require.Equal(0, nacks, "expected no nacks")
	require.Greater(acks, 5, "expected at least 5 acks")
}

func (s *serverTestSuite) TestPublisherStreamInitialization() {
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := &mock.PublisherServer{}
	s.store.OnAllowedTopics = MockAllowedTopics
	s.store.OnTopicName = MockTopicName
	defer s.store.Reset()

	// Must be authenticated and have the publisher permission.
	err := s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Create base claims to add to the stream context for authentication
	// These claims are valid but will have no topics associated with them.
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01H784KEP6F5EMW9CBYAHFB3J3",
		},
		OrgID:       "01H784KNY3GN2GC8NHW4ZKC5A9",
		ProjectID:   "01H784KXZPKMDWRX2ZRP6FSXET",
		Permissions: []string{permissions.Publisher},
	}
	stream.WithClaims(claims)

	// When no topics are associated with the project, should get a no topics available error
	err = s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "no topics available")

	// Change the claims over to a ProjectID that does contain topics and add a peer.
	claims.ProjectID = "01H6PGFTK2X53RGG2KMSGR2M61"
	stream.WithPeer(claims, MakePeer("172.92.121.6:10820"))

	// Handle stream closed before open stream message recv
	stream.WithError(mock.StreamRecv, io.EOF)
	require.NoError(s.srv.Publish(stream))

	// Handle stream crashed before open stream message recv
	stream.WithError(mock.StreamRecv, context.DeadlineExceeded)
	require.ErrorIs(s.srv.Publish(stream), context.DeadlineExceeded)

	// An OpenStream message must be the first message received by the handler
	stream.OnRecv = func() (*api.PublisherRequest, error) {
		return &api.PublisherRequest{
			Embed: &api.PublisherRequest_Event{
				Event: &api.EventWrapper{
					Event: []byte("foo"),
				},
			},
		}, nil
	}

	err = s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "an open stream message must be sent immediately after opening the stream")

	// Setup the open stream message to be returned with no events
	stream.WithEvents(&api.OpenStream{ClientId: "tester", Topics: nil})

	// Handle stream closed before stream ready message sent
	stream.WithError(mock.StreamSend, io.EOF)
	require.NoError(s.srv.Publish(stream))

	// Handle stream crashed before stream ready message sent
	stream.WithEvents(&api.OpenStream{ClientId: "tester", Topics: nil})
	stream.WithError(mock.StreamSend, context.DeadlineExceeded)
	s.GRPCErrorIs(s.srv.Publish(stream), codes.Aborted, "could not initialize publish stream")

	// Test the happy case - stream initialized but EOF before any events are sent.
	replies := make(chan *api.PublisherReply, 2)
	stream.Capture(replies)
	stream.WithEvents(&api.OpenStream{ClientId: "tester", Topics: nil})

	err = s.srv.Publish(stream)
	require.NoError(err, "the stream should have shutdown after recv EOF")

	// Should recv two message, the stream ready message and the stream closed message
	reply := <-replies
	ready := reply.GetReady()
	require.NotNil(ready, "expected the first message to be a ready message")

	require.Equal("tester", ready.ClientId)
	require.Equal("localtest", ready.ServerId)
	s.CheckTopicMap("01H6PGFTK2X53RGG2KMSGR2M61", ready.Topics)

	reply = <-replies
	closed := reply.GetCloseStream()
	require.NotNil(closed, "expected the last message to be a stream closed message")

	require.Equal(uint64(0), closed.Events)
	require.Equal(uint64(0), closed.Topics)
	require.Equal(uint64(0), closed.Acks)
	require.Equal(uint64(0), closed.Nacks)
}

func (s *serverTestSuite) TestSecondOpenStreamFail() {
	// Test cannot send a second open stream message after first open stream message
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := s.setupValidPublisher()
	defer s.store.Reset()

	// An OpenStream message must be the first message received by the handler
	stream.OnRecv = func() (*api.PublisherRequest, error) {
		return &api.PublisherRequest{
			Embed: &api.PublisherRequest_OpenStream{
				OpenStream: &api.OpenStream{
					ClientId: "tester",
				},
			},
		}, nil
	}

	err := s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.Aborted, "cannot send multiple open stream messages")
	require.Equal(2, stream.Calls(mock.StreamSend), "recv should have been called twice")
}

func (s *serverTestSuite) TestBadPublisherRequest() {
	// Test cannot send a second open stream message after first open stream message
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := s.setupValidPublisher()
	defer s.store.Reset()

	// An OpenStream message must be the first message received by the handler
	msg := 0
	stream.OnRecv = func() (*api.PublisherRequest, error) {
		if msg == 0 {
			msg++
			return &api.PublisherRequest{
				Embed: &api.PublisherRequest_OpenStream{
					OpenStream: &api.OpenStream{
						ClientId: "tester",
					},
				},
			}, nil
		}
		return &api.PublisherRequest{}, nil
	}

	err := s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "unhandled publisher request message <nil>")
	require.Equal(2, stream.Calls(mock.StreamSend), "recv should have been called twice")
}

func (s *serverTestSuite) TestPublisherStreamTopicFilter() {
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := s.setupValidPublisher()
	defer s.store.Reset()

	// Should receive an error when one of the topics is not in allowed topics
	stream.WithEvents(&api.OpenStream{ClientId: "tester", Topics: []string{"foo", "bar"}})
	err := s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "no topics available")
	require.Equal(0, stream.Calls(mock.StreamSend), "no messages should have been sent back to client")

	// Should receive a nack when an event is published to a topic that wasn't in the filter list
	event := MakeEmpty("01H6XTB5DS8YG0YZEVQ385QRTB")
	results := stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: []string{"01H6XTAVNM21F6JXNGAJF1SJ4S"}}, event)
	err = s.srv.Publish(stream)
	require.NoError(err)

	nack := results.Nack(event)
	require.NotNil(nack, "expected a nack for the given event")
	require.Equal(api.Nack_TOPIC_UNKNOWN, nack.Code)
	require.Equal("topic unknown or unhandled", nack.Error)

	// Should be able to publish an event that is in the filter list.
	event = MakeEmpty("01H6XTAVNM21F6JXNGAJF1SJ4S")
	_ = stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: []string{"01H6XTAVNM21F6JXNGAJF1SJ4S"}}, event)
	err = s.srv.Publish(stream)
	require.NoError(err)

	// NOTE: the test here is that we did not receive an ack. It might seem to make more
	// sense to check for an ack, but the problem is that we cannot synchronize with the
	// broker; so a race condition occurs if we check the ack
	require.Nil(results.Nack(event))
}

func (s *serverTestSuite) TestPublisherNackEvents() {
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := s.setupValidPublisher()
	defer s.store.Reset()

	// Should receive a nack when an event is published without a topic ID
	event := &api.EventWrapper{LocalId: []byte("abc")}
	results := stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: nil}, event)
	err := s.srv.Publish(stream)
	require.NoError(err)

	nack := results.Nack(event)
	require.NotNil(nack)
	require.Equal(api.Nack_TOPIC_UNKNOWN, nack.Code)
	require.Equal("no topic id specified", nack.Error)

	// Should receive a nack when an event is published without a valid topic ID
	event = &api.EventWrapper{TopicId: []byte("foo"), LocalId: []byte("abc")}
	results = stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: nil}, event)
	err = s.srv.Publish(stream)
	require.NoError(err)

	nack = results.Nack(event)
	require.NotNil(nack)
	require.Equal(api.Nack_TOPIC_UNKNOWN, nack.Code)
	require.Equal("invalid topic id", nack.Error)

	// Should receive a nack when an event is published to a topic not in the project
	event = MakeEmpty("01H7KAXSB26K6QRV08V7K74Z4B")
	results = stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: nil}, event)
	err = s.srv.Publish(stream)
	require.NoError(err)

	nack = results.Nack(event)
	require.NotNil(nack)
	require.Equal(api.Nack_TOPIC_UNKNOWN, nack.Code)
	require.Equal(api.CodeTopicUnknown, nack.Error)

	// Should receive a nack when an event is published with max data size exceeded
	event = MakeEvent("01H6XTAVNM21F6JXNGAJF1SJ4S", &api.Event{Data: bytes.Repeat([]byte("abc"), 1024*1024*6)})
	results = stream.WithEventResults(&api.OpenStream{ClientId: "tester", Topics: nil}, event)
	err = s.srv.Publish(stream)
	require.NoError(err)

	nack = results.Nack(event)
	require.NotNil(nack)
	require.Equal(api.Nack_MAX_EVENT_SIZE_EXCEEDED, nack.Code)
	require.Equal(api.CodeMaxEventSizeExceeded, nack.Error)
}

func TestPublisherHandler(t *testing.T) {
	store := &store.Store{}
	stream := &mock.PublisherServer{}
	handler := ensign.NewPublisherHandler(stream, store)

	t.Run("Ack", func(t *testing.T) {
		defer store.Reset()
		defer stream.Reset()

		localID := ulid.Make().Bytes()
		committed := time.Date(2023, 12, 12, 12, 12, 12, 0, time.UTC)

		stream.WithError(mock.StreamSend, io.EOF)
		err := handler.Ack(localID, committed)
		require.ErrorIs(t, err, io.EOF, "ack should return a send error")

		stream.OnSend = func(in *api.PublisherReply) error {
			ack := in.GetAck()
			if ack == nil {
				return errors.New("expected ack in publisher reply")
			}

			if !bytes.Equal(localID, ack.Id) {
				return errors.New("unexpected local ID in ack")
			}

			if !committed.Equal(ack.Committed.AsTime()) {
				return errors.New("unexpected committed timestamp in ack")
			}

			return nil
		}

		err = handler.Ack(localID, committed)
		require.NoError(t, err, "could not send correct ack")
		require.Equal(t, 2, stream.Calls(mock.StreamSend))
	})

	t.Run("Nack", func(t *testing.T) {
		defer store.Reset()
		defer stream.Reset()

		localID := ulid.Make().Bytes()
		ncode := api.Nack_MAX_EVENT_SIZE_EXCEEDED
		emsg := "this event was too big"

		// Handle error checking
		stream.WithError(mock.StreamSend, io.EOF)
		err := handler.Nack(localID, ncode, emsg)
		require.ErrorIs(t, err, io.EOF, "nack should return a send error")

		// Should use the message specified
		stream.OnSend = func(in *api.PublisherReply) error {
			nack := in.GetNack()
			if nack == nil {
				return errors.New("expected nack in publisher reply")
			}

			if !bytes.Equal(localID, nack.Id) {
				return errors.New("unexpected local ID in nack")
			}

			if ncode != nack.Code || emsg != nack.Error {
				return errors.New("unexpected nack code or error message")
			}

			return nil
		}

		err = handler.Nack(localID, ncode, emsg)
		require.NoError(t, err, "could not send correct ack")
		require.Equal(t, 2, stream.Calls(mock.StreamSend))

		// Should use default error message
		stream.OnSend = func(in *api.PublisherReply) error {
			nack := in.GetNack()
			if nack == nil {
				return errors.New("expected nack in publisher reply")
			}

			if !bytes.Equal(localID, nack.Id) {
				return errors.New("unexpected local ID in nack")
			}

			if nack.Code != ncode || nack.Error != api.CodeMaxEventSizeExceeded {
				return errors.New("unexpected nack code or error message")
			}

			return nil
		}

		err = handler.Nack(localID, ncode, "")
		require.NoError(t, err, "could not send correct ack")
		require.Equal(t, 3, stream.Calls(mock.StreamSend))
	})

	t.Run("Reply", func(t *testing.T) {
		defer store.Reset()
		defer stream.Reset()

		ack := broker.PublishResult{
			LocalID:   ulid.Make().Bytes(),
			Committed: timestamppb.Now(),
		}
		nack := broker.PublishResult{
			LocalID: ulid.Make().Bytes(),
			Code:    api.Nack_PERMISSION_DENIED,
		}

		// Handle error checking
		stream.WithError(mock.StreamSend, io.EOF)
		err := handler.Reply(ack)
		require.ErrorIs(t, err, io.EOF, "nack should return a send error")

		stream.OnSend = func(in *api.PublisherReply) error { return nil }

		// Handle ack
		err = handler.Reply(ack)
		require.NoError(t, err, "could not send ack")

		// Handle nack
		err = handler.Reply(nack)
		require.NoError(t, err, "could not send ack")
		require.Equal(t, 3, stream.Calls(mock.StreamSend))
	})

	t.Run("Publisher", func(t *testing.T) {
		defer store.Reset()
		defer stream.Reset()

		// Set up claims and peer on the stream context.
		claims := &tokens.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "01H7G592TMWS2NF995V188NZE8",
			},
			OrgID:       "01H784KNY3GN2GC8NHW4ZKC5A9",
			ProjectID:   "01H784KXZPKMDWRX2ZRP6FSXET",
			Permissions: []string{permissions.Publisher},
		}
		stream.WithPeer(claims, MakePeer("172.153.2.2:10721"))

		// Authorize must be called first otherwise this method panics
		require.Panics(t, func() { handler.Publisher() }, "expected panic when authorize was not called")

		_, err := handler.Authorize(permissions.Publisher)
		require.NoError(t, err, "expected claims to be valid")

		publisher := handler.Publisher()
		require.Equal(t, "01H7G592TMWS2NF995V188NZE8", publisher.PublisherId)
		require.Equal(t, "172.153.2.2:10721", publisher.Ipaddr)
		require.Empty(t, publisher.ClientId)
		require.Equal(t, "test-agent", publisher.UserAgent)
	})

	t.Run("CloseStream", func(t *testing.T) {
		defer store.Reset()
		defer stream.Reset()

		// A new handler needs to be created since we're testing that acks/nacks are
		// recorded on the handler and other tests may increment the acks and nacks.
		handler := ensign.NewPublisherHandler(stream, store)

		// Handle error checking
		stream.WithError(mock.StreamSend, io.EOF)
		err := handler.CloseStream("tester", 12, 2)
		require.ErrorIs(t, err, io.EOF, "close stream should return a send error")

		// PublisherHandler should record acks and nacks being sent
		var closed *api.CloseStream
		stream.OnSend = func(in *api.PublisherReply) error {
			if msg := in.GetCloseStream(); msg != nil {
				closed = msg
			}
			return nil
		}

		// Send acks, nacks, and replies to make sure we're incrementing the counts
		handler.Ack(nil, time.Now())
		handler.Ack(nil, time.Now())
		handler.Reply(broker.PublishResult{Code: -1})
		handler.Nack(nil, api.Nack_CONSENSUS_FAILURE, "")
		handler.Reply(broker.PublishResult{Code: api.Nack_SHARDING_FAILURE})
		handler.CloseStream("tester", 12, 2)

		// Check that close stream has recorded the counts
		require.Equal(t, uint64(12), closed.Events)
		require.Equal(t, uint64(2), closed.Topics)
		require.Equal(t, uint64(3), closed.Acks)
		require.Equal(t, uint64(2), closed.Nacks)
		require.Equal(t, 7, stream.Calls(mock.StreamSend))
	})

}

func (s *serverTestSuite) TestSubscriberStreamInitialization() {
	require := s.Require()

	// These tests use a mock subscriber stream server rather than creating a client
	// stream through bufconn so that we don't directly use the interceptors and
	// directly test the stream handlers instead. It also prevents concurrency issues
	// with send and recv on the stream and the possibility of EOF errors and
	// intermittent failures.
	stream := &mock.SubscribeServer{}
	s.store.OnAllowedTopics = MockAllowedTopics
	s.store.OnTopicName = MockTopicName

	// Must be authenticated and have the subscriber permission
	err := s.srv.Subscribe(stream)
	s.GRPCErrorIs(err, codes.PermissionDenied, "not authorized to perform this action")

	// Create base claims to add to the stream context for authentication
	// These claims are valid but will have no topics associated with them.
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01H784KEP6F5EMW9CBYAHFB3J3",
		},
		OrgID:       "01H784KNY3GN2GC8NHW4ZKC5A9",
		ProjectID:   "01H784KXZPKMDWRX2ZRP6FSXET",
		Permissions: []string{permissions.Subscriber},
	}
	stream.WithClaims(claims)

	// When no topics are associated with the project, should get a no topics error
	err = s.srv.Subscribe(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "no topics available")

	// Change the claims over to a project ID that does contain topics and a peer
	claims.ProjectID = "01H6PGFTK2X53RGG2KMSGR2M61"
	stream.WithPeer(claims, MakePeer("172.92.121.6:10820"))

	// Handle stream closed before open stream message recv
	stream.WithError(mock.StreamRecv, io.EOF)
	require.NoError(s.srv.Subscribe(stream))

	// Handle stream crashed before open stream message recv
	stream.WithError(mock.StreamRecv, context.DeadlineExceeded)
	require.ErrorIs(s.srv.Subscribe(stream), context.DeadlineExceeded)

	// An OpenStream message must be the first message received by the handler
	stream.OnRecv = func() (*api.SubscribeRequest, error) {
		return &api.SubscribeRequest{
			Embed: &api.SubscribeRequest_Ack{
				Ack: &api.Ack{Id: []byte("foo")},
			},
		}, nil
	}

	err = s.srv.Subscribe(stream)
	s.GRPCErrorIs(err, codes.FailedPrecondition, "must send subscription to initialize stream")

	// TODO: happy path is timing out; need way to cancel subscribe stream.
	// subscription := &api.Subscription{ClientId: "tester", Topics: nil}
	// sub := stream.WithSubscription(subscription)

	// err = s.srv.Subscribe(stream)
	// require.NoError(err, "error happened?")

	// ready := sub.Ready()
	// require.NotNil(ready, "did not get a ready response from server")
	// sub.Close()

}

func TestStreamHandler(t *testing.T) {
	meta, err := store.Open(config.StorageConfig{ReadOnly: false, Testing: true})
	require.NoError(t, err, "could not open mock store for testing")

	stream := &mock.ServerStream{}
	handler := ensign.NewStreamHandler(ensign.UnknownStream, stream, meta)

	// Should not be able to get the ProjectID or AllowedTopics without authorization.
	_, err = handler.ProjectID()
	GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")

	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")

	// When there are no claims on the context, the handler should return unauthorized
	_, err = handler.Authorize("publisher")
	GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")

	// Add claims to the context for the remainder of the tests
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "http://localhost",
			Subject:   "01H6PGFB4T34D4WWEXQMAGJNMK",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		OrgID:       "01H6PGFG71N0AFEVTK3NJB71T9",
		ProjectID:   "01H6PGFTK2X53RGG2KMSGR2M61",
		Permissions: []string{"publisher", "subscriber"},
	}

	ctx := contexts.WithClaims(context.Background(), claims)
	stream.WithContext(ctx)

	// Should return unauthorized when the claims do not have the specific permission
	_, err = handler.Authorize("cookinthekitchen")
	GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")

	// When unauthorized, should not be able to get the ProjectID or AllowedTopics
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")

	// Should be able to authorize with valid permissions
	actualClaims, err := handler.Authorize("publisher")
	require.NoError(t, err)
	require.Equal(t, claims, actualClaims)

	// After authorization, should be able to get the ProjectID
	projectID, err := handler.ProjectID()
	require.NoError(t, err)
	require.Equal(t, ulid.MustParse("01H6PGFTK2X53RGG2KMSGR2M61"), projectID)

	// When no topics are available, should get an error
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.FailedPrecondition, "no topics available")

	// Internal error should be returned if topics cannot be fetched.
	meta.UseError(store.AllowedTopics, errors.New("this is a testing error"))
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Internal, "could not open unknown stream")

	meta.OnAllowedTopics = MockAllowedTopics
	meta.UseError(store.TopicName, errors.New("this is a testing error"))
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Internal, "could not open unknown stream")

	// Should be able to fetch topics
	meta.OnTopicName = MockTopicName
	group, err := handler.AllowedTopics()
	require.NoError(t, err)
	require.Equal(t, 4, group.Length())
}

func TestStreamHandlerInvalidProjectID(t *testing.T) {
	stream := &mock.ServerStream{}
	handler := ensign.NewStreamHandler(ensign.UnknownStream, stream, nil)

	testCases := []string{
		"", "notavalidulid", "00000000000000000000000000",
	}

	for _, tc := range testCases {
		claims := &tokens.Claims{ProjectID: tc, Permissions: []string{"publisher"}}
		stream.WithContext(contexts.WithClaims(context.Background(), claims))
		_, err := handler.Authorize("publisher")
		require.NoError(t, err)

		// Empty ProjectID not allowed
		projectID, err := handler.ProjectID()
		GRPCErrorIs(t, err, codes.PermissionDenied, "not authorized to perform this action")
		require.True(t, ulids.IsZero(projectID))
	}
}

func (s *serverTestSuite) setupValidPublisher() *mock.PublisherServer {
	stream := &mock.PublisherServer{}
	s.store.OnAllowedTopics = MockAllowedTopics
	s.store.OnTopicName = MockTopicName

	// Create base claims to add to the stream context for authentication
	// These claims are valid but will have no topics associated with them.
	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01H784KEP6F5EMW9CBYAHFB3J3",
		},
		OrgID:       "01H784KNY3GN2GC8NHW4ZKC5A9",
		ProjectID:   "01H6PGFTK2X53RGG2KMSGR2M61",
		Permissions: []string{permissions.Publisher},
	}
	stream.WithPeer(claims, MakePeer("172.92.121.6:10820"))
	return stream
}

// Maps ProjectID to a map of allowed topics and their names.
var projectTopics = map[string]map[string]string{
	"01H784KXZPKMDWRX2ZRP6FSXET": {},
	"01H6PGFTK2X53RGG2KMSGR2M61": {
		"01H6XTAPN0HZ1S7KEPFBF1MMPX": "example-topic-1",
		"01H6XTAVNM21F6JXNGAJF1SJ4S": "example-topic-2",
		"01H6XTB1780D2YKMC2MBNZ4V2X": "example-topic-3",
		"01H6XTB5DS8YG0YZEVQ385QRTB": "example-topic-4",
	},
}

func MockAllowedTopics(projectID ulid.ULID) ([]ulid.ULID, error) {
	if ulids.IsZero(projectID) {
		return nil, errors.New("cannot get topics for empty ulid")
	}

	topicMap, ok := projectTopics[projectID.String()]
	if !ok {
		return nil, errors.New("not found")
	}

	topics := make([]ulid.ULID, 0, len(topicMap))
	for key := range topicMap {
		topics = append(topics, ulid.MustParse(key))
	}

	return topics, nil
}

func MockTopicName(topicID ulid.ULID) (string, error) {
	tids := topicID.String()
	for _, tmap := range projectTopics {
		if name, ok := tmap[tids]; ok {
			return name, nil
		}
	}
	return "", errors.New("unknown topic id")

}

func MakePeer(ipaddr string) *peer.Peer {
	return &peer.Peer{
		Addr:     net.TCPAddrFromAddrPort(netip.MustParseAddrPort(ipaddr)),
		AuthInfo: nil,
	}
}

func MakeEmpty(topicID string) *api.EventWrapper {
	event := &api.Event{
		Data:     []byte{},
		Metadata: make(map[string]string),
		Mimetype: mimetype.ApplicationOctetStream,
		Type:     &api.Type{Name: "Empty", MajorVersion: 1},
		Created:  timestamppb.Now(),
	}
	return MakeEvent(topicID, event)
}

func MakeEvent(topicID string, event *api.Event) *api.EventWrapper {
	env := &api.EventWrapper{
		TopicId: ulid.MustParse(topicID).Bytes(),
		LocalId: ulid.Make().Bytes(),
	}

	if err := env.Wrap(event); err != nil {
		panic(err)
	}
	return env
}
