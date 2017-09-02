# rate limit buffer

Not a super profiled data structure, but seems to be ok for an http router

    % go test -bench=. -benchmem -cpu 1,2,4,8
    BenchmarkIsRateLimitted          1000000              3131 ns/op             242 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2        1000000              1655 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4        1000000              1682 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8        1000000              1474 ns/op             243 B/op          4 allocs/op

