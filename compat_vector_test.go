package xxh3

import "testing"

func TestVectorCompat(t *testing.T) {
	check := func(b []byte) {
		t.Helper()

		for i := range b {
			b[i] = byte(i)
		}

		var avx2Sum, sse2Sum, genericSum uint64

		withAVX2(func() { avx2Sum = Hash(b) })
		withSSE2(func() { sse2Sum = Hash(b) })
		withGeneric(func() { genericSum = Hash(b) })

		if avx2Sum != sse2Sum || avx2Sum != genericSum || sse2Sum != genericSum {
			t.Errorf("data  : %d", len(b))
			t.Errorf("avx2  : %016x", avx2Sum)
			t.Errorf("sse2  : %016x", sse2Sum)
			t.Errorf("scalar: %016x", genericSum)
			t.FailNow()
		}
	}

	for _, n := range []int{
		0, 1,
		63, 64, 65,
		127, 128, 129,
		191, 192, 193,
		239, 240, 241,
		255, 256, 257,
		319, 320, 321,
		383, 384, 385,
		447, 448, 449,
		511, 512, 513,
		575, 576, 577,
		639, 640, 641,
		703, 704, 705,
		767, 768, 769,
		831, 832, 833,
		895, 896, 897,
		959, 960, 961,
		1023, 1024, 1025,
		100 * 1024,
	} {
		check(make([]byte, n))
	}
}
