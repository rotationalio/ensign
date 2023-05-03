package ensign

import (
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Cannot publish events > 5MiB long
const EventMaxDataSize int = 5.243e6

// Publish implements the streaming endpoint for the API.
func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	o11y.OnlinePublishers.Inc()
	defer o11y.OnlinePublishers.Dec()

	// Parse the context for authentication information
	var claims *tokens.Claims
	ctx := stream.Context()
	if claims, err = contexts.Authorize(ctx, permissions.Publisher); err != nil {
		sentry.Warn(ctx).Err(err).Msg("unauthorized publisher")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Str("project_id", claims.ProjectID).Msg("could not parse projectID from claims")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Get the topic IDs this user is allowed to publish to
	var projectTopics []ulid.ULID
	if projectTopics, err = s.meta.AllowedTopics(projectID); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not fetch project topics")
		return status.Error(codes.Internal, "could not open publisher stream")
	}

	allowedTopics := make(map[string]struct{})
	for _, topicID := range projectTopics {
		allowedTopics[topicID.String()] = struct{}{}
	}

	if len(allowedTopics) == 0 {
		log.Warn().Msg("publisher created with no topics")
		return status.Error(codes.FailedPrecondition, "no topics available")
	}

	// Set up the stream handlers
	nEvents := uint64(0)
	events := make(chan *api.Event, BufferSize)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the ack-back loop
	// This loop also pushes the event onto the primary buffer
	go func(events <-chan *api.Event) {
		defer wg.Done()
		for event := range events {
			// Verify the event has a topic associated with it
			if event.TopicId == "" {
				log.Warn().Msg("event published without topic id")

				// Send the nack back to the user
				stream.Send(&api.PublisherReply{
					// TODO: what are the nack error codes?
					Embed: &api.PublisherReply_Nack{
						Nack: &api.Nack{
							Id:    event.Id,
							Code:  uint32(12),
							Error: "event requires topic id",
						},
					},
				})

				// Continue processing
				continue
			}

			// Verify the event is in a topic that the user is allowed to publish to
			// TODO: this won't allow topics that were created after the stream was
			// created but are still valid. Need to unify the allowed mechanism into
			// a global topic handler check rather than in a per-stream check.
			if _, ok := allowedTopics[event.TopicId]; !ok {
				sentry.Warn(ctx).Msg("event published to topic that is not allowed")

				// Send the nack back to the user
				stream.Send(&api.PublisherReply{
					// TODO: what are the nack error codes?
					Embed: &api.PublisherReply_Nack{
						Nack: &api.Nack{
							Id:    event.Id,
							Code:  uint32(99),
							Error: "unknown topic id",
						},
					},
				})

				// Continue processing
				continue
			}

			// Check the maximum event size to prevent large events from being published.
			if len(event.Data) > EventMaxDataSize {
				sentry.Warn(ctx).Int("size", len(event.Data)).Msg("very large event published to topic and rejected")

				// Send the nack back to the user
				stream.Send(&api.PublisherReply{
					// TODO: what are the nack error codes?
					Embed: &api.PublisherReply_Nack{
						Nack: &api.Nack{
							Id:    event.Id,
							Code:  uint32(6),
							Error: "event data too large",
						},
					},
				})

				// Continue processing
				continue
			}

			// Push event on to the primary buffer
			s.pubsub.Publish(event)

			// Send ack once the event is on the primary buffer
			err = stream.Send(&api.PublisherReply{
				Embed: &api.PublisherReply_Ack{
					Ack: &api.Ack{
						Id:        event.Id,
						Committed: timestamppb.Now(),
					},
				},
			})

			if err == nil {
				nEvents++
			} else {
				log.Warn().Err(err).Msg("could not send ack")
			}
		}
	}(events)

	// Receive events from the clients
	go func(events chan<- *api.Event) {
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

			var event *api.Event
			if event = in.GetEvent(); event == nil {
				// TODO: handle control message
				continue
			}

			events <- event
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

// Subscribe implements the streaming endpoint for the API
func (s *Server) Subscribe(stream api.Ensign_SubscribeServer) (err error) {
	o11y.OnlineSubscribers.Inc()
	defer o11y.OnlineSubscribers.Dec()

	// Parse the context for authentication information
	var claims *tokens.Claims
	ctx := stream.Context()
	if claims, err = contexts.Authorize(ctx, permissions.Subscriber); err != nil {
		sentry.Warn(ctx).Err(err).Msg("unauthorized subscriber")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID = claims.ParseProjectID(); ulids.IsZero(projectID) {
		sentry.Warn(ctx).Str("project_id", claims.ProjectID).Msg("could not parse projectID from claims")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Get the topic IDs this user is allowed to subscribe to
	var projectTopics []ulid.ULID
	if projectTopics, err = s.meta.AllowedTopics(projectID); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not fetch project topics")
		return status.Error(codes.Internal, "could not open subscriber stream")
	}

	allowedTopics := make(map[string]struct{})
	for _, topicID := range projectTopics {
		allowedTopics[topicID.String()] = struct{}{}
	}

	if len(allowedTopics) == 0 {
		log.Warn().Msg("subscriber created with no topics")
		return status.Error(codes.FailedPrecondition, "no topics available")
	}

	// Setup the stream handlers
	nEvents, acks, nacks := uint64(0), uint64(0), uint64(0)
	id, events := s.pubsub.Subscribe()
	defer s.pubsub.Finish(id)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the event sending loop
	go func(events <-chan *api.Event) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Debug().Err(err).Msg("context closed in subscribe event routine")
					return
				}
			case event := <-events:
				// Filter events based on the topic ID
				if _, ok := allowedTopics[event.TopicId]; !ok {
					continue
				}

				if err = stream.Send(&api.SubscribeReply{Embed: &api.SubscribeReply_Event{Event: event}}); err != nil {
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
				// TODO: handle the consumer group
				if len(sub.Topics) > 0 {
					// Filter the allowedTopics channel
					// TODO: add some thread-safety here
					filter := make(map[string]struct{})
					for _, topic := range sub.Topics {
						filter[topic] = struct{}{}
					}

					for topic := range allowedTopics {
						if _, ok := filter[topic]; !ok {
							delete(allowedTopics, topic)
						}
					}
				}
			}
		}
	}()
	wg.Wait()
	log.Info().Uint64("nEvents", nEvents).Uint64("acks", acks).Uint64("nacks", nacks).Msg("subscribe stream terminated")
	return err
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
