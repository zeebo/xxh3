package xxh3

import "testing"

func override() (bool, bool, func()) {
	avx2Orig, sse2Orig := avx2, sse2
	return avx2Orig, sse2Orig, func() {
		avx2, sse2 = avx2Orig, sse2Orig
	}
}

func TestVectorCompat(t *testing.T) {
	check := func(b []byte) {
		t.Helper()

		avx2Orig, sse2Orig, cleanup := override()
		defer cleanup()

		avx2, sse2 = avx2Orig, false
		avx2Sum := Hash(b)

		avx2, sse2 = false, sse2Orig
		sse2Sum := Hash(b)

		avx2, sse2 = false, false
		cpuSum := Hash(b)

		if avx2Sum != sse2Sum || avx2Sum != cpuSum || sse2Sum != cpuSum {
			t.Errorf("data: %d", len(b))
			t.Errorf("avx2: %016x", avx2Sum)
			t.Errorf("sse2: %016x", sse2Sum)
			t.Errorf("cpu : %016x", cpuSum)
		}
	}

	for _, n := range []int{
		0, 1,
		63, 64, 65,
		127, 128, 129,
		191, 192, 193,
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
