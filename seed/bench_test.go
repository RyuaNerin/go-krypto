package seed

import (
	"bufio"
	"crypto/rand"
	"testing"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

func Benchmark_New(b *testing.B)     { benchNew(b, 128) }
func Benchmark_Encrypt(b *testing.B) { benchEncrypt(b, 128) }
func Benchmark_Decrypt(b *testing.B) { benchDecrypt(b, 128) }

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

	rnd.Read(src)

	b.ReportAllocs()
	b.SetBytes(BlockSize)
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

	rnd.Read(src)

	b.ReportAllocs()
	b.SetBytes(BlockSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, src)
		copy(src, dst)
	}
}
