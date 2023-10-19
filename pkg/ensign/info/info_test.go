package info_test

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/info"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestInfoGather(t *testing.T) {
	// This test works in four phases. In the first phase, there is nothing in the
	// database and gather is run. In the second phase, topics are added, some with
	// events and some without events. In the third phase, all topics have events and
	// some topics have more events. In the fourth phase, all topics have events but
	// some topics hve duplicate events.

	// Setup the database and the gatherer
	events, topics := createDatabase(t)
	gatherer := info.New(events, topics)

	// Execute phase 0: no data in the database
	wg := &sync.WaitGroup{}
	err := gatherer.Gather(wg)
	require.NoError(t, err, "couldn't execute phase 0 gather")
	wg.Wait()
	checkPhase0(t, topics)

	// Execute phase 1: initial topics and events in the database
	setupPhase1(t, events, topics)
	wg = &sync.WaitGroup{}
	err = gatherer.Gather(wg)
	require.NoError(t, err, "could not execute phase 1 gather")
	wg.Wait()
	checkPhase1(t, topics)

	// Execute phase 2: events with duplicates in the database
	setupPhase2(t, events, topics)
	wg = &sync.WaitGroup{}
	err = gatherer.Gather(wg)
	require.NoError(t, err, "could not execute phase 2 gather")
	wg.Wait()
	checkPhase2(t, topics)
}

func TestInfoGatherFatal(t *testing.T) {
	store := &mock.Store{}
	store.UseError(mock.ListAllTopics, errors.New("this should be a fatal error"))
	gatherer := info.New(store, store)

	var wg sync.WaitGroup
	err := gatherer.Gather(&wg)
	require.Error(t, err, "expected fatal error when not able to list all topics")
	wg.Wait()

	require.Equal(t, 1, store.Calls(mock.ListAllTopics))
	require.Zero(t, store.Calls(mock.List))
	require.Zero(t, store.Calls(mock.TopicInfo))
	require.Zero(t, store.Calls(mock.UpdateTopicInfo))
}

func TestInfoGatherRunShutdown(t *testing.T) {
	store := &mock.Store{}
	gatherer := info.New(store, store)

	gatherer.Run()
	err := gatherer.Shutdown()
	require.NoError(t, err)
}

func createDatabase(t *testing.T) (store.EventStore, store.MetaStore) {
	dbpath, err := os.MkdirTemp("", "infogather")
	require.NoError(t, err, "could not create temporary directory for database")
	t.Cleanup(func() { os.RemoveAll(dbpath) })

	conf := config.StorageConfig{
		ReadOnly: false,
		DataPath: dbpath,
		Testing:  false,
	}

	events, topics, err := store.Open(conf)
	require.NoError(t, err, "could not open events and topics store")
	t.Cleanup(func() {
		events.Close()
		topics.Close()
	})

	return events, topics
}

func checkPhase0(t *testing.T, topics store.TopicInfoStore) {
	count := 0
	iter := topics.ListAllTopics()
	defer iter.Release()

	for iter.Next() {
		topic, err := iter.Topic()
		require.Error(t, err, "could not parse topic")

		topicID, err := topic.ParseTopicID()
		require.NoError(t, err, "could not parse topic id")

		_, err = topics.TopicInfo(topicID)
		require.NoError(t, err, "could not fetch info for topic")
		count++
	}

	err := iter.Error()
	require.NoError(t, err, "could not iterate over topics in store")
	require.Zero(t, count, "expected no topic info in the database")
}

func setupPhaseN(t *testing.T, path string, eventDB store.EventStore, topicDB store.MetaStore) {
	topics, events := loadFixture(t, path)
	for _, topic := range topics {
		err := topicDB.CreateTopic(topic)
		require.NoError(t, err, "could not create topic in database")
	}

	for _, event := range events {
		err := eventDB.Insert(event)
		require.NoError(t, err, "could not insert event into database")
	}
}

func setupPhase1(t *testing.T, eventDB store.EventStore, topicDB store.MetaStore) {
	setupPhaseN(t, "testdata/phase1.json", eventDB, topicDB)
}

func setupPhase2(t *testing.T, eventDB store.EventStore, topicDB store.MetaStore) {
	setupPhaseN(t, "testdata/phase2.json", eventDB, topicDB)
}

func checkPhaseN(t *testing.T, expected map[string]*api.TopicInfo, topics store.TopicInfoStore) {
	count := 0
	iter := topics.ListAllTopics()
	defer iter.Release()

	for iter.Next() {
		topic, err := iter.Topic()
		require.NoError(t, err, "could not parse topic")

		topicID, err := topic.ParseTopicID()
		require.NoError(t, err, "could not parse topic id")

		info, err := topics.TopicInfo(topicID)
		require.NoError(t, err, "could not fetch info for topic")

		infoTopicID, err := info.ParseTopicID()
		require.NoError(t, err, "could not parse topicID from info")
		require.Equal(t, topicID, infoTopicID, "info topic ID and topic ID do not match")

		require.Contains(t, expected, infoTopicID.String(), "unexpected topic info written to database")
		expectedInfo := expected[infoTopicID.String()]

		require.Equal(t, expectedInfo.TopicId, info.TopicId, "topic id mismatch")
		require.Equal(t, expectedInfo.ProjectId, info.ProjectId, "project id mismatch")
		// require.Equal(t, expectedInfo.EventOffsetId, info.EventOffsetId, "event offset id mismatch")
		require.Equal(t, expectedInfo.Events, info.Events, "event count mismatch")
		require.Equal(t, expectedInfo.Duplicates, info.Duplicates, "duplicates count mismatch")
		require.Equal(t, expectedInfo.DataSizeBytes, info.DataSizeBytes, "data size mismatch")
		require.Equal(t, len(expectedInfo.Types), len(info.Types), "different numbers of event type info")
		require.False(t, info.Modified.AsTime().IsZero(), "modified timestamp is zero")

		for _, aeti := range info.Types {
			eeti := expectedInfo.FindEventTypeInfo(aeti.Type, aeti.Mimetype)
			require.Equal(t, eeti.Events, aeti.Events, "%s (%s) type events mismatch", aeti.Type.Repr(), aeti.Mimetype)
			require.Equal(t, eeti.Duplicates, aeti.Duplicates, "%s (%s) type duplicates mismatch", aeti.Type.Repr(), aeti.Mimetype)
			require.Equal(t, eeti.DataSizeBytes, aeti.DataSizeBytes, "%s (%s) type data size mismatch", aeti.Type.Repr(), aeti.Mimetype)
		}

		count++
	}

	err := iter.Error()
	require.NoError(t, err, "could not iterate over topics in store")
	require.Equal(t, len(expected), count, "expected topic infos for each topic in the database")
}

func checkPhase1(t *testing.T, topics store.TopicInfoStore) {
	expected := map[string]*api.TopicInfo{
		"01GTSMQ3V8ASAPNCFEN378T8RD": {
			TopicId:       ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        0,
			Duplicates:    0,
			DataSizeBytes: 0,
			Types:         []*api.EventTypeInfo{},
		},
		"01GTSMSX1M9G2Z45VGG4M12WC0": {
			TopicId:       ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        10,
			Duplicates:    0,
			DataSizeBytes: 0x7a1,
			Types: []*api.EventTypeInfo{
				{
					Type:          &api.Type{Name: "RandomOrder", MajorVersion: 2},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        3,
					Duplicates:    0,
					DataSizeBytes: 0x22d,
				},
				{
					Type:          &api.Type{Name: "RandomOrder", MajorVersion: 1, MinorVersion: 9, PatchVersion: 2},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        7,
					Duplicates:    0,
					DataSizeBytes: 0x574,
				},
			},
		},
		"01GTSN1139JMK1PS5A524FXWAZ": {
			TopicId:       ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        2,
			Duplicates:    0,
			DataSizeBytes: 0x145,
			Types: []*api.EventTypeInfo{
				{
					Type:          &api.Type{Name: "RandomShipment", MajorVersion: 1},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        2,
					Duplicates:    0,
					DataSizeBytes: 0x145,
				},
			},
		},
	}
	checkPhaseN(t, expected, topics)
}

func checkPhase2(t *testing.T, topics store.TopicInfoStore) {
	expected := map[string]*api.TopicInfo{
		"01GTSMQ3V8ASAPNCFEN378T8RD": {
			TopicId:       ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        0,
			Duplicates:    0,
			DataSizeBytes: 0,
			Types:         []*api.EventTypeInfo{},
		},
		"01GTSMSX1M9G2Z45VGG4M12WC0": {
			TopicId:       ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        11,
			Duplicates:    1,
			DataSizeBytes: 0x80f,
			Types: []*api.EventTypeInfo{
				{
					Type:          &api.Type{Name: "RandomOrder", MajorVersion: 2},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        3,
					Duplicates:    0,
					DataSizeBytes: 0x22d,
				},
				{
					Type:          &api.Type{Name: "RandomOrder", MajorVersion: 1, MinorVersion: 9, PatchVersion: 2},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        8,
					Duplicates:    1,
					DataSizeBytes: 0x5e2,
				},
			},
		},
		"01GTSN1139JMK1PS5A524FXWAZ": {
			TopicId:       ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        2,
			Duplicates:    0,
			DataSizeBytes: 0x145,
			Types: []*api.EventTypeInfo{
				{
					Type:          &api.Type{Name: "RandomShipment", MajorVersion: 1},
					Mimetype:      mimetype.MIME_UNSPECIFIED,
					Events:        2,
					Duplicates:    0,
					DataSizeBytes: 0x145,
				},
			},
		},
		"01GTSN1WF5BA0XCPT6ES64JVGQ": {
			TopicId:       ulid.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ").Bytes(),
			ProjectId:     ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM").Bytes(),
			EventOffsetId: []byte{0x1, 0x89, 0xec, 0x33, 0x51, 0x60, 0x0, 0x0, 0x0, 0xa},
			Events:        0,
			Duplicates:    0,
			DataSizeBytes: 0,
			Types:         []*api.EventTypeInfo{},
		},
	}
	checkPhaseN(t, expected, topics)
}

func loadFixture(t *testing.T, path string) ([]*api.Topic, []*api.EventWrapper) {
	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	data, err := os.ReadFile(path)
	require.NoError(t, err, "could not open %s", path)

	jfixtures := make(map[string][]interface{})
	err = json.Unmarshal(data, &jfixtures)
	require.NoError(t, err, "could not unmarshal json fixtures")

	topics := make([]*api.Topic, 0, len(jfixtures["topics"]))
	for _, jtopic := range jfixtures["topics"] {
		data, err := json.Marshal(jtopic)
		require.NoError(t, err, "could not marshal json topic")

		topic := &api.Topic{}
		require.NoError(t, pbjson.Unmarshal(data, topic), "could not unmarshal topic protobuf")
		topics = append(topics, topic)
	}

	events := make([]*api.EventWrapper, 0, len(jfixtures["events"]))
	for _, jevent := range jfixtures["events"] {
		data, err := json.Marshal(jevent)
		require.NoError(t, err, "could not marshal json event wrapper")

		event := &api.EventWrapper{}
		require.NoError(t, pbjson.Unmarshal(data, event), "could not unmarshal event protobuf")
		events = append(events, event)
	}

	return topics, events
}
