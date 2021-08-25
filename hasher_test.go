package xxh3

import (
	"testing"
)

func TestHasherCompat(t *testing.T) {
	buf := make([]byte, 4097)
	for i := range buf {
		buf[i] = byte(i) % 251
	}

	for n := 0; n < 4097; n++ {
		h := New()

		h.Write(buf[:n/2])
		h.Reset()

		h.Write(buf[:n])

		check := func() {
			if exp, got := Hash(buf[:n]), h.Sum64(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
			if exp, got := Hash128(buf[:n]), h.Sum128(); exp != got {
				t.Fatalf("% -4d: %016x != %016x", n, exp, got)
			}
		}

		withAVX2(check)
		withSSE2(check)
		withGeneric(check)
	}
}
