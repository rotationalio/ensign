package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipInfo struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips   map[string]*ipInfo
	mu    *sync.RWMutex
	limit rate.Limit // number of tokens allowed per second
	burst int        // maximum number of tokens allowed in a single call
}

func NewIPRateLimiter(limit rate.Limit, burst int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:   make(map[string]*ipInfo),
		mu:    &sync.RWMutex{},
		limit: limit,
		burst: burst,
	}

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.limit, i.burst)

	// add the current time
	i.ips[ip] = &ipInfo{limiter, time.Now()}

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	sync.Map
	i.mu.Lock()
	ipInfo, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	// update the last seen time associated with the IP address
	ipInfo.lastSeen = time.Now()
	return ipInfo.limiter
}

// Every minute check the map for IP addresses that haven't been seen for
// more than 3 minutes and delete the entries.
func (i *IPRateLimiter) cleanupIPInfo() {
	var mu sync.Mutex
	for {
		time.Sleep(time.Minute)

		//TODO: add the candidate keys here and lock the key value pair before the delete and unlock after
		mu.Lock()
		for ip, ipInfo := range i.ips {
			if time.Since(ipInfo.lastSeen) > 3*time.Minute {
				delete(i.ips, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimiter() http.HandlerFunc {
	var limiter = NewIPRateLimiter(1, 5)
	//TODO: add cleanup here
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if strings.TrimSpace(xForwardedFor) != "" {
			w.Header().Add("X-Rate-Limit-Request-Forwarded-For", xForwardedFor)
		}
		w.Header().Add("X-Rate-Limit-Request-Remote-Addr", r.RemoteAddr)
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%.2f", limiter.Limit()))
		w.Header().Add("X-Rate-Limit-Reset", "1")
		w.Header().Add("RateLimit-Remaining", fmt.Sprintf("%.2f", limiter.Tokens()))
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
	}
}
