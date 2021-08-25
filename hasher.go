package xxh3

import (
	"encoding/binary"
	"hash"
)

// Hasher implements the hash.Hash interface
type Hasher struct {
	acc [8]u64
	blk u64
	len u64
	buf [_block + _stripe]byte
}

var (
	_ hash.Hash   = (*Hasher)(nil)
	_ hash.Hash64 = (*Hasher)(nil)
)

// New returns a new Hasher that implements the hash.Hash interface.
func New() *Hasher {
	var h Hasher
	h.Reset()
	return &h
}

// Reset resets the Hash to its initial state.
func (h *Hasher) Reset() {
	h.acc = [8]u64{
		prime32_3, prime64_1, prime64_2, prime64_3,
		prime64_4, prime32_2, prime64_5, prime32_1,
	}
	h.blk = 0
	h.len = 0
}

// BlockSize returns the hash's underlying block size.
// The Write method will accept any amount of data, but
// it may operate more efficiently if all writes are a
// multiple of the block size.
func (h *Hasher) BlockSize() int { return _stripe }

// Size returns the number of bytes Sum will return.
func (h *Hasher) Size() int { return 8 }

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (h *Hasher) Sum(b []byte) []byte {
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], h.Sum64())
	return append(b, tmp[:]...)
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *Hasher) Write(buf []byte) (int, error) {
	h.update(buf)
	return len(buf), nil
}

// WriteString adds more data to the running hash.
// It never returns an error.
func (h *Hasher) WriteString(buf string) (int, error) {
	h.updateString(buf)
	return len(buf), nil
}

func (h *Hasher) update(buf []byte) {
	// relies on the data pointer being the first word in the string header
	h.updateString(*(*string)(ptr(&buf)))
}

func (h *Hasher) updateString(buf string) {
	for len(buf) > 0 {
		if h.len < u64(len(h.buf)) {
			n := copy(h.buf[h.len:], buf)
			h.len += u64(n)
			buf = buf[n:]
			continue
		}

		if hasAVX2 {
			accumBlockAVX2(&h.acc, ptr(&h.buf), key)
		} else if hasSSE2 {
			accumBlockSSE(&h.acc, ptr(&h.buf), key)
		} else {
			accumBlockScalar(&h.acc, ptr(&h.buf), key)
		}

		h.blk++
		h.len = _stripe
		copy(h.buf[_block:], h.buf[:_stripe])
	}
}

// Sum64 returns the 64-bit hash of the written data.
func (h *Hasher) Sum64() uint64 {
	if h.blk == 0 {
		return Hash(h.buf[:h.len])
	}

	l := h.blk*_block + h.len
	acc := l * prime64_1
	accs := h.acc

	if h.len > 0 {
		if hasAVX2 {
			accumAVX2(&accs, ptr(&h.buf[0]), key, h.len)
		} else if hasSSE2 {
			accumSSE(&accs, ptr(&h.buf[0]), key, h.len)
		} else {
			accumScalar(&accs, ptr(&h.buf[0]), key, h.len)
		}
	}

	acc += mulFold64(accs[0]^key64_011, accs[1]^key64_019)
	acc += mulFold64(accs[2]^key64_027, accs[3]^key64_035)
	acc += mulFold64(accs[4]^key64_043, accs[5]^key64_051)
	acc += mulFold64(accs[6]^key64_059, accs[7]^key64_067)

	return xxh3Avalanche(acc)
}

// Sum64 returns the 128-bit hash of the written data.
func (h *Hasher) Sum128() Uint128 {
	if h.blk == 0 {
		return Hash128(h.buf[:h.len])
	}

	l := h.blk*_block + h.len
	acc := Uint128{Lo: l * prime64_1, Hi: ^(l * prime64_2)}
	accs := h.acc

	if h.len > 0 {
		if hasAVX2 {
			accumAVX2(&accs, ptr(&h.buf[0]), key, h.len)
		} else if hasSSE2 {
			accumSSE(&accs, ptr(&h.buf[0]), key, h.len)
		} else {
			accumScalar(&accs, ptr(&h.buf[0]), key, h.len)
		}
	}

	acc.Lo += mulFold64(accs[0]^key64_011, accs[1]^key64_019)
	acc.Lo += mulFold64(accs[2]^key64_027, accs[3]^key64_035)
	acc.Lo += mulFold64(accs[4]^key64_043, accs[5]^key64_051)
	acc.Lo += mulFold64(accs[6]^key64_059, accs[7]^key64_067)
	acc.Lo = xxh3Avalanche(acc.Lo)

	acc.Hi += mulFold64(accs[0]^key64_117, accs[1]^key64_125)
	acc.Hi += mulFold64(accs[2]^key64_133, accs[3]^key64_141)
	acc.Hi += mulFold64(accs[4]^key64_149, accs[5]^key64_157)
	acc.Hi += mulFold64(accs[6]^key64_165, accs[7]^key64_173)
	acc.Hi = xxh3Avalanche(acc.Hi)

	return acc
}
