package mock_test

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	region "github.com/rotationalio/ensign/pkg/ensign/region/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnmarshalEventList(t *testing.T) {
	events, err := mock.EventListFixture("testdata/events.pb.json")
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestUnmarshalTopicList(t *testing.T) {
	topics, err := mock.TopicListFixture("testdata/topics.pb.json")
	require.NoError(t, err)
	require.Len(t, topics, 7)
}

func TestUnmarshalTopicNamesList(t *testing.T) {
	names, err := mock.TopicNamesListFixture("testdata/topicnames.pb.json")
	require.NoError(t, err)
	require.Len(t, names, 7)
}

func TestUnmarshalTopicInfosList(t *testing.T) {
	names, err := mock.TopicInfoListFixture("testdata/topicinfos.pb.json")
	require.NoError(t, err)
	require.Len(t, names, 0)
}

func TestGenerateJSONFixtures(t *testing.T) {
	if env := os.Getenv("ENSIGN_GENERATE_JSON_FIXTURES"); env != "1" {
		t.Skip("skipping json fixture generation")
	}

	fixtures := []*TopicFixtureGenerator{
		{
			TopicID:    "01GTSMQ3V8ASAPNCFEN378T8RD",
			ProjectID:  "01GTSMMC152Q95RD4TNYDFJGHT",
			Name:       "testing.testapp.alerts",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GTSMSX1M9G2Z45VGG4M12WC0",
			ProjectID:  "01GTSMMC152Q95RD4TNYDFJGHT",
			Name:       "testing.testapp.orders",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GTSN1139JMK1PS5A524FXWAZ",
			ProjectID:  "01GTSMMC152Q95RD4TNYDFJGHT",
			Name:       "testing.testapp.shipments",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GV6KXTEPSWZHZB4XW9RWDSAA",
			ProjectID:  "01GTSMMC152Q95RD4TNYDFJGHT",
			Name:       "testing.testapp.products",
			ReadOnly:   true,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GV6KYPW33RW5D800ERR3NP8S",
			ProjectID:  "01GTSMMC152Q95RD4TNYDFJGHT",
			Name:       "testing.testapp.receipts",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GTSN1WF5BA0XCPT6ES64JVGQ",
			ProjectID:  "01GTSMZNRYXNAZQF5R8NHQ14NM",
			Name:       "mock.mockapp.feed",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
		{
			TopicID:    "01GTSN2NQV61P2R4WFYF1NF1JG",
			ProjectID:  "01GTSMZNRYXNAZQF5R8NHQ14NM",
			Name:       "mock.mockapp.post",
			ReadOnly:   false,
			EventTypes: []*EventFixtureGenerator{},
		},
	}

	topics := make([]protoreflect.ProtoMessage, 0)
	names := make([]protoreflect.ProtoMessage, 0)
	infos := make([]protoreflect.ProtoMessage, 0)
	events := make([]protoreflect.ProtoMessage, 0)

	for _, fixture := range fixtures {
		topics = append(topics, fixture.Topic())
		names = append(names, fixture.TopicName())
		infos = append(infos, fixture.TopicInfo())
	}

	err := WriteFixture(topics, "testdata/topics.pb.json")
	require.NoError(t, err, "could not write topics fixture")

	err = WriteFixture(names, "testdata/topicnames.pb.json")
	require.NoError(t, err, "could not write topic names fixture")

	err = WriteFixture(infos, "testdata/topicinfos.pb.json")
	require.NoError(t, err, "could not write topic infos fixture")

	err = WriteFixture(events, "testdata/events.pb.json")
	require.NoError(t, err, "could not write events fixture")
}

type TopicFixtureGenerator struct {
	TopicID    string
	ProjectID  string
	Name       string
	ReadOnly   bool
	EventTypes []*EventFixtureGenerator
	topic      *api.Topic
}

type EventFixtureGenerator struct {
	Count         uint64
	TypeName      string
	TypeSemver    string
	Mimetype      string
	AvgEventSize  float64
	StdEventDev   float64
	MetaGenerator func(int) map[string]string
	DataGenerator func(int) []byte
	info          *api.EventTypeInfo
	sizes         []uint64
}

func (t TopicFixtureGenerator) Topic() *api.Topic {
	if t.topic == nil {
		t.topic = &api.Topic{
			Id:        ulid.MustParse(t.TopicID).Bytes(),
			ProjectId: ulid.MustParse(t.ProjectID).Bytes(),
			Name:      t.Name,
			Readonly:  t.ReadOnly,
			Offset:    0,
			Shards:    1,
			Placements: []*api.Placement{
				{
					Epoch:    1,
					Sharding: api.ShardingStrategy_NO_SHARDING,
					Regions: []region.Region{
						region.Region_STG_LKE_US_EAST_1A,
					},
					Nodes: []*api.Node{
						{
							Id:       "staging-1",
							Hostname: "staging-1.ensign.ensign.svc.local",
							Quorum:   1,
							Shard:    1,
							Region:   region.Region_STG_LKE_US_EAST_1A,
						},
					},
				},
			},
			Types: make([]*api.Type, 0, len(t.EventTypes)),
		}

		t.topic.Created = timestamppb.New(ulid.Time(ulid.MustParse(t.TopicID).Time()))
		t.topic.Modified = timestamppb.New(t.topic.Created.AsTime().Add(time.Duration(rand.Int63n(1.577e16))))

		for _, event := range t.EventTypes {
			t.topic.Offset += uint64(event.Count)
			t.topic.Types = append(t.topic.Types, event.Type())
		}
	}

	return t.topic
}

func (t *TopicFixtureGenerator) TopicName() *api.TopicName {
	return &api.TopicName{
		TopicId:   t.TopicID,
		ProjectId: t.ProjectID,
		Name:      t.Name,
	}
}

func (t *TopicFixtureGenerator) TopicInfo() *api.TopicInfo {
	topic := t.Topic()
	info := &api.TopicInfo{
		TopicId:   topic.Id,
		ProjectId: topic.ProjectId,
		Types:     make([]*api.EventTypeInfo, 0, len(t.EventTypes)),
		Modified:  topic.Modified,
	}

	for _, evt := range t.EventTypes {
		event := evt.Info()
		event.Modified = topic.Modified
		info.Types = append(info.Types, event)
		info.Events += event.Events
		info.Duplicates += event.Duplicates
		info.DataSizeBytes += event.DataSizeBytes
	}
	return info
}

func (e *TopicFixtureGenerator) Events() []*api.EventWrapper {
	// Generate a list of events for all event types
	topic := e.Topic()
	events := make([]*api.Event, 0)
	for _, etype := range e.EventTypes {
		events = append(events, etype.Events()...)
	}

	// Shuffle the events
	rand.Shuffle(len(events), func(i, j int) { events[i], events[j] = events[j], events[i] })

	// Setup the event wrappers sequence generation
	var offset uint64
	start := topic.Created.AsTime().Add(time.Duration(rand.Int63n(1.8e11)))
	sequence := rlid.Sequence(0)
	envs := make([]*api.EventWrapper, 0, len(events))

	// Create some publishers
	pubs := make([]*api.Publisher, rand.Intn(5)+1)
	for i := range pubs {
		pubs[i] = &api.Publisher{
			ClientId:    fmt.Sprintf("%06X", rand.Int31()),
			PublisherId: ulid.MustNew(ulid.Timestamp(start.Add(time.Duration(-1*rand.Int63n(1.8e11)))), crand.Reader).String(),
		}
	}

	for _, event := range events {
		offset++
		env := &api.EventWrapper{
			Id:        sequence.Next().Bytes(),
			TopicId:   topic.Id,
			Offset:    offset,
			Epoch:     1,
			Region:    region.Region_STG_LKE_US_EAST_1A,
			Publisher: pubs[rand.Intn(len(pubs))],
			Key:       nil,
			Shard:     1,
			Encryption: &api.Encryption{
				SealingAlgorithm: api.Encryption_PLAINTEXT,
			},
			Compression: &api.Compression{
				Algorithm: api.Compression_NONE,
			},
			Committed: timestamppb.New(start),
			LocalId:   nil,
		}

		event.Created = timestamppb.New(start.Add(time.Duration(-1 * rand.Int63n(1.5e9))))
		env.LocalId = ulid.MustNew(ulid.Timestamp(event.Created.AsTime()), crand.Reader).Bytes()
		if err := env.Wrap(event); err != nil {
			panic(err)
		}

		envs = append(envs, env)
		start = start.Add(time.Duration(rand.Int63n(1.8e11)))
	}

	return envs
}

func (e *EventFixtureGenerator) Events() []*api.Event {
	// Make sure info and sample sizes are generated
	e.Info()
	etype := e.Type()
	mime := mimetype.MustParse(e.Mimetype)
	events := make([]*api.Event, 0, e.info.Events)

	for i := uint64(0); i < e.info.Events; i++ {
		event := &api.Event{
			Data:     e.DataGenerator(int(e.sizes[i])),
			Metadata: e.MetaGenerator(int(i)),
			Mimetype: mime,
			Type:     etype,
		}

		events = append(events, event)
	}
	return events
}

func (e *EventFixtureGenerator) Type() *api.Type {
	parts := strings.Split(e.TypeSemver, ".")
	t := &api.Type{
		Name: e.TypeName,
	}

	if len(parts) > 0 {
		vers, _ := strconv.ParseUint(parts[0], 10, 32)
		t.MajorVersion = uint32(vers)
	}

	if len(parts) > 1 {
		vers, _ := strconv.ParseUint(parts[1], 10, 32)
		t.MinorVersion = uint32(vers)
	}

	if len(parts) > 2 {
		vers, _ := strconv.ParseUint(parts[2], 10, 32)
		t.PatchVersion = uint32(vers)
	}

	return t
}

func (e *EventFixtureGenerator) Info() *api.EventTypeInfo {
	// Less than 20% of the events should be duplicates
	duplicates := uint64(rand.Int63n(int64(float64(e.Count) * .2)))

	if e.info == nil {
		e.info = &api.EventTypeInfo{
			Type:       e.Type(),
			Mimetype:   mimetype.MustParse(e.Mimetype),
			Events:     e.Count - duplicates,
			Duplicates: duplicates,
		}
		e.sizes = make([]uint64, 0, e.info.Events)

		// Compute a data size using the averages
		for i := uint64(0); i < e.info.Events; i++ {
			sample := uint64(rand.NormFloat64()*e.StdEventDev + e.AvgEventSize)
			e.sizes = append(e.sizes, sample)
			e.info.DataSizeBytes += sample
		}
	}
	return e.info
}

func WriteFixture(fixtures []protoreflect.ProtoMessage, path string) (err error) {
	jsonpb := &protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	repr := make([]interface{}, 0, len(fixtures))
	for _, fixture := range fixtures {
		var data []byte
		if data, err = jsonpb.Marshal(fixture); err != nil {
			return err
		}

		var intermediate interface{}
		if err = json.Unmarshal(data, &intermediate); err != nil {
			return err
		}

		repr = append(repr, intermediate)
	}

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(repr)
}
