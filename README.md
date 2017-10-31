# rate limit buffer

[![Build Status](https://travis-ci.org/szuecs/rate-limit-buffer.svg)](https://travis-ci.org/szuecs/rate-limit-buffer)
[![GoDoc](https://godoc.org/github.com/szuecs/rate-limit-buffer?status.svg)](https://godoc.org/github.com/szuecs/rate-limit-buffer)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/szuecs/rate-limit-buffer)](https://goreportcard.com/report/szuecs/rate-limit-buffer)
[![Coverage Status](https://coveralls.io/repos/github/szuecs/rate-limit-buffer/badge.svg?branch=master)](https://coveralls.io/github/szuecs/rate-limit-buffer?branch=master)

RateLimiter is an interface you want to use in a protect your
backend from too many calls.

There are the following implementations:
- CircularBuffer: NewRateLimiter(int, time.Duration) RateLimiter
- ClientRateLimiter: NewClientRateLimiter(int, time.Duration) *ClientRateLimiter

CircularBuffer is a rate limiter that can only protect a backend from
maximum number of calls. It has no idea about clients or
connections. Allow(string) will ignore the string parameter, which is
only used by ClientRateLimiter to match the bucket of a client. This
CircularBuffer does not need to be refilled, because it is aware of
the oldest of call time.Time, which will be used to check the
time.Duration passed.

ClientRateLimiter implements the RateLimiter
interface. ClientRateLimiter is a rate limit data structure with a
simple interface and is able to do different kinds of rate limits. For
example it can be used to do client based rate limits, where each
client is independently counted and rate limited. The normal case in
the literature is implemented as leaky bucket data structure, which
can only implement the case were you want to protect your backend to
get only N number requests per duration.  Client based rate limits can
be used to slow down user/password enumeration attacks, protect DDoS
attacks that do not fill the pipe, but your software proxy.

## Upgrade v0.1.x to v0.2.y

There is a breaking change, which does not break the interface.
People that want to rate limit per string passed to "Allow(string)
bool" have to change to ClientRateLimiter to get the same results.

## Benchmarks

### v0.2.*

    % go test -bench=. -benchmem -cpu 1,2,4,8 | tee -a v0.1.3.txt
    goos: linux
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer
    BenchmarkRateLimiterAllow                                       10000000               155 ns/op               5 B/op          0 allocs/op
    BenchmarkRateLimiterAllow-2                                     10000000               156 ns/op               5 B/op          0 allocs/op
    BenchmarkRateLimiterAllow-4                                     10000000               153 ns/op               5 B/op          0 allocs/op
    BenchmarkRateLimiterAllow-8                                     10000000               152 ns/op               5 B/op          0 allocs/op
    BenchmarkClientRateLimiterAllow                                  1000000              1545 ns/op             143 B/op          4 allocs/op
    BenchmarkClientRateLimiterAllow-2                                1000000              1430 ns/op             143 B/op          4 allocs/op
    BenchmarkClientRateLimiterAllow-4                                1000000              1402 ns/op             143 B/op          4 allocs/op
    BenchmarkClientRateLimiterAllow-8                                1000000              1493 ns/op             143 B/op          4 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1                         5000000               353 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1-2                       5000000               344 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1-4                       5000000               334 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1-8                       5000000               328 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData10                        5000000               354 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData10-2                      5000000               326 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData10-4                      5000000               365 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData10-8                      5000000               330 ns/op              14 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData100                       5000000               362 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData100-2                     5000000               353 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData100-4                     5000000               338 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData100-8                     5000000               348 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1000                      2000000               609 ns/op              24 B/op          2 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1000-2                    2000000               512 ns/op              24 B/op          2 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1000-4                    2000000               533 ns/op              24 B/op          2 allocs/op
    BenchmarkClientRateLimiterAllowBaseData1000-8                    2000000               514 ns/op              24 B/op          2 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1                       5000000               403 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1-2                     5000000               344 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1-4                     5000000               359 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1-8                     5000000               362 ns/op              15 B/op          1 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent10                       500000              3805 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent10-2                    1000000              2190 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent10-4                    1000000              1826 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent10-8                    1000000              1982 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent100                       50000             37318 ns/op            1593 B/op        199 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent100-2                     50000             28760 ns/op            1593 B/op        199 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent100-4                    100000             18575 ns/op            1593 B/op        199 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent100-8                    100000             20966 ns/op            1594 B/op        199 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1000                       5000            350560 ns/op           15928 B/op       1990 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1000-2                    10000            252966 ns/op           15932 B/op       1990 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1000-4                    10000            192886 ns/op           15966 B/op       1990 allocs/op
    BenchmarkClientRateLimiterAllowConcurrent1000-8                    10000            220237 ns/op           15973 B/op       1990 allocs/op
    BenchmarkClientRateLimiterAllowConcurrentAddDelete10              500000              4018 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrentAddDelete10-2            500000              3018 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrentAddDelete10-4           1000000              1813 ns/op             159 B/op         19 allocs/op
    BenchmarkClientRateLimiterAllowConcurrentAddDelete10-8            500000              2181 ns/op             159 B/op         19 allocs/op
    PASS
    ok      github.com/szuecs/rate-limit-buffer     124.119s

### < v0.2.*

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

## v0.1.2 to v0.1.3

I fixed a race condition in high concurrent situations
(go test -race -v -run=Massive .), which has a small impact on
benchmarks. The system is slower under low concurrent with less CPUs
and faster in high concurrent with more CPUs:

    # Data collected with v0.1.2 and v0.1.3
    % go test -bench=. -benchmem -cpu 1,2,4,8 | tee -a v0.1.3.txt
    goos: linux
    goarch: amd64
    pkg: github.com/szuecs/rate-limit-buffer

    # Data compare
    % benchcmp -changed v0.1.2.txt v0.1.3.txt
    benchmark                                 old ns/op     new ns/op     delta
    BenchmarkAllow                            1433          1585          +10.61%
    BenchmarkAllow-2                          1334          1386          +3.90%
    BenchmarkAllow-4                          1317          1363          +3.49%
    BenchmarkAllow-8                          1262          1399          +10.86%
    BenchmarkAllowBaseData1                   340           353           +3.82%
    BenchmarkAllowBaseData1-2                 317           325           +2.52%
    BenchmarkAllowBaseData1-4                 320           338           +5.62%
    BenchmarkAllowBaseData1-8                 312           333           +6.73%
    BenchmarkAllowBaseData10                  340           355           +4.41%
    BenchmarkAllowBaseData10-2                311           324           +4.18%
    BenchmarkAllowBaseData10-4                312           333           +6.73%
    BenchmarkAllowBaseData10-8                319           335           +5.02%
    BenchmarkAllowBaseData100                 344           353           +2.62%
    BenchmarkAllowBaseData100-2               335           337           +0.60%
    BenchmarkAllowBaseData100-4               401           333           -16.96%
    BenchmarkAllowBaseData100-8               323           333           +3.10%
    BenchmarkAllowBaseData1000                525           539           +2.67%
    BenchmarkAllowBaseData1000-2              505           511           +1.19%
    BenchmarkAllowBaseData1000-4              501           552           +10.18%
    BenchmarkAllowBaseData1000-8              450           517           +14.89%
    BenchmarkAllowConcurrent1                 355           363           +2.25%
    BenchmarkAllowConcurrent1-2               336           344           +2.38%
    BenchmarkAllowConcurrent1-4               332           346           +4.22%
    BenchmarkAllowConcurrent1-8               339           348           +2.65%
    BenchmarkAllowConcurrent10                3575          3634          +1.65%
    BenchmarkAllowConcurrent10-2              2165          2296          +6.05%
    BenchmarkAllowConcurrent10-4              2163          1793          -17.11%
    BenchmarkAllowConcurrent10-8              2148          1995          -7.12%
    BenchmarkAllowConcurrent100               34943         35992         +3.00%
    BenchmarkAllowConcurrent100-2             24132         22516         -6.70%
    BenchmarkAllowConcurrent100-4             18392         18662         +1.47%
    BenchmarkAllowConcurrent100-8             19484         19960         +2.44%
    BenchmarkAllowConcurrent1000              338325        352712        +4.25%
    BenchmarkAllowConcurrent1000-2            218878        234376        +7.08%
    BenchmarkAllowConcurrent1000-4            184612        194888        +5.57%
    BenchmarkAllowConcurrent1000-8            216032        215727        -0.14%
    BenchmarkAllowConcurrentAddDelete10       3701          3687          -0.38%
    BenchmarkAllowConcurrentAddDelete10-2     2428          2427          -0.04%
    BenchmarkAllowConcurrentAddDelete10-4     1787          1812          +1.40%
    BenchmarkAllowConcurrentAddDelete10-8     2059          2126          +3.25%

    benchmark                          old bytes     new bytes     delta
    BenchmarkAllow                     128           143           +11.72%
    BenchmarkAllow-2                   127           143           +12.60%
    BenchmarkAllow-4                   127           143           +12.60%
    BenchmarkAllow-8                   127           143           +12.60%
    BenchmarkAllowBaseData1000-8       21            24            +14.29%
    BenchmarkAllowConcurrent100        1592          1593          +0.06%
    BenchmarkAllowConcurrent1000-2     15928         15933         +0.03%
    BenchmarkAllowConcurrent1000-8     16026         15973         -0.33%
