/*
Implements a go routine that collects topic info periodically outside of the broker.
*/
package info

import (
	"time"

	"github.com/rotationalio/ensign/pkg/ensign/store"
)

// InfoInterval specifies the delay between topic info gathering runs. This interval is
// approximate because the time it takes to perform an info gathering run is added to
// the interval; the more topics in the system, the longer the info gathering interval.
// NOTE: this is not configurable as it is a core system function.
const InfoInterval = 75 * time.Second

// TopicInfoGatherer runs a go routine that periodically lists all of the topics on the
// node and updates the topic info. This routine is designed to be outside of the broker
// process so that updating topic info does not slow down the evening process. It does
// mean that topic info may be behind the actual state of the server, but given a
// routine enough periodicity, should quickly be resolved with eventual consistency;
// e.g. in the absence of writes, the topic info will eventually become consistent.
//
// NOTE: This routine works as a single thread and guarantees consistency in topic info
// -- no other go routine should write to the topic info, only read from it.
type TopicInfoGatherer struct {
	events store.EventStore
	topics store.TopicInfoStore
}

func New(events store.EventStore, topics store.TopicInfoStore) *TopicInfoGatherer {
	return &TopicInfoGatherer{
		events: events,
		topics: topics,
	}
}

// Run the background go routine that collects topic info from each topic.
// NOTE: this should not be run in maintenance mode.
func (t *TopicInfoGatherer) Run() {

}
