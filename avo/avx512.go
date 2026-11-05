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
		plen := Load(Param("len"), GP64())
		prime := ZMM()
		// p are per-block partial sums that never include the carried a,
		// so the next block's stripes overlap the previous block's scramble.
		a := ZMM()
		p := [4]VecVirtual{ZMM(), ZMM(), ZMM(), ZMM()}

		Label("load")
		{
			// The caller has just written acc with 16-byte stores. Read it back in 16-byte
			// pieces so store-to-load forwarding does not fail on a load spanning several.
			VMOVDQU(acc.Offset(0x00), a.AsX())
			VINSERTI32X4(Imm(1), acc.Offset(0x10), a, a)
			VINSERTI32X4(Imm(2), acc.Offset(0x20), a, a)
			VINSERTI32X4(Imm(3), acc.Offset(0x30), a, a)
			VMOVDQU64(primeData, prime)
		}
		// Load key at 8 byte offsets in the order we use it.
		var keyReg [17]VecVirtual
		for i := range keyReg {
			keyReg[i] = ZMM()
			VMOVDQU64(key.Offset(i*8), keyReg[i])
		}
		advance := func(n int) {
			ADDQ(U32(n*64), data.Base)
			SUBQ(U32(n*64), plen)
		}

		// Key at offset 121.
		key121 := ZMM()
		VMOVDQU64(key.Offset(121), key121)

		Label("accum_large")
		{
			CMPQ(plen, U32(1024))
			JLE(LabelRef("accum"))

			for i := 0; i < 16; i += 2 {
				avx512accumSplitTwo(data, p, 64*i, keyReg[i], keyReg[i+1], true, i == 0)
			}
			VPADDQ(p[0], p[1], p[0])
			VPADDQ(p[2], p[3], p[2])
			VPADDQ(p[0], p[2], p[0])
			VPADDQ(a, p[0], a)
			advance(16)
			avx512scramble(prime, keyReg[16], a)

			JMP(LabelRef("accum_large"))
		}

		Label("accum")
		{
			// Unroll loop so we can use already loaded keys.
			for i := 0; i < 1024/64; i++ {
				CMPQ(plen, Imm(64))
				JLE(LabelRef("finalize"))

				avx512accum(data, a, 0, keyReg[i], false)
				advance(1)
			}
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

	{
		TEXT("accumBlockAVX512", NOSPLIT, "func(acc *[8]uint64, data, key *byte)")

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		prime := ZMM()
		// p are per-block partial sums that never include the carried a,
		// so the next block's stripes overlap the previous block's scramble.
		a := ZMM()
		p := [4]VecVirtual{ZMM(), ZMM(), ZMM(), ZMM()}

		Label("load")
		{
			VMOVDQU64(acc.Offset(0x00), a)
			VMOVDQU64(primeData, prime)
		}

		// Accumulate block of 1024 bytes.
		{
			for i := 0; i < 16; i += 2 {
				// Load keys just in time.
				key1, key2 := ZMM(), ZMM()
				VMOVDQU64(key.Offset(i*8), key1)
				VMOVDQU64(key.Offset(i*8+8), key2)

				avx512accumSplitTwo(data, p, 64*i, key1, key2, false, i == 0)
			}
			key128 := ZMM()
			VMOVDQU64(key.Offset(16*8), key128)
			VPADDQ(p[0], p[1], p[0])
			VPADDQ(p[2], p[3], p[2])
			VPADDQ(p[0], p[2], p[0])
			VPADDQ(a, p[0], a)
			avx512scramble(prime, key128, a)
		}

		{
			VMOVDQU64(a, acc.Offset(0x00))
			VZEROUPPER()
			RET()
		}
	}

	{
		// accumBlocksAVX512 accumulates and scrambles blocks of 1024 bytes, at least one.
		TEXT("accumBlocksAVX512", NOSPLIT, "func(acc *[8]uint64, data, key *byte, blocks uint64)")

		acc := Mem{Base: Load(Param("acc"), GP64())}
		data := Mem{Base: Load(Param("data"), GP64())}
		key := Mem{Base: Load(Param("key"), GP64())}
		blocks := Load(Param("blocks"), GP64())
		prime := ZMM()
		// p are per-block partial sums that never include the carried a,
		// so the next block's stripes overlap the previous block's scramble.
		a := ZMM()
		p := [4]VecVirtual{ZMM(), ZMM(), ZMM(), ZMM()}

		VMOVDQU64(acc.Offset(0x00), a)
		VMOVDQU64(primeData, prime)
		var keyReg [17]VecVirtual
		for i := range keyReg {
			keyReg[i] = ZMM()
			VMOVDQU64(key.Offset(i*8), keyReg[i])
		}

		Label("accum_blocks")
		{
			for i := 0; i < 16; i += 2 {
				avx512accumSplitTwo(data, p, 64*i, keyReg[i], keyReg[i+1], true, i == 0)
			}
			VPADDQ(p[0], p[1], p[0])
			VPADDQ(p[2], p[3], p[2])
			VPADDQ(p[0], p[2], p[0])
			VPADDQ(a, p[0], a)
			ADDQ(U32(1024), data.Base)
			avx512scramble(prime, keyReg[16], a)

			DECQ(blocks)
			JNZ(LabelRef("accum_blocks"))
		}

		VMOVDQU64(a, acc.Offset(0x00))
		VZEROUPPER()
		RET()
	}

	Generate()
}

func avx512scramble(prime, key, a VecVirtual) {
	y0, y1 := ZMM(), ZMM()

	VPSRLQ(Imm(0x2f), a, y0)
	// 3 way Xor, replacing VPXOR(a, y0, y0), VPXOR(key[koff], y0, y0)
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
	// Sum the stripe first so a carries one add per stripe.
	VPADDQ(y0, y1, y1)
	VPADDQ(a, y1, a)
}

func avx512accumSplitTwo(data Mem, p [4]VecVirtual, doff int, key, key2 VecVirtual, prefetch, first bool) {
	y0, y1, y2 := ZMM(), ZMM(), ZMM()
	z0, z1, z2 := ZMM(), ZMM(), ZMM()
	if first {
		// Replace destination registers with the partial sums to avoid zeroing and adding.
		y0 = p[0]
		y1 = p[1]
		z1 = p[2]
		z0 = p[3]
	}
	VMOVDQU64(data.Offset(doff+0x00), y0)
	VMOVDQU64(data.Offset(doff+0x40), z0)
	if prefetch {
		// Prefetch once per cacheline (64 bytes), 8 iterations ahead.
		PREFETCHT0(data.Offset(doff + 1024))
		PREFETCHT0(data.Offset(doff + 1024 + 64))
	}
	VPXORD(key, y0, y1)
	VPXORD(key2, z0, z1)
	VPSHUFD(Imm(49), y1, y2)
	VPSHUFD(Imm(49), z1, z2)
	VPMULUDQ(y1, y2, y1)
	VPMULUDQ(z1, z2, z1)
	VPSHUFD(Imm(78), y0, y0)
	VPSHUFD(Imm(78), z0, z0)

	if !first {
		VPADDQ(p[0], y0, p[0])
		VPADDQ(p[1], y1, p[1])
		VPADDQ(p[2], z1, p[2])
		VPADDQ(p[3], z0, p[3])
	}
}
