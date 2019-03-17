package xxh3

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

func Benchmark(b *testing.B) {
	r := func(i int) {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.SetBytes(int64(i))
			var acc uint64
			d := bytes.Repeat([]byte("x"), i)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				acc = Hash(d)
			}
			runtime.KeepAlive(acc)
		})
	}

	for i := 0; i < 24; i++ {
		r(i)
	}

	for i := 24; i < 128; i += 8 {
		r(i)
	}

	for i := 128; i <= 1024*1024; i *= 2 {
		r(i)
	}
}
