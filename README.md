# rate limit buffer

Not a super profiled data structure, but seems to be ok for an http
router. I dropped previous concurrency tests, because these were
miss leading (measuring goroutine spawn in the benchmark loop)

    # 8e146dfb7ae73ae1d80c25108d64d5431c1efbe4
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkAllow                          1000000              3192 ns/op             242 B/op          4 allocs/op
    BenchmarkAllow-2                        1000000              1715 ns/op             243 B/op          4 allocs/op
    BenchmarkAllow-4                        1000000              1661 ns/op             243 B/op          4 allocs/op
    BenchmarkAllow-8                        1000000              1430 ns/op             243 B/op          4 allocs/op


    # 4879455378dd0207c7a41ec4256023562d4529cb
    % go test -bench=. -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkAllow                          1000000              2501 ns/op             243 B/op          4 allocs/op
    BenchmarkAllow-2                        1000000              1468 ns/op             243 B/op          4 allocs/op
    BenchmarkAllow-4                        1000000              1708 ns/op             243 B/op          4 allocs/op
    BenchmarkAllow-8                        1000000              1405 ns/op             242 B/op          4 allocs/op

    # Preprovisioned data set: 67bb879964b58763a7ca3d429f10fe5f7bbe644c
    % go test -bench=BenchmarkAllowBaseData1 -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkAllowBaseData1                 5000000               315 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData1-2               5000000               296 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData1-4               5000000               290 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData1-8               5000000               288 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData10                5000000               319 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData10-2              5000000               287 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData10-4              5000000               283 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData10-8              5000000               287 ns/op              14 B/op          1 allocs/op
    BenchmarkAllowBaseData100               5000000               311 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowBaseData100-2             5000000               305 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowBaseData100-4             5000000               299 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowBaseData100-8             5000000               303 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowBaseData1000              3000000               439 ns/op              21 B/op          2 allocs/op
    BenchmarkAllowBaseData1000-2            3000000               404 ns/op              21 B/op          2 allocs/op
    BenchmarkAllowBaseData1000-4            3000000               396 ns/op              21 B/op          2 allocs/op
    BenchmarkAllowBaseData1000-8            3000000               403 ns/op              21 B/op          2 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     53.970s

I am unsure how to measure with benchmark correctly the concurrent
behavior of this data structure. It's clear that high contention will
result in less throughput. 1,10,100,1000 are the number of goroutines
working concurrently on this data. 300µs/op is ok for 1000 goroutines
reading and writing the dataset, compared to do a proxy http call it
is much faster and will almost never happen at this concurrency level.

    # Concurrent/Contention tests: efcb5016a8a33a8f271ae9a5a527a0433ed727e4
    % go test -bench=BenchmarkAllowConcurrent1 -benchmem -cpu 1,2,4,8
    goos: darwin
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkAllowConcurrent1               5000000               364 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowConcurrent1-2             5000000               319 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowConcurrent1-4             5000000               313 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowConcurrent1-8             5000000               332 ns/op              15 B/op          1 allocs/op
    BenchmarkAllowConcurrent10               500000              3383 ns/op             159 B/op         19 allocs/op
    BenchmarkAllowConcurrent10-2             500000              2204 ns/op             159 B/op         19 allocs/op
    BenchmarkAllowConcurrent10-4             500000              3325 ns/op             159 B/op         19 allocs/op
    BenchmarkAllowConcurrent10-8             500000              2582 ns/op             159 B/op         19 allocs/op
    BenchmarkAllowConcurrent100               50000             34248 ns/op            1593 B/op        199 allocs/op
    BenchmarkAllowConcurrent100-2            100000             21281 ns/op            1592 B/op        199 allocs/op
    BenchmarkAllowConcurrent100-4             50000             36870 ns/op            1593 B/op        199 allocs/op
    BenchmarkAllowConcurrent100-8             50000             35154 ns/op            1593 B/op        199 allocs/op
    BenchmarkAllowConcurrent1000               5000            319270 ns/op           15928 B/op       1990 allocs/op
    BenchmarkAllowConcurrent1000-2            10000            230893 ns/op           15936 B/op       1990 allocs/op
    BenchmarkAllowConcurrent1000-4             5000            344456 ns/op           15946 B/op       1990 allocs/op
    BenchmarkAllowConcurrent1000-8             3000            333466 ns/op           15996 B/op       1990 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     42.174s


    # Concurrent/Contention tests: 648a55836c24f8a3db74d710b666bbfa085cbe2a
    % go test -bench=Allow -cpu 1,2,4,8 >old.txt
    # Concurrent/Contention tests: 380fb9d6fe737f032ce3a1aa4f6881a60b10b273
    % go test -bench=Allow -cpu 1,2,4,8 >new.txt
    % benchcmp old.txt new.txt
    benchmark                                 old ns/op     new ns/op     delta
    BenchmarkAllow                            1475          1552          +5.22%
    BenchmarkAllow-2                          1317          1225          -6.99%
    BenchmarkAllow-4                          1488          1249          -16.06%
    BenchmarkAllow-8                          1289          1270          -1.47%
    BenchmarkAllowBaseData1                   368           316           -14.13%
    BenchmarkAllowBaseData1-2                 363           301           -17.08%
    BenchmarkAllowBaseData1-4                 326           299           -8.28%
    BenchmarkAllowBaseData1-8                 340           284           -16.47%
    BenchmarkAllowBaseData10                  339           334           -1.47%
    BenchmarkAllowBaseData10-2                319           305           -4.39%
    BenchmarkAllowBaseData10-4                311           298           -4.18%
    BenchmarkAllowBaseData10-8                318           304           -4.40%
    BenchmarkAllowBaseData100                 357           335           -6.16%
    BenchmarkAllowBaseData100-2               332           318           -4.22%
    BenchmarkAllowBaseData100-4               336           320           -4.76%
    BenchmarkAllowBaseData100-8               342           303           -11.40%
    BenchmarkAllowBaseData1000                554           509           -8.12%
    BenchmarkAllowBaseData1000-2              583           411           -29.50%
    BenchmarkAllowBaseData1000-4              510           431           -15.49%
    BenchmarkAllowBaseData1000-8              503           410           -18.49%
    BenchmarkAllowConcurrent1                 355           355           +0.00%
    BenchmarkAllowConcurrent1-2               340           333           -2.06%
    BenchmarkAllowConcurrent1-4               335           329           -1.79%
    BenchmarkAllowConcurrent1-8               336           332           -1.19%
    BenchmarkAllowConcurrent10                3836          3544          -7.61%
    BenchmarkAllowConcurrent10-2              3042          1752          -42.41%
    BenchmarkAllowConcurrent10-4              3604          1597          -55.69%
    BenchmarkAllowConcurrent10-8              3675          1856          -49.50%
    BenchmarkAllowConcurrent100               35028         35378         +1.00%
    BenchmarkAllowConcurrent100-2             36003         19574         -45.63%
    BenchmarkAllowConcurrent100-4             55464         16719         -69.86%
    BenchmarkAllowConcurrent100-8             54492         17668         -67.58%
    BenchmarkAllowConcurrent1000              389221        348280        -10.52%
    BenchmarkAllowConcurrent1000-2            392405        197976        -49.55%
    BenchmarkAllowConcurrent1000-4            849330        173836        -79.53%
    BenchmarkAllowConcurrent1000-8            1017861       187116        -81.62%
    BenchmarkAllowConcurrentAddDelete10       4174          3808          -8.77%
    BenchmarkAllowConcurrentAddDelete10-2     3425          2167          -36.73%
    BenchmarkAllowConcurrentAddDelete10-4     5709          1644          -71.20%
    BenchmarkAllowConcurrentAddDelete10-8     5133          2141          -58.29%
