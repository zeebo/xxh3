package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func AVX() {
	// Lay out the prime constant in memory, copy it so no unpack is needed.
	primeData := GLOBL("prime_avx", RODATA|NOPTR)
	for i := 0; i < 32; i += 8 {
		DATA(i, U64(2654435761))
	}

	{
		// loadEarly will fill up remaining registers with key values.
		// This can reduce L1 fetch pressure.
		const loadEarly = true

		TEXT("accumAVX2", NOSPLIT, "func(acc *[8]uint64, data, key *byte, len uint64)")
		// %rdi, %rsi, %rdx, %rcx

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		skey := Mem{Base: Load(Param("key"), GP64())}
		plen := Load(Param("len"), GP64())
		prime := YMM()

		// a[0] contains merged values
		// a[1] is a temporary accumulator used for the accum_large loop.
		a := [2][2]VecVirtual{{YMM(), YMM()}, {YMM(), YMM()}}

		advance := func(n int) {
			ADDQ(U32(n*64), data.Base)
			SUBQ(U32(n*64), plen)
		}

		// Preload a number of keys to fill regs, and interleave them.
		var loadedOffset = make(map[int]Op)

		Label("load")
		{
			VMOVDQU(acc.Offset(0x00), a[0][0])
			VMOVDQU(acc.Offset(0x20), a[0][1])
			VMOVDQU(primeData, prime)

			CMPQ(plen, U32(1024))
			JLE(LabelRef("accum"))

			used := 0
			for i := 0; i < 16; i++ {
				offset := i * 8
				nextOffset := i*8 + 32
				if loadEarly && used < 5 {
					loadedOffset[nextOffset] = YMM()
					VMOVDQU(key.Offset(nextOffset), loadedOffset[nextOffset])
					used++
				} else if loadedOffset[nextOffset] == nil {
					loadedOffset[nextOffset] = key.Offset(nextOffset)
				}
				if loadedOffset[offset] == nil {
					loadedOffset[offset] = key.Offset(offset)
				}
			}
		}

		Label("accum_large")
		{
			for i := 0; i < 16; i++ {
				avx2accumLoaded(data, a, 64*i, loadedOffset[i*8], loadedOffset[i*8+32], true, i == 0)
			}
			// add a[1] to a[0]
			VPADDQ(a[0][0], a[1][0], a[0][0])
			VPADDQ(a[0][1], a[1][1], a[0][1])
			advance(16)

			avx2scramble(prime, key, a[0], 8*16)

			CMPQ(plen, U32(1024))
			JLE(LabelRef("accum"))
			JMP(LabelRef("accum_large"))
		}

		Label("accum")
		{
			CMPQ(plen, Imm(64))
			JLE(LabelRef("finalize"))

			avx2accum(data, a[0], 0, 0, skey, false)
			advance(1)
			ADDQ(U32(8), skey.Base)

			JMP(LabelRef("accum"))
		}

		Label("finalize")
		{
			CMPQ(plen, Imm(0))
			JE(LabelRef("return"))

			SUBQ(Imm(64), data.Base)
			ADDQ(plen, data.Base)

			avx2accum(data, a[0], 0, 121, key, false)
		}

		Label("return")
		{
			VMOVDQU(a[0][0], acc.Offset(0x00))
			VMOVDQU(a[0][1], acc.Offset(0x20))
			VZEROUPPER()
			RET()
		}
	}

	{
		TEXT("accumBlockAVX2", NOSPLIT, "func(acc *[8]uint64, data, key *byte)")

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		prime := YMM()
		a := [2][2]VecVirtual{{YMM(), YMM()}, {YMM(), YMM()}}

		// Preload a number of keys to fill regs, and interleave them.
		var loadedOffset = make(map[int]Op)

		Label("load")
		{
			const loadEarly = false
			VMOVDQU(acc.Offset(0x00), a[0][0])
			VMOVDQU(acc.Offset(0x20), a[0][1])
			VMOVDQU(primeData, prime)

			used := 0
			for i := 0; i < 16; i++ {
				offset := i * 8
				nextOffset := i*8 + 32
				if loadEarly && used < 5 {
					loadedOffset[nextOffset] = YMM()
					VMOVDQU(key.Offset(nextOffset), loadedOffset[nextOffset])
					used++
				} else if loadedOffset[nextOffset] == nil {
					loadedOffset[nextOffset] = key.Offset(nextOffset)
				}
				if loadedOffset[offset] == nil {
					loadedOffset[offset] = key.Offset(offset)
				}
			}
		}

		Label("accum_block")
		{
			for i := 0; i < 16; i++ {
				avx2accumLoaded(data, a, 64*i, loadedOffset[i*8], loadedOffset[i*8+32], false, i == 0)
			}
			// add a[1] to a[0]
			VPADDQ(a[0][0], a[1][0], a[0][0])
			VPADDQ(a[0][1], a[1][1], a[0][1])
			avx2scramble(prime, key, a[0], 8*16)
		}

		Label("return")
		{
			VMOVDQU(a[0][0], acc.Offset(0x00))
			VMOVDQU(a[0][1], acc.Offset(0x20))
			VZEROUPPER()
			RET()
		}
	}

	Generate()
}

func avx2scramble(prime VecVirtual, key Mem, a [2]VecVirtual, koff int) {
	for n, offset := range []int{0x00, 0x20} {
		y0, y1 := YMM(), YMM()

		VPSRLQ(Imm(0x2f), a[n], y0)
		VPXOR(a[n], y0, y0)
		VPXOR(key.Offset(koff+offset), y0, y0)
		VPMULUDQ(prime, y0, y1)
		VPSHUFD(Imm(0xf5), y0, y0)
		VPMULUDQ(prime, y0, y0)
		VPSLLQ(Imm(0x20), y0, y0)

		VPADDQ(y1, y0, a[n])
	}
}

func avx2accum(data Mem, a [2]VecVirtual, doff, koff int, key Mem, prefetch bool) {
	y0, y1, y2, y10, y11, y12 := YMM(), YMM(), YMM(), YMM(), YMM(), YMM()
	VMOVDQU(data.Offset(doff+0x00), y0)
	VMOVDQU(data.Offset(doff+0x20), y10)
	if prefetch {
		// Prefetch once per cacheline (64 bytes), 8 iterations ahead.
		PREFETCHT0(data.Offset(doff + 512))
	}
	VPXOR(key.Offset(koff+0x00), y0, y1)
	VPXOR(key.Offset(koff+0x20), y10, y11)
	VPSHUFD(Imm(49), y1, y2)
	VPSHUFD(Imm(49), y11, y12)
	VPMULUDQ(y1, y2, y1)
	VPMULUDQ(y11, y12, y11)
	VPSHUFD(Imm(78), y0, y0)
	VPSHUFD(Imm(78), y10, y10)
	VPADDQ(a[0], y1, a[0])
	VPADDQ(a[1], y11, a[1])
	VPADDQ(a[0], y0, a[0])
	VPADDQ(a[1], y10, a[1])
}

func avx2accumLoaded(data Mem, a [2][2]VecVirtual, doff int, k0, k1 Op, prefetch, initA1 bool) {
	y0, y1, y2, y10, y11, y12 := YMM(), YMM(), YMM(), YMM(), YMM(), YMM()
	if initA1 {
		// Write to a[1] to initialize it.
		y0 = a[1][0]
		y10 = a[1][1]
	}
	VMOVDQU(data.Offset(doff+0x00), y0)
	VMOVDQU(data.Offset(doff+0x20), y10)
	if prefetch {
		// Prefetch once per cacheline (64 bytes), one block ahead.
		PREFETCHT0(data.Offset(doff + 1024))
	}
	VPXOR(k0, y0, y1)
	VPXOR(k1, y10, y11)
	VPSHUFD(Imm(49), y1, y2)
	VPSHUFD(Imm(49), y11, y12)
	VPMULUDQ(y1, y2, y1)
	VPMULUDQ(y11, y12, y11)
	VPSHUFD(Imm(78), y0, y0)
	VPSHUFD(Imm(78), y10, y10)
	VPADDQ(a[0][0], y1, a[0][0])
	VPADDQ(a[0][1], y11, a[0][1])
	if !initA1 {
		VPADDQ(a[1][0], y0, a[1][0])
		VPADDQ(a[1][1], y10, a[1][1])
	}
}
