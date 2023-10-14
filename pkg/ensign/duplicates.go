package ensign

import (
	"context"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/utils/radish"
)

func (s *Server) Rehash(ctx context.Context, topicID ulid.ULID, policy *api.Deduplication) error {
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
	s.tasks.Queue(radish.Func(func(ctx context.Context) error {
		return s.Rehash(ctx, topicID, policy)
	}))
}
