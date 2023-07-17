package hight

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	testKey = internal.Reverse(internal.HB(`00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00`))
	input   = internal.Reverse(internal.HB(`00 00 00 00 00 00 00 00`))
)

func BenchmarkNewCipher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewCipher(testKey)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {
	c, err := NewCipher(testKey)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(dst, input)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	c, err := NewCipher(testKey)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, input)
	}
}
