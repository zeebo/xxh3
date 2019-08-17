# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go.

---

# Benchmarks

Run on my `i7-6700K CPU @ 4.00GHz`

## Small Sizes

| Bytes     | Rate                                 |
|-----------|--------------------------------------|
|` 0 `      |` 0.49 ns/op `                       |
|` 1-3 `    |` 3.49 ns/op (0.28 GB/s - 0.86 GB/s) `|
|` 4-8 `    |` 4.20 ns/op (0.94 GB/s - 1.91 GB/s) `|
|` 9-16 `   |` 3.60 ns/op (2.33 GB/s - 4.45 GB/s) `|
|` 17-32 `  |` 5.01 ns/op (3.39 GB/s - 6.32 GB/s) `|
|` 33-64 `  |` 6.71 ns/op (4.92 GB/s - 9.50 GB/s) `|
|` 65-96 `  |` 9.08 ns/op (7.16 GB/s - 10.5 GB/s) `|
|` 97-128 ` |` 11.0 ns/op (8.80 GB/s - 11.6 GB/s) `|

## Large Sizes

| Bytes   | Rate                     | SSE2 Rate                | AVX2 Rate                |
|---------|--------------------------|--------------------------|--------------------------|
|` 129 `  |` 11.5 ns/op (11.2 GB/s) `|                          |                          |
|` 240 `  |` 21.3 ns/op (11.2 GB/s) `|                          |                          |
|` 241 `  |` 26.7 ns/op (9.02 GB/s) `|` 19.8 ns/op (12.2 GB/s) `|` 16.3 ns/op (14.8 GB/s) `|
|` 512 `  |` 44.4 ns/op (11.5 GB/s) `|` 28.5 ns/op (18.0 GB/s) `|` 21.1 ns/op (24.3 GB/s) `|
|` 1024 ` |` 85.1 ns/op (12.0 GB/s) `|` 45.9 ns/op (22.3 GB/s) `|` 30.1 ns/op (34.0 GB/s) `|
|` 100KB `|` 7961 ns/op (12.9 GB/s) `|` 3515 ns/op (29.1 GB/s) `|` 1856 ns/op (55.1 GB/s) `|
