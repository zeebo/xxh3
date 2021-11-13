package main

import "flag"

//go:generate go run . -avx -out ../accum_vector_avx_amd64.s -pkg xxh3
//go:generate go run . -avx512 -out ../accum_vector_avx512_amd64.s -pkg xxh3
//go:generate go run . -sse -out ../accum_vector_sse_amd64.s -pkg xxh3

var (
	avx    = flag.Bool("avx", false, "run avx generation")
	avx512 = flag.Bool("avx512", false, "run avx512 generation")
	sse    = flag.Bool("sse", false, "run sse generation")
)

func main() {
	flag.Parse()

	if *avx {
		AVX()
	} else if *sse {
		SSE()
	} else if *avx512 {
		AVX512()
	}
}
