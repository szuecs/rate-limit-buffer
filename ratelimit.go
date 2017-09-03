package circularbuffer

import (
	"sync"
	"time"
)

// RateLimitter
type RateLimitter struct {
	sync.Mutex
	bag        map[string]*CircularBuffer
	maxHits    int
	timeWindow time.Duration
}

// NewRateLimitter returns a new initialized RateLimitter with maxHits is
// the maximal number of hits per time.Duration t.
func NewRateLimitter(maxHits int, t time.Duration) *RateLimitter {
	return &RateLimitter{
		bag:        make(map[string]*CircularBuffer),
		maxHits:    maxHits,
		timeWindow: t,
	}
}

// IsRateLimitted tries to add s to a circularbuffer and returns true if we have
// a free bucket, if not it will return false, which means ratelimit.
func (rl *RateLimitter) IsRateLimitted(s string) bool {
	rl.Lock()
	if rl.bag[s] == nil {
		rl.bag[s] = NewCircularBuffer(rl.maxHits, rl.timeWindow)
	}
	res := !rl.bag[s].Add(time.Now())
	rl.Unlock()
	return res
}
