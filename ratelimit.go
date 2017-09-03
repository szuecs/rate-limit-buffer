package circularbuffer

import (
	"sync"
	"time"
)

// RateLimiter
type RateLimiter struct {
	sync.Mutex
	bag        map[string]*CircularBuffer
	maxHits    int
	timeWindow time.Duration
}

// NewRateLimiter returns a new initialized RateLimitter with maxHits is
// the maximal number of hits per time.Duration d.
func NewRateLimiter(maxHits int, d, cleanInterval time.Duration, quit chan struct{}) *RateLimiter {
	rl := &RateLimiter{
		bag:        make(map[string]*CircularBuffer),
		maxHits:    maxHits,
		timeWindow: d,
	}
	go rl.startCleanerDaemon(cleanInterval, quit)
	return rl
}

// Allow tries to add s to a circularbuffer and returns true if we have
// a free bucket, if not it will return false, which means ratelimit.
func (rl *RateLimiter) Allow(s string) bool {
	rl.Lock()
	if rl.bag[s] == nil {
		rl.bag[s] = NewCircularBuffer(rl.maxHits, rl.timeWindow)
	}
	res := rl.bag[s].Add(time.Now())
	rl.Unlock()
	return res
}

// DeleteOld removes old entries from state bag
func (rl *RateLimiter) DeleteOld() {
	rl.Lock()
	for k, cb := range rl.bag {
		if !cb.InUse() {
			delete(rl.bag, k)
		}
	}
	rl.Unlock()
}

func (rl *RateLimiter) startCleanerDaemon(d time.Duration, quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		case <-time.After(d):
			rl.DeleteOld()
		}
	}
}
