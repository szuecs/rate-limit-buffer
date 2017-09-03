# rate limit buffer

Not a super profiled data structure, but seems to be ok for an http
router. I dropped previous concurrency tests, because these were
miss leading (measuring goroutine spawn in the benchmark loop)

    # 8e146dfb7ae73ae1d80c25108d64d5431c1efbe4
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimitted                          1000000              3192 ns/op             242 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2                        1000000              1715 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4                        1000000              1661 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8                        1000000              1430 ns/op             243 B/op          4 allocs/op


    # 4879455378dd0207c7a41ec4256023562d4529cb
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimitted                          1000000              2501 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2                        1000000              1468 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4                        1000000              1708 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8                        1000000              1405 ns/op             242 B/op          4 allocs/op

    # Preprovisioned data set: 67bb879964b58763a7ca3d429f10fe5f7bbe644c
    % go test -bench=BenchmarkIsRateLimittedBaseData1 -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimittedBaseData1                 5000000               315 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData1-2               5000000               296 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData1-4               5000000               290 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData1-8               5000000               288 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData10                5000000               319 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData10-2              5000000               287 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData10-4              5000000               283 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData10-8              5000000               287 ns/op              14 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData100               5000000               311 ns/op              15 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData100-2             5000000               305 ns/op              15 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData100-4             5000000               299 ns/op              15 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData100-8             5000000               303 ns/op              15 B/op          1 allocs/op
    BenchmarkIsRateLimittedBaseData1000              3000000               439 ns/op              21 B/op          2 allocs/op
    BenchmarkIsRateLimittedBaseData1000-2            3000000               404 ns/op              21 B/op          2 allocs/op
    BenchmarkIsRateLimittedBaseData1000-4            3000000               396 ns/op              21 B/op          2 allocs/op
    BenchmarkIsRateLimittedBaseData1000-8            3000000               403 ns/op              21 B/op          2 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     53.970s

I am unsure how to measure with benchmark correctly the concurrent
behavior of this data structure.

    # Concurrent tests: 67bb879964b58763a7ca3d429f10fe5f7bbe644c
    % go test -bench=BenchmarkIsRateLimittedConcurrent10 -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimittedConcurrent10               500000              3319 ns/op             159 B/op         19 allocs/op
    BenchmarkIsRateLimittedConcurrent10-2            1000000              2097 ns/op             159 B/op         19 allocs/op
    BenchmarkIsRateLimittedConcurrent10-4             500000              3350 ns/op             159 B/op         19 allocs/op
    BenchmarkIsRateLimittedConcurrent10-8            1000000              2565 ns/op             159 B/op         19 allocs/op
    BenchmarkIsRateLimittedConcurrent100               50000             33120 ns/op            1593 B/op        199 allocs/op
    BenchmarkIsRateLimittedConcurrent100-2            100000             21246 ns/op            1592 B/op        199 allocs/op
    BenchmarkIsRateLimittedConcurrent100-4             30000             37877 ns/op            1593 B/op        199 allocs/op
    BenchmarkIsRateLimittedConcurrent100-8             50000             35584 ns/op            1594 B/op        199 allocs/op
    BenchmarkIsRateLimittedConcurrent1000               5000            315386 ns/op           15928 B/op       1990 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-2            10000            225854 ns/op           15936 B/op       1990 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-4             3000            459434 ns/op           15963 B/op       1990 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-8             5000            269482 ns/op           15974 B/op       1990 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     35.595s
