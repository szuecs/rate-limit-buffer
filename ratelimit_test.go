package circularbuffer

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimitter(t *testing.T) {
	window := 1 * time.Second
	rl := NewRateLimitter(2, window)

	if rl.IsRateLimitted("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimitted("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.IsRateLimitted("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if rl.IsRateLimitted("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.IsRateLimitted("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.IsRateLimitted("bar") {
		t.Errorf("bar should be rate limitted")
	}

	time.Sleep(window)

	if rl.IsRateLimitted("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimitted("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.IsRateLimitted("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimitted("bar") {
		t.Errorf("bar should not be rate limitted")
	}

	if !rl.IsRateLimitted("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if !rl.IsRateLimitted("bar") {
		t.Errorf("bar should be rate limitted")
	}
}

func TestRateLimitterConcurrent(t *testing.T) {
	window := 1 * time.Second
	rl := NewRateLimitter(2, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		if rl.IsRateLimitted(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if rl.IsRateLimitted(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if !rl.IsRateLimitted(s) {
			t.Errorf("%s should be rate limitted", s)
		}

		time.Sleep(window)

		if rl.IsRateLimitted(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if rl.IsRateLimitted(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if !rl.IsRateLimitted(s) {
			t.Errorf("%s should be rate limitted", s)
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
}

func BenchmarkIsRateLimitted(b *testing.B) {
	window := 10 * time.Millisecond
	rl := NewRateLimitter(2, window)

	for n := 0; n < b.N; n++ {
		if rl.IsRateLimitted("foo") && rl.IsRateLimitted(fmt.Sprintf("foo%d", n)) {
			b.Errorf("Failed 2nd should never be limitted")
		}
	}
}

func BenchmarkIsRateLimittedConcurrent1(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)

	for n := 0; n < b.N; n++ {
		go func(j int) {
			rl.IsRateLimitted(fmt.Sprintf("foo%d", j))
		}(n)
	}
}
func BenchmarkIsRateLimittedConcurrent10(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			go func(j int) {
				rl.IsRateLimitted(fmt.Sprintf("foo%d", j))
			}(i + n)
		}
	}
}
func BenchmarkIsRateLimittedConcurrent100(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)

	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go func(j int) {
				rl.IsRateLimitted(fmt.Sprintf("foo%d", j))
			}(i + n)
		}
	}
}
func BenchmarkIsRateLimittedConcurrent1000(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)

	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000; i++ {
			go func(j int) {
				rl.IsRateLimitted(fmt.Sprintf("foo%d", j))
			}(i + n)
		}
	}
}
