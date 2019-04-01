package xxh3

import (
	"math/bits"
	"unsafe"
)

type (
	ptr = unsafe.Pointer
	ui  = uintptr
)

const (
	_stripe = 64
	_block  = 1024

	prime32_1 = 2654435761

	prime64_1 = 11400714785074694791
	prime64_2 = 14029467366897019727
	prime64_3 = 1609587929392839161
	prime64_4 = 9650029242287828579
	prime64_5 = 2870177450012600261
)

var key = ptr(&[...]uint32{
	0xb8fe6c39, 0x23a44bbe, 0x7c01812c, 0xf721ad1c,
	0xded46de9, 0x839097db, 0x7240a4a4, 0xb7b3671f,
	0xcb79e64e, 0xccc0e578, 0x825ad07d, 0xccff7221,
	0xb8084674, 0xf743248e, 0xe03590e6, 0x813a264c,
	0x3c2852bb, 0x91c300cb, 0x88d0658b, 0x1b532ea3,
	0x71644897, 0xa20df94e, 0x3819ef46, 0xa9deacd8,
	0xa8fa763f, 0xe39c343f, 0xf9dcbbc7, 0xc70b4f1d,
	0x8a51e04b, 0xcdb45931, 0xc89f7ec9, 0xd9787364,
	0xeac5ac83, 0x34d3ebc3, 0xc581a0ff, 0xfa1363eb,
	0x170ddd51, 0xb7f0da49, 0xd3165526, 0x29d4689e,
	0x2b16be58, 0x7d47a1fc, 0x8ff8b8d1, 0x7ad031ce,
	0x45cb3a8f, 0x95160428, 0xafd7fbca, 0xbb4b407e,
})

// HashString returns the hash of the byte slice.
func HashString(s string) uint64 {
	fn := hash
	if len(s) > 128 {
		fn = hashLarge
	}
	return fn(s)
}

// Hash returns the hash of the byte slice.
func Hash(b []byte) uint64 {
	fn := hash
	if len(b) > 128 {
		fn = hashLarge
	}
	return fn(*(*string)(ptr(&b)))
}

func hash(s string) (acc uint64) {
	p, l := *(*ptr)(ptr(&s)), uint64(len(s))

	if l <= 16 {
		if l > 8 {
			ll1 := *(*uint64)(p) ^ *(*uint64)(key)
			ll2 := *(*uint64)(ptr(ui(p) + ui(l) - 8)) ^ *(*uint64)(ptr(ui(key) + 8))
			hi, lo := bits.Mul64(ll1, ll2)
			acc = l + ll1 + ll2 + (hi ^ lo)

			// avalanche
			acc ^= acc >> 37
			acc *= prime64_3
			acc ^= acc >> 32

			return acc
		} else if l >= 4 {
			in1 := *(*uint32)(p)
			in2 := *(*uint32)(ptr(ui(p) + ui(l) - 4))
			in64 := uint64(in1) + uint64(in2)<<32
			keyed := in64 ^ *(*uint64)(key)
			hi, lo := bits.Mul64(keyed, prime64_1)
			acc = l + hi ^ lo

			// avalanche
			acc ^= acc >> 37
			acc *= prime64_3
			acc ^= acc >> 32

			return acc
		} else if l > 0 {
			c1 := *(*uint8)(p)
			c2 := *(*uint8)(ptr(ui(p) + (ui(l) >> 1)))
			c3 := *(*uint8)(ptr(ui(p) + ui(l) - 1))
			l1 := uint32(c1) + (uint32(c2) << 8)
			l2 := uint32(l) + (uint32(c3) << 2)
			m1 := uint64(l1 + *(*uint32)(ptr(ui(key) + 0)))
			m2 := uint64(l2 + *(*uint32)(ptr(ui(key) + 4)))
			acc = m1 * m2

			// avalanche
			acc ^= acc >> 37
			acc *= prime64_3
			acc ^= acc >> 32

			return acc
		}

		return 0
	}

	if l > 32 {
		if l > 64 {
			if l > 96 {
				hi1, lo1 := bits.Mul64(
					*(*uint64)(ptr(ui(p) + 48))^*(*uint64)(ptr(ui(key) + 96)),
					*(*uint64)(ptr(ui(p) + 48 + 8))^*(*uint64)(ptr(ui(key) + 96 + 8)),
				)
				acc += hi1 ^ lo1

				hi2, lo2 := bits.Mul64(
					*(*uint64)(ptr(ui(p) + ui(l) - 64))^*(*uint64)(ptr(ui(key) + 112)),
					*(*uint64)(ptr(ui(p) + ui(l) - 64 + 8))^*(*uint64)(ptr(ui(key) + 112 + 8)),
				)
				acc += hi2 ^ lo2
			} // 96

			hi1, lo1 := bits.Mul64(
				*(*uint64)(ptr(ui(p) + 32))^*(*uint64)(ptr(ui(key) + 64)),
				*(*uint64)(ptr(ui(p) + 32 + 8))^*(*uint64)(ptr(ui(key) + 64 + 8)),
			)
			acc += hi1 ^ lo1

			hi2, lo2 := bits.Mul64(
				*(*uint64)(ptr(ui(p) + ui(l) - 48))^*(*uint64)(ptr(ui(key) + 80)),
				*(*uint64)(ptr(ui(p) + ui(l) - 48 + 8))^*(*uint64)(ptr(ui(key) + 80 + 8)),
			)
			acc += hi2 ^ lo2
		} // 64

		hi1, lo1 := bits.Mul64(
			*(*uint64)(ptr(ui(p) + 16))^*(*uint64)(ptr(ui(key) + 32)),
			*(*uint64)(ptr(ui(p) + 16 + 8))^*(*uint64)(ptr(ui(key) + 32 + 8)),
		)
		acc += hi1 ^ lo1

		hi2, lo2 := bits.Mul64(
			*(*uint64)(ptr(ui(p) + ui(l) - 32))^*(*uint64)(ptr(ui(key) + 48)),
			*(*uint64)(ptr(ui(p) + ui(l) - 32 + 8))^*(*uint64)(ptr(ui(key) + 48 + 8)),
		)
		acc += hi2 ^ lo2
	} // 32

	hi1, lo1 := bits.Mul64(
		*(*uint64)(ptr(ui(p) + 0))^*(*uint64)(ptr(ui(key) + 0)),
		*(*uint64)(ptr(ui(p) + 0 + 8))^*(*uint64)(ptr(ui(key) + 0 + 8)),
	)
	acc += hi1 ^ lo1

	hi2, lo2 := bits.Mul64(
		*(*uint64)(ptr(ui(p) + ui(l) - 16))^*(*uint64)(ptr(ui(key) + 16)),
		*(*uint64)(ptr(ui(p) + ui(l) - 16 + 8))^*(*uint64)(ptr(ui(key) + 16 + 8)),
	)
	acc += hi2 ^ lo2

	// avalanche
	acc ^= acc >> 37
	acc *= prime64_3
	acc ^= acc >> 32

	return acc
}

// hashLarge handles lengths greater than 128.
func hashLarge(s string) uint64 {
	p, l := *(*ptr)(ptr(&s)), uint64(len(s))

	if avx2 || sse2 {
		return hashVector(p, l)
	}

	acc := l * prime64_1
	accs := [8]uint64{0, prime64_1, prime64_2, prime64_3, prime64_4, prime64_5, 0, 0}

	for l >= _block {
		k := key

		// accs
		for i := 0; i < 16; i++ {
			dl0, dr0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			kl0, kr0 := *(*uint32)(ptr(ui(k) + 0)), *(*uint32)(ptr(ui(k) + 4))
			accs[0] += uint64(dl0^kl0)*uint64(dr0^kr0) + uint64(dl0) + uint64(dr0)<<32

			dl1, dr1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			kl1, kr1 := *(*uint32)(ptr(ui(k) + 8)), *(*uint32)(ptr(ui(k) + 12))
			accs[1] += uint64(dl1^kl1)*uint64(dr1^kr1) + uint64(dl1) + uint64(dr1)<<32

			dl2, dr2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			kl2, kr2 := *(*uint32)(ptr(ui(k) + 16)), *(*uint32)(ptr(ui(k) + 20))
			accs[2] += uint64(dl2^kl2)*uint64(dr2^kr2) + uint64(dl2) + uint64(dr2)<<32

			dl3, dr3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			kl3, kr3 := *(*uint32)(ptr(ui(k) + 24)), *(*uint32)(ptr(ui(k) + 28))
			accs[3] += uint64(dl3^kl3)*uint64(dr3^kr3) + uint64(dl3) + uint64(dr3)<<32

			dl4, dr4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			kl4, kr4 := *(*uint32)(ptr(ui(k) + 32)), *(*uint32)(ptr(ui(k) + 36))
			accs[4] += uint64(dl4^kl4)*uint64(dr4^kr4) + uint64(dl4) + uint64(dr4)<<32

			dl5, dr5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			kl5, kr5 := *(*uint32)(ptr(ui(k) + 40)), *(*uint32)(ptr(ui(k) + 44))
			accs[5] += uint64(dl5^kl5)*uint64(dr5^kr5) + uint64(dl5) + uint64(dr5)<<32

			dl6, dr6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			kl6, kr6 := *(*uint32)(ptr(ui(k) + 48)), *(*uint32)(ptr(ui(k) + 52))
			accs[6] += uint64(dl6^kl6)*uint64(dr6^kr6) + uint64(dl6) + uint64(dr6)<<32

			dl7, dr7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			kl7, kr7 := *(*uint32)(ptr(ui(k) + 56)), *(*uint32)(ptr(ui(k) + 60))
			accs[7] += uint64(dl7^kl7)*uint64(dr7^kr7) + uint64(dl7) + uint64(dr7)<<32

			p, k = ptr(ui(p)+_stripe), ptr(ui(k)+8)
		}

		// scramble accs
		accs[0] ^= accs[0] >> 47
		accs[0] ^= *(*uint64)(ptr(ui(k) + 0*8))
		accs[0] *= prime32_1

		accs[1] ^= accs[1] >> 47
		accs[1] ^= *(*uint64)(ptr(ui(k) + 1*8))
		accs[1] *= prime32_1

		accs[2] ^= accs[2] >> 47
		accs[2] ^= *(*uint64)(ptr(ui(k) + 2*8))
		accs[2] *= prime32_1

		accs[3] ^= accs[3] >> 47
		accs[3] ^= *(*uint64)(ptr(ui(k) + 3*8))
		accs[3] *= prime32_1

		accs[4] ^= accs[4] >> 47
		accs[4] ^= *(*uint64)(ptr(ui(k) + 4*8))
		accs[4] *= prime32_1

		accs[5] ^= accs[5] >> 47
		accs[5] ^= *(*uint64)(ptr(ui(k) + 5*8))
		accs[5] *= prime32_1

		accs[6] ^= accs[6] >> 47
		accs[6] ^= *(*uint64)(ptr(ui(k) + 6*8))
		accs[6] *= prime32_1

		accs[7] ^= accs[7] >> 47
		accs[7] ^= *(*uint64)(ptr(ui(k) + 7*8))
		accs[7] *= prime32_1

		l -= 16 * _stripe
	}

	if l > 0 {
		t, k := (l%_block)/_stripe, key

		for i := uint64(0); i < t; i++ {
			dl0, dr0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			kl0, kr0 := *(*uint32)(ptr(ui(k) + 0)), *(*uint32)(ptr(ui(k) + 4))
			accs[0] += uint64(dl0^kl0)*uint64(dr0^kr0) + uint64(dl0) + uint64(dr0)<<32

			dl1, dr1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			kl1, kr1 := *(*uint32)(ptr(ui(k) + 8)), *(*uint32)(ptr(ui(k) + 12))
			accs[1] += uint64(dl1^kl1)*uint64(dr1^kr1) + uint64(dl1) + uint64(dr1)<<32

			dl2, dr2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			kl2, kr2 := *(*uint32)(ptr(ui(k) + 16)), *(*uint32)(ptr(ui(k) + 20))
			accs[2] += uint64(dl2^kl2)*uint64(dr2^kr2) + uint64(dl2) + uint64(dr2)<<32

			dl3, dr3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			kl3, kr3 := *(*uint32)(ptr(ui(k) + 24)), *(*uint32)(ptr(ui(k) + 28))
			accs[3] += uint64(dl3^kl3)*uint64(dr3^kr3) + uint64(dl3) + uint64(dr3)<<32

			dl4, dr4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			kl4, kr4 := *(*uint32)(ptr(ui(k) + 32)), *(*uint32)(ptr(ui(k) + 36))
			accs[4] += uint64(dl4^kl4)*uint64(dr4^kr4) + uint64(dl4) + uint64(dr4)<<32

			dl5, dr5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			kl5, kr5 := *(*uint32)(ptr(ui(k) + 40)), *(*uint32)(ptr(ui(k) + 44))
			accs[5] += uint64(dl5^kl5)*uint64(dr5^kr5) + uint64(dl5) + uint64(dr5)<<32

			dl6, dr6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			kl6, kr6 := *(*uint32)(ptr(ui(k) + 48)), *(*uint32)(ptr(ui(k) + 52))
			accs[6] += uint64(dl6^kl6)*uint64(dr6^kr6) + uint64(dl6) + uint64(dr6)<<32

			dl7, dr7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			kl7, kr7 := *(*uint32)(ptr(ui(k) + 56)), *(*uint32)(ptr(ui(k) + 60))
			accs[7] += uint64(dl7^kl7)*uint64(dr7^kr7) + uint64(dl7) + uint64(dr7)<<32

			p, k, l = ptr(ui(p)+_stripe), ptr(ui(k)+8), l-_stripe
		}

		if l > 0 {
			p = ptr(ui(p) - uintptr(_stripe-l))

			dl0, dr0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			kl0, kr0 := *(*uint32)(ptr(ui(k) + 0)), *(*uint32)(ptr(ui(k) + 4))
			accs[0] += uint64(dl0^kl0)*uint64(dr0^kr0) + uint64(dl0) + uint64(dr0)<<32

			dl1, dr1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			kl1, kr1 := *(*uint32)(ptr(ui(k) + 8)), *(*uint32)(ptr(ui(k) + 12))
			accs[1] += uint64(dl1^kl1)*uint64(dr1^kr1) + uint64(dl1) + uint64(dr1)<<32

			dl2, dr2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			kl2, kr2 := *(*uint32)(ptr(ui(k) + 16)), *(*uint32)(ptr(ui(k) + 20))
			accs[2] += uint64(dl2^kl2)*uint64(dr2^kr2) + uint64(dl2) + uint64(dr2)<<32

			dl3, dr3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			kl3, kr3 := *(*uint32)(ptr(ui(k) + 24)), *(*uint32)(ptr(ui(k) + 28))
			accs[3] += uint64(dl3^kl3)*uint64(dr3^kr3) + uint64(dl3) + uint64(dr3)<<32

			dl4, dr4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			kl4, kr4 := *(*uint32)(ptr(ui(k) + 32)), *(*uint32)(ptr(ui(k) + 36))
			accs[4] += uint64(dl4^kl4)*uint64(dr4^kr4) + uint64(dl4) + uint64(dr4)<<32

			dl5, dr5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			kl5, kr5 := *(*uint32)(ptr(ui(k) + 40)), *(*uint32)(ptr(ui(k) + 44))
			accs[5] += uint64(dl5^kl5)*uint64(dr5^kr5) + uint64(dl5) + uint64(dr5)<<32

			dl6, dr6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			kl6, kr6 := *(*uint32)(ptr(ui(k) + 48)), *(*uint32)(ptr(ui(k) + 52))
			accs[6] += uint64(dl6^kl6)*uint64(dr6^kr6) + uint64(dl6) + uint64(dr6)<<32

			dl7, dr7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			kl7, kr7 := *(*uint32)(ptr(ui(k) + 56)), *(*uint32)(ptr(ui(k) + 60))
			accs[7] += uint64(dl7^kl7)*uint64(dr7^kr7) + uint64(dl7) + uint64(dr7)<<32
		}

	}

	// merge accs
	hi1, lo1 := bits.Mul64(accs[0]^*(*uint64)(ptr(ui(key) + 0*8)), accs[1]^*(*uint64)(ptr(ui(key) + 1*8)))
	acc += hi1 ^ lo1
	hi2, lo2 := bits.Mul64(accs[2]^*(*uint64)(ptr(ui(key) + 2*8)), accs[3]^*(*uint64)(ptr(ui(key) + 3*8)))
	acc += hi2 ^ lo2
	hi3, lo3 := bits.Mul64(accs[4]^*(*uint64)(ptr(ui(key) + 4*8)), accs[5]^*(*uint64)(ptr(ui(key) + 5*8)))
	acc += hi3 ^ lo3
	hi4, lo4 := bits.Mul64(accs[6]^*(*uint64)(ptr(ui(key) + 6*8)), accs[7]^*(*uint64)(ptr(ui(key) + 7*8)))
	acc += hi4 ^ lo4

	// avalanche
	acc ^= acc >> 37
	acc *= prime64_3
	acc ^= acc >> 32

	return acc
}
