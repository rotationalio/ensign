package tokens

import "errors"

var (
	ErrCacheMiss    = errors.New("requested key is not in the cache")
	ErrCacheExpired = errors.New("requested key is expired")
)
