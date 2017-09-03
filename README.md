# rate limit buffer

Not a super profiled data structure, but seems to be ok for an http router

    # 8e146dfb7ae73ae1d80c25108d64d5431c1efbe4
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimitted                          1000000              3192 ns/op             242 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2                        1000000              1715 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4                        1000000              1661 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8                        1000000              1430 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimittedConcurrent1               1000000              2136 ns/op             548 B/op          4 allocs/op
    BenchmarkIsRateLimittedConcurrent1-2             1000000              5041 ns/op             469 B/op          3 allocs/op
    BenchmarkIsRateLimittedConcurrent1-4             1000000              9084 ns/op             393 B/op          2 allocs/op
    BenchmarkIsRateLimittedConcurrent1-8              500000              2685 ns/op             374 B/op          3 allocs/op
    BenchmarkIsRateLimittedConcurrent10                30000             47515 ns/op             681 B/op         22 allocs/op
    BenchmarkIsRateLimittedConcurrent10-2             100000             15181 ns/op             807 B/op         23 allocs/op
    BenchmarkIsRateLimittedConcurrent10-4             200000             13900 ns/op             954 B/op         25 allocs/op
    BenchmarkIsRateLimittedConcurrent10-8             200000              8321 ns/op            1076 B/op         26 allocs/op
    BenchmarkIsRateLimittedConcurrent100               10000            182509 ns/op            5287 B/op        236 allocs/op
    BenchmarkIsRateLimittedConcurrent100-2             20000             90301 ns/op            1989 B/op        202 allocs/op
    BenchmarkIsRateLimittedConcurrent100-4             20000             70309 ns/op            2219 B/op        204 allocs/op
    BenchmarkIsRateLimittedConcurrent100-8             20000             68141 ns/op            3735 B/op        220 allocs/op
    BenchmarkIsRateLimittedConcurrent1000               1000           1519939 ns/op           16818 B/op       2004 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-2             2000            858816 ns/op           16561 B/op       2003 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-4             2000            826575 ns/op           20100 B/op       2039 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-8             2000            704271 ns/op           49505 B/op       2346 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     274.247s


    # 4879455378dd0207c7a41ec4256023562d4529cb
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkIsRateLimitted                          1000000              2501 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-2                        1000000              1468 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-4                        1000000              1708 ns/op             243 B/op          4 allocs/op
    BenchmarkIsRateLimitted-8                        1000000              1405 ns/op             242 B/op          4 allocs/op
    BenchmarkIsRateLimittedConcurrent1               1000000              2062 ns/op             548 B/op          4 allocs/op
    BenchmarkIsRateLimittedConcurrent1-2             1000000              4720 ns/op             465 B/op          3 allocs/op
    BenchmarkIsRateLimittedConcurrent1-4             1000000              9356 ns/op             434 B/op          2 allocs/op
    BenchmarkIsRateLimittedConcurrent1-8             1000000             15033 ns/op             243 B/op          2 allocs/op
    BenchmarkIsRateLimittedConcurrent10               100000             51367 ns/op            2982 B/op         27 allocs/op
    BenchmarkIsRateLimittedConcurrent10-2             100000             16291 ns/op             692 B/op         23 allocs/op
    BenchmarkIsRateLimittedConcurrent10-4             100000             14418 ns/op             947 B/op         25 allocs/op
    BenchmarkIsRateLimittedConcurrent10-8             200000             11745 ns/op            1252 B/op         25 allocs/op
    BenchmarkIsRateLimittedConcurrent100               10000            435849 ns/op            5287 B/op        236 allocs/op
    BenchmarkIsRateLimittedConcurrent100-2             20000             82162 ns/op            1991 B/op        202 allocs/op
    BenchmarkIsRateLimittedConcurrent100-4             30000             61462 ns/op            2176 B/op        203 allocs/op
    BenchmarkIsRateLimittedConcurrent100-8             20000             70651 ns/op            3057 B/op        213 allocs/op
    BenchmarkIsRateLimittedConcurrent1000               1000           1306187 ns/op           16820 B/op       2004 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-2             2000            810523 ns/op           16557 B/op       2003 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-4             2000            673209 ns/op           20582 B/op       2044 allocs/op
    BenchmarkIsRateLimittedConcurrent1000-8             3000            718642 ns/op           40715 B/op       2254 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     337.055s
