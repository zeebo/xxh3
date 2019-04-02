# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go. 

**Important note**: I have no idea if it matches the output, yet, and upstream is still iterating on the design, so don't use this in production.

---

# Benchmarks

Run on my `i7-6700K CPU @ 4.00GHz`

## Small Sizes

| Bytes    | Rate                                 |
|----------|--------------------------------------|
|` 0 `     |` 1.92 ns/op `                        |
|` 1-3 `   |` 3.13 ns/op (0.32 GB/s - 0.96 GB/s) `|
|` 4-8 `   |` 3.08 ns/op (1.30 GB/s - 2.61 GB/s) `|
|` 9-16 `  |` 3.00 ns/op (3.00 GB/s - 5.30 GB/s) `|
|` 17-32 ` |` 3.82 ns/op (4.45 GB/s - 8.36 GB/s) `|
|` 33-64 ` |` 5.65 ns/op (5.84 GB/s - 11.6 GB/s) `|
|` 65-96 ` |` 7.20 ns/op (8.91 GB/s - 13.3 GB/s) `|
|` 97-128 `|` 8.70 ns/op (11.0 GB/s - 14.7 GB/s) `|

## Large Sizes

| Bytes   | Rate                      | SSE2 Rate                | AVX2 Rate                |
|---------|---------------------------|--------------------------|--------------------------|
|` 129 `  |` 24.7 ns/op (5.22 GB/s) ` |` 16.6 ns/op (7.76 GB/s) `|` 14.8 ns/op (8.74 GB/s) `|
|` 256 `  |` 30.9 ns/op (8.30 GB/s) ` |` 19.0 ns/op (13.4 GB/s) `|` 15.7 ns/op (16.3 GB/s) `|
|` 512 `  |` 55.1 ns/op (9.30 GB/s) ` |` 27.5 ns/op (18.6 GB/s) `|` 19.9 ns/op (25.8 GB/s) `|
|` 1024 ` |` 107 ns/op (9.56 GB/s) `  |` 43.4 ns/op (23.6 GB/s) `|` 29.4 ns/op (34.9 GB/s) `|
|` 100KB `|` 10172 ns/op (10.1 GB/s) `|` 3479 ns/op (29.4 GB/s) `|` 1837 ns/op (55.7 GB/s) `|
