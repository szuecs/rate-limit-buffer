package circularbuffer

import (
	"sync"
	"time"
)

// RateLimiter is an interface which can be used to implement
// rate limiting.
type RateLimiter interface {
	// Allow returns true if call should be allowed, false in case
	// you should rate limit.
	Allow(string) bool
	// Close cleans up the RateLimiter implementation.
	Close()
}

// NewRateLimiter returns a new initialized RateLimitter with maxHits
// as the maximal number of hits per time.Duration d. This can be used
// to implement maximum number of requests for a backend to protect
// from a known scaling limit.
func NewRateLimiter(maxHits int, d time.Duration) RateLimiter {
	return NewCircularBuffer(maxHits, d)
}

// Allow returns true if there is a free bucket and we should not rate
// limit, if not it will return false, which means ratelimit.
func (cb *CircularBuffer) Allow(s string) bool {
	return cb.Add(time.Now())
}

// Close implements the RateLimiter interface to shutdown, nothing to
// do.
func (cb *CircularBuffer) Close() {
}

// ClientRateLimiter implements the RateLimiter interface and does
// rate limiting based on the the String passed to Allow(). This can
// be used to limit per client calls to the backend. For example you
// can slow down user enumeration or dictionary attacks to /login
// APIs.
type ClientRateLimiter struct {
	sync.RWMutex
	bag        map[string]*CircularBuffer
	maxHits    int
	timeWindow time.Duration
	quitCH     chan struct{}
}

// NewRateLimiter returns a new initialized RateLimitter with maxHits is
// the maximal number of hits per time.Duration d.
func NewClientRateLimiter(maxHits int, d, cleanInterval time.Duration) *ClientRateLimiter {
	quit := make(chan struct{})
	crl := &ClientRateLimiter{
		bag:        make(map[string]*CircularBuffer),
		maxHits:    maxHits,
		timeWindow: d,
		quitCH:     quit,
	}
	go crl.startCleanerDaemon(cleanInterval)
	return crl
}

// Allow tries to add s to a circularbuffer and returns true if we have
// a free bucket, if not it will return false, which means ratelimit.
func (rl *ClientRateLimiter) Allow(s string) bool {
	var source *CircularBuffer
	var present bool

	rl.RLock()
	if source, present = rl.bag[s]; !present {
		rl.RUnlock()
		rl.Lock()
		source = NewCircularBuffer(rl.maxHits, rl.timeWindow)
		rl.bag[s] = source
		rl.Unlock()
	} else {
		rl.RUnlock()
	}
	present = source.Add(time.Now())
	return present
}

// DeleteOld removes old entries from state bag
func (rl *ClientRateLimiter) DeleteOld() {
	rl.Lock()
	for k, cb := range rl.bag {
		if !cb.InUse() {
			delete(rl.bag, k)
		}
	}
	rl.Unlock()
}

// Close will stop the cleanup goroutine
func (rl *ClientRateLimiter) Close() {
	close(rl.quitCH)
}

func (rl *ClientRateLimiter) startCleanerDaemon(d time.Duration) {
	for {
		select {
		case <-rl.quitCH:
			return
		case <-time.After(d):
			rl.DeleteOld()
		}
	}
}
