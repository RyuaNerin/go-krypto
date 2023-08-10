package has160

import (
	"testing"
)

func Benchmark_Hash_8_Go(b *testing.B)  { benchmarkSize(b, 8) }
func Benchmark_Hash_1K_Go(b *testing.B) { benchmarkSize(b, 1024) }
func Benchmark_Hash_8K_Go(b *testing.B) { benchmarkSize(b, 8192) }

var benchBuf = make([]byte, 8192)

func benchmarkSize(b *testing.B, size int) {
	sum := make([]byte, Size)

	var h has160Context
	h.Reset()

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(benchBuf[:size])
		h.Sum(sum[:0])
	}
}
