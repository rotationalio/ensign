package ensign

import (
	"io"
	"sync"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/o11y"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Publish implements the streaming endpoint for the API.
func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	o11y.OnlinePublishers.Inc()
	defer o11y.OnlinePublishers.Dec()

	// Set up the stream handlers
	nEvents := uint64(0)
	ctx := stream.Context()
	events := make(chan *api.Event, BufferSize)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the ack-back loop
	// This loop also pushes the event onto the primary buffer
	go func(events <-chan *api.Event) {
		defer wg.Done()
		for event := range events {
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
				close(events)
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
