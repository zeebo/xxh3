// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	// Lay out the prime constant in memory
	primeData := GLOBL("prime_avx", RODATA|NOPTR)
	DATA(0, U32(2654435761))

	TEXT("accum_avx", NOSPLIT, "func(acc *[8]uint64, data, key *byte, len uint64)")
	// %rdi, %rsi, %rdx, %rcx

	acc := Mem{Base: Load(Param("acc"), GP64())}
	data := Mem{Base: Load(Param("data"), GP64())}
	key := Mem{Base: Load(Param("key"), GP64())}
	len := Load(Param("len"), GP64())
	prime := YMM()
	a := [2]VecVirtual{YMM(), YMM()}

	advance := func(n int) {
		ADDQ(U32(n*64), data.Base)
		ADDQ(U32(n*8), key.Base)
		SUBQ(U32(n*64), len)
	}

	accum := func(n int) {
		doff, koff := 64*n, 8*n

		for n, offset := range []int{0x00, 0x20} {
			y0, y1, y2 := YMM(), YMM(), YMM()

			VMOVDQU(data.Offset(offset+doff), y0)
			VPADDD(key.Offset(offset+koff), y0, y1)
			VPSHUFD(Imm(245), y1, y2)
			VPMULUDQ(y2, y1, y1)
			VPADDQ(y0, y1, y0)

			VPADDQ(a[n], y0, a[n])
		}
	}

	scramble := func() {
		for n, offset := range []int{0x00, 0x20} {
			y0, y1 := YMM(), YMM()

			VPSRLQ(Imm(0x2f), a[n], y0)
			VPXOR(a[n], y0, y0)
			VPXOR(key.Offset(offset), y0, y0)
			VPMULUDQ(prime, y0, y1)
			VPSHUFD(Imm(0xf5), y0, y0)
			VPMULUDQ(prime, y0, y0)
			VPSLLQ(Imm(0x20), y0, y0)

			VPADDQ(y1, y0, a[n])
		}
	}

	Label("load")
	{
		VMOVDQU(acc.Offset(0x00), a[0])
		VMOVDQU(acc.Offset(0x20), a[1])
		VPBROADCASTQ(primeData, prime)
	}

	Label("accum_large")
	{
		CMPQ(len, U32(1024))
		JLT(LabelRef("accum"))

		for i := 0; i < 16; i++ {
			accum(i)
		}
		advance(16)

		scramble()
		Load(Param("key"), key.Base)

		JMP(LabelRef("accum_large"))
	}

	Label("accum")
	{
		CMPQ(len, Imm(64))
		JLT(LabelRef("finalize"))

		accum(0)
		advance(1)

		JMP(LabelRef("accum"))
	}

	Label("finalize")
	{
		CMPQ(len, Imm(0))
		JE(LabelRef("return"))

		SUBQ(Imm(64), data.Base)
		ADDQ(len, data.Base)

		accum(0)
	}

	Label("return")
	{
		VMOVDQU(a[0], acc.Offset(0x00))
		VMOVDQU(a[1], acc.Offset(0x20))
		RET()
	}

	Generate()
}
