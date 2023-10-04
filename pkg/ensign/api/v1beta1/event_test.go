package api_test

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEventWrapper(t *testing.T) {
	// Create an event wrapper and random event data
	wrap := &api.EventWrapper{
		Id:      ulid.Make().Bytes(),
		TopicId: ulid.Make().Bytes(),
		Offset:  421,
		Epoch:   23,
	}

	evt := &api.Event{
		Data: make([]byte, 128),
	}
	rand.Read(evt.Data)

	err := wrap.Wrap(evt)
	require.NoError(t, err, "should be able to wrap an event in an event wrapper")

	cmp, err := wrap.Unwrap()
	require.NoError(t, err, "should be able to unwrap an event in an event wrapper")
	require.NotNil(t, cmp, "the unwrapped event should not be nil")
	require.True(t, proto.Equal(evt, cmp), "the unwrapped event should match the original")
	require.NotSame(t, evt, cmp, "a pointer to the same event should not be returned")

	wrap.Event = nil
	empty, err := wrap.Unwrap()
	require.EqualError(t, err, "event wrapper contains no event")
	require.Empty(t, empty, "no data event should be zero-valued")

	wrap.Event = []byte("foo")
	_, err = wrap.Unwrap()
	require.Error(t, err, "should not be able to unwrap non-protobuf data")
}

func TestEventWrapperIDParsing(t *testing.T) {
	testCases := []struct {
		eventID  []byte
		expected rlid.RLID
		err      error
	}{
		{nil, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte{}, rlid.RLID{}, rlid.ErrDataSize},
		{[]byte("foo"), rlid.RLID{}, rlid.ErrDataSize},
		{rlid.MustParse("064zwbj8vg00000n").Bytes(), rlid.MustParse("064zwbj8vg00000n"), nil},
	}

	for i, tc := range testCases {
		event := &api.EventWrapper{Id: tc.eventID}
		eventID, err := event.ParseEventID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, rlid.Null, eventID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, eventID, "test case %d failed", i)
		}
	}
}

func TestEventWrapperTopicIDParsing(t *testing.T) {
	testCases := []struct {
		topicID  []byte
		expected ulid.ULID
		err      error
	}{
		{nil, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte{}, ulid.ULID{}, ulid.ErrDataSize},
		{[]byte("foo"), ulid.ULID{}, ulid.ErrDataSize},
		{ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ").Bytes(), ulid.MustParse("01H7Z2XD3VDEV0VKG8PF699ZGQ"), nil},
	}

	for i, tc := range testCases {
		event := &api.EventWrapper{TopicId: tc.topicID}
		topicID, err := event.ParseTopicID()
		if tc.err != nil {
			require.Error(t, err, "test case %d failed", i)
			require.Equal(t, ulid.ULID{}, topicID, "test case %d failed", i)
		} else {
			require.NoError(t, err, "test case %d failed", i)
			require.Equal(t, tc.expected, topicID, "test case %d failed", i)
		}
	}
}

func TestEventEquality(t *testing.T) {
	testCases := []struct {
		name   string
		alpha  *api.Event
		bravo  *api.Event
		assert require.BoolAssertionFunc
	}{
		{
			"nil events", nil, nil, require.True,
		},
		{
			"zero valued events", &api.Event{}, &api.Event{}, require.True,
		},
		{
			"zero valued and nil events", &api.Event{}, nil, require.False,
		},
		{
			"nil and zero valued events", nil, &api.Event{}, require.False,
		},
		{
			"equal, fully populated events",
			mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.True,
		},
		{
			"equal, fully populated events with different created timestamps",
			mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-11-13T21:56:01-05:00"),
			require.True,
		},
		{
			"different data",
			mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different mimetype",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different type name",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different type major version",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v2.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different type minor version",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.8.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different type patch version",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.19", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different key value 1",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:zap,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different key value 2",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:blue", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"extra key",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red,age:42", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"missing key",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"different keys",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "bar:foo,name:strangeloop", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"empty key values",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo,color", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
		{
			"no keys",
			mkevt("8sS9TZqAY33MOj9RMJytyg", "foo:bar,color:red", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			mkevt("8sS9TZqAY33MOj9RMJytyg", "", mimetype.ApplicationAvro, "MockEvent v1.2.3", "2023-10-04T08:17:22-05:00"),
			require.False,
		},
	}

	for i, tc := range testCases {
		tc.assert(t, tc.alpha.Equals(tc.bravo), "test case %s (%d) failed", tc.name, i)
	}
}

func TestResolveType(t *testing.T) {
	testCases := []struct {
		event    *api.Event
		expected *api.Type
	}{
		{&api.Event{}, api.UnspecifiedType},
		{&api.Event{Type: nil}, api.UnspecifiedType},
		{&api.Event{Type: &api.Type{}}, api.UnspecifiedType},
		{&api.Event{Type: &api.Type{Name: "TestType", MajorVersion: 1}}, &api.Type{Name: "TestType", MajorVersion: 1}},
	}

	for i, tc := range testCases {
		actual := tc.event.ResolveType()
		require.Equal(t, tc.expected, actual, "test case %d failed", i)
	}
}

func TestResolveClientID(t *testing.T) {
	testCases := []struct {
		pub      *api.Publisher
		expected string
	}{
		{&api.Publisher{}, ""},
		{&api.Publisher{PublisherId: "testpub1"}, "testpub1"},
		{&api.Publisher{ClientId: "testclient1"}, "testclient1"},
		{&api.Publisher{PublisherId: "testpub1", ClientId: "testclient1"}, "testclient1"},
	}

	for i, tc := range testCases {
		actual := tc.pub.ResolveClientID()
		require.Equal(t, tc.expected, actual, "test case %d failed", i)
	}
}

func TestTypeEquality(t *testing.T) {
	testCases := []struct {
		in      *api.Type
		require require.BoolAssertionFunc
	}{
		{nil, require.False},
		{&api.Type{}, require.False},
		{&api.Type{Name: "TESTTYPE", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "testtype", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "testType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 4, MinorVersion: 2, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 7, PatchVersion: 3}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 14}, require.False},
		{&api.Type{Name: "TestType", MajorVersion: 1, MinorVersion: 2, PatchVersion: 3}, require.True},
	}

	etype := &api.Type{
		Name:         "TestType",
		MajorVersion: 1,
		MinorVersion: 2,
		PatchVersion: 3,
	}

	for i, tc := range testCases {
		tc.require(t, etype.Equals(tc.in), "test case %d failed", i)
	}

	// Two zero valued types should compare to true
	zero := &api.Type{}
	require.True(t, zero.Equals(&api.Type{}), "two zero valued types should compare to true")

	// Zero valued type should not equal the unspecified type
	require.False(t, zero.Equals(api.UnspecifiedType), "zero valued type should not equal unspecified type")

	// Nil types should equal each other
	var nilype *api.Type
	require.True(t, nilype.Equals(nil), "nil types should equal each other")
}

func TestTypeIsZero(t *testing.T) {
	testCases := []struct {
		in      *api.Type
		require require.BoolAssertionFunc
	}{
		{&api.Type{}, require.True},
		{api.UnspecifiedType, require.True},
		{&api.Type{Name: "TestType", MajorVersion: 0, MinorVersion: 0, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 1, MinorVersion: 0, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 0, MinorVersion: 2, PatchVersion: 0}, require.False},
		{&api.Type{Name: "", MajorVersion: 0, MinorVersion: 0, PatchVersion: 3}, require.False},
	}

	for i, tc := range testCases {
		tc.require(t, tc.in.IsZero(), "test case %d failed", i)
	}
}

func TestMkEvtHelper(t *testing.T) {
	event := mkevt("x/Xvi+2nnU8lfETEZ4C7YQ", "foo:bar,color:red", mimetype.ApplicationParquet, "TestEvent v1.2.3", "2023-10-04T08:17:22-05:00")
	require.Equal(t, event.Data, []byte{0xc7, 0xf5, 0xef, 0x8b, 0xed, 0xa7, 0x9d, 0x4f, 0x25, 0x7c, 0x44, 0xc4, 0x67, 0x80, 0xbb, 0x61}, "ensure event data is b64 decoded")
	require.Equal(t, event.Metadata, map[string]string{"foo": "bar", "color": "red"}, "metadata should be equal")
	require.Equal(t, event.Mimetype, mimetype.ApplicationParquet)
	require.Equal(t, event.Type.Name, "TestEvent")
	require.Equal(t, event.Type.MajorVersion, uint32(1))
	require.Equal(t, event.Type.MinorVersion, uint32(2))
	require.Equal(t, event.Type.PatchVersion, uint32(3))
	require.Equal(t, event.Created.AsTime(), time.Date(2023, 10, 4, 13, 17, 22, 0, time.UTC))
}

// Helper to quickly make events for testing purposes. Data should be base64 encoded
// data, kvs should be key:val,key:val pairs, etype should be a semvar for the type,
// e.g. Generic v1.2.3. and finally created should be an RFC3339 string or empty string
// to use a constant timestamp. If data or kvs is empty then empty byte array and empty
// metadata maps will be created. If etype is empty then the unspecified type is used.
func mkevt(data, kvs string, mime mimetype.MIME, etype, created string) *api.Event {
	e := &api.Event{
		Data:     make([]byte, 0),
		Metadata: make(map[string]string),
		Mimetype: mime,
		Type:     api.UnspecifiedType,
		Created:  timestamppb.New(time.Date(2023, 10, 11, 11, 14, 23, 0, time.UTC)),
	}

	if data != "" {
		var err error
		if e.Data, err = base64.RawStdEncoding.DecodeString(data); err != nil {
			e.Data = []byte(data)
		}
	}

	if kvs != "" {
		for _, pair := range strings.Split(kvs, ",") {
			parts := strings.Split(pair, ":")
			switch len(parts) {
			case 0:
				continue
			case 1:
				e.Metadata[parts[0]] = ""
			case 2:
				e.Metadata[parts[0]] = parts[1]
			default:
				panic("could not parse kvs")
			}
		}
	}

	if etype != "" {
		e.Type = &api.Type{}
		parts := strings.Split(etype, " ")
		if len(parts) == 2 {
			e.Type.Name = parts[0]

			semver := strings.Split(strings.TrimPrefix(parts[1], "v"), ".")
			if len(semver) == 3 {
				e.Type.MajorVersion = parseuint32(semver[0])
				e.Type.MinorVersion = parseuint32(semver[1])
				e.Type.PatchVersion = parseuint32(semver[2])
			} else {
				panic("could not parse etype version")
			}
		} else {
			panic("could not parse etype")
		}
	}

	if created != "" {
		if ts, err := time.Parse(time.RFC3339, created); err == nil {
			e.Created = timestamppb.New(ts)
		}
	}

	return e
}

func parseuint32(s string) uint32 {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(num)
}
