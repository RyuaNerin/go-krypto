package aria

import (
	"bufio"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

var (
	rnd = bufio.NewReaderSize(rand.Reader, 1<<15)
)

func Benchmark_ARIA_New(b *testing.B) {
	benchmarkAllSizes(
		b,
		func(b *testing.B, keySize int) {
			rnd := bufio.NewReaderSize(rand.Reader, 1<<15)
			k := make([]byte, keySize/8)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rnd.Read(k)
				_, err := NewCipher(k)
				if err != nil {
					b.Error(err)
				}
			}
		},
	)
}

func Benchmark_ARIA_Encrypt(b *testing.B) {
	benchmarkAllSizesBlock(b, func(c cipher.Block, dst, src []byte) { c.Encrypt(dst, src) })
}

func Benchmark_Decrypt(b *testing.B) {
	benchmarkAllSizesBlock(b, func(c cipher.Block, dst, src []byte) { c.Decrypt(dst, src) })
}

func benchmarkAllSizes(b *testing.B, f func(*testing.B, int)) {
	tests := []struct {
		name    string
		keySize int
	}{
		{"128", 128},
		{"196", 196},
		{"256", 256},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			f(b, test.keySize)
		})
	}
}

func benchmarkAllSizesBlock(b *testing.B, do func(c cipher.Block, dst []byte, src []byte)) {
	benchmarkAllSizes(
		b,
		func(b *testing.B, keySize int) {
			k := make([]byte, keySize/8)
			c, err := NewCipher(k)
			if err != nil {
				b.Error(err)
			}

			src := make([]byte, BlockSize)
			dst := make([]byte, BlockSize)

			rnd.Read(src)

			b.ReportAllocs()
			b.SetBytes(BlockSize)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				do(c, dst, src)
				copy(src, dst)
			}
		},
	)
}
