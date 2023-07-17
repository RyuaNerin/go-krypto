package lea

import (
	"testing"
)

func Benchmark_LEA128_New(b *testing.B) { benchNewCipher(b, 128) }
func Benchmark_LEA196_New(b *testing.B) { benchNewCipher(b, 196) }
func Benchmark_LEA256_New(b *testing.B) { benchNewCipher(b, 256) }

func Benchmark_LEA128_Encrypt_1Block_Go(b *testing.B) { block(b, 128, 1, leaEnc1Go, true) }
func Benchmark_LEA196_Encrypt_1Block_Go(b *testing.B) { block(b, 196, 1, leaEnc1Go, true) }
func Benchmark_LEA256_Encrypt_1Block_Go(b *testing.B) { block(b, 256, 1, leaEnc1Go, true) }

func Benchmark_LEA128_Decrypt_1Block_Go(b *testing.B) { block(b, 128, 1, leaDec1Go, true) }
func Benchmark_LEA196_Decrypt_1Block_Go(b *testing.B) { block(b, 196, 1, leaDec1Go, true) }
func Benchmark_LEA256_Decrypt_1Block_Go(b *testing.B) { block(b, 256, 1, leaDec1Go, true) }

func benchNewCipher(b *testing.B, keySize int) {
	k := make([]byte, keySize/8)

	for i := 0; i < b.N; i++ {
		_, err := NewCipher(k)
		if err != nil {
			b.Error(err)
		}
	}
}

func block(b *testing.B, keySize int, blocks int, f funcBlock, do bool) {
	if !do {
		b.Skip()
		return
	}

	k := make([]byte, keySize/8)

	src := make([]byte, BlockSize*blocks)
	dst := make([]byte, BlockSize*blocks)

	var ctx leaContextGo
	err := ctx.initContext(k)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(ctx.round, ctx.rk, dst, src)
	}
}
