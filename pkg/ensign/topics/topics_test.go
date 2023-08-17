package topics_test

import (
	"encoding/base64"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/topics"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func TestNameGroup(t *testing.T) {
	group := &topics.NameGroup{}
	name1 := "testing.123"
	tids1 := "01H0TF3S4ES3VSWSGY3RH4E0H5"
	hash1 := base64.RawURLEncoding.EncodeToString(api.TopicNameHash(name1))

	require.False(t, group.Contains(name1))
	require.False(t, group.Contains(tids1))
	require.False(t, group.Contains(hash1))

	// Add name and ID
	require.NoError(t, group.Add(name1, ulid.MustParse(tids1)))
	require.True(t, group.Contains(name1))
	require.True(t, group.Contains(tids1))
	require.True(t, group.Contains(hash1))

	name2 := "simple-topic.testing"
	tids2 := "01H0TFBZ90B2JV7YG06Z66PTXN"
	hash2 := base64.RawURLEncoding.EncodeToString(api.TopicNameHash(name2))

	require.False(t, group.Contains(name2))
	require.False(t, group.Contains(tids2))
	require.False(t, group.Contains(hash2))

	// Add topic with name and ID
	group.AddTopic(&api.Topic{Id: ulid.MustParse(tids2).Bytes(), Name: name2})
	require.True(t, group.Contains(name2))
	require.True(t, group.Contains(tids2))
	require.True(t, group.Contains(hash2))

	expected := map[string][]byte{
		name1: ulid.MustParse(tids1).Bytes(),
		name2: ulid.MustParse(tids2).Bytes(),
	}

	topics := group.TopicMap()
	require.Len(t, topics, group.Length())
	require.Equal(t, expected, topics)

}

func TestNameGroupFilter(t *testing.T) {
	fixtures := []struct {
		name    string
		topicID ulid.ULID
	}{
		{"testing.testapp.alerts", ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD")},
		{"testing.testapp.orders", ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0")},
		{"testing.testapp.shipments", ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")},
		{"mock.mockapp.feed", ulid.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ")},
		{"mock.mockapp.post", ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG")},
		{"mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM")},
		{"snake-case", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT")},
	}

	group := &topics.NameGroup{}
	for i, fixture := range fixtures {
		err := group.Add(fixture.name, fixture.topicID)
		require.NoError(t, err, "could not add fixture %d to group", i)
	}

	// Filter nothing should return an empty filter
	bravo := group.Filter()
	require.Equal(t, 0, bravo.Length())

	// Filtering for things that are not in the group should return an empty group
	bravo = group.Filter("foo", "bar", "01H7GJDB3ETF5F06CHNEAABM4F", "vdKucRbIWkV0olXguvjWrw")
	require.Equal(t, 0, bravo.Length())

	// Should be able to filter with topic names
	bravo = group.Filter("mock.mockapp.feed", "snake-case", "testing.testapp.shipments")
	require.Equal(t, 3, bravo.Length())
	require.True(t, bravo.ContainsTopicName("snake-case"))
	require.False(t, bravo.ContainsTopicName("mock1"))

	// Should be able to filter with topic IDs
	bravo = group.Filter("01GTSN1WF5BA0XCPT6ES64JVGQ", "01GTSMQ3V8ASAPNCFEN378T8RD", "01GTSMMC152Q95RD4TNYDFJGHT", "01GTSN1139JMK1PS5A524FXWAZ")
	require.Equal(t, 4, bravo.Length())
	require.True(t, bravo.ContainsTopicName("testing.testapp.shipments"))
	require.False(t, bravo.ContainsTopicName("mock.mockapp.post"))

	// Should be able to filter with name hashes
	bravo = group.Filter("aGdsSRMLeh-urMLeQu2XRQ", "lVEWQE8IkIM1MokVF0V1mw")
	require.Equal(t, 2, bravo.Length())
	require.True(t, bravo.ContainsTopicName("testing.testapp.shipments"))
	require.False(t, bravo.ContainsTopicName("snake-case"))

	// Should be able to filter with multiple types
	bravo = group.Filter("aGdsSRMLeh-urMLeQu2XRQ", "01GTSMQ3V8ASAPNCFEN378T8RD", "testing.testapp.orders")
	require.Equal(t, 3, bravo.Length())
	require.True(t, bravo.ContainsTopicName("testing.testapp.orders"))
	require.False(t, bravo.ContainsTopicName("mock.mockapp.feed"))

	// Should ignore duplicates when filtering
	bravo = group.Filter("aGdsSRMLeh-urMLeQu2XRQ", "testing.testapp.shipments", "01GTSN1139JMK1PS5A524FXWAZ")
	require.Equal(t, 1, bravo.Length())
	require.True(t, bravo.ContainsTopicName("testing.testapp.shipments"))
	require.False(t, bravo.ContainsTopicName("mock.mockapp.post"))

	// Should ignore things that aren't in the input list when filtering
	bravo = group.Filter("aGdsSRMLeh-urMLeQu2XRQ", "", "01GTSMQ3V8ASAPNCFEN378T8RD", "foo", "testing.testapp.orders", "01H7GKJ581VE32QWQGE5JXRP8K")
	require.Equal(t, 3, bravo.Length())
	require.True(t, bravo.ContainsTopicName("testing.testapp.orders"))
	require.False(t, bravo.ContainsTopicName("mock.mockapp.feed"))

}

func TestNameGroupLookup(t *testing.T) {
	fixtures := []struct {
		name    string
		topicID ulid.ULID
	}{
		{"testing.testapp.alerts", ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD")},
		{"testing.testapp.orders", ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0")},
		{"testing.testapp.shipments", ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")},
		{"mock.mockapp.feed", ulid.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ")},
		{"mock.mockapp.post", ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG")},
		{"mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM")},
		{"snake-case", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT")},
	}

	group := &topics.NameGroup{}
	for i, fixture := range fixtures {
		err := group.Add(fixture.name, fixture.topicID)
		require.NoError(t, err, "could not add fixture %d to group", i)
	}

	t.Run("Generic", func(t *testing.T) {
		testCases := []struct {
			input   string
			name    string
			topicID ulid.ULID
			ok      require.BoolAssertionFunc
		}{
			{"", "", ulids.Null, require.False},
			{"00000000000000000000000000", "", ulids.Null, require.False},
			{"notinthegroup", "", ulids.Null, require.False},
			{"01H7GH6EA425FQTQANZT21DZFN", "", ulids.Null, require.False},
			{"gZ_jm2rUedIoDI6RUt4pgg", "", ulids.Null, require.False},
			{"01GTSN1139JMK1PS5A524FXWAZ", "testing.testapp.shipments", ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), require.True},
			{"snake-case", "snake-case", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"), require.True},
			{"lVEWQE8IkIM1MokVF0V1mw", "mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM"), require.True},
		}

		for i, tc := range testCases {
			name, topicID, ok := group.Lookup(tc.input)
			tc.ok(t, ok, "expected retrieval for test case %d", i)
			require.Equal(t, tc.name, name, "test case %d name comparision failed", i)
			require.Equal(t, tc.topicID, topicID, "test case %d topic id comparision failed", i)
		}
	})

	t.Run("TopicID", func(t *testing.T) {
		testCases := []struct {
			input ulid.ULID
			name  string
			ok    require.BoolAssertionFunc
		}{
			{ulids.Null, "", require.False},
			{ulid.MustParse("01H7GH6EA425FQTQANZT21DZFN"), "", require.False},
			{ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), "testing.testapp.shipments", require.True},
		}

		for i, tc := range testCases {
			name, ok := group.LookupTopicID(tc.input)
			tc.ok(t, ok, "expected retrieval for test case %d", i)
			require.Equal(t, tc.name, name, "test case %d name comparision failed", i)
		}
	})

	t.Run("Name", func(t *testing.T) {
		testCases := []struct {
			input   string
			topicID ulid.ULID
			ok      require.BoolAssertionFunc
		}{
			{"", ulids.Null, require.False},
			{"00000000000000000000000000", ulids.Null, require.False},
			{"notinthegroup", ulids.Null, require.False},
			{"01H7GH6EA425FQTQANZT21DZFN", ulids.Null, require.False},
			{"snake-case", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"), require.True},
			{"mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM"), require.True},
		}

		for i, tc := range testCases {
			topicID, ok := group.LookupTopicName(tc.input)
			tc.ok(t, ok, "expected retrieval for test case %d", i)
			require.Equal(t, tc.topicID, topicID, "test case %d topic id comparision failed", i)
		}
	})

	t.Run("Hash", func(t *testing.T) {
		testCases := []struct {
			input   string
			name    string
			topicID ulid.ULID
			ok      require.BoolAssertionFunc
		}{
			{"gZ_jm2rUedIoDI6RUt4pgg", "", ulids.Null, require.False},
			{"lVEWQE8IkIM1MokVF0V1mw", "mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM"), require.True},
		}

		for i, tc := range testCases {
			hash, err := base64.RawURLEncoding.DecodeString(tc.input)
			require.NoError(t, err, "could not decode input for test case %d", i)

			name, topicID, ok := group.LookupTopicHash(hash)
			tc.ok(t, ok, "expected retrieval for test case %d", i)
			require.Equal(t, tc.name, name, "test case %d name comparision failed", i)
			require.Equal(t, tc.topicID, topicID, "test case %d topic id comparision failed", i)
		}

		_, _, ok := group.LookupTopicHash(nil)
		require.False(t, ok, "expected nil to return nothing")
		_, _, ok = group.LookupTopicHash([]byte{})
		require.False(t, ok, "expected empty byte slice to return nothing")
	})
}

func TestNameGroupFilterID(t *testing.T) {
	fixtures := []struct {
		name string
		id   ulid.ULID
	}{
		{"testing.1", ulid.MustParse("01H0TB6F9MAMMQ8T9DZAZGQ5RH")},
		{"testing.2", ulid.MustParse("01H0TF3S4ES3VSWSGY3RH4E0H5")},
		{"testing.3", ulid.MustParse("01H0TFBZ90B2JV7YG06Z66PTXN")},
		{"testing.4", ulid.MustParse("01H0TFNFT9439MQR3EF11X6ES1")},
		{"testing.5", ulid.MustParse("01H0TFP3Q2WRE2QZXBDNW4HR2Z")},
		{"testing.6", ulid.MustParse("01H0TFP8S2QQ38GZX17X5RGB3J")},
		{"testing.7", ulid.MustParse("01H0TFPDTKHZP4R416GG1C61GP")},
		{"testing.8", ulid.MustParse("01H0TFPJVD2NCBX2AHP310MN0S")},
	}

	group := &topics.NameGroup{}
	for _, fixture := range fixtures {
		require.NoError(t, group.Add(fixture.name, fixture.id))
	}

	filtered := group.FilterTopicID(fixtures[1].id, fixtures[3].id, fixtures[5].id, fixtures[7].id)
	require.NotSame(t, group, filtered)
	require.Equal(t, 4, filtered.Length())

	for i, fixture := range fixtures {
		require.True(t, group.Contains(fixture.name))
		require.True(t, group.Contains(fixture.id.String()))

		if i%2 == 0 {
			require.False(t, filtered.Contains(fixture.name))
			require.False(t, filtered.Contains(fixture.id.String()))
		} else {
			require.True(t, filtered.Contains(fixture.name))
			require.True(t, filtered.Contains(fixture.id.String()))
		}
	}
}

func TestNameGroupFilterName(t *testing.T) {
	fixtures := []struct {
		name string
		id   ulid.ULID
	}{
		{"testing.1", ulid.MustParse("01H0TB6F9MAMMQ8T9DZAZGQ5RH")},
		{"testing.2", ulid.MustParse("01H0TF3S4ES3VSWSGY3RH4E0H5")},
		{"testing.3", ulid.MustParse("01H0TFBZ90B2JV7YG06Z66PTXN")},
		{"testing.4", ulid.MustParse("01H0TFNFT9439MQR3EF11X6ES1")},
		{"testing.5", ulid.MustParse("01H0TFP3Q2WRE2QZXBDNW4HR2Z")},
		{"testing.6", ulid.MustParse("01H0TFP8S2QQ38GZX17X5RGB3J")},
		{"testing.7", ulid.MustParse("01H0TFPDTKHZP4R416GG1C61GP")},
		{"testing.8", ulid.MustParse("01H0TFPJVD2NCBX2AHP310MN0S")},
	}

	group := &topics.NameGroup{}
	for _, fixture := range fixtures {
		require.NoError(t, group.Add(fixture.name, fixture.id))
	}

	filtered := group.FilterTopicName(fixtures[1].name, fixtures[3].name, fixtures[5].name, fixtures[7].name)
	require.NotSame(t, group, filtered)
	require.Equal(t, 4, filtered.Length())

	for i, fixture := range fixtures {
		require.True(t, group.Contains(fixture.name))
		require.True(t, group.Contains(fixture.id.String()))

		if i%2 == 0 {
			require.False(t, filtered.Contains(fixture.name))
			require.False(t, filtered.Contains(fixture.id.String()))
		} else {
			require.True(t, filtered.Contains(fixture.name))
			require.True(t, filtered.Contains(fixture.id.String()))
		}
	}
}

func TestEmptyNameGroup(t *testing.T) {
	group := &topics.NameGroup{}
	require.Equal(t, 0, group.Length())

	err := group.Add("", ulids.Null)
	require.ErrorIs(t, err, topics.ErrEmptyReference)
	require.Equal(t, 0, group.Length())

	err = group.Add("", ulid.MustParse("01H78XH126J1XHRR2CAQBBT7RC"))
	require.ErrorIs(t, err, topics.ErrEmptyReference)
	require.Equal(t, 0, group.Length())

	err = group.Add("example", ulids.Null)
	require.ErrorIs(t, err, topics.ErrEmptyReference)
	require.Equal(t, 0, group.Length())
}

func TestAddTopicIDTwice(t *testing.T) {
	topicID := ulid.MustParse("01H78XH126J1XHRR2CAQBBT7RC")
	group := &topics.NameGroup{}

	require.NoError(t, group.Add("foo", topicID))
	require.Equal(t, 1, group.Length())

	err := group.Add("bar", topicID)
	require.ErrorIs(t, err, topics.ErrAlreadyExists)
	require.Equal(t, 1, group.Length())
}

func TestAddTopicNameTwice(t *testing.T) {
	group := &topics.NameGroup{}

	require.NoError(t, group.Add("foo", ulid.MustParse("01H78XH126J1XHRR2CAQBBT7RC")))
	require.Equal(t, 1, group.Length())

	err := group.Add("foo", ulid.MustParse("01H78XT88RRYKX9SFQNN47B7WK"))
	require.ErrorIs(t, err, topics.ErrAlreadyExists)
	require.Equal(t, 1, group.Length())
}

func TestTopicIDs(t *testing.T) {
	fixtures := []struct {
		name    string
		topicID ulid.ULID
	}{
		{"testing.testapp.alerts", ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD")},
		{"testing.testapp.orders", ulid.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0")},
		{"testing.testapp.shipments", ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")},
		{"mock.mockapp.feed", ulid.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ")},
		{"mock.mockapp.post", ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG")},
		{"mock1", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM")},
		{"snake-case", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT")},
	}

	group := &topics.NameGroup{}
	for i, fixture := range fixtures {
		err := group.Add(fixture.name, fixture.topicID)
		require.NoError(t, err, "could not add fixture %d to group", i)
	}

	topicIDs := group.TopicIDs()
	require.Len(t, topicIDs, len(fixtures))
	for _, fixture := range fixtures {
		require.Contains(t, topicIDs, fixture.topicID)
	}
}

func TestNameHash(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "AAAAAAAAAAAAAAAAAAAAAA"},
		{"mock1", "lVEWQE8IkIM1MokVF0V1mw"},
		{"testing.testapp.shipments", "aGdsSRMLeh-urMLeQu2XRQ"},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, topics.NameHash(tc.input))
	}
}
