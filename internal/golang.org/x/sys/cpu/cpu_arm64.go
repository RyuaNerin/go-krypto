package cpu

// cacheLineSize is used to prevent false sharing of cache lines.
// We choose 128 because Apple Silicon, a.k.a. M1, has 128-byte cache line size.
// It doesn't cost much and is much more future-proof.
const cacheLineSize = 128

func initOptions() {
	options = []option{}
}

func archInit() {
}
