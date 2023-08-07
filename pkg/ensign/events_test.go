package ensign_test

import (
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
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/mock"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

func (s *serverTestSuite) TestPublisherStreamInitialization() {
	require := s.Require()

	// These tests use a mock publisher server rather than creating a client stream
	// through bufconn so that we don't directly use the interceptors and directly test
	// the stream handlers instead. It also prevents concurrency issues with send and
	// recv on a stream and the possibility of EOF errors and intermittent failures.
	stream := &mock.PublisherServer{}
	s.store.OnAllowedTopics = MockAllowedTopics
	s.store.OnTopicName = MockTopicName

	// Must be authenticated and have the publisher permission.
	err := s.srv.Publish(stream)
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

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
	require.ErrorIs(s.srv.Publish(stream), context.DeadlineExceeded)

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

	require.Equal(uint64(0), closed.Consumers)
	require.Equal(uint64(0), closed.Topics)
	require.Equal(uint64(0), closed.Events)

	// TODO: test cannot send a second open stream message after first open stream message
}

func TestStreamHandler(t *testing.T) {
	meta, err := store.Open(config.StorageConfig{ReadOnly: false, Testing: true})
	require.NoError(t, err, "could not open mock store for testing")

	stream := &mock.ServerStream{}
	handler := ensign.NewStreamHandler(ensign.UnknownStream, stream, meta)

	// Should not be able to get the ProjectID or AllowedTopics without authorization.
	_, err = handler.ProjectID()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// When there are no claims on the context, the handler should return unauthorized
	_, err = handler.Authorize("publisher")
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

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
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

	// When unauthorized, should not be able to get the ProjectID or AllowedTopics
	_, err = handler.AllowedTopics()
	GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")

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
	s.GRPCErrorIs(err, codes.Unauthenticated, "not authorized to perform this action")

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
		GRPCErrorIs(t, err, codes.Unauthenticated, "not authorized to perform this action")
		require.True(t, ulids.IsZero(projectID))
	}
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
