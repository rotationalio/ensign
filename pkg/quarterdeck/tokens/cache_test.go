package tokens_test

import (
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *TokenTestSuite) TestCache() {
	require := s.Require()

	// Create a new cache
	cache := tokens.NewCache(10)

	// Get should return an error when the key is not in the cache
	_, err := cache.Get(ulids.New(), ulids.New())
	require.ErrorIs(err, tokens.ErrCacheMiss, "expected error when key is not in the cache")

	// Create an expired token string
	expiredTokens, err := tokens.New(s.expiredConf)
	require.NoError(err, "could not create expired token manager")
	claims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}
	expired, _, err := expiredTokens.CreateTokenPair(claims)
	require.NoError(err, "could not create expired token")

	// Add the expired token to the cache
	userID := ulids.New()
	projectID := ulids.New()
	require.NoError(cache.Add(userID, projectID, expired), "could not add token to cache")

	// Attempting to retrieve an expired token should return an error
	_, err = cache.Get(userID, projectID)
	require.ErrorIs(err, tokens.ErrCacheExpired, "expected error when token is expired")

	// Create a valid token string
	validTokens, err := tokens.New(s.conf)
	require.NoError(err, "could not create valid token manager")
	valid, _, err := validTokens.CreateTokenPair(claims)
	require.NoError(err, "could not create valid token")

	// Add the valid token to the cache, should be able to retrieve it
	require.NoError(cache.Add(userID, projectID, valid), "could not add token to cache")
	tks, err := cache.Get(userID, projectID)
	require.NoError(err, "could not retrieve token from cache")
	require.Equal(valid, tks, "token retrieved from cache does not match")

	// Add should return an error when given an unparseable token
	err = cache.Add(userID, projectID, "not a token")
	require.EqualError(err, "token contains an invalid number of segments", "expected error when given an invalid token")

	// Add tokens to exceed the cache size and trigger eviction
	userIDs := make([]ulid.ULID, 0)
	projectIDs := make([]ulid.ULID, 0)
	accessTokens := make([]string, 0)

	for i := 0; i < 5; i++ {
		expired, _, err := expiredTokens.CreateTokenPair(claims)
		require.NoError(err, "could not create expired token")
		userID := ulids.New()
		projectID := ulids.New()
		userIDs = append(userIDs, userID)
		projectIDs = append(projectIDs, projectID)
		accessTokens = append(accessTokens, expired)
		err = cache.Add(userID, projectID, expired)
		require.NoError(err, "could not add token to cache")
	}

	for i := 0; i < 5; i++ {
		valid, _, err := validTokens.CreateTokenPair(claims)
		require.NoError(err, "could not create valid token")
		userID := ulids.New()
		projectID := ulids.New()
		userIDs = append(userIDs, userID)
		projectIDs = append(projectIDs, projectID)
		accessTokens = append(accessTokens, valid)
		err = cache.Add(userID, projectID, valid)
		require.NoError(err, "could not add token to cache")
	}

	// First token should no longer be in the cache
	_, err = cache.Get(userID, projectID)
	require.ErrorIs(err, tokens.ErrCacheMiss, "expected first token to be evicted from cache")

	// Attempt to retrieve an expired token, should remove that and all previous tokens
	_, err = cache.Get(userIDs[3], projectIDs[3])
	require.ErrorIs(err, tokens.ErrCacheExpired, "expected expired token to be evicted from cache")
	require.Equal(6, cache.Size(), "expected cache to have 6 items after evicting expired tokens")

	// The last expired token should also not be retrievable
	_, err = cache.Get(userIDs[4], projectIDs[4])
	require.ErrorIs(err, tokens.ErrCacheExpired, "expected expired token to be evicted from cache")

	// All other tokens should still be in the cache
	for i := 5; i < 10; i++ {
		tks, err := cache.Get(userIDs[i], projectIDs[i])
		require.NoError(err, "could not retrieve token from cache")
		require.Equal(accessTokens[i], tks, "wrong token retrieved from cache")
	}

	// Remove should not panic when the key is not in the cache
	cache.Remove(ulids.New(), ulids.New())

	// Remove the rest of the tokens from the cache
	for i := 5; i < 10; i++ {
		cache.Remove(userIDs[i], projectIDs[i])
		_, err = cache.Get(userIDs[i], projectIDs[i])
		require.ErrorIs(err, tokens.ErrCacheMiss, "token should no longer be in the cache after removal")
	}

	// Cache should be empty
	require.Equal(0, cache.Size(), "expected cache to be empty")
}

func (s *TokenTestSuite) TestCacheConcurrency() {
	require := s.Require()

	// Create a new cache
	cache := tokens.NewCache(100)

	// Create a valid token string
	validTokens, err := tokens.New(s.conf)
	require.NoError(err, "could not create valid token manager")
	claims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Generate some user IDs, project IDs, and tokens
	userIDs := make([]ulid.ULID, 0)
	projectIDs := make([]ulid.ULID, 0)
	accessTokens := make([]string, 0)
	for i := 0; i < 100; i++ {
		userIDs = append(userIDs, ulids.New())
		projectIDs = append(projectIDs, ulids.New())
		valid, _, err := validTokens.CreateTokenPair(claims)
		require.NoError(err, "could not create valid token")
		accessTokens = append(accessTokens, valid)
	}

	// Go routine to add tokens to the cache
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			require.NoError(cache.Add(userIDs[i], projectIDs[i], accessTokens[i]), "could not add token to cache")
		}
	}()

	// Go routine to add more tokens to the cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 50; i < 100; i++ {
			require.NoError(cache.Add(userIDs[i], projectIDs[i], accessTokens[i]), "could not add token to cache")
		}
	}()

	// Wait for the go routines to finish
	wg.Wait()

	// Go routine to retrieve tokens from the cache
	retrieve := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			tks, err := cache.Get(userIDs[i], projectIDs[i])
			require.NoError(err, "could not retrieve token from cache")
			require.Equal(accessTokens[i], tks, "wrong token retrieved from cache")
		}
	}

	// Test retrieving tokens from the cache concurrently
	wg.Add(2)
	go retrieve()
	go retrieve()
	wg.Wait()

	// Go routine to add more tokens to the cache
	addTokens := func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			claims := &tokens.Claims{
				Name:  "Leopold Wentzel",
				Email: "leopold.wentzel@gmail.com",
			}

			valid, _, err := validTokens.CreateTokenPair(claims)
			require.NoError(err, "could not create valid token")
			require.NoError(cache.Add(ulids.New(), ulids.New(), valid), "could not add token to cache")
		}
	}

	// Test cache evictions happening concurrently
	wg.Add(2)
	go addTokens()
	go addTokens()

	// Wait for the go routines to finish
	wg.Wait()

	// Should still be 100 tokens in the cache
	require.Equal(100, cache.Size(), "expected cache to have 100 items")
}
