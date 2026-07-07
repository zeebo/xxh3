# XXH3
[![GoDoc](https://godoc.org/github.com/zeebo/xxh3?status.svg)](https://godoc.org/github.com/zeebo/xxh3)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/xxh3/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/xxh3?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/xxh3)](https://goreportcard.com/report/github.com/zeebo/xxh3)

This package is a port of the [xxh3](https://github.com/Cyan4973/xxHash) library to Go.

Upstream has fixed the output as of v0.8.0, and this package matches that.

---

# Benchmarks

## AMD Ryzen 9 9950X

### Small Sizes

| Bytes      | 64-bit                       | 64-bit Seed                  | 128-bit                      | 128-bit Seed                 |
|------------|------------------------------|------------------------------|------------------------------|------------------------------|
| ` 0 `      | ` 1.16 ns `                  | ` 1.72 ns `                  | ` 1.31 ns `                  | ` 2.13 ns `                  |
| ` 1-3 `    | ` 1.58 ns (0.63-1.90 GB/s) ` | ` 1.79 ns (0.56-1.68 GB/s) ` | ` 1.92 ns (0.52-1.56 GB/s) ` | ` 2.32 ns (0.43-1.29 GB/s) ` |
| ` 4-8 `    | ` 1.56 ns (2.56-5.13 GB/s) ` | ` 1.80 ns (2.22-4.44 GB/s) ` | ` 1.81 ns (2.21-4.42 GB/s) ` | ` 2.15 ns (1.86-3.72 GB/s) ` |
| ` 9-16 `   | ` 1.38 ns (6.52-11.6 GB/s) ` | ` 1.71 ns (5.26-9.36 GB/s) ` | ` 2.25 ns (4.00-7.11 GB/s) ` | ` 2.50 ns (3.60-6.40 GB/s) ` |
| ` 17-32 `  | ` 1.72 ns (9.88-18.6 GB/s) ` | ` 2.07 ns (8.21-15.5 GB/s) ` | ` 2.69 ns (6.32-11.9 GB/s) ` | ` 3.05 ns (5.57-10.5 GB/s) ` |
| ` 33-64 `  | ` 2.39 ns (13.8-26.8 GB/s) ` | ` 2.76 ns (12.0-23.2 GB/s) ` | ` 3.47 ns (9.51-18.4 GB/s) ` | ` 3.87 ns (8.53-16.5 GB/s) ` |
| ` 65-96 `  | ` 3.16 ns (20.6-30.4 GB/s) ` | ` 3.63 ns (17.9-26.4 GB/s) ` | ` 4.46 ns (14.6-21.5 GB/s) ` | ` 4.95 ns (13.1-19.4 GB/s) ` |
| ` 97-128 ` | ` 3.89 ns (24.9-32.9 GB/s) ` | ` 4.41 ns (22.0-29.0 GB/s) ` | ` 5.22 ns (18.6-24.5 GB/s) ` | ` 5.75 ns (16.9-22.3 GB/s) ` |

### Large Sizes

The seed and 128-bit variants share the same bulk loop and land within ~1% of
the figures below on these sizes, so only the 64-bit default is shown.

| Bytes     | Rate                       | SSE2 Rate                  | AVX2 Rate                  | AVX512 Rate                 |
|-----------|----------------------------|----------------------------|----------------------------|-----------------------------|
| ` 129 `   | ` 4.52 ns/op (28.5 GB/s) ` |                            |                            |                             |
| ` 240 `   | ` 7.39 ns/op (32.5 GB/s) ` |                            |                            |                             |
| ` 241 `   | ` 13.6 ns/op (17.7 GB/s) ` | ` 7.95 ns/op (30.3 GB/s) ` | ` 11.3 ns/op (21.4 GB/s) ` |                             |
| ` 512 `   | ` 23.1 ns/op (22.2 GB/s) ` | ` 11.2 ns/op (45.7 GB/s) ` | ` 14.1 ns/op (36.4 GB/s) ` |                             |
| ` 1024 `  | ` 42.1 ns/op (24.3 GB/s) ` | ` 20.7 ns/op (49.5 GB/s) ` | ` 20.9 ns/op (49.0 GB/s) ` | ` 23.9 ns/op (42.9 GB/s) `  |
| ` 100KB ` | ` 4.03 us/op (25.4 GB/s) ` | ` 2.10 us/op (48.7 GB/s) ` | ` 1.07 us/op (96.0 GB/s) ` | ` 685 ns/op (149 GB/s) `    |

## Intel i7-8850H @ 2.60GHz

### Small Sizes

| Bytes      | Rate                                   |
|------------|----------------------------------------|
| ` 0 `      | ` 0.74 ns/op `                         |
| ` 1-3 `    | ` 4.19 ns/op (0.24 GB/s - 0.71 GB/s) ` |
| ` 4-8 `    | ` 4.16 ns/op (0.97 GB/s - 1.98 GB/s) ` |
| ` 9-16 `   | ` 4.46 ns/op (2.02 GB/s - 3.58 GB/s) ` |
| ` 17-32 `  | ` 6.22 ns/op (2.76 GB/s - 5.15 GB/s) ` |
| ` 33-64 `  | ` 8.00 ns/op (4.13 GB/s - 8.13 GB/s) ` |
| ` 65-96 `  | ` 11.0 ns/op (5.91 GB/s - 8.84 GB/s) ` |
| ` 97-128 ` | ` 12.8 ns/op (7.68 GB/s - 10.0 GB/s) ` |

### Large Sizes

| Bytes     | Rate                       | SSE2 Rate                  | AVX2 Rate                  |
|-----------|----------------------------|----------------------------|----------------------------|
| ` 129 `   | ` 13.6 ns/op (9.45 GB/s) ` |                            |                            |
| ` 240 `   | ` 23.8 ns/op (10.1 GB/s) ` |                            |                            |
| ` 241 `   | ` 40.5 ns/op (5.97 GB/s) ` | ` 23.3 ns/op (10.4 GB/s) ` | ` 20.1 ns/op (12.0 GB/s) ` |
| ` 512 `   | ` 69.8 ns/op (7.34 GB/s) ` | ` 30.4 ns/op (16.9 GB/s) ` | ` 24.7 ns/op (20.7 GB/s) ` |
| ` 1024 `  | ` 132  ns/op (7.77 GB/s) ` | ` 48.9 ns/op (20.9 GB/s) ` | ` 37.7 ns/op (27.2 GB/s) ` |
| ` 100KB ` | ` 13.0 us/op (7.88 GB/s) ` | ` 4.05 us/op (25.3 GB/s) ` | ` 2.31 us/op (44.3 GB/s) ` |
