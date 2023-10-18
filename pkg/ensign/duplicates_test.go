package ensign_test

import (
	"crypto/rand"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/iterator"
	store "github.com/rotationalio/ensign/pkg/ensign/store/mock"
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
