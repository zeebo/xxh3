# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go. It is (currently) pure Go and performs best when strings are 128 bytes or shorter.

---

Lower bytes | Upper bytes | ns/op | Lower rate | Upper rate
-----------------------------
0 | 0 | 1.91 | - | -
1 | 3 | 3.00 | 333 MB/s | 999 MB/s
4 | 8 | 2.75 | 1.45 GB/s | 2.9 GB/s
9 | 16 | 2.88 | 3.13 GB/s | 5.55 GB/s
17 | 32 | 4.31 | 3.95 GB/s | 7.43 GB/s
33 | 64 | 5.98 | - | 10.7 GB/s
65 | 96 | 7.78 | - | 12.4 GB/s
97 | 128 | 9.25 | - | 13.8 GB/s
129 | - | - | 7.73 GB/s | 9.8 GB/s