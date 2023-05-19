package topics_test

import (
	"encoding/base64"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/topics"
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
	group.Add(name1, ulid.MustParse(tids1))
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
		group.Add(fixture.name, fixture.id)
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
		group.Add(fixture.name, fixture.id)
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
