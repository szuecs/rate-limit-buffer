# rate limit buffer

Not a super profiled data structure, but seems to be ok for an http router

    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimitted                          1000000              2410 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2                        1000000              1584 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4                        1000000              1603 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8                        1000000              1606 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimittedConcurrent10                50000             22425 ns/op            4581 B/op         38 allocs/op
    BenchmarkIsRateLimittedConcurrent10-2             200000             47161 ns/op            3607 B/op         28 allocs/op
    BenchmarkIsRateLimittedConcurrent10-4             100000             13162 ns/op             795 B/op         25 allocs/op
    BenchmarkIsRateLimittedConcurrent10-8             200000             10312 ns/op             821 B/op         25 allocs/op
    BenchmarkIsRateLimittedConcurrent100               10000            162948 ns/op            1988 B/op        202 allocs/op
    BenchmarkIsRateLimittedConcurrent100-2             20000             78815 ns/op            2004 B/op        202 allocs/op
    BenchmarkIsRateLimittedConcurrent100-4             20000             67206 ns/op            2709 B/op        209 allocs/op
    BenchmarkIsRateLimittedConcurrent100-8             20000             71729 ns/op            3003 B/op        212 allocs/op
    BenchmarkIsRateLimittedConcurrent1000               1000           1375966 ns/op           16819 B/op       2004 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-2             2000            818376 ns/op           16855 B/op       2006 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-4             2000            776018 ns/op           32331 B/op       2167 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-8             2000            681790 ns/op           48675 B/op       2337 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     118.149s
