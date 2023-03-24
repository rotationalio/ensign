package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/config"
	"golang.org/x/time/rate"
)

const (
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"
)

/*
IPRateLimiter is an IP address based limiter that controls how frequently requests
can be made from a single IP address.
limit: represents the number of tokens that can be added to the token bucket per second
burst: maximum number of tokens/requests in a "token bucket" and is initially full
each request consumes tokens from the token bucket and if the bucket is empty
when the request is made, the request is rejected
*/
type IPRateLimiter struct {
	sync.RWMutex
	ips   map[string]*ipInfo //contains a map of IP address to information about its request pattern
	limit rate.Limit
	burst int
}
type ipInfo struct {
	sync.RWMutex
	limiter  *rate.Limiter
	lastSeen time.Time //this is used to delete entries from the "ips" map
}

func NewIPRateLimiter(limit rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*ipInfo),
		limit: limit,
		burst: burst,
	}
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.Lock()
	defer i.Unlock()

	// Here is the second check, e.g. the double check
	if ipInfo, exists := i.ips[ip]; exists {
		return ipInfo.limiter
	}

	// Otherwise the condition from the RLock is still true, so create the limiter
	limiter := rate.NewLimiter(i.limit, i.burst)
	i.ips[ip] = &ipInfo{sync.RWMutex{}, limiter, time.Now()}
	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.RLock()
	ipInfo, exists := i.ips[ip]

	if !exists {
		i.RUnlock()
		return i.AddIP(ip)
	}

	i.RUnlock()

	// update the last seen time associated with the IP address
	ipInfo.Lock()
	ipInfo.lastSeen = time.Now()
	ipInfo.Unlock()
	return ipInfo.limiter
}

// This method checks the "ips" map for IP addresses that haven't been seen for
// more than ttl minutes and deletes those entries.
func (i *IPRateLimiter) cleanupIPInfo(ttl time.Duration) {
	for {
		time.Sleep(time.Minute)
		deleteList := []string{}
		for ip, ipInfo := range i.ips {
			ipInfo.RLock()
			if time.Since(ipInfo.lastSeen) > ttl*time.Minute {
				deleteList = append(deleteList, ip)
			}
			ipInfo.RUnlock()
		}
		for _, ip := range deleteList {
			i.Lock()
			delete(i.ips, ip)
			i.Unlock()
		}
	}
}

func RateLimiter(conf config.RateLimitConfig) gin.HandlerFunc {
	var limiter = NewIPRateLimiter(rate.Limit(conf.PerSecond), conf.Burst)
	//run `cleanupIPInfo` in a go routine to periodically remove entries from the map
	go limiter.cleanupIPInfo(conf.TTL)
	return func(c *gin.Context) {
		// c.ClientIP() does a more thorough check to return the real client IP
		// it also strips out the port, which ensures that we don't create multiple
		// limiters for requests coming from the same IP address due to different port values
		lim := limiter.GetLimiter(c.ClientIP())
		if !lim.Allow() {
			// Add rate limiter headers
			c.Header(HeaderRateLimitLimit, fmt.Sprintf("%d", lim.Burst()))
			c.Header(HeaderRateLimitRemaining, fmt.Sprintf("%0.0f", lim.Tokens()))
			c.Header(HeaderRateLimitReset, "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, api.ErrorResponse(ErrRateLimit))
			return
		}

		// Add rate limiter headers
		// NOTE: these have to be set after the call to lim.Allow()
		c.Header(HeaderRateLimitLimit, fmt.Sprintf("%d", lim.Burst()))
		c.Header(HeaderRateLimitRemaining, fmt.Sprintf("%0.0f", lim.Tokens()))
		c.Header(HeaderRateLimitReset, "1")
		c.Next()
	}
}
