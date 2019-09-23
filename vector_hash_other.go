// +build !amd64

package xxh3

const (
	avx2 = false
	sse2 = false
)

func hash_vector(_ ptr, _ u64) u64 { return 0 }
