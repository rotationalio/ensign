package tokens

import (
	"container/list"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

const defaultCacheSize = 128

// Cache is an in-memory cache that stores issued token claims by ULID.
type Cache struct {
	sync.RWMutex

	// The doubly linked list allows us to remove expired tokens without iterating over
	// the entire list.
	items *list.List

	// A map to quickly find tokens in the cache.
	index map[string]*list.Element

	// Maximum number of tokens in the cache.
	maxTokens uint32
}

// SessionToken is a string token issued to users that has a limited lifetime. The
// cache is responsible for checking whether tokens are expired and removing expired
// tokens to maintain the cache size.
type SessionToken struct {
	// ID of the user.
	userID ulid.ULID

	// Project ID of the project the user is accessing.
	projectID ulid.ULID

	// Expiration time of the token.
	expiresAt time.Time

	// The actual token string issued to the user.
	tks string
}

// The key is used to lookup the token in the cache by user ID and project ID.
func (t *SessionToken) Key() string {
	return t.userID.String() + t.projectID.String()
}

// Create a new token cache with a maximum number of items.
func NewCache(size uint32) (c *Cache) {
	c = &Cache{
		items: list.New(),
		index: make(map[string]*list.Element, 0),
	}

	if size == 0 {
		c.maxTokens = defaultCacheSize
	} else {
		c.maxTokens = size
	}

	return c
}

// Get the token from the user ID and project ID or return an error if it does not
// exist.
func (c *Cache) Get(userID, projectID ulid.ULID) (tks string, err error) {
	var (
		element *list.Element
		ok      bool
	)

	// Lookup the token from the cache.
	token := &SessionToken{
		userID:    userID,
		projectID: projectID,
	}

	c.Lock()
	defer c.Unlock()
	if element, ok = c.index[token.Key()]; !ok {
		return "", ErrCacheMiss
	}

	// Check if the token has expired. This will panic if the element in the cache is
	// not a SessionToken so external packages should not access the cache directly.
	token = element.Value.(*SessionToken)
	if token.expiresAt.Before(time.Now()) {
		// If expired, remove this token and all tokens created before it. There's no
		// need to check if they are expired because the items are naturally sorted by
		// expiration time.
		c.trunc(element)
		return "", ErrCacheExpired
	}

	return token.tks, nil
}

// Add a token to the cache by user ID and project ID.
func (c *Cache) Add(userID, projectID ulid.ULID, tks string) (err error) {
	// Create a new token
	token := &SessionToken{
		userID:    userID,
		projectID: projectID,
		tks:       tks,
	}
	if token.expiresAt, err = ExpiresAt(token.tks); err != nil {
		return err
	}

	// Evict the oldest tokens if the cache is full.
	c.Lock()
	defer c.Unlock()
	c.prune()

	// Add the token to the cache.
	element := c.items.PushBack(token)
	c.index[token.Key()] = element

	return nil
}

// Remove a token from the cache by user ID and project ID.
func (c *Cache) Remove(userID, projectID ulid.ULID) {
	// Lookup the token from the cache. Do not return an error if the token does not
	// exist.
	token := &SessionToken{
		userID:    userID,
		projectID: projectID,
	}
	key := token.Key()

	c.Lock()
	defer c.Unlock()
	if element, ok := c.index[key]; ok {
		c.remove(element)
	}
}

// Clear the entire cache.
func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.items.Init()
	c.index = make(map[string]*list.Element, 0)
}

// Return the current size of the cache for profiling or testing purposes.
func (c *Cache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return c.items.Len()
}

// Remove oldest tokens from the cache until the capacity is below the maximum.
func (c *Cache) prune() {
	for c.items.Len() >= int(c.maxTokens) {
		element := c.items.Front()
		if element == nil {
			break
		}

		c.remove(element)
	}
}

// Remove an item from the cache and all items before it.
func (c *Cache) trunc(element *list.Element) {
	for element != nil {
		prev := element.Prev()
		c.remove(element)
		element = prev
	}
}

// Remove a single item from the cache.
func (c *Cache) remove(element *list.Element) {
	c.items.Remove(element)
	delete(c.index, element.Value.(*SessionToken).Key())
}
