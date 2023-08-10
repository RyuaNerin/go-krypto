package lsh512

import (
	"hash"
	"testing"
)

func Benchmark_Hash_8_Go(b *testing.B)  { benchmarkSize(b, newContext, 8, true) }
func Benchmark_Hash_1K_Go(b *testing.B) { benchmarkSize(b, newContext, 1024, true) }
func Benchmark_Hash_8K_Go(b *testing.B) { benchmarkSize(b, newContext, 8192, true) }

var benchBuf = make([]byte, 8192)

func benchmarkSize(b *testing.B, newHash func(size int) hash.Hash, size int, do bool) {
	tests := []struct {
		name string
		size int
	}{
		{"512", Size},
		{"384", Size384},
		{"256", Size256},
		{"224", Size224},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			sum := make([]byte, Size)
			h := newHash(test.size)

			b.ResetTimer()
			b.ReportAllocs()
			b.SetBytes(int64(size))
			for i := 0; i < b.N; i++ {
				h.Reset()
				h.Write(benchBuf[:size])
				h.Sum(sum[:0])
			}
		})
	}
}
