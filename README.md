# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go. 

*Important note*: I have no idea if it matches the output, yet, and upstream is still iterating on the design, so don't use this in production.

---

## Small Sizes

| Lower bytes | Upper bytes | ns/op | Lower GB/s | Upper GB/s |
|-------------|-------------|-------|------------|------------|
| 0           | 0           | 1.92  | -          | -          |
| 1           | 3           | 3.13  | 0.32       | 9.56       |
| 4           | 8           | 3.08  | 1.30       | 2.61       |
| 9           | 16          | 3.00  | 3.00       | 5.30       |
| 17          | 32          | 3.82  | 4.45       | 8.36       |
| 33          | 64          | 5.65  | 5.84       | 11.6       |
| 65          | 96          | 7.20  | 8.91       | 13.3       |
| 97          | 128         | 8.70  | 11.0       | 14.7       |

## Large Sizes

| Bytes | ns/op (GB/s)      | SSE2             | AVX2             |
|-------|-------------------|------------------|------------------|
| 129   | 24.7 (5.22 GB/s)  | 16.6 (7.76 GB/s) | 14.8 (8.74 GB/s) |
| 256   | 30.9 (8.30 GB/s)  | 19.0 (13.4 GB/s) | 15.7 (16.3 GB/s) |
| 512   | 55.1 (9.30 GB/s)  | 27.5 (18.6 GB/s) | 19.9 (25.8 GB/s) |
| 1024  | 107 (9.56 GB/s)   | 43.4 (23.6 GB/s) | 29.4 (34.9 GB/s) | 
| 100KB | 10172 (10.1 GB/s) | 3479 (29.4 GB/s) | 1837 (55.7 GB/s) |
