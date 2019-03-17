# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go. It is (currently) pure Go and performs best when strings are 128 bytes or shorter.

---

```
goos: linux
goarch: amd64

Benchmark/0-8              1000000000          2.29 ns/op
Benchmark/1-8               500000000          3.60 ns/op       277.46 MB/s
Benchmark/2-8               500000000          3.60 ns/op       555.01 MB/s
Benchmark/3-8               500000000          3.64 ns/op       825.30 MB/s
Benchmark/4-8               500000000          3.26 ns/op      1226.62 MB/s
Benchmark/5-8               500000000          3.25 ns/op      1537.03 MB/s
Benchmark/6-8               500000000          3.25 ns/op      1843.53 MB/s
Benchmark/7-8               500000000          3.29 ns/op      2125.94 MB/s
Benchmark/8-8               500000000          3.24 ns/op      2468.29 MB/s
Benchmark/9-8               500000000          3.46 ns/op      2598.92 MB/s
Benchmark/10-8              500000000          3.47 ns/op      2880.75 MB/s
Benchmark/11-8              500000000          3.46 ns/op      3174.78 MB/s
Benchmark/12-8              500000000          3.47 ns/op      3459.33 MB/s
Benchmark/13-8              500000000          3.46 ns/op      3760.17 MB/s
Benchmark/14-8              500000000          3.45 ns/op      4063.50 MB/s
Benchmark/15-8              500000000          3.46 ns/op      4332.27 MB/s
Benchmark/16-8              500000000          3.45 ns/op      4633.76 MB/s
Benchmark/17-8              300000000          4.49 ns/op      3786.75 MB/s
Benchmark/18-8              300000000          4.46 ns/op      4032.63 MB/s
Benchmark/19-8              300000000          4.49 ns/op      4233.76 MB/s
Benchmark/20-8              300000000          4.50 ns/op      4445.07 MB/s
Benchmark/21-8              300000000          4.45 ns/op      4723.27 MB/s
Benchmark/22-8              300000000          4.47 ns/op      4924.00 MB/s
Benchmark/23-8              300000000          4.45 ns/op      5167.80 MB/s
Benchmark/24-8              300000000          4.51 ns/op      5325.97 MB/s
Benchmark/32-8              300000000          4.47 ns/op      7164.10 MB/s
Benchmark/40-8              200000000          6.10 ns/op      6553.21 MB/s
Benchmark/48-8              200000000          6.15 ns/op      7805.37 MB/s
Benchmark/56-8              200000000          6.15 ns/op      9102.46 MB/s
Benchmark/64-8              200000000          6.16 ns/op     10394.88 MB/s
Benchmark/72-8              200000000          7.99 ns/op      9016.52 MB/s
Benchmark/80-8              200000000          8.01 ns/op      9985.19 MB/s
Benchmark/88-8              200000000          7.98 ns/op     11034.41 MB/s
Benchmark/96-8              200000000          7.96 ns/op     12063.13 MB/s
Benchmark/104-8             200000000          9.93 ns/op     10468.24 MB/s
Benchmark/112-8             200000000          9.88 ns/op     11339.29 MB/s
Benchmark/120-8             200000000          9.88 ns/op     12142.79 MB/s
Benchmark/128-8             200000000          9.91 ns/op     12917.24 MB/s
Benchmark/256-8             50000000           31.9 ns/op      8018.31 MB/s
Benchmark/512-8             30000000           55.8 ns/op      9168.85 MB/s
Benchmark/1024-8            20000000            111 ns/op      9167.42 MB/s
Benchmark/2048-8            10000000            215 ns/op      9499.29 MB/s
Benchmark/4096-8             3000000            430 ns/op      9513.85 MB/s
Benchmark/8192-8             2000000            840 ns/op      9750.47 MB/s
Benchmark/16384-8            1000000           1670 ns/op      9806.26 MB/s
Benchmark/32768-8             500000           3404 ns/op      9625.70 MB/s
Benchmark/65536-8             200000           6664 ns/op      9833.05 MB/s
Benchmark/131072-8            100000          13430 ns/op      9759.49 MB/s
Benchmark/262144-8             50000          27001 ns/op      9708.65 MB/s
Benchmark/524288-8             30000          54225 ns/op      9668.70 MB/s
Benchmark/1048576-8            10000         107726 ns/op      9733.72 MB/s
```