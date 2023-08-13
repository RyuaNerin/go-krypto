package lsh256

import (
	"hash"
	"testing"
)

var benchBuf = make([]byte, 8192)

func init() {
	for idx := range benchBuf {
		benchBuf[idx] = byte(idx % 256)
	}
}

func benchmarkSize(b *testing.B, newHash func(size int) hash.Hash, size int, do bool) {
	tests := []struct {
		name string
		int  int
	}{
		{"224", Size224},
		{"256", Size},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			sum := make([]byte, Size)
			h := newHash(test.int)

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
