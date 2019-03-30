// +build !amd64

package xxh3

const (
	avx2 = false
	sse2 = false
)

func hashVector(p ptr, l uint64) uint64 { return 0 }
