package circularbuffer

import (
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
