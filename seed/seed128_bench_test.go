package seed

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	testKey128 = internal.HB(`00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00`)
	input      = internal.HB(`00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F`)
)

func BenchmarkNewCipher128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewCipher(testKey128)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEncrypt128(b *testing.B) {
	c, err := NewCipher(testKey128)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(dst, input)
	}
}

func BenchmarkDecrypt128(b *testing.B) {
	c, err := NewCipher(testKey128)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, input)
	}
}
