package circularbuffer

import (
	"fmt"
	"sync"
	"time"
)

// CircularBuffer has slots to store times as int64 and an offset,
// which marks the next free entry. Slots are fixed in size.
type CircularBuffer struct {
	sync.RWMutex
	slots      []time.Time
	offset     int
	timeWindow time.Duration
}

func NewCircularBuffer(l int, t time.Duration) *CircularBuffer {
	return &CircularBuffer{
		slots:      make([]time.Time, l, l),
		offset:     0,
		timeWindow: t,
	}
}

func (cb *CircularBuffer) Cap() int {
	return len(cb.slots)
}

func (cb *CircularBuffer) Len() int {
	n := 0
	for i := 0; i < len(cb.slots); i++ {
		cb.RLock()
		slot := cb.slots[i]
		cb.RUnlock()
		if slot.Add(cb.timeWindow).After(time.Now()) {
			n++
		}
	}
	return n
}

func (cb *CircularBuffer) InUse() bool {
	cb.RLock()
	l := len(cb.slots)
	offset := cb.offset
	cb.RUnlock()

	newestOffset := (offset - 1) % l
	if newestOffset < 0 {
		newestOffset = l + newestOffset
	}

	cb.RLock()
	slot := cb.slots[newestOffset]
	cb.RUnlock()

	return slot.Add(cb.timeWindow).After(time.Now())
}

// Free returns if there is space or the bucket is full for the current time.
// Example:
//   time.Now(): 5
//   timeWindow: 2
//   [1 2 3 4]
//          ^
//   5-2 = 3 --> 2 free slots [1,2] are too old and are Free already
func (cb *CircularBuffer) Free() bool {
	cb.RLock()
	slot := cb.slots[cb.offset]
	cb.RUnlock()
	return slot.Add(cb.timeWindow).Before(time.Now())
}

// Add adds an element to the next free bucket in the buffer and
// returns true. It returns false if there is no free bucket.
// Example
//   [_ _ _ _]
//    ^
//   [1 _ _ _]
//      ^
//   [1 2 _ _]
//        ^
//   [1 2 3 _]
//          ^
//   [1 2 3 4]
//    ^
func (cb *CircularBuffer) Add(t time.Time) bool {
	if cb.Free() {
		cb.Lock()
		cb.slots[cb.offset] = t
		cb.offset = (cb.offset + 1) % len(cb.slots)
		cb.Unlock()
		return true
	}
	return false
}

func (cb *CircularBuffer) current() time.Time {
	cb.RLock()
	curOff := cb.offset - 1
	if curOff < 0 {
		curOff += len(cb.slots)
	}
	cur := cb.slots[curOff]
	cb.RUnlock()
	return cur
}

func (cb *CircularBuffer) delta() time.Duration {
	cb.RLock()
	cur := cb.current()
	next := cb.Next()
	cb.RUnlock()
	return cur.Sub(next)
}

func (cb *CircularBuffer) Next() time.Time {
	cb.RLock()
	next := cb.slots[cb.offset]
	cb.RUnlock()
	return next
}

func (cb *CircularBuffer) retryAfter() time.Duration {
	if cb.Free() {
		return 0
	}
	first := cb.Next()
	next := first.Add(cb.timeWindow)
	now := time.Now()
	return next.Sub(now)
}

// needs to be called with Lock() held by caller
func (cb *CircularBuffer) resize(n int) {
	if n <= 0 {
		fmt.Printf("resize(): refusing to resize circular buffer to %d", n)
		return
	}
	cur := len(cb.slots)
	if cur == n {
		return
	}
	newSlots := make([]time.Time, n, n)
	if cur < n {
		copy(newSlots, cb.slots)
		cb.slots = newSlots
	} else {
		cb.offset = cb.offset - n
		if cb.offset < 0 {
			cb.offset += cur
		}
		for i := 0; i < n; i++ {
			off := (cb.offset + i) % cur
			newSlots[i] = cb.slots[off]
		}
		cb.slots = newSlots
		cb.offset = n - 1
	}
}
