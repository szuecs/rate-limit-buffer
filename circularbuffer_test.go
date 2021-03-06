package circularbuffer

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	l := 2
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)

	if cb.Len() != 0 {
		t.Errorf("buffer is not used expected 0, but is %v", cb.Len())
	}
	cb.Add(time.Now())
	if cb.Len() != 1 {
		t.Errorf("expected 1, but is %v", cb.Len())
	}
	cb.Add(time.Now())
	if cb.Len() != 2 {
		t.Errorf("expected 2, but is %v", cb.Len())
	}
	cb.Add(time.Now())
	if cb.Len() != 2 {
		t.Errorf("expected 2, but is %v", cb.Len())
	}

	time.Sleep(window)
	if cb.Len() != 0 {
		t.Errorf("expected 0, but is %v", cb.Len())
	}
}

func TestInUse(t *testing.T) {
	l := 2
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)
	if cb.InUse() {
		t.Errorf("buffer should not be in use")
	}
	cb.Add(time.Now())
	if !cb.InUse() {
		t.Errorf("buffer should be in use")
	}
	cb.Add(time.Now())
	if !cb.InUse() {
		t.Errorf("buffer should be in use")
	}

	time.Sleep(window)
	if cb.InUse() {
		t.Errorf("buffer should not be in use anymore")
	}
}

func TestFree(t *testing.T) {
	l := 2
	cb := NewCircularBuffer(l, 1*time.Second)
	if !cb.Free() {
		t.Errorf("empty buffer should be Free")
	}
	cb.Add(time.Now())
	if !cb.Free() {
		t.Errorf("buffer has one slot Free")
	}
	cb.Add(time.Now())
	if cb.Free() {
		t.Errorf("buffer is full and should not be Free")
	}
}

func TestCurrent(t *testing.T) {
	l := 2
	window := 100 * time.Millisecond
	cb := NewCircularBuffer(l, window)
	start := time.Now()
	for i := 0; i < 2*l; i++ {
		new := start.Add(time.Duration(i) * window)
		cb.Add(new)
		if cb.Current("") != new {
			t.Errorf("current position should be the last one added")
		}
		time.Sleep(window)
	}
}

func TestDelta(t *testing.T) {
	l := 4
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)
	start := time.Now()
	for i := 0; i < l; i++ {
		cb.Add(start.Add(time.Duration(i) * window))
	}
	want := time.Duration(3) * window
	delta := cb.delta()
	if delta != want {
		t.Errorf("want != delta => %s / %s", want, delta)
	}
}

func TestAdd(t *testing.T) {
	l := 4
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)
	if !cb.Add(time.Now()) {
		t.Errorf("empty buffer Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 1 entry not full Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 2 entries not full Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 3 entries not full Add() should return true")
	}

	if cb.Add(time.Now()) {
		t.Errorf("buffer is full Add() should return false")
	}

	time.Sleep(window)

	if !cb.Add(time.Now()) {
		t.Errorf("empty buffer Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 1 entry not full Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 2 entries not full Add() should return true")
	}
	if !cb.Add(time.Now()) {
		t.Errorf("buffer with 3 entries not full Add() should return true")
	}

	if cb.Add(time.Now()) {
		t.Errorf("buffer is full Add() should return false")
	}
}

func TestCicularBufferMassiveConcurrentUse(t *testing.T) {
	l := 1 << 21
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)
	if !cb.Add(time.Now()) {
		t.Errorf("empty buffer Add() should return true")
	}
	if !cb.InUse() {
		t.Errorf("buffer should be in use")
	}

	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		for i := 0; i < 1<<10; i++ {
			if !cb.Add(time.Now()) {
				t.Errorf("%s Add should return true", s)
			}
			if !cb.InUse() {
				t.Errorf("%s buffer should be in use", s)
			}
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()

}

func TestResizeBufferIncrease(t *testing.T) {
	l := 4
	window := 1 * time.Second
	cb := NewCircularBuffer(l, window)
	start := time.Now()
	for i := 0; i < l; i++ {
		cb.Add(start)
	}
	cb.resize(2 * l)
	for i := 0; i < l; i++ {
		if !cb.slots[i].Equal(start) {
			t.Errorf("invalid value found in slot %d: %s", i, cb.slots[i])
		}
	}
	for i := l; i < 2*l; i++ {
		if !cb.slots[i].IsZero() {
			t.Errorf("invalid value found in slot %d: %s", i, cb.slots[i])
		}
	}
}

func TestResizeBufferDecreaseFullVaryingOffset(t *testing.T) {
	l := 8
	window := 1 * time.Second
	for off := 0; off < l; off++ {
		for newSize := 1; newSize < l; newSize++ {
			cb := NewCircularBuffer(l, window)
			cb.offset = off
			start := time.Time{}
			for i := 0; i < l; i++ {
				added := cb.Add(start.Add(time.Duration(i) * window))
				if ! added {
					t.Errorf("%v not added", time.Duration(i)*window)
				}
			}
			cb.resize(newSize)
			for i := 0; i < newSize; i++ {
				if !cb.slots[i].Equal(start.Add(window * time.Duration(l-newSize+i))) {
					t.Errorf("invalid value found for new size %d in slot %d: %s", newSize, i, cb.slots[i])
				}
			}
			if cb.offset != 0 {
				t.Errorf("offset is not 0. Is: %d", cb.offset)
			}
		}
	}
}

func TestResizeBufferDecreaseFullOverwritten(t *testing.T) {
	for newSize := 1; newSize < 5; newSize++ {
		for writes := 6; writes >= 10; writes++ {
			b := NewCircularBuffer(5, 5*time.Second)
			ts := time.Now()
			for i := 0; i < writes; i++ {
				b.Add(ts.Add(time.Duration(i) * time.Second))
				time.Sleep(time.Second)
			}

			b.Resize("", newSize)

			if b.Len() != 3 {
				t.Errorf("length is not 3")
			}

			for i := 0; i < 3; i++ {
				expectedTs := ts.Add(time.Duration(i+(writes-newSize)) * time.Second)
				if b.slots[i] != expectedTs {
					t.Errorf("(%d) Expected %v got %v", i, expectedTs, b.slots[i])
				}
			}
			if b.offset != 0 {
				t.Errorf("unexpected offset: %d", b.offset)
			}
		}
	}
}

func TestResizeBufferDecreaseNonFull(t *testing.T) {
	for newSize := 2; newSize < 5; newSize++ {
		for writes := newSize + 1; writes <= 5; writes ++ {
			b := NewCircularBuffer(5, 1*time.Minute)
			timestamp := time.Now()
			for i := 0; i < writes; i++ {
				b.Add(timestamp.Add(time.Duration(i) * time.Second))
			}

			b.Resize("", newSize)
			newLen := int(math.Min(float64(writes), float64(newSize)))
			if b.Len() != newLen {
				t.Errorf("length is not %d, is %d", newLen, b.Len())
			}

			for i := 0; i < newLen; i++ {
				diff := writes - newSize
				expected := timestamp.Add(time.Duration(diff+i) * time.Second)
				if b.slots[i] != expected {
					t.Errorf("unexpected time: %v expected, %v", b.slots[i].Second(), expected.Second())
				}
			}

			if b.offset != 0 {
				t.Errorf("unexpected offset: %d", b.offset)
			}
		}
	}
}

func TestRetryAfter(t *testing.T) {
	l := 2
	window := 100 * time.Millisecond
	cb := NewCircularBuffer(l, window)
	if cb.retryAfter() != 0 {
		t.Errorf("wait time should be zero if there are free slots")
	}
	start := time.Now()
	for i := 0; i < l*10; i++ {
		new := start.Add(time.Duration(i) * window)
		cb.Add(new)
		time.Sleep(10 * time.Millisecond)
		if !cb.Free() && cb.retryAfter() == 0 {
			t.Errorf("wait time should not be zero if there are no free slots")
		}
		if cb.retryAfter() < 0 {
			t.Errorf("wait time should not be less than zero")
		}
	}
}
