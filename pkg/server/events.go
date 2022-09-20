package server

import (
	"fmt"
	"io"
	"sync"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/server/o11y"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const BufferSize = 10000

// PubSub is a simple dispatcher that has a publish queue that fans in events from
// different publisher streams and assigns event ids then fans the events out to send
// them to one or more subscriber streams with an outgoing buffer. Backpressure is
// applied to the publisher streams when the buffers get full.
type PubSub struct {
	inQ     chan inQ
	outQ    chan *api.Event
	counter uint64
	subs    []chan<- *api.Event
}

type inQ struct {
	res   chan uint64
	event *api.Event
}

func NewPubSub() (ps *PubSub) {
	ps = &PubSub{
		inQ:     make(chan inQ, BufferSize),
		outQ:    make(chan *api.Event, BufferSize),
		counter: 0,
		subs:    make([]chan<- *api.Event, 0),
	}
	go ps.pub()
	go ps.sub()
	return ps
}

// Handles incoming events being published; loops on the incoming queue, assigns a
// monotonically increasing event ID then puts the even on the outgoing queue.
func (ps *PubSub) pub() {
	for i := range ps.inQ {
		ps.counter++
		i.event.Id = fmt.Sprintf("%04x", ps.counter)
		ps.outQ <- i.event
		i.res <- ps.counter
	}
}

// Handles outgoing events being sent to subscribers; loops on the outgoing queue and
// and sends the event to all connected subscribers. If there are no subscribers then
// the event is dropped.
func (ps *PubSub) sub() {
	for e := range ps.outQ {

		log.Debug().Str("id", e.Id).Msg("event handled")
	}
}

func (ps *PubSub) Publish(event *api.Event) uint64 {
	q := inQ{
		event: event,
		res:   make(chan uint64, 1),
	}
	ps.inQ <- q
	return <-q.res
}

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
