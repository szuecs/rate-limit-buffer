package circularbuffer

import (
	"context"
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
			if !rl.AllowContext(context.Background(), "") {
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

func BenchmarkRateLimiterAllowContext(b *testing.B) {
	window := 1 * time.Second
	rl := NewRateLimiter(1<<21, window)
	for n := 0; n < b.N; n++ {
		rl.AllowContext(context.Background(), "")
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

	rl.AllowContext(context.Background(), "foo")
	rl.AllowContext(context.Background(), "bar")
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

func TestClientRateLimiterOldest(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(2, window) // [ , ]
	zero := time.Time{}

	if !rl.Oldest("foo").Equal(zero) { // [0, 0]
		t.Errorf("0foo should return zero")
	}

	rl.AllowContext(context.Background(), "foo") // [t0, 0]
	if !rl.Oldest("foo").Equal(zero) {
		t.Errorf("1foo should return zero")
	}

	t0 := rl.Current("foo")
	rl.AllowContext(context.Background(), "foo") // [t0, t1]
	if !rl.Oldest("foo").Equal(t0) {
		t.Errorf("2foo should return t0")
	}
}

func TestClientRateLimiterCurrent(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(2, window)

	t0 := time.Now()
	rl.AllowContext(context.Background(), "foo")
	rl.AllowContext(context.Background(), "foo")

	t1 := time.Now()
	rl.AllowContext(context.Background(), "bar")
	rl.AllowContext(context.Background(), "bar")
	t2 := time.Now()

	if rl.Current("foo").Before(t0) || rl.Current("foo").After(t1) {
		t.Errorf("foo should be in the time frame t0 < current < t1")
	}
	if rl.Current("bar").Before(t1) || rl.Current("bar").After(t2) {
		t.Errorf("bar should be in the time frame t1 < current < t2")
	}
}

func TestClientRateLimiterAllowContext(t *testing.T) {
	window := 1 * time.Second
	rl := newClientRateLimiter(2, window)

	if !rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if rl.AllowContext(context.Background(), "bar") {
		t.Errorf("bar should be rate limitted")
	}

	time.Sleep(window)

	if !rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "bar") {
		t.Errorf("bar should not be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should not be rate limitted")
	}
	if !rl.AllowContext(context.Background(), "bar") {
		t.Errorf("bar should not be rate limitted")
	}

	if rl.AllowContext(context.Background(), "foo") {
		t.Errorf("foo should be rate limitted")
	}
	if rl.AllowContext(context.Background(), "bar") {
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
		if !rl.AllowContext(context.Background(), s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if !rl.AllowContext(context.Background(), s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if rl.AllowContext(context.Background(), s) {
			t.Errorf("%s should be rate limitted", s)
		}

		time.Sleep(window)

		if !rl.AllowContext(context.Background(), s) {
			t.Errorf("%s should not be rate limitted", s)
		}
		if !rl.AllowContext(context.Background(), s) {
			t.Errorf("%s should not be rate limitted", s)
		}

		if rl.AllowContext(context.Background(), s) {
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
			if !rl.AllowContext(context.Background(), s) {
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

func BenchmarkClientRateLimiterAllowContext(b *testing.B) {
	window := 10 * time.Millisecond
	rl := newClientRateLimiter(2, window)

	for n := 0; n < b.N; n++ {
		if !rl.AllowContext(context.Background(), "foo") && !rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", n)) {
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
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData10(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 10
	for i := 0; i < m*m; i++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData100(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < m*m; i++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", n%m))
	}
	rl.Close()
}
func BenchmarkClientRateLimiterAllowBaseData1000(b *testing.B) {
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 1000
	for i := 0; i < m*m; i++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", i%m))
	}

	for n := 0; n < b.N; n++ {
		rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", n%m))
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
				rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", (j+n)%m))
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
				rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", (j+n)%m))
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
				rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", (j+n)%m))
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
				rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", (j+n)%m))
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
				rl.AllowContext(context.Background(), fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}

func TestClientRateLimiterRetryAfter(t *testing.T) {
	l := 10
	window := 5 * time.Second
	rl := newClientRateLimiter(l, window)

	for i := 0; i < l; i++ {
		if rl.RetryAfter("foo") != 0 {
			t.Errorf("shouldn't have waiting time")
		}
		rl.AllowContext(context.Background(), "foo")
	}
	if rl.RetryAfter("foo") != 5 {
		t.Errorf("should wait for the exact time window")
	}
	if rl.RetryAfter("bar") != 0 {
		t.Errorf("shouldn't have waiting time")
	}
	rl.Close()
}

func TestClientRateLimiterRetryAfterConcurrent(t *testing.T) {
	l := 3
	window := 1 * time.Second
	rl := newClientRateLimiter(l, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) == 0 {
			t.Errorf("%v should return waiting time", s)
		}

		time.Sleep(window)

		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) != 0 {
			t.Errorf("%v shouldn't have waiting time", s)
		}
		rl.AllowContext(context.Background(), s)
		if rl.RetryAfter(s) == 0 {
			t.Errorf("%v should return waiting time", s)
		}

		wg.Done()
	}
	go f("foo")
	go f("bar")
	go f("baz")
	wg.Wait()
	rl.Close()
}

func TestClientRateLimiterRetryAfterMassiveConcurrent(t *testing.T) {
	l := 1 << 21
	n := 1 << 10
	window := 1 * time.Second
	rl := newClientRateLimiter(l, window)
	var wg sync.WaitGroup
	wg.Add(3)
	f := func(s string) {
		for i := 0; i < n; i++ {
			if rl.RetryAfter(s) != 0 {
				t.Errorf("%v shouldn't have waiting time", s)
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

func BenchmarkClientRateLimiterRetryAfter(b *testing.B) {
	window := 1 * time.Second
	rl := newClientRateLimiter(1<<21, window)
	for n := 0; n < b.N; n++ {
		rl.RetryAfter("foo")
	}
	rl.Close()
}

func BenchmarkClientRateLimiterRetryAfterConcurrent1(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.RetryAfter(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}

func BenchmarkClientRateLimiterRetryAfterConcurrent10(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.RetryAfter(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}

func BenchmarkClientRateLimiterRetryAfterConcurrent100(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.RetryAfter(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}

func BenchmarkClientRateLimiterRetryAfterConcurrent1000(b *testing.B) {
	var wg sync.WaitGroup
	window := time.Second
	rl := newClientRateLimiter(10, window)
	m := 100
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(j int) {
			for n := 0; n < b.N; n++ {
				rl.RetryAfter(fmt.Sprintf("foo%d", (j+n)%m))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	rl.Close()
}
