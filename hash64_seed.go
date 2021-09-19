package xxh3

import "math/bits"

const secret_size = 192

// HashSeed returns the hash of the byte slice with given seed.
func HashSeed(b []byte, seed uint64) uint64 {
	if len(b) <= 16 {
		return hashSmallSeed(*(*ptr)(ptr(&b)), len(b), seed)
	}
	return hashMedSeed(*(*ptr)(ptr(&b)), len(b), seed)
}

// HashStringSeed returns the hash of the string slice with given seed.
func HashStringSeed(s string, seed uint64) uint64 {
	if len(s) <= 16 {
		return hashSmallSeed(*(*ptr)(ptr(&s)), len(s), seed)
	}
	return hashMedSeed(*(*ptr)(ptr(&s)), len(s), seed)
}

func hashSmallSeed(p ptr, l int, seed uint64) (acc u64) {
	switch {
	case l > 8:
		inputlo := readU64(p, 0) ^ (key64_024 ^ key64_032 + seed)
		inputhi := readU64(p, ui(l)-8) ^ (key64_040 ^ key64_048 - seed)
		folded := mulFold64(inputlo, inputhi)
		return xxh3Avalanche(u64(l) + bits.ReverseBytes64(inputlo) + inputhi + folded)

	case l >= 4:
		seed ^= u64(bits.ReverseBytes32(u32(seed))) << 32
		input1 := readU32(p, 0)
		input2 := readU32(p, ui(l)-4)
		input64 := u64(input2) + u64(input1)<<32
		keyed := input64 ^ (key64_008 ^ key64_016 - seed)
		return rrmxmx(keyed, u64(l))

	case l == 3:
		c12 := u64(readU16(p, 0))
		c3 := u64(readU8(p, 2))
		acc = c12<<16 + c3 + 3<<8

	case l == 2:
		c12 := u64(readU16(p, 0))
		acc = c12*(1<<24+1)>>8 + 2<<8

	case l == 1:
		c1 := u64(readU8(p, 0))
		acc = c1*(1<<24+1<<16+1) + 1<<8

	case l == 0:
		return xxhAvalancheSmall(seed ^ key64_056 ^ key64_064)
	}

	return xxhAvalancheSmall(acc ^ (u64(key32_000^key32_004) + seed))
}

func hashMedSeed(p ptr, l int, seed uint64) (acc u64) {
	switch {
	case l <= 128:
		acc = u64(l) * prime64_1

		if l > 32 {
			if l > 64 {
				if l > 96 {
					acc += mulFold64(
						readU64(p, 6*8)^(key64_096+seed),
						readU64(p, 7*8)^(key64_104-seed))

					acc += mulFold64(
						readU64(p, ui(l)-8*8)^(key64_112+seed),
						readU64(p, ui(l)-7*8)^(key64_120-seed))
				} // 96

				acc += mulFold64(
					readU64(p, 4*8)^(key64_064+seed),
					readU64(p, 5*8)^(key64_072-seed))

				acc += mulFold64(
					readU64(p, ui(l)-6*8)^(key64_080+seed),
					readU64(p, ui(l)-5*8)^(key64_088-seed))
			} // 64

			acc += mulFold64(
				readU64(p, 2*8)^(key64_032+seed),
				readU64(p, 3*8)^(key64_040-seed))

			acc += mulFold64(
				readU64(p, ui(l)-4*8)^(key64_048+seed),
				readU64(p, ui(l)-3*8)^(key64_056-seed))
		} // 32

		acc += mulFold64(
			readU64(p, 0*8)^(key64_000+seed),
			readU64(p, 1*8)^(key64_008-seed))

		acc += mulFold64(
			readU64(p, ui(l)-2*8)^(key64_016+seed),
			readU64(p, ui(l)-1*8)^(key64_024-seed))

		return xxh3Avalanche(acc)

	case l <= 240:
		acc = u64(l) * prime64_1

		acc += mulFold64(
			readU64(p, 0*16+0)^(key64_000+seed),
			readU64(p, 0*16+8)^(key64_008-seed))

		acc += mulFold64(
			readU64(p, 1*16+0)^(key64_016+seed),
			readU64(p, 1*16+8)^(key64_024-seed))

		acc += mulFold64(
			readU64(p, 2*16+0)^(key64_032+seed),
			readU64(p, 2*16+8)^(key64_040-seed))

		acc += mulFold64(
			readU64(p, 3*16+0)^(key64_048+seed),
			readU64(p, 3*16+8)^(key64_056-seed))

		acc += mulFold64(
			readU64(p, 4*16+0)^(key64_064+seed),
			readU64(p, 4*16+8)^(key64_072-seed))

		acc += mulFold64(
			readU64(p, 5*16+0)^(key64_080+seed),
			readU64(p, 5*16+8)^(key64_088-seed))

		acc += mulFold64(
			readU64(p, 6*16+0)^(key64_096+seed),
			readU64(p, 6*16+8)^(key64_104-seed))

		acc += mulFold64(
			readU64(p, 7*16+0)^(key64_112+seed),
			readU64(p, 7*16+8)^(key64_120-seed))

		// avalanche
		acc = xxh3Avalanche(acc)

		// trailing groups after 128
		top := ui(l) &^ 15
		for i := ui(8 * 16); i < top; i += 16 {
			acc += mulFold64(
				readU64(p, i+0)^(readU64(key, i-125)+seed),
				readU64(p, i+8)^(readU64(key, i-117)-seed))
		}

		// last 16 bytes
		acc += mulFold64(
			readU64(p, ui(l)-16)^(key64_119+seed),
			readU64(p, ui(l)-8)^(key64_127-seed))

		return xxh3Avalanche(acc)

	default:
		secret := key
		if seed != 0 {
			secret = ptr(&[secret_size]byte{})
			initSecret(secret, seed)
		}
		return hashLargeSeed(p, u64(l), secret)
	}
}

func initSecret(secret ptr, seed u64) {
	for i := ui(0); i < secret_size/16; i++ {
		lo := readU64(key, 16*i) + seed
		hi := readU64(key, 16*i+8) - seed
		writeU64(secret, 16*i, lo)
		writeU64(secret, 16*i+8, hi)
	}
}

func hashLargeSeed(p ptr, l u64, secret ptr) (acc u64) {
	acc = l * prime64_1
	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	if hasAVX2 {
		accumAVX2(&accs, p, secret, l)
	} else if hasSSE2 {
		accumSSE(&accs, p, secret, l)
	} else {
		accumScalarSeed(&accs, p, secret, l)
	}

	// merge accs
	acc += mulFold64(accs[0]^readU64(secret, 11), accs[1]^readU64(secret, 19))
	acc += mulFold64(accs[2]^readU64(secret, 27), accs[3]^readU64(secret, 35))
	acc += mulFold64(accs[4]^readU64(secret, 43), accs[5]^readU64(secret, 51))
	acc += mulFold64(accs[6]^readU64(secret, 59), accs[7]^readU64(secret, 67))

	return xxh3Avalanche(acc)
}

func xxhAvalancheSmall(x u64) u64 {
	x ^= x >> 33
	x *= prime64_2
	x ^= x >> 29
	x *= prime64_3
	x ^= x >> 32
	return x
}
