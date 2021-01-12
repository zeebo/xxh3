// +build !amd64

package xxh3

var avx2, sse2 bool

func override() (bool, bool, func()) {
	return false, false, func() {}
}
