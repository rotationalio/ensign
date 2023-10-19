/*
Implements a go routine that collects topic info periodically outside of the broker.
*/
package info

import (
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	// InfoInterval specifies the delay between topic info gathering runs. This interval
	// is approximate because the time it takes to perform an info gathering run is
	// added to the interval; the more topics in the system, the longer the info
	// gathering interval.
	InfoInterval = 75 * time.Second

	// InfoWorkers specifes the number of workers that loop through each topic. The goal
	// of parallelization is to ensure that a large topic does not dominate the info
	// gathering process, but to minimize the amount of CPU needed to perform topic info
	// gathering in favor of publisher and subscriber routines.
	InfoWorkers = 4
)

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
	sync.Mutex
	events  store.EventStore
	topics  store.TopicInfoStore
	done    chan struct{}
	running bool
}

func New(events store.EventStore, topics store.TopicInfoStore) *TopicInfoGatherer {
	return &TopicInfoGatherer{
		events:  events,
		topics:  topics,
		done:    make(chan struct{}),
		running: false,
	}
}

// Run the background go routine that collects topic info from each topic.
// NOTE: this should not be run in maintenance mode.
// WARNING: Do not call this method more than once per process!
func (t *TopicInfoGatherer) Run() {
	go func() {
		ticker := time.NewTicker(InfoInterval)
		log.Info().Dur("interval", InfoInterval).Msg("topic info gatherer started")

		for {
			select {
			case <-t.done:
				log.Info().Msg("topic info gatherer stopped")
				return
			case <-ticker.C:
				var wg sync.WaitGroup
				if err := t.Gather(&wg); err != nil {
					sentry.Fatal(nil).Err(err).Msg("topic info gatherer terminated")
					return
				}
				wg.Wait()
			}
		}
	}()

	t.Lock()
	t.running = true
	t.Unlock()
}

// Shutdown the topic info gatherer; blocks until the topic info has completed.
// WARNING: Do not call this method more than once per process!
func (t *TopicInfoGatherer) Shutdown() error {
	t.Lock()
	defer t.Unlock()

	// This check ensures that we can call shutdown even if the info gatherer is not
	// running, that way we don't block on the send to the done channel. This is
	// important for server tests that don't run the info gatherer.
	if !t.running {
		return nil
	}

	t.done <- struct{}{}
	t.running = false
	return nil
}

// Loops through all topics currently stored in the database and gathers topic info
// from the events, then saves that topic info back to disk. Any errors returned from
// this method are fatal; e.g. the gatherer cannot access the database, otherwise if
// there is a transient failure, then the error is logged.
func (t *TopicInfoGatherer) Gather(wg *sync.WaitGroup) error {
	// Start off the topic info gathering workers
	// NOTE: the wait group is passed in for the outer level so that the gather function
	// can defer and error without waiting for the workers to complete.
	topics := make(chan *api.Topic, InfoWorkers)
	defer close(topics)

	for i := 0; i < InfoWorkers; i++ {
		wg.Add(1)
		go t.worker(wg, topics)
	}

	// Iterate over all topics in the database
	iter := t.topics.ListAllTopics()
	defer iter.Release()

	nTopics := 0
	for iter.Next() {
		topic, err := iter.Topic()
		if err != nil {
			sentry.Error(nil).Bytes("objectKey", iter.Key()).Msg("could not parse topic")
			continue
		}
		topics <- topic
		nTopics++
	}

	if err := iter.Error(); err != nil {
		return err
	}

	log.Debug().Int("workers", InfoWorkers).Int("topics", nTopics).Msg("topic infos gathered")
	return nil
}

func (t *TopicInfoGatherer) worker(wg *sync.WaitGroup, topics <-chan *api.Topic) {
	defer wg.Done()
	for topic := range topics {
		log.Debug().Bytes("topic", topic.Id).Msg("gathering topic info")
		if err := t.handleTopic(topic); err != nil {
			sentry.Warn(nil).Err(err).Bytes("topicID", topic.Id).Bytes("projectID", topic.ProjectId).Msg("could not gather topic info")
		}
	}
}

func (t *TopicInfoGatherer) handleTopic(topic *api.Topic) (err error) {
	var (
		topicID ulid.ULID
		info    *api.TopicInfo
		events  iterator.EventIterator
	)

	if topicID, err = topic.ParseTopicID(); err != nil {
		return fmt.Errorf("could not parse topicID: %w", err)
	}

	if info, err = t.topics.TopicInfo(topicID); err != nil {
		return fmt.Errorf("could not fetch topic info: %w", err)
	}

	events = t.events.List(topicID)
	defer events.Release()

	// Seek over any events that have already been processed.
	if len(info.EventOffsetId) != 0 {
		var eventID rlid.RLID
		if eventID, err = info.ParseEventOffsetID(); err != nil {
			return fmt.Errorf("could not unmarshal event offset id: %w", err)
		}
		events.Seek(eventID)
	}

eventLoop:
	for events.Next() {
		// Fetch the raw data from the iterator rather than parsing the event wrapper
		// so that we can compute the raw data size of the event.
		data := events.Value()
		dataSize := uint64(len(data))

		info.Events++
		info.DataSizeBytes += dataSize

		// Parse the event wrapper from the data
		event := &api.EventWrapper{}
		if err = proto.Unmarshal(data, event); err != nil {
			sentry.Warn(nil).Err(err).Bytes("eventKey", events.Key()).Msg("could not unmarshal event")
			continue eventLoop
		}

		// Store the last ID on the topic info so that we can seek to the next event
		info.EventOffsetId = event.Id

		// Check if the event is a duplicate
		if event.IsDuplicate {
			info.Duplicates++

			// Rehydrate the event from the original to ensure mime and event type
			// correctly compute the duplication if needed.
			var target *api.EventWrapper
			if target, err = t.events.Retrieve(topicID, rlid.RLID(event.DuplicateId)); err != nil {
				sentry.Warn(nil).Err(err).Bytes("targetKey", event.DuplicateId).Msg("could not retreive target of duplicate event")
				continue eventLoop
			}

			if err = event.DuplicateFrom(target); err != nil {
				sentry.Warn(nil).Err(err).Bytes("targetKey", event.DuplicateId).Msg("could not dereferecnce duplicate")
				continue eventLoop
			}
		}

		// Unwrap the event to perform type checking.
		var e *api.Event
		if e, err = event.Unwrap(); err != nil {
			sentry.Warn(nil).Err(err).Bytes("eventKey", events.Key()).Msg("could not unwrap event from wrapper")
			continue
		}

		// Update the event type info on the info
		// NOTE: ResolveType() returns Unspecified if the event does not have a type.
		etypeinfo := info.FindEventTypeInfo(e.ResolveType(), e.Mimetype)
		etypeinfo.Events++
		etypeinfo.DataSizeBytes += dataSize
		if event.IsDuplicate {
			etypeinfo.Duplicates++
		}
	}

	if err = events.Error(); err != nil {
		return fmt.Errorf("could not fetch events: %w", err)
	}

	// Save the topic info back to disk
	if err = t.topics.UpdateTopicInfo(info); err != nil {
		return err
	}
	return nil
}
