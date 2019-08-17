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
func accum_avx(acc *[8]u64, data, key unsafe.Pointer, len u64)

//go:noescape
func accum_sse(acc *[8]u64, data, key unsafe.Pointer, len u64)

func hash_vector(p ptr, l u64) (acc u64) {
	acc = l * prime64_1
	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if avx2 {
		accum_avx(&accs, p, key, l)
	} else {
		accum_sse(&accs, p, key, l)
	}

	// merge accs
	hi1, lo1 := bits.Mul64(accs[0]^0x6dd4de1cad21f72c, accs[1]^0xa44072db979083e9)
	acc += hi1 ^ lo1
	hi2, lo2 := bits.Mul64(accs[2]^0xe679cb1f67b3b7a4, accs[3]^0xd05a8278e5c0cc4e)
	acc += hi2 ^ lo2
	hi3, lo3 := bits.Mul64(accs[4]^0x4608b82172ffcc7d, accs[5]^0x9035e08e2443f774)
	acc += hi3 ^ lo3
	hi4, lo4 := bits.Mul64(accs[6]^0x52283c4c263a81e6, accs[7]^0x65d088cb00c391bb)
	acc += hi4 ^ lo4

	// avalanche
	acc ^= acc >> 37
	acc *= prime64_3
	acc ^= acc >> 32

	return acc
}
