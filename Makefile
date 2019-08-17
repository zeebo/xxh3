.PHONY: all
all: vector_avx_amd64.s vector_sse_amd64.s

vector_avx_amd64.s: _avo/avx.go
	GO111MODULE=on go run _avo/avx.go > vector_avx_amd64.s

vector_sse_amd64.s: _avo/sse.go
	GO111MODULE=on go run _avo/sse.go > vector_sse_amd64.s

clean:
	rm vector_avx_amd64.s
	rm vector_sse_amd64.s
