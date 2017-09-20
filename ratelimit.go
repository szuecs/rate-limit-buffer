package circularbuffer

import (
	"sync"
	"time"
)

// RateLimiter
type RateLimiter struct {
	sync.RWMutex
	bag        map[string]*CircularBuffer
	maxHits    int
	timeWindow time.Duration
	quitCH     chan struct{}
}

// NewRateLimiter returns a new initialized RateLimitter with maxHits is
// the maximal number of hits per time.Duration d.
func NewRateLimiter(maxHits int, d, cleanInterval time.Duration, quit chan struct{}) *RateLimiter {
	rl := &RateLimiter{
		bag:        make(map[string]*CircularBuffer),
		maxHits:    maxHits,
		timeWindow: d,
		quitCH:     quit,
	}
	go rl.startCleanerDaemon(cleanInterval)
	return rl
}

// Allow tries to add s to a circularbuffer and returns true if we have
// a free bucket, if not it will return false, which means ratelimit.
func (rl *RateLimiter) Allow(s string) bool {
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
func (rl *RateLimiter) DeleteOld() {
	rl.Lock()
	for k, cb := range rl.bag {
		if !cb.InUse() {
			delete(rl.bag, k)
		}
	}
	rl.Unlock()
}

// Close will stop the cleanup goroutine
func (rl *RateLimiter) Close() {
	close(rl.quitCH)
}

func (rl *RateLimiter) startCleanerDaemon(d time.Duration) {
	for {
		select {
		case <-rl.quitCH:
			return
		case <-time.After(d):
			rl.DeleteOld()
		}
	}
}
