package xxh3

import (
	"math/bits"
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	avx2 = cpu.X86.HasAVX2
	sse2 = cpu.X86.HasSSE2
)

//go:noescape
func accum_avx(acc *[8]uint64, data, key unsafe.Pointer, len uint64)

//go:noescape
func accum_sse(acc *[8]uint64, data, key unsafe.Pointer, len uint64)

func hashVector(p ptr, l uint64) uint64 {
	acc := [8]uint64{0, prime64_1, prime64_2, prime64_3, prime64_4, prime64_5, 0, 0}

	if avx2 && false {
		accum_avx(&acc, p, key, l)
	} else {
		accum_sse(&acc, p, key, l)
	}

	// merge_accs
	result := l * prime64_1
	hi1, lo1 := bits.Mul64(acc[0]^*(*uint64)(ptr(ui(key) + 0*8)), acc[1]^*(*uint64)(ptr(ui(key) + 1*8)))
	result += hi1 ^ lo1
	hi2, lo2 := bits.Mul64(acc[2]^*(*uint64)(ptr(ui(key) + 2*8)), acc[3]^*(*uint64)(ptr(ui(key) + 3*8)))
	result += hi2 ^ lo2
	hi3, lo3 := bits.Mul64(acc[4]^*(*uint64)(ptr(ui(key) + 4*8)), acc[5]^*(*uint64)(ptr(ui(key) + 5*8)))
	result += hi3 ^ lo3
	hi4, lo4 := bits.Mul64(acc[6]^*(*uint64)(ptr(ui(key) + 6*8)), acc[7]^*(*uint64)(ptr(ui(key) + 7*8)))
	result += hi4 ^ lo4

	result ^= result >> 37
	result *= prime64_3
	result ^= result >> 32
	return result
}
