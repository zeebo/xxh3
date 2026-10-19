package xxh3

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	chunkedWriteSize      = 1 << 20
	chunkedWriteChunkSize = 32 << 10
)

func TestHasherCompat(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}

	for n := range buf {
		check := func() {
			h := New()
			h.Write(buf[:n/2])
			h.Reset()
			h.Write(buf[:n])
			if exp, got := Hash(buf[:n]), h.Sum64(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
			if exp, got := Hash128(buf[:n]), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
			seed := uint64(n)
			h.ResetSeed(seed)
			h.Write(buf[:n])
			if exp, got := HashSeed(buf[:n], seed), h.Sum64(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
			if exp, got := Hash128Seed(buf[:n], seed), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
		}

		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
		withGeneric(check)
	}

	chunked := make([]byte, _block+1+3*chunkedWriteChunkSize)
	for i := range chunked {
		chunked[i] = byte(uint64(i+1) * 2654435761)
	}
	check := func() {
		h := New()
		h.Write(chunked[:_block+1])
		for start := _block + 1; start < len(chunked); start += chunkedWriteChunkSize {
			h.Write(chunked[start:min(start+chunkedWriteChunkSize, len(chunked))])
		}
		if want, got := Hash(chunked), h.Sum64(); want != got {
			t.Fatalf("chunked Sum64: %016x != %016x", got, want)
		}
		if want, got := Hash128(chunked), h.Sum128(); want != got {
			t.Fatalf("chunked Sum128: %016x != %016x", got, want)
		}
	}
	withAVX512(check)
	withAVX2(check)
	withSSE2(check)
	withGeneric(check)
}

func TestHasher128Compat(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}

	for n := range buf {
		check := func() {
			h := New128()
			h.Write(buf[:n/2])
			h.Reset()
			h.Write(buf[:n])
			if exp, got := Hash128(buf[:n]), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}

			seed := uint64(n)
			h.ResetSeed(seed)
			h.Write(buf[:n])
			if exp, got := Hash128Seed(buf[:n], seed), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
		}

		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
		withGeneric(check)
	}
}

func TestHasherZeroValue(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}

	for n := range buf {
		check := func() {
			var h Hasher
			h.Write(buf[:n])
			if exp, got := Hash(buf[:n]), h.Sum64(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
			if exp, got := Hash128(buf[:n]), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
		}

		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
		withGeneric(check)
	}
}

func TestHasher128ZeroValue(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}

	for n := range buf {
		check := func() {
			var h Hasher128
			h.Write(buf[:n])
			if exp, got := Hash128(buf[:n]), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
		}

		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
		withGeneric(check)
	}
}

func TestHasherCompatSeed(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}
	rng := rand.New(rand.NewSource(42))

	for n := range buf {
		seed := rng.Uint64()

		check := func() {
			h := NewSeed(seed)

			h.Write(buf[:n/2])
			h.Reset()
			h.Write(buf[:n])

			if exp, got := HashSeed(buf[:n], seed), h.Sum64(); exp != got {
				t.Fatalf("Sum64: % -4d: %016x != %016x, seed:%x", n, exp, got, seed)
				return
			}
			if exp, got := Hash128Seed(buf[:n], seed), h.Sum128(); exp != got {
				t.Errorf("Sum128: % -4d: %016x != %016x", n, exp, got)
			}
		}

		withGeneric(check)
		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
	}
}

func TestHasher128CompatSeed(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}
	rng := rand.New(rand.NewSource(42))

	for n := range buf {
		seed := rng.Uint64()

		check := func() {
			h := NewSeed128(seed)

			h.Write(buf[:n/2])
			h.Reset()
			h.Write(buf[:n])

			if exp, got := Hash128Seed(buf[:n], seed), h.Sum128(); exp != got {
				t.Errorf("Sum128: % -4d: %016x != %016x", n, exp, got)
			}
		}

		withGeneric(check)
		withAVX512(check)
		withAVX2(check)
		withSSE2(check)
	}
}

func BenchmarkHasher64(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	for n := uint(4); n <= 28; n += 2 {
		size := 1 << n
		buf := make([]byte, size)
		for i := range buf {
			buf[i] = byte(uint64(i+int(n)+1) * 2654435761)
		}
		seed := rng.Uint64() | 1
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			var bn *testing.B
			check := func() {
				bn.Run("plain", func(b *testing.B) {
					h := New()
					b.ReportAllocs()
					b.ResetTimer()
					b.SetBytes(int64(size))
					for i := 0; i < b.N; i++ {
						h.Reset()
						h.Write(buf[:size])
						_ = h.Sum64()
					}
				})
				bn.Run("seed", func(b *testing.B) {
					h := NewSeed(seed)
					b.ReportAllocs()
					b.ResetTimer()
					b.SetBytes(int64(size))
					for i := 0; i < b.N; i++ {
						h.Reset()
						h.Write(buf[:size])
						_ = h.Sum64()
					}
				})
			}

			b.Run("go", func(b *testing.B) {
				bn = b
				withGeneric(check)
			})
			if hasAVX512 {
				b.Run("avx512", func(b *testing.B) {
					bn = b
					withAVX512(check)
				})
			}
			if hasAVX2 {
				b.Run("avx2", func(b *testing.B) {
					bn = b
					withAVX2(check)
				})
			}
			if hasSSE2 {
				b.Run("sse2", func(b *testing.B) {
					bn = b
					withSSE2(check)
				})
			}
		})
	}
}

func BenchmarkHasher64Chunked(b *testing.B) {
	buf := make([]byte, chunkedWriteSize)
	for i := range buf {
		buf[i] = byte(uint64(i+1) * 2654435761)
	}
	text := string(buf)

	for _, prefix := range []int{1, _stripe, _block + 1} {
		b.Run(fmt.Sprintf("prefix-%d", prefix), func(b *testing.B) {
			run := func(b *testing.B) {
				b.Run("bytes", func(b *testing.B) {
					benchmarkHasher64Chunked(b, buf, text, prefix, false)
				})
				b.Run("string", func(b *testing.B) {
					benchmarkHasher64Chunked(b, buf, text, prefix, true)
				})
			}

			b.Run("go", func(b *testing.B) {
				withGeneric(func() { run(b) })
			})
			if hasAVX512 {
				b.Run("avx512", func(b *testing.B) {
					withAVX512(func() { run(b) })
				})
			}
			if hasAVX2 {
				b.Run("avx2", func(b *testing.B) {
					withAVX2(func() { run(b) })
				})
			}
			if hasSSE2 {
				b.Run("sse2", func(b *testing.B) {
					withSSE2(func() { run(b) })
				})
			}
		})
	}
}

func benchmarkHasher64Chunked(
	b *testing.B,
	buf []byte,
	text string,
	prefix int,
	stringInput bool,
) {
	h := New()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for range b.N {
		h.Reset()
		if stringInput {
			h.WriteString(text[:prefix])
			for start := prefix; start < len(text); start += chunkedWriteChunkSize {
				h.WriteString(text[start:min(start+chunkedWriteChunkSize, len(text))])
			}
		} else {
			h.Write(buf[:prefix])
			for start := prefix; start < len(buf); start += chunkedWriteChunkSize {
				h.Write(buf[start:min(start+chunkedWriteChunkSize, len(buf))])
			}
		}
		_ = h.Sum64()
	}
}

func BenchmarkHasher128(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	for n := uint(4); n <= 28; n += 2 {
		size := 1 << n
		buf := make([]byte, size)
		for i := range buf {
			buf[i] = byte(uint64(i+int(n)+1) * 2654435761)
		}
		seed := rng.Uint64() | 1
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			var bn *testing.B
			check := func() {
				bn.Run("plain", func(b *testing.B) {
					h := New128()
					b.ReportAllocs()
					b.ResetTimer()
					b.SetBytes(int64(size))
					for i := 0; i < b.N; i++ {
						h.Reset()
						h.Write(buf[:size])
						_ = h.Sum128()
					}
				})
				bn.Run("seed", func(b *testing.B) {
					h := NewSeed128(seed)
					b.ReportAllocs()
					b.ResetTimer()
					b.SetBytes(int64(size))
					for i := 0; i < b.N; i++ {
						h.Reset()
						h.Write(buf[:size])
						_ = h.Sum128()
					}
				})
			}

			b.Run("go", func(b *testing.B) {
				bn = b
				withGeneric(check)
			})
			if hasAVX512 {
				b.Run("avx512", func(b *testing.B) {
					bn = b
					withAVX512(check)
				})
			}
			if hasAVX2 {
				b.Run("avx2", func(b *testing.B) {
					bn = b
					withAVX2(check)
				})
			}
			if hasSSE2 {
				b.Run("sse2", func(b *testing.B) {
					bn = b
					withSSE2(check)
				})
			}
		})
	}
}
