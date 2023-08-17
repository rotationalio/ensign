package mock_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestIterator(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		emptyTest := func(init makeIterator) func(t *testing.T) {
			return func(t *testing.T) {
				it := init()
				require.False(t, it.Next())
				require.False(t, it.Prev())
				require.NoError(t, it.Error(), "expected no error after calls to next and prev")

				// Calling key and value panic for empty iterators because Next returned false
				require.Panics(t, func() { it.Key() })
				require.Panics(t, func() { it.Value() })

				// If the iterator is a paginator, check creating an empty page without error
				if pg, ok := it.(iterator.Paginator); ok {
					// Should create an empty page without an error
					out, npt, err := pg.Page(&api.PageInfo{})
					require.NoError(t, err)
					require.Empty(t, npt)
					require.Empty(t, out)
				}
			}
		}

		t.Run("Event", emptyTest(func() Iterator { return makeEmptyEventIterator() }))
		t.Run("Topic", emptyTest(func() Iterator { return makeEmptyTopicIterator() }))
		t.Run("TopicName", emptyTest(func() Iterator { return makeEmptyTopicNamesIterator() }))
	})

	t.Run("Release", func(t *testing.T) {
		releaseTest := func(init makeIterator) func(t *testing.T) {
			return func(t *testing.T) {
				it := init()
				it.Release()
				require.False(t, it.Next())
				require.ErrorIs(t, it.Error(), leveldb.ErrIterReleased)

				it = init()
				it.Release()
				require.False(t, it.Prev())
				require.ErrorIs(t, it.Error(), leveldb.ErrIterReleased)

				it = init()
				it.Release()
				require.Nil(t, it.Key())
				require.ErrorIs(t, it.Error(), leveldb.ErrIterReleased)

				it = init()
				it.Release()
				_, err := it.Object()
				require.ErrorIs(t, err, leveldb.ErrIterReleased)
				require.ErrorIs(t, it.Error(), leveldb.ErrIterReleased)

				it = init()
				it.Release()
				if pg, ok := it.(iterator.Paginator); ok {
					_, _, err = pg.Page(&api.PageInfo{})
					require.ErrorIs(t, err, leveldb.ErrIterReleased)
					require.ErrorIs(t, it.Error(), leveldb.ErrIterReleased)
				}
			}
		}

		t.Run("Event", releaseTest(func() Iterator { return makeEmptyEventIterator() }))
		t.Run("Topic", releaseTest(func() Iterator { return makeEmptyTopicIterator() }))
		t.Run("TopicName", releaseTest(func() Iterator { return makeEmptyTopicNamesIterator() }))
	})

	t.Run("Error", func(t *testing.T) {
		errorTest := func(init makeIterator) func(t *testing.T) {
			return func(t *testing.T) {
				it := init()
				require.ErrorIs(t, it.Error(), errTestIterator)

				it.Release()
				require.ErrorIs(t, it.Error(), errTestIterator)

				require.False(t, it.Next())
				require.False(t, it.Prev())
				require.Nil(t, it.Key())
				require.ErrorIs(t, it.Error(), errTestIterator)

				_, err := it.Object()
				require.ErrorIs(t, err, errTestIterator)

				if pg, ok := it.(iterator.Paginator); ok {
					_, _, err = pg.Page(&api.PageInfo{})
					require.ErrorIs(t, err, errTestIterator)
				}
			}
		}

		t.Run("Event", errorTest(func() Iterator { return makeEventErrorIterator() }))
		t.Run("Topic", errorTest(func() Iterator { return makeTopicErrorIterator() }))
		t.Run("TopicName", errorTest(func() Iterator { return makeTopicNamesErrorIterator() }))
	})
}

func TestEventIterator(t *testing.T) {
	fixture, err := mock.EventListFixture("testdata/events.pb.json")
	require.NoError(t, err, "could not load testdata/events.pb.json")

	it := mock.NewEventIterator(fixture)

	eventIDs := make([][]byte, 0, len(fixture))
	for it.Next() {
		event, err := it.Event()
		require.NoError(t, err)
		require.Len(t, event.Id, 10)
		eventIDs = append(eventIDs, event.Id)

		// TODO: test key is correct
	}
	require.Len(t, eventIDs, len(fixture))

}

func TestTopicIterator(t *testing.T) {
	fixture, err := mock.TopicListFixture("testdata/topics.pb.json")
	require.NoError(t, err, "could not load testdata/topics.pb.json")

	it := mock.NewTopicIterator(fixture)

	topics := make([]string, 0, len(fixture))
	for it.Next() {
		topic, err := it.Topic()
		require.NoError(t, err)
		topics = append(topics, topic.Name)

		key := it.Key()
		expected := meta.TopicKey(topic)
		require.True(t, bytes.Equal(key, expected[:]))
	}
	require.Len(t, topics, len(fixture))

}

func TestTopicPagination(t *testing.T) {
	fixture, err := mock.TopicListFixture("testdata/topics.pb.json")
	require.NoError(t, err, "could not load testdata/topics.pb.json")

	t.Run("All", func(t *testing.T) {
		it := mock.NewTopicIterator(fixture)
		page, err := it.NextPage(&api.PageInfo{PageSize: 100})
		require.NoError(t, err)
		require.Len(t, page.Topics, len(fixture))
		require.Empty(t, page.NextPageToken)
	})

	t.Run("Split", func(t *testing.T) {
		it := mock.NewTopicIterator(fixture)
		nPages, nTopics := 0, 0
		page := &api.PageInfo{PageSize: 3}

		// Maximum number of iterations is 100 to prevent test timeout
		for i := 0; i < 100; i++ {
			topics, err := it.NextPage(page)
			require.NoError(t, err)
			page.NextPageToken = topics.NextPageToken

			nPages++
			nTopics += len(topics.Topics)

			if page.NextPageToken == "" {
				break
			}
		}

		require.Equal(t, len(fixture), nTopics)
		require.Equal(t, (len(fixture)/3)+1, nPages)
	})

}

func TestTopicNamesIterator(t *testing.T) {
	fixture, err := mock.TopicNamesListFixture("testdata/topic_names.pb.json")
	require.NoError(t, err, "could not load testdata/topic_names.pb.json")

	it := mock.NewTopicNamesIterator(fixture)

	topics := make([]string, 0, len(fixture))
	for it.Next() {
		topic, err := it.TopicName()
		require.NoError(t, err)
		topics = append(topics, topic.Name)

		key := it.Key()
		expected := meta.TopicNameKey(&api.Topic{ProjectId: ulid.MustParse(topic.ProjectId).Bytes(), Name: topic.Name})
		require.True(t, bytes.Equal(key, expected[:]))
	}
	require.Len(t, topics, len(fixture))

}

func TestTopicNamesPagination(t *testing.T) {
	fixture, err := mock.TopicNamesListFixture("testdata/topic_names.pb.json")
	require.NoError(t, err, "could not load testdata/topic_names.pb.json")

	t.Run("All", func(t *testing.T) {
		it := mock.NewTopicNamesIterator(fixture)
		page, err := it.NextPage(&api.PageInfo{PageSize: 100})
		require.NoError(t, err)
		require.Len(t, page.TopicNames, len(fixture))
		require.Empty(t, page.NextPageToken)
	})

	t.Run("Split", func(t *testing.T) {
		it := mock.NewTopicNamesIterator(fixture)
		nPages, nTopics := 0, 0
		page := &api.PageInfo{PageSize: 3}

		// Maximum number of iterations is 100 to prevent test timeout
		for i := 0; i < 100; i++ {
			topics, err := it.NextPage(page)
			require.NoError(t, err)
			page.NextPageToken = topics.NextPageToken

			nPages++
			nTopics += len(topics.TopicNames)

			if page.NextPageToken == "" {
				break
			}
		}

		require.Equal(t, len(fixture), nTopics)
		require.Equal(t, (len(fixture)/3)+1, nPages)
	})

}

type Iterator interface {
	iterator.Iterator
	iterator.ObjectIterator
}

var errTestIterator = errors.New("this is a test iterator error")

type makeIterator func() Iterator

func makeEmptyEventIterator() *mock.EventIterator {
	return mock.NewEventIterator(nil)
}

func makeEmptyTopicIterator() *mock.TopicIterator {
	return mock.NewTopicIterator(nil)
}

func makeEmptyTopicNamesIterator() *mock.TopicNamesIterator {
	return mock.NewTopicNamesIterator(nil)
}

func makeEventErrorIterator() *mock.EventIterator {
	return mock.NewEventErrorIterator(errTestIterator)
}

func makeTopicErrorIterator() *mock.TopicIterator {
	return mock.NewTopicErrorIterator(errTestIterator)
}

func makeTopicNamesErrorIterator() *mock.TopicNamesIterator {
	return mock.NewTopicNamesErrorIterator(errTestIterator)
}
