module github.com/zeebo/xxh3/internal/compare

go 1.22

require (
	github.com/cespare/xxhash v1.1.0
	github.com/zeebo/xxh3 v0.0.0-20190706061215-e3a4ded7d14f
)

require (
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/zeebo/xxh3 => ../../
