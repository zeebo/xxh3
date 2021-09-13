package xxh3

import (
	"math/rand"
	"testing"
)

func TestHasherCompat(t *testing.T) {
	buf := make([]byte, 40970)
	for i := range buf {
		buf[i] = byte((i + 1) * 2654435761)
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
		buf[i] = byte((i + 1) * 2654435761)
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
