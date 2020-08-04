package xxh3

import "math/bits"

// Hash returns the hash of the byte slice.
func Hash(b []byte) uint64 {
	fn := hashMed
	if len(b) <= 16 {
		fn = hashSmall
	}
	return fn(*(*ptr)(ptr(&b)), len(b))

}

// Hash returns the hash of the string slice.
func HashString(s string) uint64 {
	fn := hashMed
	if len(s) <= 16 {
		fn = hashSmall
	}
	return fn(*(*ptr)(ptr(&s)), len(s))
}

func hashSmall(p ptr, l int) (acc u64) {
	const seed = 0

	switch {
	case l > 8:
		inputlo := readU64(p, 0) ^ (key64_024 ^ key64_032 + seed)
		inputhi := readU64(p, ui(l)-8) ^ (key64_040 ^ key64_048 - seed)
		folded := mulFold64(inputlo, inputhi)
		return xxh3Avalanche(u64(l) + bits.ReverseBytes64(inputlo) + inputhi + folded)

	case l >= 4:
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
		return 0x2d06800538d394c2 // xxh_avalanche(key64_056 ^ key64_064)
	}

	return xxhAvalancheSmall(acc)
}

func hashMed(p ptr, l int) (acc u64) {
	const seed = 0

	switch {
	case l <= 128:
		acc = u64(l) * prime64_1

		if l > 32 {
			if l > 64 {
				if l > 96 {
					acc += mulFold64(
						readU64(p, 6*8)^key64_096,
						readU64(p, 7*8)^key64_104)

					acc += mulFold64(
						readU64(p, ui(l)-8*8)^key64_112,
						readU64(p, ui(l)-7*8)^key64_120)
				} // 96

				acc += mulFold64(
					readU64(p, 4*8)^key64_064,
					readU64(p, 5*8)^key64_072)

				acc += mulFold64(
					readU64(p, ui(l)-6*8)^key64_080,
					readU64(p, ui(l)-5*8)^key64_088)
			} // 64

			acc += mulFold64(
				readU64(p, 2*8)^key64_032,
				readU64(p, 3*8)^key64_040)

			acc += mulFold64(
				readU64(p, ui(l)-4*8)^key64_048,
				readU64(p, ui(l)-3*8)^key64_056)
		} // 32

		acc += mulFold64(
			readU64(p, 0*8)^key64_000,
			readU64(p, 1*8)^key64_008)

		acc += mulFold64(
			readU64(p, ui(l)-2*8)^key64_016,
			readU64(p, ui(l)-1*8)^key64_024)

		return xxh3Avalanche(acc)

	case l <= 240:
		acc = u64(l) * prime64_1

		acc += mulFold64(
			readU64(p, 0*16+0)^key64_000,
			readU64(p, 0*16+8)^key64_008)

		acc += mulFold64(
			readU64(p, 1*16+0)^key64_016,
			readU64(p, 1*16+8)^key64_024)

		acc += mulFold64(
			readU64(p, 2*16+0)^key64_032,
			readU64(p, 2*16+8)^key64_040)

		acc += mulFold64(
			readU64(p, 3*16+0)^key64_048,
			readU64(p, 3*16+8)^key64_056)

		acc += mulFold64(
			readU64(p, 4*16+0)^key64_064,
			readU64(p, 4*16+8)^key64_072)

		acc += mulFold64(
			readU64(p, 5*16+0)^key64_080,
			readU64(p, 5*16+8)^key64_088)

		acc += mulFold64(
			readU64(p, 6*16+0)^key64_096,
			readU64(p, 6*16+8)^key64_104)

		acc += mulFold64(
			readU64(p, 7*16+0)^key64_112,
			readU64(p, 7*16+8)^key64_120)

		// avalanche
		acc = xxh3Avalanche(acc)

		// trailing groups after 128
		top := ui(l) &^ 15
		for i := ui(8 * 16); i < top; i += 16 {
			acc += mulFold64(
				readU64(p, i+0)^readU64(key, i-125),
				readU64(p, i+8)^readU64(key, i-117))
		}

		// last 16 bytes
		acc += mulFold64(
			readU64(p, ui(l)-16)^key64_127,
			readU64(p, ui(l)-8)^key64_135)

		return xxh3Avalanche(acc)

	case avx2, sse2:
		return hashVector(p, u64(l))

	default:
		return hashLarge(p, u64(l))
	}
}

func hashLarge(p ptr, l u64) (acc u64) {
	acc = l * prime64_1
	accs := [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1}

	for l > _block {
		k := key

		// accs
		for i := 0; i < 16; i++ {
			dv0 := readU64(p, 8*0)
			dk0 := dv0 ^ readU64(k, 8*0)
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := readU64(p, 8*1)
			dk1 := dv1 ^ readU64(k, 8*1)
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := readU64(p, 8*2)
			dk2 := dv2 ^ readU64(k, 8*2)
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := readU64(p, 8*3)
			dk3 := dv3 ^ readU64(k, 8*3)
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := readU64(p, 8*4)
			dk4 := dv4 ^ readU64(k, 8*4)
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := readU64(p, 8*5)
			dk5 := dv5 ^ readU64(k, 8*5)
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := readU64(p, 8*6)
			dk6 := dv6 ^ readU64(k, 8*6)
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := readU64(p, 8*7)
			dk7 := dv7 ^ readU64(k, 8*7)
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)

			l -= _stripe
			if l > 0 {
				p, k = ptr(ui(p)+_stripe), ptr(ui(k)+8)
			}
		}

		// scramble accs
		accs[0] ^= accs[0] >> 47
		accs[0] ^= key64_128
		accs[0] *= prime32_1

		accs[1] ^= accs[1] >> 47
		accs[1] ^= key64_136
		accs[1] *= prime32_1

		accs[2] ^= accs[2] >> 47
		accs[2] ^= key64_144
		accs[2] *= prime32_1

		accs[3] ^= accs[3] >> 47
		accs[3] ^= key64_152
		accs[3] *= prime32_1

		accs[4] ^= accs[4] >> 47
		accs[4] ^= key64_160
		accs[4] *= prime32_1

		accs[5] ^= accs[5] >> 47
		accs[5] ^= key64_168
		accs[5] *= prime32_1

		accs[6] ^= accs[6] >> 47
		accs[6] ^= key64_176
		accs[6] *= prime32_1

		accs[7] ^= accs[7] >> 47
		accs[7] ^= key64_184
		accs[7] *= prime32_1
	}

	if l > 0 {
		t, k := (l-1)/_stripe, key

		for i := u64(0); i < t; i++ {
			dv0 := readU64(p, 8*0)
			dk0 := dv0 ^ readU64(k, 8*0)
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := readU64(p, 8*1)
			dk1 := dv1 ^ readU64(k, 8*1)
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := readU64(p, 8*2)
			dk2 := dv2 ^ readU64(k, 8*2)
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := readU64(p, 8*3)
			dk3 := dv3 ^ readU64(k, 8*3)
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := readU64(p, 8*4)
			dk4 := dv4 ^ readU64(k, 8*4)
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := readU64(p, 8*5)
			dk5 := dv5 ^ readU64(k, 8*5)
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := readU64(p, 8*6)
			dk6 := dv6 ^ readU64(k, 8*6)
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := readU64(p, 8*7)
			dk7 := dv7 ^ readU64(k, 8*7)
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)

			l -= _stripe
			if l > 0 {
				p, k = ptr(ui(p)+_stripe), ptr(ui(k)+8)
			}
		}

		if l > 0 {
			p = ptr(ui(p) - uintptr(_stripe-l))

			dv0 := readU64(p, 8*0)
			dk0 := dv0 ^ key64_121
			accs[1] += dv0
			accs[0] += (dk0 & 0xffffffff) * (dk0 >> 32)

			dv1 := readU64(p, 8*1)
			dk1 := dv1 ^ key64_129
			accs[0] += dv1
			accs[1] += (dk1 & 0xffffffff) * (dk1 >> 32)

			dv2 := readU64(p, 8*2)
			dk2 := dv2 ^ key64_137
			accs[3] += dv2
			accs[2] += (dk2 & 0xffffffff) * (dk2 >> 32)

			dv3 := readU64(p, 8*3)
			dk3 := dv3 ^ key64_145
			accs[2] += dv3
			accs[3] += (dk3 & 0xffffffff) * (dk3 >> 32)

			dv4 := readU64(p, 8*4)
			dk4 := dv4 ^ key64_153
			accs[5] += dv4
			accs[4] += (dk4 & 0xffffffff) * (dk4 >> 32)

			dv5 := readU64(p, 8*5)
			dk5 := dv5 ^ key64_161
			accs[4] += dv5
			accs[5] += (dk5 & 0xffffffff) * (dk5 >> 32)

			dv6 := readU64(p, 8*6)
			dk6 := dv6 ^ key64_169
			accs[7] += dv6
			accs[6] += (dk6 & 0xffffffff) * (dk6 >> 32)

			dv7 := readU64(p, 8*7)
			dk7 := dv7 ^ key64_177
			accs[6] += dv7
			accs[7] += (dk7 & 0xffffffff) * (dk7 >> 32)
		}
	}

	// merge accs
	acc += mulFold64(accs[0]^key64_011, accs[1]^key64_019)
	acc += mulFold64(accs[2]^key64_027, accs[3]^key64_035)
	acc += mulFold64(accs[4]^key64_043, accs[5]^key64_051)
	acc += mulFold64(accs[6]^key64_059, accs[7]^key64_067)

	return xxh3Avalanche(acc)
}
