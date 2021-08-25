.PHONY: all vet
all: accum_vector_avx_amd64.s accum_vector_sse_amd64.s _compat

accum_vector_avx_amd64.s: avo/avx.go
	cd ./avo; go run . -avx > ../accum_vector_avx_amd64.s

accum_vector_sse_amd64.s: avo/sse.go
	cd ./avo; go run . -sse > ../accum_vector_sse_amd64.s

clean:
	rm accum_vector_avx_amd64.s
	rm accum_vector_sse_amd64.s
	rm _compat

upstream/xxhash.o: upstream/xxhash.h
	( cd upstream && make )

_compat: _compat.c upstream/xxhash.o
	gcc -o _compat _compat.c ./upstream/xxhash.o

vet:
	GOARCH=amd64 go vet
	GOARCH=386 go vet
	GOARCH=arm go vet