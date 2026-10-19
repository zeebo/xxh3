package xxh3

import (
	"fmt"
	"testing"
)

func TestHasherSeededSnapshots(t *testing.T) {
	data := make([]byte, 1153)
	for i := range data {
		data[i] = byte(i*131 + 17)
	}
	for _, size := range []int{240, 241, 512, 1024, 1088, 1089} {
		for _, seed := range []uint64{0, 1, 42, ^uint64(0)} {
			t.Run(fmt.Sprintf("size-%d/seed-%d", size, seed), func(t *testing.T) {
				check := func() {
					hasher := NewSeed(seed)
					hasher.Write(data[:size])
					for snapshot := 0; snapshot < 2; snapshot++ {
						if got, want := hasher.Sum64(), HashSeed(data[:size], seed); got != want {
							t.Fatalf("snapshot %d Sum64: %016x != %016x", snapshot, got, want)
						}
						if got, want := hasher.Sum128(), Hash128Seed(data[:size], seed); got != want {
							t.Fatalf("snapshot %d Sum128: %032x != %032x", snapshot, got, want)
						}
					}
					hasher.WriteString(string(data[size : size+64]))
					if got, want := hasher.Sum64(), HashSeed(data[:size+64], seed); got != want {
						t.Fatalf("after write Sum64: %016x != %016x", got, want)
					}
					if got, want := hasher.Sum128(), Hash128Seed(data[:size+64], seed); got != want {
						t.Fatalf("after write Sum128: %032x != %032x", got, want)
					}
				}
				withGeneric(check)
				withSSE2(check)
				withAVX2(check)
				withAVX512(check)
			})
		}
	}
}
