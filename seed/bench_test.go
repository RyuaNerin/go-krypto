package seed

import (
	"testing"
)

func Benchmark_SEED128_New(b *testing.B)     { benchNew(b, 128) }
func Benchmark_SEED128_Encrypt(b *testing.B) { benchEncrypt(b, 128) }
func Benchmark_SEED128_Decrypt(b *testing.B) { benchDecrypt(b, 128) }

func benchNew(b *testing.B, keySize int) {
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
		copy(src, dst)
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
		copy(src, dst)
	}
}
