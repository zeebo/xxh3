package xxh3

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2
	hasSSE2 = cpu.X86.HasSSE2
)

//go:noescape
func accumAVX2(acc *[8]u64, data, key unsafe.Pointer, len u64)

//go:noescape
func accumSSE(acc *[8]u64, data, key unsafe.Pointer, len u64)

//go:noescape
func accumBlockAVX2(acc *[8]u64, data, key unsafe.Pointer)

//go:noescape
func accumBlockSSE(acc *[8]u64, data, key unsafe.Pointer)

func withOverrides(avx2, sse2 bool, cb func()) {
	avx2Orig, sse2Orig := hasAVX2, hasSSE2
	hasAVX2, hasSSE2 = avx2, sse2
	defer func() { hasAVX2, hasSSE2 = avx2Orig, sse2Orig }()
	cb()
}

func withAVX2(cb func())    { withOverrides(hasAVX2, false, cb) }
func withSSE2(cb func())    { withOverrides(false, hasSSE2, cb) }
func withGeneric(cb func()) { withOverrides(false, false, cb) }
