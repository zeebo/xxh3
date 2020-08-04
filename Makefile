.PHONY: all
all: vector_avx_amd64.s vector_sse_amd64.s _compat

vector_avx_amd64.s: avo/avx.go
	cd ./avo; go run . -avx > ../vector_avx_amd64.s

vector_sse_amd64.s: avo/sse.go
	cd ./avo; go run . -sse > ../vector_sse_amd64.s

clean:
	rm vector_avx_amd64.s
	rm vector_sse_amd64.s
	rm _compat

upstream/xxhash.o: upstream/xxhash.h
	( cd upstream && make )

_compat: _compat.c upstream/xxhash.o
	gcc -o _compat _compat.c ./upstream/xxhash.o
