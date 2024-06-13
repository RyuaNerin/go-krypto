package kipher_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"

	. "github.com/RyuaNerin/testingutil"
)

func TestCTRWithStd(t *testing.T) {
	type ctr struct {
		c, k cipher.Stream
	}

	BTTC(
		t,
		128,
		aes.BlockSize, // iv
		aes.BlockSize*16,
		1,
		func(key, iv []byte) (interface{}, error) {
			b, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}

			data := &ctr{
				c: cipher.NewCTR(ikipher.WrapCipher(b), iv),
				k: kipher.NewCTR(ikipher.WrapKipher(b), iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).c.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).k.XORKeyStream(dst, src) },
		false,
	)
}

func BenchmarkCTR(b *testing.B) {
	b.Run("AES/std", benchCTR(aes.NewCipher, false))
	b.Run("AES/krypto", benchCTR(aes.NewCipher, true))
	b.Run("LEA/std", benchCTR(lea.NewCipher, false))
	b.Run("LEA/krypto", benchCTR(lea.NewCipher, true))
}

func benchCTR(
	newCipher func([]byte) (cipher.Block, error),
	useKipher bool,
) func(b *testing.B) {
	b, _ := newCipher(make([]byte, keySize))
	blockSize := b.BlockSize()

	return func(b *testing.B) {
		BBD(
			b,
			keySize*8,
			blockSize,
			blocks*blockSize,
			func(key, iv []byte) (interface{}, error) {
				block, err := newCipher(key)
				if err != nil {
					return nil, err
				}

				if useKipher {
					return kipher.NewCTR(ikipher.WrapKipher(block), iv), nil
				} else {
					return cipher.NewCTR(ikipher.WrapCipher(block), iv), nil
				}
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.Stream).XORKeyStream(dst, src)
			},
			false,
		)
	}
}
