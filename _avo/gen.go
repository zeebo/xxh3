package main

import "flag"

var (
	avx = flag.Bool("avx", false, "run avx generation")
	sse = flag.Bool("sse", false, "run sse generation")
)

func main() {
	flag.Parse()

	if *avx {
		AVX()
	} else if *sse {
		SSE()
	}
}
