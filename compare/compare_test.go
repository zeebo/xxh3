package compare

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/zeebo/xxh3"
)

func BenchmarkZeebo(b *testing.B) {
	sizes := []int{
		0, 1, 3, 4, 8, 9, 16, 17, 32,
		33, 64, 65, 96, 97, 128, 129, 240, 241,
		512, 1024, 100 * 1024,
	}

	hashes := []struct {
		name string
		fn   func([]byte) uint64
	}{
		{"Zeebo", xxh3.Hash},
		{"Cespare", xxhash.Sum64},
	}

	for _, hash := range hashes {
		b.Run(hash.name, func(b *testing.B) {
			for _, size := range sizes {
				b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
					b.SetBytes(int64(size))
					var acc uint64
					d := make([]byte, size)
					b.ReportAllocs()
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						acc = hash.fn(d)
					}
					runtime.KeepAlive(acc)
				})
			}
		})
	}
}
