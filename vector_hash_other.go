// +build !amd64

package xxh3

const (
	avx2 = false
	sse2 = false
)

func hash_vector(s string) u64 { return 0 }
