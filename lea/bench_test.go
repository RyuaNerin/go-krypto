package lea

import (
	"bufio"
	"crypto/rand"
	"testing"
)

func benchAll(b *testing.B, f func(*testing.B, int)) {
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

func Benchmark_New(b *testing.B) { benchAll(b, benchNewCipher) }

func Benchmark_Encrypt_1Block_Go(b *testing.B) { benchAll(b, block(b, 1, leaEnc1Go, true)) }
func Benchmark_Decrypt_1Block_Go(b *testing.B) { benchAll(b, block(b, 1, leaDec1Go, true)) }

func benchNewCipher(b *testing.B, keySize int) {
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
}

func block(b *testing.B, blocks int, f funcBlock, do bool) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		if !do {
			b.Skip()
			return
		}

		rnd := bufio.NewReaderSize(rand.Reader, 1<<15)

		k := make([]byte, keySize/8)
		rnd.Read(k)

		src := make([]byte, BlockSize*blocks)
		dst := make([]byte, BlockSize*blocks)

		rnd.Read(src)

		var ctx leaContext
		err := ctx.initContext(k)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(len(src)))
		for i := 0; i < b.N; i++ {
			f(&ctx, dst, src)
			copy(dst, src)
		}
	}
}
