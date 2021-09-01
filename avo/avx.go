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
		TEXT("accumAVX2", NOSPLIT, "func(acc *[8]uint64, data, key *byte, len uint64)")
		// %rdi, %rsi, %rdx, %rcx

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		skey := Mem{Base: Load(Param("key"), GP64())}
		plen := Load(Param("len"), GP64())
		prime := YMM()
		a := [...]VecVirtual{YMM(), YMM()}

		advance := func(n int) {
			ADDQ(U32(n*64), data.Base)
			SUBQ(U32(n*64), plen)
		}

		Label("load")
		{
			VMOVDQU(acc.Offset(0x00), a[0])
			VMOVDQU(acc.Offset(0x20), a[1])
			VMOVDQU(primeData, prime)
		}

		Label("accum_large")
		{
			CMPQ(plen, U32(1024))
			JLE(LabelRef("accum"))

			for i := 0; i < 16; i++ {
				avx2accum(data, a, 64*i, 8*i, key, true)
			}
			advance(16)
			avx2scramble(prime, key, a, 8*16)

			JMP(LabelRef("accum_large"))
		}

		Label("accum")
		{
			CMPQ(plen, Imm(64))
			JLE(LabelRef("finalize"))

			avx2accum(data, a, 0, 0, skey, false)
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

			avx2accum(data, a, 0, 121, key, false)
		}

		Label("return")
		{
			VMOVDQU(a[0], acc.Offset(0x00))
			VMOVDQU(a[1], acc.Offset(0x20))
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
		a := [...]VecVirtual{YMM(), YMM()}

		Label("load")
		{
			VMOVDQU(acc.Offset(0x00), a[0])
			VMOVDQU(acc.Offset(0x20), a[1])
			VMOVDQU(primeData, prime)
		}

		Label("accum_block")
		{
			for i := 0; i < 16; i++ {
				avx2accum(data, a, 64*i, 8*i, key, false)
			}
			avx2scramble(prime, key, a, 8*16)
		}

		Label("return")
		{
			VMOVDQU(a[0], acc.Offset(0x00))
			VMOVDQU(a[1], acc.Offset(0x20))
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
	VPADDQ(a[0], y0, a[0])
	VPADDQ(a[0], y1, a[0])
	VPADDQ(a[1], y10, a[1])
	VPADDQ(a[1], y11, a[1])
}
