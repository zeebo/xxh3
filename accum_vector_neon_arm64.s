#include "textflag.h"

DATA prime_neon<>+0(SB)/4, $0x9e3779b1
DATA prime_neon<>+4(SB)/4, $0x9e3779b1
DATA prime_neon<>+8(SB)/4, $0x9e3779b1
DATA prime_neon<>+12(SB)/4, $0x9e3779b1
GLOBL prime_neon<>(SB), RODATA|NOPTR, $16

// XTN Vd.2S, Vn.2D - Narrow 64-bit to 32-bit (low part)
// Encoding: 0|0|0|01110|10|100001|001010|Rn|Rd = 0x0EA12800
#define XTN_2S_2D(Vd, Vn) WORD $(0x0EA12800 | ((Vn) << 5) | (Vd))

// SHRN #32, Vn.2D, Vd.2S - Shift right 32 and narrow (high part)
#define SHRN_32_2D_2S(Vd, Vn) WORD $(0x0F208400 | ((Vn) << 5) | (Vd))

// UMULL Vd.2D, Vn.2S, Vm.2S - Widening multiply 32x32→64
#define UMULL_2D_2S_2S(Vd, Vn, Vm) WORD $(0x2EA0C000 | ((Vm) << 16) | ((Vn) << 5) | (Vd))

// UMLAL Vd.2D, Vn.2S, Vm.2S - Widening multiply 32x32→64 of the low halves, added to Vd
#define UMLAL_2D_2S_2S(Vd, Vn, Vm) WORD $(0x2EA08000 | ((Vm) << 16) | ((Vn) << 5) | (Vd))

// UMLAL2 Vd.2D, Vn.4S, Vm.4S - Widening multiply 32x32→64 of the high halves, added to Vd
#define UMLAL2_2D_4S_4S(Vd, Vn, Vm) WORD $(0x6EA08000 | ((Vm) << 16) | ((Vn) << 5) | (Vd))

// Register layout shared by both functions. The eight 64-bit accumulators are
// split the way upstream xxHash splits them on AArch64: six in NEON registers
// and two in general purpose registers. The vector pipes are the bottleneck,
// and the scalar lanes run beside them on the integer pipes.
//
//   V1, V2, V3      acc[0:2], acc[2:4], acc[4:6], without the swapped-data term
//   V11, V12, V13   sum of the raw data added to those lanes since the last FOLD
//   R8, R9          acc[6], acc[7]
//   V0, R19         prime32_1, for the scramble
//
// A round is acc[i] += data[i^1] + lo32(data[i]^key[i]) * hi32(data[i]^key[i]).
// Swapping the data words costs one vector instruction per register per
// stripe, so the NEON lanes add the data unswapped into V11-V13 and FOLD swaps
// that sum once, before the accumulators are next read.

// LOAD reads the accumulators at (R0) and clears the pending data sums.
#define LOAD \
	VLD1  (R0), [V1.D2, V2.D2, V3.D2] \
	LDP   48(R0), (R8, R9)            \
	VMOVI $0, V11.B16                 \
	VMOVI $0, V12.B16                 \
	VMOVI $0, V13.B16

// FOLD adds the pending data sums, swapped, into V1-V3.
#define FOLD \
	VEXT $8, V11.B16, V11.B16, V11.B16 \
	VEXT $8, V12.B16, V12.B16, V12.B16 \
	VEXT $8, V13.B16, V13.B16, V13.B16 \
	VADD V11.D2, V1.D2, V1.D2          \
	VADD V12.D2, V2.D2, V2.D2          \
	VADD V13.D2, V3.D2, V3.D2

// STORE folds and writes the accumulators back to (R0).
#define STORE \
	FOLD                              \
	VST1 [V1.D2, V2.D2, V3.D2], (R0) \
	STP  (R8, R9), 48(R0)

// ROUND accumulates the 64-byte stripe at doff(dbase). Its 64 bytes of key are
// already in registers: K0-K2 for the NEON lanes, KA and KB for the scalar ones.
//
// V5-V7 hold the data and V8-V10 data XOR key. Viewing V8/V9 as 4x32 bits,
// UZP1 gathers the low half of their four 64-bit lanes in V14 and UZP2 the
// high half in V15, so one UMLAL/UMLAL2 pair multiplies both registers. V10
// has no partner and is narrowed on its own. R10/R11, R14/R15 and R20/R21
// play the same roles for the scalar lanes.
#define ROUND(dbase, doff, K0, K1, K2, KA, KB) \
	FMOVQ  (doff)(dbase), F5            \
	FMOVQ  (doff+16)(dbase), F6         \
	FMOVQ  (doff+32)(dbase), F7         \
	LDP    (doff+48)(dbase), (R10, R11) \
	VEOR   V5.B16, K0.B16, V8.B16       \
	VEOR   V6.B16, K1.B16, V9.B16       \
	VEOR   V7.B16, K2.B16, V10.B16      \
	EOR    R10, KA, R14                 \
	EOR    R11, KB, R15                 \
	VUZP1  V9.S4, V8.S4, V14.S4         \
	VUZP2  V9.S4, V8.S4, V15.S4         \
	XTN_2S_2D(16, 10)                   \
	SHRN_32_2D_2S(17, 10)               \
	LSR    $32, R14, R20                \
	LSR    $32, R15, R21                \
	UMLAL_2D_2S_2S(1, 14, 15)           \
	UMLAL2_2D_4S_4S(2, 14, 15)          \
	UMLAL_2D_2S_2S(3, 16, 17)           \
	ADD    R11, R8, R8                  \
	ADD    R10, R9, R9                  \
	UMADDL R20, R8, R14, R8             \
	UMADDL R21, R9, R15, R9             \
	VADD   V5.D2, V11.D2, V11.D2        \
	VADD   V6.D2, V12.D2, V12.D2        \
	VADD   V7.D2, V13.D2, V13.D2

// STRIPE is ROUND with the key read from koff(kbase).
#define STRIPE(dbase, doff, kbase, koff) \
	FMOVQ (koff)(kbase), F20            \
	FMOVQ (koff+16)(kbase), F21         \
	FMOVQ (koff+32)(kbase), F22         \
	LDP   (koff+48)(kbase), (R12, R13)  \
	ROUND(dbase, doff, V20, V21, V22, R12, R13)

// SCRAMBLE for one NEON accumulator pair: acc = (acc ^ acc>>47 ^ key) * prime32_1
// Uses V5 for data, V6/V7/V8 temps
// V_acc_num = numeric reg for macros (1-3), V_acc = symbolic (V1-V3)
#define SCRAMBLE(key_off, V_acc_num, V_acc) \
	VUSHR $47, V_acc.D2, V5.D2         \
	VEOR  V_acc.B16, V5.B16, V_acc.B16 \
	ADD   $key_off, R2, R5             \
	VLD1  (R5), [V5.D2]                \
	VEOR  V_acc.B16, V5.B16, V_acc.B16 \
	XTN_2S_2D(6, V_acc_num)            \
	SHRN_32_2D_2S(7, V_acc_num)        \
	UMULL_2D_2S_2S(V_acc_num, 6, 0)    \
	UMULL_2D_2S_2S(8, 7, 0)            \
	VSHL  $32, V8.D2, V8.D2            \
	VADD  V8.D2, V_acc.D2, V_acc.D2

// SCRAMBLE_SCALAR is SCRAMBLE for one accumulator held in Racc.
#define SCRAMBLE_SCALAR(key_off, Racc) \
	MOVD key_off(R2), R10     \
	EOR  Racc>>47, Racc, Racc \
	EOR  R10, Racc, Racc      \
	MUL  R19, Racc, Racc

// SCRAMBLE_ALL folds the pending data sums and scrambles all eight
// accumulators with key[128:192].
#define SCRAMBLE_ALL \
	FOLD                     \
	VMOVI $0, V11.B16        \
	VMOVI $0, V12.B16        \
	VMOVI $0, V13.B16        \
	SCRAMBLE(128, 1, V1)     \
	SCRAMBLE(144, 2, V2)     \
	SCRAMBLE(160, 3, V3)     \
	SCRAMBLE_SCALAR(176, R8) \
	SCRAMBLE_SCALAR(184, R9)

// BLOCK accumulates the 1024 bytes (16 stripes) at (R1) and scrambles.
//
// Stripe n reads data[64n:64n+64] and key[8n:8n+64]. The key advances by only
// 8 bytes per stripe, so two stripes later the same key bytes are used one
// 16-byte lane further down. Reading all 64 key bytes for every stripe costs
// about 15% on Neoverse V2, so each stripe reads only what is new to it: one
// vector into V20-V22 (even stripes) or V23-V25 (odd stripes), and one word
// into R12/R13, with the registers rotating roles.
//
// The second half is addressed from R7 to keep the LDP offsets encodable.
#define BLOCK \
	FMOVQ  0(R2), F20                       \
	FMOVQ  16(R2), F21                      \
	FMOVQ  32(R2), F22                      \
	LDP    48(R2), (R12, R13)               \
	ROUND(R1, 0, V20, V21, V22, R12, R13)   \
	FMOVQ  8(R2), F23                       \
	FMOVQ  24(R2), F24                      \
	FMOVQ  40(R2), F25                      \
	MOVD   64(R2), R12                      \
	ROUND(R1, 64, V23, V24, V25, R13, R12)  \
	FMOVQ  48(R2), F20                      \
	MOVD   72(R2), R13                      \
	ROUND(R1, 128, V21, V22, V20, R12, R13) \
	FMOVQ  56(R2), F23                      \
	MOVD   80(R2), R12                      \
	ROUND(R1, 192, V24, V25, V23, R13, R12) \
	FMOVQ  64(R2), F21                      \
	MOVD   88(R2), R13                      \
	ROUND(R1, 256, V22, V20, V21, R12, R13) \
	FMOVQ  72(R2), F24                      \
	MOVD   96(R2), R12                      \
	ROUND(R1, 320, V25, V23, V24, R13, R12) \
	FMOVQ  80(R2), F22                      \
	MOVD   104(R2), R13                     \
	ROUND(R1, 384, V20, V21, V22, R12, R13) \
	FMOVQ  88(R2), F25                      \
	MOVD   112(R2), R12                     \
	ROUND(R1, 448, V23, V24, V25, R13, R12) \
	ADD    $512, R1, R7                     \
	FMOVQ  96(R2), F20                      \
	MOVD   120(R2), R13                     \
	ROUND(R7, 0, V21, V22, V20, R12, R13)   \
	FMOVQ  104(R2), F23                     \
	MOVD   128(R2), R12                     \
	ROUND(R7, 64, V24, V25, V23, R13, R12)  \
	FMOVQ  112(R2), F21                     \
	MOVD   136(R2), R13                     \
	ROUND(R7, 128, V22, V20, V21, R12, R13) \
	FMOVQ  120(R2), F24                     \
	MOVD   144(R2), R12                     \
	ROUND(R7, 192, V25, V23, V24, R13, R12) \
	FMOVQ  128(R2), F22                     \
	MOVD   152(R2), R13                     \
	ROUND(R7, 256, V20, V21, V22, R12, R13) \
	FMOVQ  136(R2), F25                     \
	MOVD   160(R2), R12                     \
	ROUND(R7, 320, V23, V24, V25, R13, R12) \
	FMOVQ  144(R2), F20                     \
	MOVD   168(R2), R13                     \
	ROUND(R7, 384, V21, V22, V20, R12, R13) \
	FMOVQ  152(R2), F23                     \
	MOVD   176(R2), R12                     \
	ROUND(R7, 448, V24, V25, V23, R13, R12) \
	SCRAMBLE_ALL

// func accumNEON(acc *[8]uint64, data *byte, key *byte, len uint64)
TEXT ·accumNEON(SB), NOSPLIT, $0-32
	MOVD acc+0(FP), R0
	MOVD data+8(FP), R1
	MOVD key+16(FP), R2
	MOVD len+24(FP), R3

	LOAD

	// Inputs of at most one block never scramble
	CMP $1024, R3
	BLE accum

	// Load prime constant for scramble
	MOVD $prime_neon<>(SB), R4
	VLD1 (R4), [V0.D2]
	MOVD $0x9e3779b1, R19

accum_large:
	BLOCK

	ADD $1024, R1, R1
	SUB $1024, R3, R3
	CMP $1024, R3
	BGT accum_large

accum:
	// If no remaining bytes, we're done
	CBZ R3, done

	// Compute number of full stripes: R6 = (R3-1) / 64
	SUB $1, R3, R6
	LSR $6, R6, R6
	CBZ R6, finalize

	MOVD $0, R4 // key offset

accum_loop:
	// Process one 64-byte stripe
	ADD R4, R2, R5
	STRIPE(R1, 0, R5, 0)

	ADD  $64, R1, R1
	SUB  $64, R3, R3
	ADD  $8, R4, R4
	AND  $127, R4, R4
	SUBS $1, R6, R6
	BNE  accum_loop

finalize:
	// Always process final stripe if there's remaining data (R3 > 0)
	CBZ R3, done

	// Adjust data pointer to last 64 bytes: R1 = R1 + R3 - 64
	ADD R3, R1, R6
	SUB $64, R6, R6

	// Final key offset is 121
	ADD $121, R2, R5
	STRIPE(R6, 0, R5, 0)

done:
	STORE
	RET

// func accumBlockNEON(acc *[8]uint64, data *byte, key *byte)
// Processes exactly 1024 bytes (16 stripes) and scrambles
TEXT ·accumBlockNEON(SB), NOSPLIT, $0-24
	MOVD acc+0(FP), R0
	MOVD data+8(FP), R1
	MOVD key+16(FP), R2

	LOAD

	// Load prime constant for scramble
	MOVD $prime_neon<>(SB), R4
	VLD1 (R4), [V0.D2]
	MOVD $0x9e3779b1, R19

	BLOCK

	STORE
	RET
