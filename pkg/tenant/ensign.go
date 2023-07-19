package tenant

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/metatopic"
	sdk "github.com/rotationalio/go-ensign"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TopicSubscriber is a struct with a go routine that subscribes to the Ensign
// "meta topic" topic and applies asynchronous updates to the topics in the Tenant
// database.
type TopicSubscriber struct {
	client *EnsignClient
	stop   chan struct{}
}

func NewTopicSubscriber(client *EnsignClient) *TopicSubscriber {
	return &TopicSubscriber{
		client: client,
	}
}

// Run the topic subscriber under the waitgroup. This allows the caller to wait for the
// subscriber to graacefully exit after calling Stop().
func (s *TopicSubscriber) Run(wg *sync.WaitGroup) error {
	if s.stop != nil {
		return errors.New("topic subscriber is already running")
	}

	if wg == nil {
		return errors.New("waitgroup must be provided to run the topic subscriber")
	}

	s.stop = make(chan struct{})
	wg.Add(1)
	go func() {
		s.Subscribe()
		s.stop = nil
		wg.Done()
	}()
	return nil
}

// Stop the topic subscriber.
func (s *TopicSubscriber) Stop() {
	if s.stop != nil {
		close(s.stop)
	}
}

func (s *TopicSubscriber) Subscribe() {
	var (
		err error
		sub *sdk.Subscription
	)

	// Subscribe to the meta topic
	if sub, err = s.client.Subscribe(); err != nil {
		// Note: Using WithLevel with FatalLevel does not exit the program but this is
		// likely a critical configuration error that we want to fix immediately.
		log.WithLevel(zerolog.FatalLevel).Err(err).Msg("failed to subscribe to meta topic")
		return
	}
	defer sub.Close()

	// Handle events from the stream
	// This assumes that the SDK properly handles connection issues and will reconnect
	// to the stream if necessary.
	for {
		select {
		case <-s.stop:
			return
		case event := <-sub.C:
			if event.Mimetype != mimetype.ApplicationMsgPack {
				log.Warn().Str("mimetype", event.Mimetype.String()).Msg("unexpected mimetype for metatopic event")
				event.Nack(api.Nack_UNHANDLED_MIMETYPE)
				continue
			}

			if event.Type == nil {
				log.Warn().Msg("missing event type")
				event.Nack(api.Nack_UNKNOWN_TYPE)
				continue
			}

			if event.Type.Name != metatopic.SchemaName {
				log.Warn().Str("type", event.Type.Name).Msg("unexpected event type")
				event.Nack(api.Nack_UNKNOWN_TYPE)
				continue
			}

			// TODO: Ensure version is correct

			// Attempt to decode the event data
			topicUpdate := &metatopic.TopicUpdate{}
			if err = topicUpdate.Unmarshal(event.Data); err != nil {
				log.Warn().Err(err).Msg("failed to unmarshal topic update event")
				event.Nack(api.Nack_UNKNOWN_TYPE)
				continue
			}

			// Validate the event schema
			if err = topicUpdate.Validate(); err != nil {
				log.Warn().Err(err).Msg("failed to validate topic update event")
				event.Nack(api.Nack_UNPROCESSED)
				continue
			}

			// Handle the event. This will return an error if the event includes an
			// update or change to a topic that does not exist in the database.
			if err = s.performUpdate(topicUpdate); err != nil {
				log.Error().Err(err).Msg("failed to perform topic update")
				event.Nack(api.Nack_UNPROCESSED)
				continue
			}

			// Acknowledge the event
			event.Ack()
		}
	}
}

// Write a topic update to the database. This assumes that the update schema has
// already been validated.
func (s *TopicSubscriber) performUpdate(update *metatopic.TopicUpdate) (err error) {
	var topic *db.Topic

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch update.UpdateType {
	case metatopic.TopicUpdateCreated:
		// Create the topic
		topic := &db.Topic{
			OrgID:       update.OrgID,
			ProjectID:   update.ProjectID,
			ID:          update.TopicID,
			Name:        update.Topic.Name,
			Events:      update.Topic.Events,
			Storage:     update.Topic.Storage,
			Publishers:  update.Topic.Publishers,
			Subscribers: update.Topic.Subscribers,
			Created:     update.Topic.Created,
			Modified:    update.Topic.Modified,
		}

		if err = db.CreateTopic(ctx, topic); err != nil {
			return err
		}
	case metatopic.TopicUpdateModified:
		// Get the existing topic
		if topic, err = db.RetrieveTopic(ctx, update.TopicID); err != nil {
			return err
		}

		// Update the modifiable fields on the topic
		topic.Name = update.Topic.Name
		topic.Events = update.Topic.Events
		topic.Storage = update.Topic.Storage
		topic.Publishers = update.Topic.Publishers
		topic.Subscribers = update.Topic.Subscribers

		if err = db.UpdateTopic(ctx, topic); err != nil {
			return err
		}
	case metatopic.TopicUpdateStateChange:
		// Retrieve the topic
		if topic, err = db.RetrieveTopic(ctx, update.TopicID); err != nil {
			return err
		}

		// Currently the only valid state change is read-only
		topic.State = api.TopicTombstone_READONLY

		if err = db.UpdateTopic(ctx, topic); err != nil {
			return err
		}
	case metatopic.TopicUpdateDeleted:
		// Delete the topic
		if err = db.DeleteTopic(ctx, update.TopicID); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown topic update type: %s", update.UpdateType)
	}

	return nil
}
