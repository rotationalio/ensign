package ensign

import (
	"io"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/broker"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/topics"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Cannot publish events > 5MiB long
const EventMaxDataSize int = 5.243e6

// Publish implements a streaming endpoint that allows users to publish events into a
// topic or topics that are managed by the current broker.
//
// The Publish stream has two phases and operates three go routines. The first phase is
// the initialization phase where the stream is opened and authorized, then the topics
// are loaded from the database. The handler waits for an OpenStream message to be
// received from the client before proceeding, and sends either an error or a
// StreamReady message if successful. Clients should ensure that they send an OpenStream
// message and then recv StreamReady before publishing any events to a topic.
//
// The second phase initializes two go routines: the primary go routine receives events
// from the client and extracts them. It then sends them via a channel to a second go
// routine that handles pre-broker processing and any acks/nacks received from the
// broker. The handler go routine waits for these routines to complete before returning
// any error from the recv routine.
//
// Permissions: publisher
func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	o11y.OnlinePublishers.Inc()
	defer o11y.OnlinePublishers.Dec()

	// Create a Publish handler for the stream
	ctx := stream.Context()
	handler := NewPublisherHandler(stream, s.meta)

	// Authorize the user to ensure they are allowed to publish.
	if _, err = handler.Authorize(permissions.Publisher); err != nil {
		// NOTE: Authorize() returns a status error that can be returned directly.
		return err
	}

	// Get the allowed topics based on the claims
	var allowedTopics *topics.NameGroup
	if allowedTopics, err = handler.AllowedTopics(); err != nil {
		// NOTE: AllowedTopics() returns a status error that can be returned directly.
		return err
	}

	// Get publisher information from the stream
	// Sets the publisherID as the API Key ID from the claims subject
	publisher := handler.Publisher()

	// Recv the OpenStream message from the client
	var in *api.PublisherRequest
	if in, err = stream.Recv(); err != nil {
		if streamClosed(err) {
			log.Info().Msg("publish stream closed")
			return nil
		}
		sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
		return err
	}

	var open *api.OpenStream
	if open = in.GetOpenStream(); open == nil {
		return status.Error(codes.FailedPrecondition, "an open stream message must be sent immediately after opening the stream")
	}

	// Verify topics sent in open stream message and filter allowed topics.
	if len(open.Topics) > 0 {
		// Only the topics specified by the publisher will be allowed to be published to.
		allowedTopics = allowedTopics.Filter(open.Topics...)
		if allowedTopics.Length() == 0 {
			log.Warn().Int("topic_filter", len(open.Topics)).Msg("publish stream opened with no topics")
			return status.Error(codes.FailedPrecondition, "no topics available")
		}
	}

	// Send back topic mapping and stream ready notification.
	publisher.ClientId = open.ClientId
	ready := &api.StreamReady{
		ClientId: publisher.ResolveClientID(),
		ServerId: s.conf.Monitoring.NodeID,
		Topics:   allowedTopics.TopicMap(),
	}

	if err = stream.Send(&api.PublisherReply{Embed: &api.PublisherReply_Ready{Ready: ready}}); err != nil {
		if streamClosed(err) {
			log.Info().Msg("publish stream closed")
			return nil
		}
		sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
		return status.Error(codes.Aborted, "could not initialize publish stream")
	}

	// Set up the stream handlers
	streamID, results, err := s.broker.Register()
	if err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not register publisher with broker")
		return status.Error(codes.Unavailable, "ensign broker is not available")
	}
	defer s.broker.Close(streamID)

	var nEvents uint64
	publishedTo := make(map[ulid.ULID]struct{})
	events := make(chan *api.EventWrapper, BufferSize)

	var wg sync.WaitGroup
	wg.Add(2)

	// Receive events from the clients
	// This is the primary routine for the publisher since we want to ensure that we
	// receive all events that the publisher publishes. If an error occurs or the stream
	// closes during this routine, then we signal all other go routines to stop.
	// This go routine sets the external err so that it can be sent back to the user if
	// something goes wrong, no errors are sent back from send errors.
	// NOTE: this go routine cannot send messages since it calls recv!
	go func(events chan<- *api.EventWrapper) {
		defer wg.Done()
		defer close(events)
		for {
			select {
			case <-ctx.Done():
				// NOTE: If the context errors (e.g. deadline exceeded) then no error is
				// returned, which is why err := in the next line to prevent shadowing.
				if err := ctx.Err(); err != nil {
					log.Debug().Err(err).Msg("context closed")
				}
				return
			default:
			}

			var in *api.PublisherRequest
			if in, err = stream.Recv(); err != nil {
				if streamClosed(err) {
					log.Info().Msg("publish stream closed")
					err = nil
					return
				}

				// Set the error message to aborted if we cannot recv a message
				err = status.Error(codes.Aborted, "could not recv event from client")
				sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
				return
			}

			// Handle the different types of messages the publisher will send
			switch msg := in.Embed.(type) {
			case *api.PublisherRequest_Event:
				events <- msg.Event
			case *api.PublisherRequest_OpenStream:
				// We have already processed an open stream request and cannot accept
				// additional open stream messages, error closing the stream.
				err = status.Errorf(codes.Aborted, "cannot send multiple open stream messages")
				sentry.Warn(ctx).Err(err).Msg("multiple open stream messages received")
				return
			default:
				// If an unknown message type comes in, return an error to the user.
				err = status.Errorf(codes.FailedPrecondition, "unhandled publisher request message %T", msg)
				sentry.Warn(ctx).Err(err).Msg("could not handle publisher request")
				return
			}
		}
	}(events)

	// Handle events and send acks/nacks back to the client.
	// This go routine handles the events sent from the client, including any pre-broker
	// preprocessing such as checking the topic or the event size. The pre-broker checks
	// happen here to ensure that nacks can be sent back to the user. This go routine
	// also listens for acks/nacks from the broker and returns them to the user as well.
	// This routine should only log errors, not return from them, and should stop when
	// the events receiving go routine is concluded.
	// NOTE: this go routine cannot recv messages since it calls send!
	go func(events <-chan *api.EventWrapper, results <-chan broker.PublishResult) {
		defer wg.Done()
		for {
			select {
			// Handle events coming from the client.
			case event, open := <-events:
				// If the events channel has closed, that is the signal to stop this routine
				if !open {
					return
				}

				// Verify the event has a topic associated with it
				if len(event.TopicId) == 0 {
					log.Warn().Msg("event published without topic id")

					// Send the nack back to the user and log send error if any
					handler.Nack(event.LocalId, api.Nack_TOPIC_UNKNOWN, "no topic id specified")
					continue
				}

				// Verify the event is in a topic that the user is allowed to publish to
				// TODO: this won't allow topics that were created after the stream was
				// created but are still valid. Need to unify the allowed mechanism into
				// a global topic handler check rather than in a per-stream check.
				var topicID ulid.ULID
				if topicID, err = event.ParseTopicID(); err != nil {
					sentry.Debug(ctx).Err(err).Msg("could not parse topic id from user")

					// Send the nack back to the user and log send error if any
					handler.Nack(event.LocalId, api.Nack_TOPIC_UNKNOWN, "invalid topic id")
					continue
				}

				if ok := allowedTopics.ContainsTopicID(topicID); !ok {
					sentry.Warn(ctx).Msg("event published to topic that is not allowed")

					// Send the nack back to the user and log send error if any
					handler.Nack(event.LocalId, api.Nack_TOPIC_UNKNOWN, "")
					continue
				}

				// Check the maximum event size to prevent large events from being published.
				if len(event.Event) > EventMaxDataSize {
					sentry.Warn(ctx).Int("size", len(event.Event)).Msg("very large event published to topic and rejected")

					// Send the nack back to the user and log send error if any
					handler.Nack(event.LocalId, api.Nack_MAX_EVENT_SIZE_EXCEEDED, "")
					continue
				}

				// Push event on to the primary buffer
				event.Publisher = publisher
				s.broker.Publish(streamID, event)

				// Increment counters for sending back closed stream message
				nEvents++
				if _, ok := publishedTo[topicID]; !ok {
					publishedTo[topicID] = struct{}{}
				}

			// Handle acks/nacks coming from the broker
			case result := <-results:
				handler.Reply(result)
			}
		}
	}(events, results)

	// Wait for client to close the event stream and to handle remaining events.
	wg.Wait()

	// Exhaust results stream to ensure all results are sent back
	s.broker.Close(streamID)
	for result := range results {
		handler.Reply(result)
	}

	// Send close stream message and log that the stream has been closed
	handler.CloseStream(nEvents, uint64(len(publishedTo)))
	return err
}

type PublisherHandler struct {
	StreamHandler
	stream api.Ensign_PublishServer
	nAcks  uint64
	nNacks uint64
}

func NewPublisherHandler(stream api.Ensign_PublishServer, meta store.MetaStore) *PublisherHandler {
	return &PublisherHandler{
		StreamHandler: StreamHandler{
			stype:  PublisherStream,
			stream: stream,
			meta:   meta,
		},
		stream: stream,
	}
}

// Sends an ack back to the client, logging any send errors that occur.
func (p *PublisherHandler) Ack(eventID []byte, committed time.Time) error {
	err := p.stream.Send(&api.PublisherReply{
		Embed: &api.PublisherReply_Ack{
			Ack: &api.Ack{
				Id:        eventID,
				Committed: timestamppb.New(committed),
			},
		},
	})

	if err != nil {
		log.Warn().Err(err).Bytes("localID", eventID).Msg("could not send ack")
	}

	p.nAcks++
	return err
}

// Sends a nack back to the client, logging any send errors that occur.
func (p *PublisherHandler) Nack(eventID []byte, code api.Nack_Code, message string) error {
	// Use default nack message if none is specified.
	if message == "" {
		message = api.DefaultNackMessage(code)
	}

	err := p.stream.Send(&api.PublisherReply{
		Embed: &api.PublisherReply_Nack{
			Nack: &api.Nack{
				Id:    eventID,
				Code:  code,
				Error: message,
			},
		},
	})

	if err != nil {
		log.Warn().Err(err).Bytes("localID", eventID).Msg("could not send nack")
	}

	p.nNacks++
	return err
}

// Sends close stream message and logs stream closed along with any send errors.
func (p PublisherHandler) CloseStream(events, topics uint64) error {
	err := p.stream.Send(&api.PublisherReply{
		Embed: &api.PublisherReply_CloseStream{
			CloseStream: &api.CloseStream{
				Events: events,
				Topics: topics,
				Acks:   p.nAcks,
				Nacks:  p.nNacks,
			},
		},
	})

	ctx := log.With().Uint64("events", events).Uint64("topics", topics).
		Uint64("acks", p.nAcks).Uint64("nacks", p.nNacks).
		Logger()

	if err != nil {
		ctx.Warn().Err(err).Msg("could not send close stream message")
	} else {
		ctx.Info().Msg("publisher stream closed")
	}
	return err
}

// Handles the publisher reply from the broker
func (p *PublisherHandler) Reply(msg broker.PublishResult) error {
	if msg.IsNack() {
		p.nNacks++
	} else {
		p.nAcks++
	}

	if err := p.stream.Send(msg.Reply()); err != nil {
		log.Warn().Err(err).
			Bool("is_ack", msg.IsAck()).Bool("is_nack", msg.IsNack()).
			Bytes("localID", msg.LocalID).
			Msg("could not send broker reply")
		return err
	}
	return nil
}

// Publisher gathers publisher info from the claims and from the request.
// NOTE: authorize must be called before this method can be called.
func (p PublisherHandler) Publisher() *api.Publisher {
	ctx := p.stream.Context()
	publisher := &api.Publisher{PublisherId: p.claims.Subject}
	if remote, ok := peer.FromContext(ctx); ok {
		publisher.Ipaddr = remote.Addr.String()
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		publisher.UserAgent = strings.Join(md.Get("user-agent"), ",")
	}

	return publisher
}

// Subscribe implements the streaming endpoint for the API
//
// Permissions: subscriber
func (s *Server) Subscribe(stream api.Ensign_SubscribeServer) (err error) {
	o11y.OnlineSubscribers.Inc()
	defer o11y.OnlineSubscribers.Dec()

	// Create a Subscriber handler for the stream
	handler := NewSubscribeHandler(stream, s.meta)

	// Parse the context for authentication information
	if _, err = handler.Authorize(permissions.Subscriber); err != nil {
		return err
	}

	// Get the allowed topics based on the claims
	var allowedTopics *topics.NameGroup
	if allowedTopics, err = handler.AllowedTopics(); err != nil {
		return err
	}

	// Recv the subscription message from the client to initialize the stream
	ctx := stream.Context()
	var in *api.SubscribeRequest
	if in, err = stream.Recv(); err != nil {
		if streamClosed(err) {
			log.Info().Msg("publish stream closed")
			return nil
		}
		sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
		return err
	}

	var sub *api.Subscription
	if sub = in.GetSubscription(); sub == nil {
		return status.Error(codes.FailedPrecondition, "must send subscription to initialize stream")
	}

	// Handle the subscription stream initialization
	// TODO: handle the consumer group
	if len(sub.Topics) > 0 {
		allowedTopics = allowedTopics.Filter(sub.Topics...)
		if allowedTopics.Length() == 0 {
			log.Warn().Int("topic_filter", len(sub.Topics)).Msg("subscribe stream opened with no topics")
			return status.Error(codes.FailedPrecondition, "no topics available")
		}
	}

	// Send back topic mapping
	ready := &api.StreamReady{
		ClientId: sub.ClientId,
		ServerId: s.conf.Monitoring.NodeID,
		Topics:   allowedTopics.TopicMap(),
	}

	if err = stream.Send(&api.SubscribeReply{Embed: &api.SubscribeReply_Ready{Ready: ready}}); err != nil {
		if streamClosed(err) {
			log.Info().Msg("subscribe stream closed")
			return nil
		}
		sentry.Warn(ctx).Err(err).Msg("subscribe stream crashed")
		return err
	}

	// Setup the stream handlers
	nEvents, acks, nacks := uint64(0), uint64(0), uint64(0)
	streamID, events, err := s.broker.Subscribe(allowedTopics.TopicIDs()...)
	defer s.broker.Close(streamID)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the event sending loop
	go func(events <-chan *api.EventWrapper) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Debug().Err(err).Msg("context closed in subscribe event routine")
					return
				}
			case event := <-events:

				var topicID ulid.ULID
				if topicID, err = event.ParseTopicID(); err != nil {
					sentry.Warn(ctx).Err(err).Bytes("topicID", event.TopicId).Bytes("event", event.Id).Msg("could not parse topic id on event in log")
					continue
				}

				// Filter events based on the topic ID
				if ok := allowedTopics.ContainsTopicID(topicID); !ok {
					continue
				}

				if err = handler.Send(event); err != nil {
					if streamClosed(err) {
						log.Info().Msg("subscribe stream closed")
						err = nil
						return
					}
					sentry.Warn(ctx).Err(err).Msg("subscribe stream crashed")
					return
				}
				nEvents++
			}
		}
	}(events)

	// Receive acks from the clients
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Debug().Err(err).Msg("context closed in subscribe ack routine")
					return
				}
			default:
			}

			var in *api.SubscribeRequest
			if in, err = stream.Recv(); err != nil {
				if streamClosed(err) {
					log.Info().Msg("subscribe stream closed")
					err = nil
					return
				}
				sentry.Warn(ctx).Err(err).Msg("subscribe stream crashed")
				return
			}

			// TODO: handle acks and nacks
			if ack := in.GetAck(); ack != nil {
				acks++
			} else if nack := in.GetNack(); nack != nil {
				nacks++
			}
		}
	}()
	wg.Wait()
	log.Info().Uint64("nEvents", nEvents).Uint64("acks", acks).Uint64("nacks", nacks).Msg("subscribe stream terminated")
	return err
}

type SubscriberHandler struct {
	StreamHandler
	stream api.Ensign_SubscribeServer
}

func NewSubscribeHandler(stream api.Ensign_SubscribeServer, meta store.MetaStore) *SubscriberHandler {
	return &SubscriberHandler{
		StreamHandler: StreamHandler{
			stype:  SubscriberStream,
			stream: stream,
			meta:   meta,
		},
		stream: stream,
	}
}

func (s SubscriberHandler) Send(event *api.EventWrapper) error {
	return s.stream.Send(&api.SubscribeReply{
		Embed: &api.SubscribeReply_Event{
			Event: event,
		},
	})
}

// StreamHandler provides some common functionality to both the Publisher and Subscriber
// stream handlers, for example providing authentication and collecting allowed topics.
type StreamHandler struct {
	stype     StreamType
	stream    grpc.ServerStream
	meta      store.MetaStore
	claims    *tokens.Claims
	projectID ulid.ULID
}

// Creates a new StreamHandler -- primarily used for testing.
func NewStreamHandler(stype StreamType, stream grpc.ServerStream, meta store.MetaStore) *StreamHandler {
	return &StreamHandler{
		stype:  stype,
		stream: stream,
		meta:   meta,
	}
}

// Authorize fetches the claims from the stream context, returning an error if the user
// is not authorized. The claims are cached on the stream handler and returned without
// error after the first time they are correctly fetched. Fetching claims requires a
// permission (e.g. subscribe or publish). If not authorized a status error is returned.
// Authorize MUST be called before projectID or projectTopics is called.
func (s *StreamHandler) Authorize(permission string) (_ *tokens.Claims, err error) {
	if s.claims == nil {
		ctx := s.stream.Context()
		if s.claims, err = contexts.Authorize(ctx, permission); err != nil {
			sentry.Warn(ctx).Err(err).Str("permission", permission).Msgf("unauthorized %s stream", s.stype)
			return nil, status.Error(codes.Unauthenticated, "not authorized to perform this action")
		}
	}
	return s.claims, nil
}

// Returns the ProjectID associated with the claims. Authorize must be called first or
// this method will error. Returns a status error if no project ID is on the claims.
// The projectID is cached on the stream handler and will be returned without error.
func (s *StreamHandler) ProjectID() (ulid.ULID, error) {
	if ulids.IsZero(s.projectID) {
		if s.claims == nil {
			sentry.Error(s.stream.Context()).Msg("project ID fetched without authorization")
			return ulids.Null, status.Error(codes.Unauthenticated, "not authorized to perform this action")
		}

		if s.projectID = s.claims.ParseProjectID(); ulids.IsZero(s.projectID) {
			sentry.Warn(s.stream.Context()).Str("project_id", s.claims.ProjectID).Msg("could not parse projectID from claims")
			return ulids.Null, status.Error(codes.Unauthenticated, "not authorized to perform this action")
		}
	}
	return s.projectID, nil
}

// AllowedTopics returns a set of topic IDs and hashed topic names that are allowed to
// be accessed by the given claims. This set can be filtered to further restrict the
// stream based on user input. A specialized data structure is used to make it easy to
// perform content filtering based on name and ID.
func (s *StreamHandler) AllowedTopics() (group *topics.NameGroup, err error) {
	var projectID ulid.ULID
	if projectID, err = s.ProjectID(); err != nil {
		return nil, err
	}

	var projectTopics []ulid.ULID
	if projectTopics, err = s.meta.AllowedTopics(projectID); err != nil {
		sentry.Error(s.stream.Context()).Err(err).Msg("could not fetch project topics")
		return nil, status.Errorf(codes.Internal, "could not open %s stream", s.stype)
	}

	group = &topics.NameGroup{}
	for _, topicID := range projectTopics {
		var name string
		if name, err = s.meta.TopicName(topicID); err != nil {
			sentry.Error(s.stream.Context()).Err(err).Str("topicID", topicID.String()).Msg("could not get topic name from ID")
			return nil, status.Errorf(codes.Internal, "could not open %s stream", s.stype)
		}

		if err := group.Add(name, topicID); err != nil {
			sentry.Error(s.stream.Context()).Err(err).Str("topicID", topicID.String()).Str("topic", name).Msg("could not add topic name and ID to group")
			return nil, status.Errorf(codes.Internal, "could not open %s stream", s.stype)
		}
	}

	if group.Length() == 0 {
		log.Warn().Msgf("%s stream opened with no topics", s.stype)
		return nil, status.Error(codes.FailedPrecondition, "no topics available")
	}

	return group, nil
}

func streamClosed(err error) bool {
	if err == io.EOF {
		return true
	}

	if serr, ok := status.FromError(err); ok {
		if serr.Code() == codes.Canceled {
			return true
		}
	}

	return false
}

type StreamType uint8

const (
	UnknownStream StreamType = iota
	PublisherStream
	SubscriberStream
)

const (
	unknownStream    = "unknown"
	publisherStream  = "publisher"
	subscriberStream = "subscriber"
)

func (s StreamType) String() string {
	switch s {
	case PublisherStream:
		return publisherStream
	case SubscriberStream:
		return subscriberStream
	default:
		return unknownStream
	}
}
