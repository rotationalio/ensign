package ensign

import (
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
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
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Cannot publish events > 5MiB long
const EventMaxDataSize int = 5.243e6

// Publish implements the streaming endpoint for the API.
func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	o11y.OnlinePublishers.Inc()
	defer o11y.OnlinePublishers.Dec()

	// Create a Publish handler for the stream
	handler := &PublisherHandler{
		StreamHandler: StreamHandler{
			stype:  "publisher",
			stream: stream,
			meta:   s.meta,
		},
		stream: stream,
	}

	// Authorize the user to ensure they are allowed to publish.
	var claims *tokens.Claims
	if claims, err = handler.Authorize(permissions.Publisher); err != nil {
		return err
	}

	// Get the allowed topics based on the claims
	var allowedTopics *topics.NameGroup
	if allowedTopics, err = handler.AllowedTopics(); err != nil {
		return err
	}

	// Publisher information
	ctx := stream.Context()
	publisher := &api.Publisher{
		PublisherId: claims.Subject,
	}
	if remote, ok := peer.FromContext(ctx); ok {
		publisher.Ipaddr = remote.Addr.String()
	}

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

	// TODO: verify topics sent in open stream message
	publisher.ClientId = open.ClientId

	// Send back topic mapping
	ready := &api.StreamReady{
		ClientId: publisher.ClientId,
		ServerId: s.conf.Monitoring.NodeID,
		Topics:   allowedTopics.TopicMap(),
	}

	if err = stream.Send(&api.PublisherReply{Embed: &api.PublisherReply_Ready{Ready: ready}}); err != nil {
		if streamClosed(err) {
			log.Info().Msg("publish stream closed")
			return nil
		}
		sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
		return err
	}

	// Set up the stream handlers
	nEvents := uint64(0)
	events := make(chan *api.EventWrapper, BufferSize)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the ack-back loop
	// This loop also pushes the event onto the primary buffer
	go func(events <-chan *api.EventWrapper) {
		defer wg.Done()
		for event := range events {
			// Verify the event has a topic associated with it
			if len(event.TopicId) == 0 {
				log.Warn().Msg("event published without topic id")

				// Send the nack back to the user
				// TODO: handle the error if any
				handler.Nack(event.Id, api.Nack_TOPIC_UKNOWN, "")
				continue
			}

			// Verify the event is in a topic that the user is allowed to publish to
			// TODO: this won't allow topics that were created after the stream was
			// created but are still valid. Need to unify the allowed mechanism into
			// a global topic handler check rather than in a per-stream check.
			var topicID ulid.ULID
			if topicID, err = event.ParseTopicID(); err != nil {
				sentry.Debug(ctx).Err(err).Msg("could not parse topic id from user")

				// Send the nack back to the user
				// TODO: handle the error if any
				handler.Nack(event.Id, api.Nack_TOPIC_UKNOWN, "")
				continue
			}

			if ok := allowedTopics.ContainsTopicID(topicID); !ok {
				sentry.Warn(ctx).Msg("event published to topic that is not allowed")

				// Send the nack back to the user
				// TODO: handle the error if any
				handler.Nack(event.Id, api.Nack_TOPIC_UKNOWN, "")
				continue
			}

			// Check the maximum event size to prevent large events from being published.
			if len(event.Event) > EventMaxDataSize {
				sentry.Warn(ctx).Int("size", len(event.Event)).Msg("very large event published to topic and rejected")

				// Send the nack back to the user
				// TODO: handle the error if any
				handler.Nack(event.Id, api.Nack_MAX_EVENT_SIZE_EXCEEDED, "")
				continue
			}

			// Push event on to the primary buffer
			event.Publisher = publisher
			s.pubsub.Publish(event)

			// Send ack once the event is on the primary buffer
			err = handler.Ack(event.Id, time.Now())

			if err == nil {
				nEvents++
			} else {
				log.Warn().Err(err).Msg("could not send ack")
			}
		}
	}(events)

	// Receive events from the clients
	go func(events chan<- *api.EventWrapper) {
		defer wg.Done()
		defer close(events)
		for {
			select {
			case <-ctx.Done():
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
				sentry.Warn(ctx).Err(err).Msg("publish stream crashed")
				return
			}

			// Handle the different types of messages the publisher will send
			switch msg := in.Embed.(type) {
			case *api.PublisherRequest_Event:
				events <- msg.Event
			case *api.PublisherRequest_OpenStream:
				// TODO: verify topics that are sent in the open stream message
				// TODO: this should be more general in the recv loop
				publisher.ClientId = msg.OpenStream.ClientId
			default:
				// TODO: how do we send errors from here?
				err = status.Errorf(codes.FailedPrecondition, "unhandled publisher request message %T", msg)
				sentry.Warn(ctx).Err(err).Msg("could not handle publisher request")
			}
		}
	}(events)

	wg.Wait()
	stream.Send(&api.PublisherReply{
		Embed: &api.PublisherReply_CloseStream{
			CloseStream: &api.CloseStream{
				Events: nEvents,
			},
		},
	})
	return err
}

type PublisherHandler struct {
	StreamHandler
	stream api.Ensign_PublishServer
}

func (p PublisherHandler) Ack(eventID []byte, committed time.Time) error {
	return p.stream.Send(&api.PublisherReply{
		Embed: &api.PublisherReply_Ack{
			Ack: &api.Ack{
				Id:        eventID,
				Committed: timestamppb.New(committed),
			},
		},
	})
}

func (p PublisherHandler) Nack(eventID []byte, code api.Nack_Code, err string) error {
	return p.stream.Send(&api.PublisherReply{
		// TODO: what are the nack error codes?
		Embed: &api.PublisherReply_Nack{
			Nack: &api.Nack{
				Id:    eventID,
				Code:  code,
				Error: err,
			},
		},
	})
}

// Subscribe implements the streaming endpoint for the API
func (s *Server) Subscribe(stream api.Ensign_SubscribeServer) (err error) {
	o11y.OnlineSubscribers.Inc()
	defer o11y.OnlineSubscribers.Dec()

	// Create a Subscriber handler for the stream
	handler := &SubscriberHandler{
		StreamHandler: StreamHandler{
			stype:  "subscriber",
			stream: stream,
			meta:   s.meta,
		},
		stream: stream,
	}

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

	// HACK: this is super messy, clean up!
	// Handle the subscription stream initialization
	// TODO: handle the consumer group
	if len(sub.Topics) > 0 {
		// Filter the allowedTopics channel
		// TODO: add some thread-safety here
		topicIDs := make([]ulid.ULID, 0, len(sub.Topics))
		for _, topic := range sub.Topics {
			// TODO: don't just ignore unparsable topics
			if tid, err := ulids.Parse(topic); err == nil && !ulids.IsZero(tid) {
				topicIDs = append(topicIDs, tid)
			}
		}

		allowedTopics = allowedTopics.FilterTopicID(topicIDs...)
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
	id, events := s.pubsub.Subscribe()
	defer s.pubsub.Finish(id)

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

			if ack := in.GetAck(); ack != nil {
				acks++
			} else if nack := in.GetNack(); nack != nil {
				nacks++
			}

			// Set up the topic filter if one has come in
			// TODO: make this a prerequisite
			if sub := in.GetSubscription(); sub != nil {
				// Filter the allowedTopics channel
				topicIDs := make([]ulid.ULID, 0, len(sub.Topics))
				for _, topic := range sub.Topics {
					// TODO: don't just ignore unparsable topics
					if tid, err := ulids.Parse(topic); err == nil && !ulids.IsZero(tid) {
						topicIDs = append(topicIDs, tid)
					}
				}

				allowedTopics = allowedTopics.FilterTopicID(topicIDs...)
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
	stype     string
	stream    grpc.ServerStream
	meta      store.MetaStore
	claims    *tokens.Claims
	projectID ulid.ULID
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
		group.Add(name, topicID)
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
