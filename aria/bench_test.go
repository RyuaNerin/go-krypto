package aria

import (
	"testing"
)

func Benchmark_ARIA128_New(b *testing.B) { benchNewCipher(b, 128) }
func Benchmark_ARIA196_New(b *testing.B) { benchNewCipher(b, 196) }
func Benchmark_ARIA256_New(b *testing.B) { benchNewCipher(b, 256) }

func Benchmark_ARIA128_Encrypt(b *testing.B) { benchEncrypt(b, 128) }
func Benchmark_ARIA196_Encrypt(b *testing.B) { benchEncrypt(b, 196) }
func Benchmark_ARIA256_Encrypt(b *testing.B) { benchEncrypt(b, 256) }

func Benchmark_ARIA128_Decrypt(b *testing.B) { benchDecrypt(b, 128) }
func Benchmark_ARIA196_Decrypt(b *testing.B) { benchDecrypt(b, 196) }
func Benchmark_ARIA256_Decrypt(b *testing.B) { benchDecrypt(b, 256) }

func benchNewCipher(b *testing.B, keySize int) {
	k := make([]byte, keySize/8)

	for i := 0; i < b.N; i++ {
		_, err := NewCipher(k)
		if err != nil {
			b.Error(err)
		}
	}
}

func benchEncrypt(b *testing.B, keySize int) {
	k := make([]byte, keySize/8)
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(dst, src)
	}
}

func benchDecrypt(b *testing.B, keySize int) {
	k := make([]byte, keySize/8)
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, src)
	}
}
