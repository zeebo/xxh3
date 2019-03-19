// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("accum", NOSPLIT, "func(acc *[8]uint64, data, key *byte, len uint64)")

	acc := Mem{Base: Load(Param("acc"), GP64())}
	data := Mem{Base: Load(Param("data"), GP64())}
	key := Mem{Base: Load(Param("key"), GP64())}
	len := Load(Param("len"), GP64())
	ctr := GP64()
	XORQ(ctr, ctr)

	x0, x1 := XMM(), XMM()
	y0, y1, y2 := x0.AsY(), x1.AsY(), YMM()

	accum := func() {
		// Load the data into registers
		VMOVDQU(data, x0)
		VMOVDQU(key, x1)
		VINSERTI128(Imm(1), data.Offset(0x10), y0, y0)
		VINSERTI128(Imm(1), key.Offset(0x10), y1, y1)

		// Do the math and store into acc
		VPADDD(y0, y1, y1)
		VPSHUFD(Imm(0x31), y1, y2)
		VPADDQ(acc, y0, y0)
		VPMULUDQ(y2, y1, y1)
		VPADDQ(y1, y0, y0)
		VMOVDQU(y0, acc)

		// Load the next half of data into registers
		VMOVDQU(data.Offset(0x20), x0)
		VINSERTI128(Imm(1), data.Offset(0x30), y0, y0)
		VMOVDQU(key.Offset(0x20), x1)
		VINSERTI128(Imm(1), key.Offset(0x30), y1, y1)

		// Do the math and store into acc
		VPADDD(y0, y1, y1)
		VPSHUFD(Imm(0x31), y1, y2)
		VPADDQ(acc.Offset(0x20), y0, y0)
		VPMULUDQ(y2, y1, y1)
		VPADDQ(y1, y0, y0)
		VMOVDQU(y0, acc.Offset(0x20))
	}

	Label("loop")
	{
		CMPQ(len, Imm(64))
		JLT(LabelRef("finalize"))

		// Accumulate 64 bytes of data.
		accum()

		// Update our pointers and jump to loop head
		ADDQ(Imm(64), data.Base)
		ADDQ(Imm(8), key.Base)
		SUBQ(Imm(64), len)

		// Check if we've done 16 iterations. If so, we need to mix.
		INCQ(ctr)
		CMPQ(ctr, Imm(16))
		JE(LabelRef("mix"))
		JMP(LabelRef("loop"))
	}

	Label("mix")
	{
		// Load the data into registers
		VMOVDQU(key, x1)
		VMOVDQU(acc, y0)
		VINSERTI128(Imm(1), key.Offset(0x10), y1, y1)

		// Do the math and store into acc
		VPSRLQ(Imm(0x2f), y0, y2)
		VPXOR(y2, y0, y0)
		VPMULUDQ(y1, y0, y2)
		VPSHUFD(Imm(0x31), y1, y1)
		VPSHUFD(Imm(0x31), y0, y0)
		VPMULUDQ(y1, y0, y0)
		VPXOR(y0, y2, y0)
		VMOVDQU(y0, acc)

		// Load the next half of data into registers
		VMOVDQU(key.Offset(0x20), x1)
		VMOVDQU(acc.Offset(0x20), y0)
		VINSERTI128(Imm(1), key.Offset(0x30), y1, y1)

		// Do the math and store into acc
		VPSRLQ(Imm(0x2f), y0, y2)
		VPXOR(y2, y0, y0)
		VPMULUDQ(y1, y0, y2)
		VPSHUFD(Imm(0x31), y1, y1)
		VPSHUFD(Imm(0x31), y0, y0)
		VPMULUDQ(y1, y0, y0)
		VPXOR(y0, y2, y0)
		VMOVDQU(y0, acc.Offset(0x20))
		XORQ(ctr, ctr)
		Load(Param("key"), key.Base)
		JMP(LabelRef("loop"))
	}

	Label("finalize")
	{
		CMPQ(len, Imm(0))
		JE(LabelRef("ret"))

		SUBQ(Imm(64), data.Base)
		ADDQ(len, data.Base)
		accum()
	}

	Label("ret")
	RET()

	Generate()
}
