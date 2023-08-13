package hight

import (
	"bufio"
	"crypto/rand"
	"testing"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

func Benchmark_New(b *testing.B) {
	k := make([]byte, KeySize)

	for i := 0; i < b.N; i++ {
		_, err := NewCipher(k)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_Encrypt(b *testing.B) {
	k := make([]byte, KeySize)
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)
	src := make([]byte, BlockSize)

	rand.Read(src)

	b.ReportAllocs()
	b.SetBytes(BlockSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(dst, src)
		copy(dst, src)
	}
}

func Benchmark_Decrypt(b *testing.B) {
	k := make([]byte, KeySize)
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)
	src := make([]byte, BlockSize)

	rand.Read(src)

	b.ReportAllocs()
	b.SetBytes(BlockSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, src)
		copy(dst, src)
	}
}
