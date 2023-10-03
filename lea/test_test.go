package lea

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

var (
	as = []CipherSize{
		{Name: "128", Size: 128},
		{Name: "196", Size: 196},
		{Name: "256", Size: 256},
	}
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) {
	BA(b, as, func(b *testing.B, keySize int) {
		BBNew(b, keySize, 0, BIW(NewCipher))
	})
}

func Benchmark_Encrypt_1Blocks_Go(b *testing.B) { BA(b, as, bb(b, 1, leaEnc1Go, true)) }
func Benchmark_Encrypt_4Blocks_Go(b *testing.B) { BA(b, as, bb(b, 4, leaEnc4Go, true)) }
func Benchmark_Encrypt_8Blocks_Go(b *testing.B) { BA(b, as, bb(b, 8, leaEnc8Go, true)) }

func Benchmark_Decrypt_1Blocks_Go(b *testing.B) { BA(b, as, bb(b, 1, leaDec1Go, true)) }
func Benchmark_Decrypt_4Blocks_Go(b *testing.B) { BA(b, as, bb(b, 4, leaDec4Go, true)) }
func Benchmark_Decrypt_8Blocks_Go(b *testing.B) { BA(b, as, bb(b, 8, leaDec8Go, true)) }

func bb(b *testing.B, blocks int, f funcBlock, do bool) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		if !do {
			b.Skip()
			return
		}

		BBDo(
			b,
			keySize,
			0,
			blocks*BlockSize,
			BIW(NewCipher),
			func(c interface{}, dst, src []byte) {
				f(c.(*leaContext), dst, src)
			},
		)
	}
}
