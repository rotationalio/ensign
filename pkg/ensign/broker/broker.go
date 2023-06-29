package broker

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
)

// Every Ensign node is composed of a single broker routine that collects events from
// publisher handlers, commits the events through consensus, writes the events to disk,
// and ensures that registered consumer groups receive any published events they are
// subscribed to.
type Broker interface {
	// Run the broker; any fatal errors will be sent on the specified channel.
	Run(errc chan<- error)

	// Register a publisher to receive an ack/nack channel for events that are
	// published using the publisher ID specified.
	Register(publisherID rlid.RLID) <-chan bool

	// Publish an event from the specified publisher. When the event is committed, an
	// acknowledgement or error is sent on the channel specified when registering.
	Publish(publisherID rlid.RLID, event *api.EventWrapper) error

	// Subscribe to events filtered by topic ids. All recent events will be sent on the
	// event wrapper channel once they are committed.
	Subscribe(subscriberID rlid.RLID, topics ...ulid.ULID) <-chan *api.EventWrapper

	// Close either a publisher or subscriber so no events will be sent from the broker.
	Close(id rlid.RLID) error
}
