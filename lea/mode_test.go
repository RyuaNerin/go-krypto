//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lea

import (
	"crypto/cipher"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BlockMode_CBC_Decrypt_Std(b *testing.B) {
	BA(b, as, benchCBC(newCipherSimple, cipher.NewCBCDecrypter), false)
}

func Benchmark_BlockMode_CBC_Decrypt_Asm(b *testing.B) {
	BA(b, as, benchCBC(NewCipher, cipher.NewCBCDecrypter), false)
}

func Benchmark_BlockMode_CTR_Std(b *testing.B) {
	BA(b, as, benchCTR(newCipherSimple), false)
}

func Benchmark_BlockMode_CTR_Krypto(b *testing.B) {
	BA(b, as, benchCTR(NewCipher), false)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

func newCipherSimple(key []byte) (cipher.Block, error) {
	block, err := NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &blockWrap{block}, nil
}

type blockWrap struct {
	b cipher.Block
}

func (bw *blockWrap) BlockSize() int {
	return bw.b.BlockSize()
}

func (bw *blockWrap) Encrypt(dst, src []byte) {
	bw.b.Encrypt(dst, src)
}

func (bw *blockWrap) Decrypt(dst, src []byte) {
	bw.b.Decrypt(dst, src)
}
