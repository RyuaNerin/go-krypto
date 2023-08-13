package lsh256

import (
	"bufio"
	"crypto/rand"
	"hash"
	"testing"
)

func Benchmark_Hash_8_Go(b *testing.B)  { benchmarkSize(b, newContextGo, 8, true) }
func Benchmark_Hash_1K_Go(b *testing.B) { benchmarkSize(b, newContextGo, 1024, true) }
func Benchmark_Hash_8K_Go(b *testing.B) { benchmarkSize(b, newContextGo, 8192, true) }

var benchBuf = make([]byte, 8192)

func init() {
	rnd := bufio.NewReaderSize(rand.Reader, 1<<15)
	rnd.Read(benchBuf)
}

func benchmarkSize(b *testing.B, newHash func(size int) hash.Hash, size int, do bool) {
	tests := []struct {
		name string
		size int
	}{
		{"256", Size},
		{"224", Size224},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			sum := make([]byte, Size)
			h := newHash(test.size)

			b.ReportAllocs()
			b.SetBytes(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Reset()
				h.Write(benchBuf[:size])
				h.Sum(sum[:0])
			}
		})
	}
}
