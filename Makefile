.PHONY: all
all: vector_avx_amd64.s vector_sse_amd64.s

vector_avx_amd64.s: avo/avx.go
	cd ./avo; go run . -avx > ../vector_avx_amd64.s

vector_sse_amd64.s: avo/sse.go
	cd ./avo; go run . -sse > ../vector_sse_amd64.s

clean:
	rm vector_avx_amd64.s
	rm vector_sse_amd64.s
