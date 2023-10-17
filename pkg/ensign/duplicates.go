package ensign

import (
	"context"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/cenkalti/backoff/v4"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/utils/radish"
)

const (
	filterFPRate       = 0.01
	filterMinSize uint = 10000
)

func (s *Server) TopicFilter(topicID ulid.ULID) (_ *bloom.BloomFilter, err error) {
	// Load the topic info to determine the bloom filter size.
	var info *api.TopicInfo
	if info, err = s.meta.TopicInfo(topicID); err != nil {
		return nil, err
	}

	// The filter size should be the larger of twice the number of events in the topic
	// or the minimum filter size (10k hashes by default).
	filterSize := filterMinSize
	if uint(info.Events*2) > filterSize {
		filterSize = uint(info.Events * 2)
	}

	// Create the bloom filter with index hashes from the database.
	filter := bloom.NewWithEstimates(filterSize, filterFPRate)
	iter := s.data.LoadIndash(topicID)
	defer iter.Release()

	for iter.Next() {
		var hash []byte
		if hash, err = iter.Hash(); err != nil {
			// NOTE: we are not skipping bad hashes because this would make it possible
			// to miss duplicates -- however, it would be possible to relax this.
			return nil, err
		}
		filter.Add(hash)
	}

	if err = iter.Error(); err != nil {
		return nil, err
	}
	return filter, nil
}

func (s *Server) Rehash(ctx context.Context, topicID ulid.ULID, policy *api.Deduplication) error {
	// TODO: clear old hashes from the database.

	iter := s.data.List(topicID)
	defer iter.Release()

	for iter.Next() {
		// Respect context cancellation and deadlines
		if err := ctx.Err(); err != nil {
			return err
		}

		event, err := iter.Event()
		if err != nil {
			return err
		}

		hash, err := event.Hash(policy)
		if err != nil {
			return err
		}

		// TODO: check if the hash is a duplicate of another event already.
		if err := s.data.Indash(topicID, hash, rlid.RLID(event.Id)); err != nil {
			return err
		}
	}

	return iter.Error()
}

func (s *Server) QueueRehash(topicID ulid.ULID, policy *api.Deduplication) {
	s.tasks.Queue(radish.TaskFunc(func(ctx context.Context) error {
		return s.Rehash(ctx, topicID, policy)
	}), radish.WithErrorf("could not complete rehash of %s", topicID),
		radish.WithRetries(1),
		radish.WithBackoff(backoff.NewConstantBackOff(1*time.Minute)),
		radish.WithTimeout(30*time.Minute),
	)
}
