package ensign

import (
	"io"
	"sync"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Publish(stream api.Ensign_PublishServer) (err error) {
	// Set up the stream handlers
	nEvents := uint64(0)
	ctx := stream.Context()
	events := make(chan *api.Event, 10000)

	var wg sync.WaitGroup
	wg.Add(2)

	// Execute the ack-back loop
	go func(events <-chan *api.Event) {
		defer wg.Done()
		for event := range events {
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
				err = ctx.Err()
				return
			default:
			}

			var in *api.Event
			if in, err = stream.Recv(); err != nil {
				if err == io.EOF {
					log.Info().Msg("publish stream closed")
					err = nil
					return
				}
				log.Error().Err(err).Msg("publish stream crashed")
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
