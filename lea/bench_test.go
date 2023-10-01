package lea

import (
	"bufio"
	"crypto/rand"
	"testing"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

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

func Benchmark_Encrypt_4Blocks_Go(b *testing.B) { benchAll(b, block(b, 4, leaEnc4Go, true)) }
func Benchmark_Encrypt_8Blocks_Go(b *testing.B) { benchAll(b, block(b, 8, leaEnc8Go, true)) }

func Benchmark_Decrypt_4Blocks_Go(b *testing.B) { benchAll(b, block(b, 4, leaDec4Go, true)) }
func Benchmark_Decrypt_8Blocks_Go(b *testing.B) { benchAll(b, block(b, 8, leaDec8Go, true)) }

func benchNewCipher(b *testing.B, keySize int) {
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

		key := make([]byte, keySize/8)
		src := make([]byte, BlockSize*blocks)
		dst := make([]byte, BlockSize*blocks)

		rnd.Read(key)
		rnd.Read(src)

		var ctx leaContext
		err := ctx.initContext(key)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f(&ctx, dst, src)
			copy(dst, src)
		}
	}
}
