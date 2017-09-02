package circularbuffer

import (
	"time"
)

// CircularBuffer has slots to store times as int64 and an offset,
// which marks the next free entry. Slots are fixed in size.
type CircularBuffer struct {
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

// Free returns if there is space or the bucket is full for the current time.
// Example:
// time.Now(): 5
// timeWindow: 2
// [1 2 3 4]
//        ^
// 5-2 = 3 --> 2 free slots [1,2] are too old and are Free already
func (cb *CircularBuffer) Free() bool {
	return cb.slots[cb.offset].Add(cb.timeWindow).Before(time.Now())
}

// Add adds an element to the next free bucket in the buffer and
// returns true. It returns false if there is no free bucket.
// Example
// [_ _ _ _]
//  ^
// [1 _ _ _]
//    ^
// [1 2 _ _]
//      ^
// [1 2 3 _]
//        ^
// [1 2 3 4]
//  ^
func (cb *CircularBuffer) Add(t time.Time) bool {
	if cb.Free() {
		cb.slots[cb.offset] = t
		cb.offset = (cb.offset + 1) % len(cb.slots)
		return true
	}
	return false
}
