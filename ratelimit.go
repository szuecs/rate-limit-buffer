package circularbuffer

import (
	"sync"
	"time"
)

// RateLimitter
type RateLimitter struct {
	sync.Mutex
	bag        map[string]*CircularBuffer
	len        int
	timeWindow time.Duration
}

func NewRateLimitter(l int, t time.Duration) *RateLimitter {
	return &RateLimitter{
		bag:        make(map[string]*CircularBuffer),
		len:        l,
		timeWindow: t,
	}
}

// IsRateLimitted tries to add s to a circularbuffer and returns true if we have
// a free bucket, if not it will return false, which means ratelimit.
func (rl *RateLimitter) IsRateLimitted(s string) bool {
	rl.Lock()
	if rl.bag[s] == nil {
		rl.bag[s] = NewCircularBuffer(rl.len, rl.timeWindow)
	}
	res := !rl.bag[s].Add(time.Now())
	rl.Unlock()
	return res
}
