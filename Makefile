.PHONY: all vet
all: genasm _compat

genasm: avo/avx.go avo/sse.go
	cd ./avo; go generate gen.go

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