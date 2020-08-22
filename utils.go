package xxh3

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

type (
	ptr = unsafe.Pointer
	ui  = uintptr

	u8  = uint8
	u32 = uint32
	u64 = uint64
)

var le = binary.LittleEndian

func readU8(p ptr, o ui) uint8   { return *(*uint8)(ptr(ui(p) + o)) }
func readU16(p ptr, o ui) uint16 { return le.Uint16((*[2]byte)(ptr(ui(p) + o))[:]) }
func readU32(p ptr, o ui) uint32 { return le.Uint32((*[4]byte)(ptr(ui(p) + o))[:]) }
func readU64(p ptr, o ui) uint64 { return le.Uint64((*[8]byte)(ptr(ui(p) + o))[:]) }

func writeU64(p ptr, o ui, v u64) { le.PutUint64((*[8]byte)(ptr(ui(p) + o))[:], v)}

func xxhAvalancheSmall(x u64) u64 {
	x ^= x >> 33
	x *= prime64_2
	x ^= x >> 29
	x *= prime64_3
	x ^= x >> 32
	return x
}

func xxh3Avalanche(x u64) u64 {
	x ^= x >> 37
	x *= 0x165667919e3779f9
	x ^= x >> 32
	return x
}

func rrmxmx(h64 u64, len u64) u64 {
	h64 ^= bits.RotateLeft64(h64, 49) ^ bits.RotateLeft64(h64, 24)
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 35) + len
	h64 *= 0x9fb21c651e98df25
	h64 ^= (h64 >> 28)
	return h64
}

func mulFold64(x, y u64) u64 {
	hi, lo := bits.Mul64(x, y)
	return hi ^ lo
}
