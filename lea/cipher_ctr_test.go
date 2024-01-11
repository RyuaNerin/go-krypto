//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lea

import (
	"crypto/cipher"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_BlockMode_CTR(t *testing.T) { TA(t, as, testCTR, false) }

func testCTR(t *testing.T, keySize int) {
	type ctr struct {
		std, asm cipher.Stream
	}

	BTTC(
		t,
		keySize,
		BlockSize, // iv
		BlockSize*8,
		1,
		func(key, iv []byte) (interface{}, error) {
			cStd, err := NewCipher(key)
			if err != nil {
				return nil, err
			}
			cAsm, err := newCipherSimple(key)
			if err != nil {
				return nil, err
			}

			data := &ctr{
				std: cipher.NewCTR(cStd, iv),
				asm: cipher.NewCTR(cAsm, iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).std.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).asm.XORKeyStream(dst, src) },
		false,
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BlockMode_CTR_Std(b *testing.B)    { BA(b, as, benchCTR(newCipherSimple), false) }
func Benchmark_BlockMode_CTR_Krypto(b *testing.B) { BA(b, as, benchCTR(NewCipher), false) }

func benchCTR(
	newCipher func(key []byte) (cipher.Block, error),
) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		BBD(
			b,
			keySize,
			BlockSize,
			BlockSize*8,
			func(key, additional []byte) (interface{}, error) {
				ctx, err := newCipher(key)
				if err != nil {
					return nil, err
				}

				ctr := cipher.NewCTR(ctx, make([]byte, BlockSize))
				return ctr, nil
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.Stream).XORKeyStream(dst, src)
			},
			false,
		)
	}
}
