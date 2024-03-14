package ensign_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (s *serverTestSuite) TestTopicFilter() {
	require := s.Require()

	s.Run("Notfound", func() {
		s.store.UseError(store.TopicInfo, errors.ErrNotFound)
		_, err := s.srv.TopicFilter(ulid.MustParse("01HCZHJ1DP6W0WVQHXXRHVAMSH"))
		require.ErrorIs(err, errors.ErrNotFound)
	})

	s.Run("Empty", func() {
		s.store.OnTopicInfo = func(u ulid.ULID) (*api.TopicInfo, error) {
			return &api.TopicInfo{TopicId: u.Bytes(), Events: 0}, nil
		}
		s.store.OnLoadIndash = func(u ulid.ULID) iterator.IndashIterator {
			return store.NewIndashIterator(nil)
		}

		filter, err := s.srv.TopicFilter(ulid.MustParse("01HCZHJ1DP6W0WVQHXXRHVAMSH"))
		require.NoError(err, "expected filter to be returned")

		m, k := bloom.EstimateParameters(10000, 0.01)
		require.Equal(m, filter.Cap())
		require.Equal(k, filter.K())
	})

	s.Run("Small", func() {
		hashes := make([][]byte, 0, 5000)
		for i := 0; i < 5000; i++ {
			hash := make([]byte, 16)
			_, err := rand.Read(hash)
			require.NoError(err, "could not create random hash")
			hashes = append(hashes, hash)
		}

		s.store.OnTopicInfo = func(u ulid.ULID) (*api.TopicInfo, error) {
			return &api.TopicInfo{TopicId: u.Bytes(), Events: uint64(len(hashes))}, nil
		}
		s.store.OnLoadIndash = func(u ulid.ULID) iterator.IndashIterator {
			return store.NewIndashIterator(hashes)
		}

		filter, err := s.srv.TopicFilter(ulid.MustParse("01HCZHJ1DP6W0WVQHXXRHVAMSH"))
		require.NoError(err, "expected filter to be returned")

		m, k := bloom.EstimateParameters(10000, 0.01)
		require.Equal(m, filter.Cap())
		require.Equal(k, filter.K())

		for _, hash := range hashes {
			require.True(filter.Test(hash))
		}
	})

	s.Run("Large", func() {
		if testing.Short() {
			s.T().Skip("skipping large topic filter test")
			return
		}

		hashes := make([][]byte, 0, 50000)
		for i := 0; i < 50000; i++ {
			hash := make([]byte, 16)
			_, err := rand.Read(hash)
			require.NoError(err, "could not create random hash")
			hashes = append(hashes, hash)
		}

		s.store.OnTopicInfo = func(u ulid.ULID) (*api.TopicInfo, error) {
			return &api.TopicInfo{TopicId: u.Bytes(), Events: uint64(len(hashes))}, nil
		}
		s.store.OnLoadIndash = func(u ulid.ULID) iterator.IndashIterator {
			return store.NewIndashIterator(hashes)
		}

		filter, err := s.srv.TopicFilter(ulid.MustParse("01HCZHJ1DP6W0WVQHXXRHVAMSH"))
		require.NoError(err, "expected filter to be returned")

		m, k := bloom.EstimateParameters(100000, 0.01)
		require.Equal(m, filter.Cap())
		require.Equal(k, filter.K())

		for _, hash := range hashes {
			require.True(filter.Test(hash))
		}
	})
}

func (s *serverTestSuite) TestRehash() {
	require := s.Require()
	if testing.Short() {
		s.T().Skip("rehashing tests take a long time to run")
		return
	}

	// Load fixtures from disk
	topic, topicInfo, events, err := loadDuplicates()
	require.NoError(err, "could not load duplicates dataset from testdata/duplicates.pb.json")
	require.Equal(topic.Id, topicInfo.TopicId, "fixture verification failed")
	require.Equal(int(topic.Offset), len(events), "fixture verification failed")

	// Create an in-memory mock store for the tests
	indash := make(map[string]int, len(events))
	s.store.OnClearIndash = func(ulid.ULID) error {
		indash = nil
		indash = make(map[string]int, len(events))
		return nil
	}

	s.store.OnTopicInfo = func(topicID ulid.ULID) (*api.TopicInfo, error) {
		if bytes.Equal(topicInfo.TopicId, topicID.Bytes()) {
			return topicInfo, nil
		}
		return nil, errors.ErrNotFound
	}

	s.store.OnUnhash = func(_ ulid.ULID, hash []byte) (*api.EventWrapper, error) {
		key := base64.RawStdEncoding.EncodeToString(hash)
		return events[indash[key]], nil
	}

	s.store.OnIndash = func(_ ulid.ULID, hash []byte, eid rlid.RLID) error {
		key := base64.RawStdEncoding.EncodeToString(hash)
		indash[key] = int(eid.Sequence())
		return nil
	}

	s.store.OnInsert = func(e *api.EventWrapper) error {
		id := rlid.RLID(e.Id)
		offset := id.Sequence()
		events[offset] = e
		return nil
	}

	s.store.OnList = func(ulid.ULID) iterator.EventIterator { return store.NewEventIterator(events) }
	s.store.OnUpdateTopicInfo = func(ti *api.TopicInfo) error { topicInfo = ti; return nil }

	// TEST 1: DEDUPLICATION NONE --> DATAGRAM
	policy := &api.Deduplication{Strategy: api.Deduplication_DATAGRAM}
	err = s.srv.Rehash(context.Background(), ulid.ULID(topic.Id), policy.Normalize())
	require.NoError(err, "could not rehash from deduplication none to deduplication datagram")

	// TODO: make assertions about duplicates from dataset
	// TODO: why was only one duplicate found?
	require.Equal(uint64(0x1), topicInfo.Duplicates)

	s.store.Reset()
	// END TEST 1

	// TODO: test restoration and rehashing to other policies.

}

func loadDuplicates() (topic *api.Topic, info *api.TopicInfo, events []*api.EventWrapper, err error) {
	var data map[string]interface{}
	if data, err = loadJSONData("testdata/duplicates.pb.json"); err != nil {
		return nil, nil, nil, err
	}

	topic = &api.Topic{}
	if err = loadProtoJSON(data["topic"], topic); err != nil {
		return nil, nil, nil, err
	}

	info = &api.TopicInfo{}
	if err = loadProtoJSON(data["topic_info"], info); err != nil {
		return nil, nil, nil, err
	}

	if events, err = loadEvents(data["events"]); err != nil {
		return nil, nil, nil, err
	}

	return topic, info, events, nil
}

func loadJSONData(path string) (obj map[string]interface{}, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	obj = make(map[string]interface{})
	if err = json.NewDecoder(f).Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func loadProtoJSON(data interface{}, obj protoreflect.ProtoMessage) (err error) {
	var raw []byte
	if raw, err = json.Marshal(data.(map[string]interface{})); err != nil {
		return err
	}

	if err = protojson.Unmarshal(raw, obj); err != nil {
		return fmt.Errorf("could not unmarshal protobuf from json: %w", err)
	}
	return nil
}

func loadEvents(data interface{}) (events []*api.EventWrapper, err error) {
	items := data.([]interface{})
	events = make([]*api.EventWrapper, 0, len(items))

	for _, item := range items {
		event := &api.EventWrapper{}
		if err = loadProtoJSON(item, event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
