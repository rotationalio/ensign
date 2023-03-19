package ensign

import (
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/contexts"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Publish implements the streaming endpoint for the API.
func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	o11y.OnlinePublishers.Inc()
	defer o11y.OnlinePublishers.Dec()

	// Parse the context for authentication information
	ctx := stream.Context()
	claims, ok := contexts.ClaimsFrom(ctx)
	if !ok {
		// NOTE: this should never happen because the interceptor will catch it, but
		// this check prevents nil panics and guards against regressions.
		sentry.Fatal(ctx).Msg("no claims available in publish stream")
		return status.Error(codes.Unauthenticated, "missing credentials")
	}

	// Verify that the user has permissions to publish.
	if !claims.HasPermission(permissions.Publisher) {
		log.Warn().Msg("attempt to open publisher stream without publisher permission")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	var projectID ulid.ULID
	if projectID, err = ulids.Parse(claims.ProjectID); err != nil || ulids.IsZero(projectID) {
		sentry.Warn(ctx).Err(err).Msg("could not parse projectID from claims")
		return status.Error(codes.Unauthenticated, "not authorized to perform this action")
	}

	// Get the topic IDs this user is allowed to publish to
	// TODO: make this more efficient and globalize to prevent inconsistencies
	allowedTopics := make(map[string]struct{})
	topics := s.meta.ListTopics(projectID)

	for topics.Next() {
		// Get and parse the key for the topic ID
		var key meta.ObjectKey
		copy(key[:], topics.Key())

		// Parse the ULID:
		var topicID ulid.ULID
		if topicID, err = key.ObjectID(); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not parse topic id")
			topics.Release()
			return status.Error(codes.Internal, "could not open publisher stream")
		}

		allowedTopics[topicID.String()] = struct{}{}
	}

	if err := topics.Error(); err != nil {
		sentry.Warn(ctx).Err(err).Msg("could not fetch topics")
		topics.Release()
		return status.Error(codes.Internal, "could not open publisher stream")
	}
	topics.Release()

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
				// TODO: should we error here and close the stream?
				log.Warn().Msg("event published without topic id")

				// Send the nack back to the user
				stream.Send(&api.Publication{
					// TODO: what are the nack error codes?
					Embed: &api.Publication_Nack{
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
				// TODO: should we error here and close the stream?
				log.Warn().Msg("event published to topic that is not allowed")

				// Send the nack back to the user
				stream.Send(&api.Publication{
					// TODO: what are the nack error codes?
					Embed: &api.Publication_Nack{
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

			// Push event on to the primary buffer
			s.pubsub.Publish(event)

			// Send ack once the event is on the primary buffer
			err = stream.Send(&api.Publication{
				Embed: &api.Publication_Ack{
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

			var in *api.Event
			if in, err = stream.Recv(); err != nil {
				if streamClosed(err) {
					log.Info().Msg("publish stream closed")
					err = nil
					return
				}
				log.Warn().Err(err).Msg("publish stream crashed")
				return
			}

			events <- in
		}
	}(events)

	wg.Wait()
	stream.Send(&api.Publication{
		Embed: &api.Publication_CloseStream{
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

	// Setup the stream handlers
	nEvents, acks, nacks := uint64(0), uint64(0), uint64(0)
	ctx := stream.Context()
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
				if err = stream.Send(event); err != nil {
					if streamClosed(err) {
						log.Info().Msg("subscribe stream closed")
						err = nil
						return
					}
					log.Warn().Err(err).Msg("subscribe stream crashed")
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

			var in *api.Subscription
			if in, err = stream.Recv(); err != nil {
				if streamClosed(err) {
					log.Info().Msg("subscribe stream closed")
					err = nil
					return
				}
				log.Warn().Err(err).Msg("subscribe stream crashed")
				return
			}

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
