//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_BlockMode_CBC_Decrypt(t *testing.T) { TA(t, as, testCBC(cipher.NewCBCDecrypter), false) }

func testCBC(newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		type ctr struct {
			std, asm cipher.BlockMode
		}

		BTTC(
			t,
			keySize,
			BlockSize, // iv
			BlockSize*8,
			BlockSize,
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
					std: newBlockMode(cStd, iv),
					asm: newBlockMode(cAsm, iv),
				}
				return data, nil
			},
			func(data interface{}, dst, src []byte) { data.(*ctr).std.CryptBlocks(dst, src) },
			func(data interface{}, dst, src []byte) { data.(*ctr).asm.CryptBlocks(dst, src) },
			false,
		)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BlockMode_CBC_Decrypt_Std(b *testing.B) {
	BA(b, as, benchCBC(newCipherSimple, cipher.NewCBCDecrypter), false)
}
func Benchmark_BlockMode_CBC_Decrypt_Asm(b *testing.B) {
	BA(b, as, benchCBC(NewCipher, cipher.NewCBCDecrypter), false)
}

func benchCBC(newCipher func(key []byte) (cipher.Block, error), newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		BBD(
			b,
			keySize,
			BlockSize,
			BlockSize*8,
			func(key, iv []byte) (interface{}, error) {
				ctx, err := newCipher(key)
				if err != nil {
					return nil, err
				}

				cbc := newBlockMode(ctx, iv)
				return cbc, nil
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.BlockMode).CryptBlocks(dst, src)
			},
			false,
		)
	}
}
