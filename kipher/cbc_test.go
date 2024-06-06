package kipher

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	. "github.com/RyuaNerin/testingutil"
)

func Test_CBC(t *testing.T) { TA(t, as, testCBC, false) }

func testCBC(t *testing.T, inputBlocks int) {
	type ctr struct {
		c, k cipher.BlockMode
	}

	BTTC(
		t,
		128,
		aes.BlockSize, // iv
		aes.BlockSize*inputBlocks,
		aes.BlockSize,
		func(key, iv []byte) (interface{}, error) {
			bc, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}
			bk := internal.WrapBlock(bc)

			data := &ctr{
				c: cipher.NewCBCDecrypter(bc, iv),
				k: NewCBCDecrypter(bk, iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).c.CryptBlocks(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).k.CryptBlocks(dst, src) },
		false,
	)
}
