package circularbuffer

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterAllowMassiveConcurrent(t *testing.T) {
	window := 1 * time.Second
	rl := NewRateLimiter(1<<21, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		for i := 0; i < 1<<10; i++ {
			if !rl.Allow("") {
				t.Errorf("%s should not be rate limitted", s)
			}
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
	rl.Close()
}

func BenchmarkRateLimiterAllow(b *testing.B) {
	window := 1 * time.Second
	rl := NewRateLimiter(1<<21, window)
	for n := 0; n < b.N; n++ {
		rl.Allow("")
	}
	rl.Close()
}

func newClientRateLimiter(maxHits int, d time.Duration) *ClientRateLimiter {
	return NewClientRateLimiter(maxHits, d, 5*d)
}

func TestClientRateLimiterClose(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("failed to close %v", r)
		}
	}()

	window := 1 * time.Second
	rl := newClientRateLimiter(5, window)
	rl.Close()

}

func TestClientRateLimiterDeleteOld(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(5, window)

	rl.Allow("foo")
	rl.Allow("bar")
	rl.DeleteOld()
	if _, ok := rl.bag["foo"]; !ok {
		t.Errorf("foo should be found")
	}
	if _, ok := rl.bag["bar"]; !ok {
		t.Errorf("bar should be found")
	}

	time.Sleep(window)
	rl.DeleteOld()
	if _, ok := rl.bag["foo"]; ok {
		t.Errorf("foo should not be found")
	}
	if _, ok := rl.bag["bar"]; ok {
		t.Errorf("bar should not be found")
	}
	rl.Close()
}

func TestClientRateLimiterAllow(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(2, window)

	if !rl.Allow("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.Allow("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.Allow("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if !rl.Allow("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.Allow("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.Allow("bar") {
		t.Errorf("bar should be rate limitted")
	}

	time.Sleep(window)

	if !rl.Allow("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.Allow("bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.Allow("foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.Allow("bar") {
		t.Errorf("bar should not be rate limitted")
	}

	if rl.Allow("foo") {
		t.Errorf("foo should be rate limitted")
	}
	if rl.Allow("bar") {
		t.Errorf("bar should be rate limitted")
	}
	rl.Close()
}

func TestClientRateLimiterAllowConcurrent(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(2, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		if !rl.Allow(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if !rl.Allow(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if rl.Allow(s) {
			t.Errorf("%s should be rate limitted", s)
		}

		time.Sleep(window)

		if !rl.Allow(s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if !rl.Allow(s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if rl.Allow(s) {
			t.Errorf("%s should be rate limitted", s)
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
	rl.Close()
}

func TestClientRateLimiterAllowMassiveConcurrent(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(1<<21, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		for i := 0; i < 1<<10; i++ {
			if !rl.Allow(s) {
				t.Errorf("%s should not be rate limitted", s)
			}
		}
		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
	rl.Close()
}

func BenchmarkClientRateLimiterAllow(b *testing.B) {
	window := 10 * time.Millisecond
	rl := newClientRateLimiter(2, window)

	for n := 0; n < b.N; n++ {
		if !rl.Allow("foo") && !rl.Allow(fmt.Sprintf("foo%d", n)) {
			b.Errorf("Failed 2nd should never be limitted")
		}
	}
	rl.Close()
}

func BenchmarkClientRateLimiterAllowBaseData1(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 10
	for i := 0; i < m; i++ {
		rl.Allow(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.Allow(fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData10(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 10
	for i := 0; i < m*m; i++ {
		rl.Allow(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.Allow(fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData100(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < m*m; i++ {
		rl.Allow(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.Allow(fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData1000(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 1000
	for i := 0; i < m*m; i++ {
		rl.Allow(fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.Allow(fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}

func BenchmarkClientRateLimiterAllowConcurrent1(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.Allow(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}
func BenchmarkClientRateLimiterAllowConcurrent10(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.Allow(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}
func BenchmarkClientRateLimiterAllowConcurrent100(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.Allow(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}
func BenchmarkClientRateLimiterAllowConcurrent1000(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.Allow(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}

func BenchmarkClientRateLimiterAllowConcurrentAddDelete10(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := NewClientRateLimiter(10, window, window)
	m := 100

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.Allow(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}
