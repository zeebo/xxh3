package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// AVX512 can be used for very long blocks.
func AVX512() {
	// Lay out the prime constant in memory, copy it so no unpack is needed.
	primeData := GLOBL("prime_avx512", RODATA|NOPTR)
	for i := 0; i < 64; i += 8 {
		DATA(i, U64(2654435761))
	}

	{

		TEXT("accumAVX512", NOSPLIT, "func(acc *[8]uint64, data, key *byte, len uint64)")

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		skey := Mem{Base: Load(Param("key"), GP64())}
		plen := Load(Param("len"), GP64())
		prime := ZMM()
		a := ZMM()

		var keyReg [17]VecVirtual
		for i := range keyReg {
			keyReg[i] = ZMM()
			VMOVDQU64(key.Offset(i*8), keyReg[i])
		}

		advance := func(n int) {
			ADDQ(U32(n*64), data.Base)
			SUBQ(U32(n*64), plen)
		}

		Label("load")
		{
			VMOVDQU64(acc.Offset(0x00), a)
			VMOVDQU64(primeData, prime)
		}

		key121 := ZMM()
		VMOVDQU64(key.Offset(121), key121)

		Label("accum_large")
		{
			CMPQ(plen, U32(1024))
			JLE(LabelRef("accum"))

			for i := 0; i < 16; i++ {
				avx512accum(data, a, 64*i, keyReg[i], true)
			}
			advance(16)
			avx512scramble(prime, keyReg[16], a)

			JMP(LabelRef("accum_large"))
		}

		Label("accum")
		{
			CMPQ(plen, Imm(64))
			JLE(LabelRef("finalize"))

			avx512accum(data, a, 0, keyReg[0], false)
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

			avx512accum(data, a, 0, key121, false)
		}

		Label("return")
		{
			VMOVDQU64(a, acc.Offset(0x00))
			VZEROUPPER()
			RET()
		}
	}

	Generate()
}

func avx512scramble(prime, key, a VecVirtual) {
	y0, y1 := ZMM(), ZMM()

	VPSRLQ(Imm(0x2f), a, y0)
	//VPXOR(a, y0, y0)
	//VPXOR(key[koff], y0, y0)
	// 3 way xor:
	VPTERNLOGD(U8(0x96), a, key, y0)
	VPMULUDQ(prime, y0, y1)
	VPSHUFD(Imm(0xf5), y0, y0)
	VPMULUDQ(prime, y0, y0)
	VPSLLQ(Imm(0x20), y0, y0)

	VPADDQ(y1, y0, a)
}

func avx512accum(data Mem, a VecVirtual, doff int, key VecVirtual, prefetch bool) {
	y0, y1, y2 := ZMM(), ZMM(), ZMM()
	VMOVDQU64(data.Offset(doff+0x00), y0)
	if prefetch {
		// Prefetch once per cacheline (64 bytes), 8 iterations ahead.
		PREFETCHT0(data.Offset(doff + 1024))
	}
	VPXORD(key, y0, y1)
	VPSHUFD(Imm(49), y1, y2)
	VPMULUDQ(y1, y2, y1)
	VPSHUFD(Imm(78), y0, y0)
	VPADDQ(a, y0, a)
	VPADDQ(a, y1, a)
}
