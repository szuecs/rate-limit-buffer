package circularbuffer

import (
	"testing"
	"time"
)

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
