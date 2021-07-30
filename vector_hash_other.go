//go:build !amd64
// +build !amd64

package xxh3

func hashLarge(p ptr, l u64) (acc u64) {
	acc = l * prime64_1

	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	accumScalar(&accs, p, key, l)

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

	accumScalar(&accs, p, key, l)

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
