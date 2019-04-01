package xxh3

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkFixed(b *testing.B) {
	r := func(i int) {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.SetBytes(int64(i))
			var acc uint64
			d := bytes.Repeat([]byte("x"), i)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				acc = Hash(d)
			}
			runtime.KeepAlive(acc)
		})
	}

	r(0)
	r(1)
	r(3)
	r(4)
	r(8)
	r(9)
	r(16)
	r(17)
	r(32)
	r(33)
	r(64)
	r(65)
	r(96)
	r(97)
	r(128)
	r(129)
	r(256)
	r(512)
	r(1024)
	r(100 * 1024)
}
