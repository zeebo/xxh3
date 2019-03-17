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

// HashString returns the hash of the string.
func HashString(s string) uint64 {
	return hash(*(*ptr)(ptr(&s)), uint64(len(s)))
}

// Hash returns the hash of the byte slice.
func Hash(b []byte) uint64 {
	return hash(*(*ptr)(ptr(&b)), uint64(len(b)))
}

func hash(p ptr, l uint64) uint64 {
	ol := l

	if l <= 16 {
		if l > 8 {
			acc := prime64_1 * l
			l1 := *(*uint64)(p) + *(*uint64)(key)
			l2 := *(*uint64)(ptr(ui(p) + ui(l) - 8)) + *(*uint64)(ptr(ui(key) + 8))
			hi, lo := bits.Mul64(l1, l2)
			acc += hi + lo

			acc ^= acc >> 29
			acc *= prime64_3
			acc ^= acc >> 32
			return acc

		} else if l >= 4 {
			acc := prime64_1 * l
			l1 := *(*uint32)(p) + *(*uint32)(key)
			l2 := *(*uint32)(ptr(ui(p) + ui(l) - 4)) + *(*uint32)(ptr(ui(key) + 4))
			acc += uint64(l1) * uint64(l2)

			acc ^= acc >> 29
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
			acc := m1 * m2

			acc ^= acc >> 29
			acc *= prime64_3
			acc ^= acc >> 32
			return acc

		}

		return 0
	}

	var hi1, lo1, hi2, lo2, hi3, lo3, hi4, lo4 uint64

	accs := prime64_1 * l
	if l > 32 {
		if l > 64 {
			if l > 96 {
				if l > 128 {
					goto long
				}
				hi1, lo1 = bits.Mul64(
					*(*uint64)(ptr(ui(p) + 48))^*(*uint64)(ptr(ui(key) + 96)),
					*(*uint64)(ptr(ui(p) + 48 + 8))^*(*uint64)(ptr(ui(key) + 96 + 8)),
				)
				accs += hi1 + lo1

				hi2, lo2 = bits.Mul64(
					*(*uint64)(ptr(ui(p) + ui(l) - 64))^*(*uint64)(ptr(ui(key) + 112)),
					*(*uint64)(ptr(ui(p) + ui(l) - 64 + 8))^*(*uint64)(ptr(ui(key) + 112 + 8)),
				)
				accs += hi2 + lo2
			}
			hi1, lo1 = bits.Mul64(
				*(*uint64)(ptr(ui(p) + 32))^*(*uint64)(ptr(ui(key) + 64)),
				*(*uint64)(ptr(ui(p) + 32 + 8))^*(*uint64)(ptr(ui(key) + 64 + 8)),
			)
			accs += hi1 + lo1

			hi2, lo2 = bits.Mul64(
				*(*uint64)(ptr(ui(p) + ui(l) - 48))^*(*uint64)(ptr(ui(key) + 80)),
				*(*uint64)(ptr(ui(p) + ui(l) - 48 + 8))^*(*uint64)(ptr(ui(key) + 80 + 8)),
			)
			accs += hi2 + lo2
		}
		hi1, lo1 = bits.Mul64(
			*(*uint64)(ptr(ui(p) + 16))^*(*uint64)(ptr(ui(key) + 32)),
			*(*uint64)(ptr(ui(p) + 16 + 8))^*(*uint64)(ptr(ui(key) + 32 + 8)),
		)
		accs += hi1 + lo1

		hi2, lo2 = bits.Mul64(
			*(*uint64)(ptr(ui(p) + ui(l) - 32))^*(*uint64)(ptr(ui(key) + 48)),
			*(*uint64)(ptr(ui(p) + ui(l) - 32 + 8))^*(*uint64)(ptr(ui(key) + 48 + 8)),
		)
		accs += hi2 + lo2
	}
	hi1, lo1 = bits.Mul64(
		*(*uint64)(ptr(ui(p) + 0))^*(*uint64)(ptr(ui(key) + 0)),
		*(*uint64)(ptr(ui(p) + 0 + 8))^*(*uint64)(ptr(ui(key) + 0 + 8)),
	)
	accs += hi1 + lo1

	hi2, lo2 = bits.Mul64(
		*(*uint64)(ptr(ui(p) + ui(l) - 16))^*(*uint64)(ptr(ui(key) + 16)),
		*(*uint64)(ptr(ui(p) + ui(l) - 16 + 8))^*(*uint64)(ptr(ui(key) + 16 + 8)),
	)
	accs += hi2 + lo2

	accs ^= accs >> 29
	accs *= prime64_3
	accs ^= accs >> 32
	return accs

long:
	acc := [8]uint64{0, prime64_1, prime64_2, prime64_3, prime64_4, prime64_5, 0, 0}
	blocks := l / _block

	for n := uint64(0); n < blocks; n++ {
		k := key

		// acc
		for i := 0; i < 16; i++ {
			l0, r0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			acc[0] += uint64(l0+*(*uint32)(ptr(ui(k) + 0)))*uint64(r0+*(*uint32)(ptr(ui(k) + 4))) + uint64(l0) + (uint64(r0) << 32)

			l1, r1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			acc[1] += uint64(l1+*(*uint32)(ptr(ui(k) + 8)))*uint64(r1+*(*uint32)(ptr(ui(k) + 12))) + uint64(l1) + (uint64(r1) << 32)

			l2, r2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			acc[2] += uint64(l2+*(*uint32)(ptr(ui(k) + 16)))*uint64(r2+*(*uint32)(ptr(ui(k) + 20))) + uint64(l2) + (uint64(r2) << 32)

			l3, r3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			acc[3] += uint64(l3+*(*uint32)(ptr(ui(k) + 24)))*uint64(r3+*(*uint32)(ptr(ui(k) + 28))) + uint64(l3) + (uint64(r3) << 32)

			l4, r4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			acc[4] += uint64(l4+*(*uint32)(ptr(ui(k) + 32)))*uint64(r4+*(*uint32)(ptr(ui(k) + 36))) + uint64(l4) + (uint64(r4) << 32)

			l5, r5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			acc[5] += uint64(l5+*(*uint32)(ptr(ui(k) + 40)))*uint64(r5+*(*uint32)(ptr(ui(k) + 44))) + uint64(l5) + (uint64(r5) << 32)

			l6, r6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			acc[6] += uint64(l6+*(*uint32)(ptr(ui(k) + 48)))*uint64(r6+*(*uint32)(ptr(ui(k) + 52))) + uint64(l6) + (uint64(r6) << 32)

			l7, r7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			acc[7] += uint64(l7+*(*uint32)(ptr(ui(k) + 56)))*uint64(r7+*(*uint32)(ptr(ui(k) + 60))) + uint64(l7) + (uint64(r7) << 32)

			p, k, l = ptr(ui(p)+_stripe), ptr(ui(k)+8), l-_stripe
		}

		// scramble acc
		acc[0] ^= acc[0] >> 47
		acc[0] = (uint64(uint32(acc[0])) * uint64(*(*uint32)(ptr(ui(k) + 0)))) ^ ((acc[0] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 4))))

		acc[1] ^= acc[1] >> 47
		acc[1] = (uint64(uint32(acc[1])) * uint64(*(*uint32)(ptr(ui(k) + 8)))) ^ ((acc[1] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 12))))

		acc[2] ^= acc[2] >> 47
		acc[2] = (uint64(uint32(acc[2])) * uint64(*(*uint32)(ptr(ui(k) + 16)))) ^ ((acc[2] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 20))))

		acc[3] ^= acc[3] >> 47
		acc[3] = (uint64(uint32(acc[3])) * uint64(*(*uint32)(ptr(ui(k) + 24)))) ^ ((acc[3] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 28))))

		acc[4] ^= acc[4] >> 47
		acc[4] = (uint64(uint32(acc[4])) * uint64(*(*uint32)(ptr(ui(k) + 32)))) ^ ((acc[4] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 36))))

		acc[5] ^= acc[5] >> 47
		acc[5] = (uint64(uint32(acc[5])) * uint64(*(*uint32)(ptr(ui(k) + 40)))) ^ ((acc[5] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 44))))

		acc[6] ^= acc[6] >> 47
		acc[6] = (uint64(uint32(acc[6])) * uint64(*(*uint32)(ptr(ui(k) + 48)))) ^ ((acc[6] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 52))))

		acc[7] ^= acc[7] >> 47
		acc[7] = (uint64(uint32(acc[7])) * uint64(*(*uint32)(ptr(ui(k) + 56)))) ^ ((acc[7] >> 32) * uint64(*(*uint32)(ptr(ui(k) + 60))))
	}

	if l > 0 {
		t, k := (l%_block)/_stripe, key
		for i := uint64(0); i < t; i++ {
			l0, r0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			acc[0] += uint64(l0+*(*uint32)(ptr(ui(k) + 0)))*uint64(r0+*(*uint32)(ptr(ui(k) + 4))) + uint64(l0) + (uint64(r0) << 32)

			l1, r1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			acc[1] += uint64(l1+*(*uint32)(ptr(ui(k) + 8)))*uint64(r1+*(*uint32)(ptr(ui(k) + 12))) + uint64(l1) + (uint64(r1) << 32)

			l2, r2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			acc[2] += uint64(l2+*(*uint32)(ptr(ui(k) + 16)))*uint64(r2+*(*uint32)(ptr(ui(k) + 20))) + uint64(l2) + (uint64(r2) << 32)

			l3, r3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			acc[3] += uint64(l3+*(*uint32)(ptr(ui(k) + 24)))*uint64(r3+*(*uint32)(ptr(ui(k) + 28))) + uint64(l3) + (uint64(r3) << 32)

			l4, r4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			acc[4] += uint64(l4+*(*uint32)(ptr(ui(k) + 32)))*uint64(r4+*(*uint32)(ptr(ui(k) + 36))) + uint64(l4) + (uint64(r4) << 32)

			l5, r5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			acc[5] += uint64(l5+*(*uint32)(ptr(ui(k) + 40)))*uint64(r5+*(*uint32)(ptr(ui(k) + 44))) + uint64(l5) + (uint64(r5) << 32)

			l6, r6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			acc[6] += uint64(l6+*(*uint32)(ptr(ui(k) + 48)))*uint64(r6+*(*uint32)(ptr(ui(k) + 52))) + uint64(l6) + (uint64(r6) << 32)

			l7, r7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			acc[7] += uint64(l7+*(*uint32)(ptr(ui(k) + 56)))*uint64(r7+*(*uint32)(ptr(ui(k) + 60))) + uint64(l7) + (uint64(r7) << 32)

			p, k, l = ptr(ui(p)+_stripe), ptr(ui(k)+8), l-_stripe
		}

		if l > 0 {
			p = ptr(ui(p) - uintptr(_stripe-l))

			l0, r0 := *(*uint32)(ptr(ui(p) + 0)), *(*uint32)(ptr(ui(p) + 4))
			acc[0] += uint64(l0+*(*uint32)(ptr(ui(k) + 0)))*uint64(r0+*(*uint32)(ptr(ui(k) + 4))) + uint64(l0) + (uint64(r0) << 32)

			l1, r1 := *(*uint32)(ptr(ui(p) + 8)), *(*uint32)(ptr(ui(p) + 12))
			acc[1] += uint64(l1+*(*uint32)(ptr(ui(k) + 8)))*uint64(r1+*(*uint32)(ptr(ui(k) + 12))) + uint64(l1) + (uint64(r1) << 32)

			l2, r2 := *(*uint32)(ptr(ui(p) + 16)), *(*uint32)(ptr(ui(p) + 20))
			acc[2] += uint64(l2+*(*uint32)(ptr(ui(k) + 16)))*uint64(r2+*(*uint32)(ptr(ui(k) + 20))) + uint64(l2) + (uint64(r2) << 32)

			l3, r3 := *(*uint32)(ptr(ui(p) + 24)), *(*uint32)(ptr(ui(p) + 28))
			acc[3] += uint64(l3+*(*uint32)(ptr(ui(k) + 24)))*uint64(r3+*(*uint32)(ptr(ui(k) + 28))) + uint64(l3) + (uint64(r3) << 32)

			l4, r4 := *(*uint32)(ptr(ui(p) + 32)), *(*uint32)(ptr(ui(p) + 36))
			acc[4] += uint64(l4+*(*uint32)(ptr(ui(k) + 32)))*uint64(r4+*(*uint32)(ptr(ui(k) + 36))) + uint64(l4) + (uint64(r4) << 32)

			l5, r5 := *(*uint32)(ptr(ui(p) + 40)), *(*uint32)(ptr(ui(p) + 44))
			acc[5] += uint64(l5+*(*uint32)(ptr(ui(k) + 40)))*uint64(r5+*(*uint32)(ptr(ui(k) + 44))) + uint64(l5) + (uint64(r5) << 32)

			l6, r6 := *(*uint32)(ptr(ui(p) + 48)), *(*uint32)(ptr(ui(p) + 52))
			acc[6] += uint64(l6+*(*uint32)(ptr(ui(k) + 48)))*uint64(r6+*(*uint32)(ptr(ui(k) + 52))) + uint64(l6) + (uint64(r6) << 32)

			l7, r7 := *(*uint32)(ptr(ui(p) + 56)), *(*uint32)(ptr(ui(p) + 60))
			acc[7] += uint64(l7+*(*uint32)(ptr(ui(k) + 56)))*uint64(r7+*(*uint32)(ptr(ui(k) + 60))) + uint64(l7) + (uint64(r7) << 32)
		}
	}

	// merge_accs
	hi1, lo1 = bits.Mul64(acc[0]^*(*uint64)(ptr(ui(key) + 0)), acc[1]^*(*uint64)(ptr(ui(key) + 8)))
	hi2, lo2 = bits.Mul64(acc[2]^*(*uint64)(ptr(ui(key) + 16)), acc[3]^*(*uint64)(ptr(ui(key) + 24)))
	hi3, lo3 = bits.Mul64(acc[4]^*(*uint64)(ptr(ui(key) + 32)), acc[5]^*(*uint64)(ptr(ui(key) + 40)))
	hi4, lo4 = bits.Mul64(acc[6]^*(*uint64)(ptr(ui(key) + 48)), acc[7]^*(*uint64)(ptr(ui(key) + 56)))
	result := ol*prime64_1 + hi1 + lo1 + hi2 + lo2 + hi3 + lo3 + hi4 + lo4

	result ^= result >> 29
	result *= prime64_3
	result ^= result >> 32
	return result
}
