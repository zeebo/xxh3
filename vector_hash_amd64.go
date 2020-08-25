package xxh3

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

//go:generate bash -c "go run github.com/zeebo/xxh3/avo -avx > vector_avx_amd64.s"
//go:generate bash -c "go run github.com/zeebo/xxh3/avo -sse > vector_sse_amd64.s"

var (
	avx2 = cpu.X86.HasAVX2
	sse2 = cpu.X86.HasSSE2
)

//go:noescape
func accumAVX2(acc *[8]u64, data, key unsafe.Pointer, len u64)

//go:noescape
func accumSSE(acc *[8]u64, data, key unsafe.Pointer, len u64)

func hashVector(p ptr, l u64, secret ptr) (acc u64) {
	acc = l * prime64_1
	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if avx2 {
		accumAVX2(&accs, p, secret, l)
	} else {
		accumSSE(&accs, p, secret, l)
	}

	// merge accs
	acc += mulFold64(accs[0]^readU64(secret, 11), accs[1]^readU64(secret, 19))
	acc += mulFold64(accs[2]^readU64(secret, 27), accs[3]^readU64(secret, 35))
	acc += mulFold64(accs[4]^readU64(secret, 43), accs[5]^readU64(secret, 51))
	acc += mulFold64(accs[6]^readU64(secret, 59), accs[7]^readU64(secret, 67))
	
	return xxh3Avalanche(acc)
}
