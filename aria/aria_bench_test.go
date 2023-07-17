package aria

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	testKey128 = internal.HB(`00000000000000000000000000000000`)
	testKey196 = internal.HB(`000000000000000000000000000000000000000000000000`)
	testKey256 = internal.HB(`0000000000000000000000000000000000000000000000000000000000000000`)
	input      = internal.HB(`80000000000000000000000000000000`)
)

func BenchmarkNewCipher128(b *testing.B) { benchNewCipher(b, testKey128) }
func BenchmarkNewCipher196(b *testing.B) { benchNewCipher(b, testKey196) }
func BenchmarkNewCipher256(b *testing.B) { benchNewCipher(b, testKey256) }

func BenchmarkEncrypt128(b *testing.B) { benchEncrypt(b, testKey128) }
func BenchmarkEncrypt196(b *testing.B) { benchEncrypt(b, testKey196) }
func BenchmarkEncrypt256(b *testing.B) { benchEncrypt(b, testKey256) }

func BenchmarkDecrypt128(b *testing.B) { benchDescript(b, testKey128) }
func BenchmarkDecrypt196(b *testing.B) { benchDescript(b, testKey196) }
func BenchmarkDecrypt256(b *testing.B) { benchDescript(b, testKey256) }

func benchNewCipher(b *testing.B, k []byte) {
	for i := 0; i < b.N; i++ {
		_, err := NewCipher(k)
		if err != nil {
			b.Error(err)
		}
	}
}

func benchEncrypt(b *testing.B, k []byte) {
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(dst, input)
	}
}

func benchDescript(b *testing.B, k []byte) {
	c, err := NewCipher(k)
	if err != nil {
		b.Error(err)
	}

	dst := make([]byte, BlockSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(dst, input)
	}
}
