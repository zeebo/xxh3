package xxh3

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkHasherSeededSmall(b *testing.B) {
	for _, size := range []int{64, 240, 241, 256, 512, 1024, 1088, 1089, 4096} {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i*131 + 17)
		}
		for _, seed := range []uint64{0, 42} {
			b.Run(fmt.Sprintf("%d/64/seed-%d", size, seed), func(b *testing.B) {
				hasher := NewSeed(seed)
				var result uint64
				b.ReportAllocs()
				b.SetBytes(int64(size))
				b.ResetTimer()
				for range b.N {
					hasher.Reset()
					hasher.Write(data)
					result = hasher.Sum64()
				}
				runtime.KeepAlive(result)
			})
			b.Run(fmt.Sprintf("%d/128/seed-%d", size, seed), func(b *testing.B) {
				hasher := NewSeed128(seed)
				var result Uint128
				b.ReportAllocs()
				b.SetBytes(int64(size))
				b.ResetTimer()
				for range b.N {
					hasher.Reset()
					hasher.Write(data)
					result = hasher.Sum128()
				}
				runtime.KeepAlive(result)
			})
		}
	}
}
