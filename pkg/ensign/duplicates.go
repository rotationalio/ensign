package ensign

import (
	"bytes"
	"context"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
)

const (
	filterFPRate       = 0.01
	filterMinSize uint = 10000
)

// TopicFilter loads a bloom filter with all of the event hashes for the events in the
// specified topic. The TopicInfo for the event is read to determine how many events are
// in the topic. The bloom filter is constructed as the larger of either 10k events or
// twice the number of events in the topic and with a false positive rate of 1%. The
// filter can be tested and modified as needed to detect duplicates.
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
			// to miss duplicates -- however, it could be possible to relax this.
			return nil, err
		}
		filter.Add(hash)
	}

	if err = iter.Error(); err != nil {
		return nil, err
	}
	return filter, nil
}

// Rehash clears the old event hashes and recomputes the hashes with the new policy.
// TODO: this method operates on a snapshot of the database and is not concurrency safe.
// TODO: rehash does not take offset policies into account.
func (s *Server) Rehash(ctx context.Context, topicID ulid.ULID, policy *api.Deduplication) (err error) {
	// Clear old hashes from the database.
	if err = s.data.ClearIndash(topicID); err != nil {
		return err
	}

	// Load topic info to build bloom filter
	// Will return not found if there is no associated topic.
	var info *api.TopicInfo
	if info, err = s.meta.TopicInfo(topicID); err != nil {
		return err
	}

	// Build the bloom filter for deduplication
	filter := bloom.NewWithEstimates(uint(info.Events), filterFPRate)
	policy = policy.Normalize()

	// Reset the topicInfo duplicate counts now that we've created the filter
	// NOTE: the next time the topic info gatherer is run, it will seek to the event ID
	// specified by this topic info so the count and sizes should not change.
	info.Duplicates = 0
	for _, etype := range info.Types {
		etype.Duplicates = 0
	}

	// Respect context cancellation before moving into iteration
	if err = ctx.Err(); err != nil {
		return err
	}

	// Iterate over all of the events in the database and re-compute hashes.
	iter := s.data.List(topicID)
	defer iter.Release()

deduplication:
	for iter.Next() {
		// Respect context cancellation and deadlines
		if err := ctx.Err(); err != nil {
			return err
		}

		event, err := iter.Event()
		if err != nil {
			return err
		}

		// If we've reached the end of the events specified by the topic info snapshot
		// then stop looping otherwise we may inject a consistency issue
		if bytes.Equal(event.Id, info.EventOffsetId) {
			break deduplication
		}

		// If none then skip over the deduplication checking step.
		if policy.Strategy != api.Deduplication_NONE {
			// Compute the hash of the event given the deduplication policy
			hash, err := event.Hash(policy)
			if err != nil {
				return err
			}

			// Check if the event is a duplicate of another event already
			if filter.TestOrAdd(hash) {
				// Load the identified duplicate, verify that it is a duplicate
				var target *api.EventWrapper
				if target, err = s.data.Unhash(topicID, hash); err != nil {
					return err
				}

				var isDuplicate bool
				if isDuplicate, err = event.Duplicates(target, policy); err != nil {
					return err
				}

				if isDuplicate {
					// Mark the event as a duplicate and save back to database
					// TODO: handle offset -- e.g. is the target the duplicate or the event?
					if err = event.DuplicateOf(target, policy); err != nil {
						return err
					}

					// Save the duplicate back to the database
					if err = s.data.Insert(event); err != nil {
						return err
					}

					// Update the duplicate counts on the topic info
					info.Duplicates++
					if e, err := event.Unwrap(); err == nil {
						etype := info.FindEventTypeInfo(e.ResolveType(), e.Mimetype)
						etype.Duplicates++
					}

					// NOTE: We assume that everything from this point on in the loop that
					// the event is not a duplicate and should be treated as original.
					continue deduplication
				}
			}

			// If the topic is not a duplicate store the hash in the database.
			if err := s.data.Indash(topicID, hash, rlid.RLID(event.Id)); err != nil {
				return err
			}
		}

		// From this point on, we assume that the event is not a duplicate.
		// Ensure the topic is not marked as a duplicate; if it is, mark it as a non-duplicate.
		if event.IsDuplicate {
			var orig *api.EventWrapper
			if orig, err = s.data.Retrieve(topicID, rlid.RLID(event.DuplicateId)); err != nil {
				return err
			}

			// TODO: do we need to handle the offset here?
			if err = event.DuplicateFrom(orig); err != nil {
				return err
			}

			// Save the duplicate back to the database
			if err = s.data.Insert(event); err != nil {
				return err
			}
		}
	}

	if err = iter.Error(); err != nil {
		return err
	}

	// Save the topic info back to disk so that it can be carried on later.
	if err = s.meta.UpdateTopicInfo(info); err != nil {
		return err
	}
	return nil
}
