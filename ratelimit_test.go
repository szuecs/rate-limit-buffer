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

	if rl.IsRateLimited("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimited("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.IsRateLimited("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if rl.IsRateLimited("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.IsRateLimited("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.IsRateLimited("bar") {
		t.Errorf("bar should be rate limitted")
	}

	time.Sleep(window)

	if rl.IsRateLimited("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimited("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.IsRateLimited("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.IsRateLimited("bar") {
		t.Errorf("bar should not be rate limitted")
	}

	if !rl.IsRateLimited("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if !rl.IsRateLimited("bar") {
		t.Errorf("bar should be rate limitted")
	}
}

func TestRateLimitterConcurrent(t *testing.T) {
	window := 1 * time.Second
	rl := NewRateLimitter(2, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		if rl.IsRateLimited(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if rl.IsRateLimited(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if !rl.IsRateLimited(s) {
			t.Errorf("%s should be rate limitted", s)
		}

		time.Sleep(window)

		if rl.IsRateLimited(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if rl.IsRateLimited(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if !rl.IsRateLimited(s) {
			t.Errorf("%s should be rate limitted", s)
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
}

func BenchmarkIsRateLimited(b *testing.B) {
	window := 10 * time.Millisecond
	rl := NewRateLimitter(2, window)

	for n := 0; n < b.N; n++ {
		if rl.IsRateLimited("foo") && rl.IsRateLimited(fmt.Sprintf("foo%d", n)) {
			b.Errorf("Failed 2nd should never be limitted")
		}
	}
}

func BenchmarkIsRateLimitedBaseData1(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 10
	for i := 0; i < m; i++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", n%m))
	}
}
func BenchmarkIsRateLimitedBaseData10(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 10
	for i := 0; i < m*m; i++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", n%m))
	}
}
func BenchmarkIsRateLimitedBaseData100(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 100
	for i := 0; i < m*m; i++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", n%m))
	}
}
func BenchmarkIsRateLimitedBaseData1000(b *testing.B) {
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 1000
	for i := 0; i < m*m; i++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.IsRateLimited(fmt.Sprintf("foo%d", n%m))
	}
}

func BenchmarkIsRateLimitedConcurrent1(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 100

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.IsRateLimited(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func BenchmarkIsRateLimitedConcurrent10(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 100

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.IsRateLimited(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func BenchmarkIsRateLimitedConcurrent100(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 100

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.IsRateLimited(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func BenchmarkIsRateLimitedConcurrent1000(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := NewRateLimitter(10, window)
	m := 100

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.IsRateLimited(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
