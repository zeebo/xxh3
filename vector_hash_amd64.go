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

func hashLarge(p ptr, l u64) (acc u64) {
	acc = l * prime64_1

	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if avx2 {
		accumAVX2(&accs, p, key, l)
	} else if sse2 {
		accumSSE(&accs, p, key, l)
	} else {
		accumScalar(&accs, p, key, l)
	}

	// merge accs
	acc += mulFold64(accs[0]^key64_011, accs[1]^key64_019)
	acc += mulFold64(accs[2]^key64_027, accs[3]^key64_035)
	acc += mulFold64(accs[4]^key64_043, accs[5]^key64_051)
	acc += mulFold64(accs[6]^key64_059, accs[7]^key64_067)

	return xxh3Avalanche(acc)
}

func hashLarge128(p ptr, l u64) (acc u128) {
	acc.Lo = l * prime64_1
	acc.Hi = ^(l * prime64_2)

	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if avx2 {
		accumAVX2(&accs, p, key, l)
	} else if sse2 {
		accumSSE(&accs, p, key, l)
	} else {
		accumScalar(&accs, p, key, l)
	}

	// merge accs
	acc.Lo += mulFold64(accs[0]^key64_011, accs[1]^key64_019)
	acc.Lo += mulFold64(accs[2]^key64_027, accs[3]^key64_035)
	acc.Lo += mulFold64(accs[4]^key64_043, accs[5]^key64_051)
	acc.Lo += mulFold64(accs[6]^key64_059, accs[7]^key64_067)
	acc.Lo = xxh3Avalanche(acc.Lo)

	acc.Hi += mulFold64(accs[0]^key64_117, accs[1]^key64_125)
	acc.Hi += mulFold64(accs[2]^key64_133, accs[3]^key64_141)
	acc.Hi += mulFold64(accs[4]^key64_149, accs[5]^key64_157)
	acc.Hi += mulFold64(accs[6]^key64_165, accs[7]^key64_173)
	acc.Hi = xxh3Avalanche(acc.Hi)

	return acc
}
