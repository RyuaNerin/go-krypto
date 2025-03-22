package kipher_test

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/kipherimpl"
	"github.com/RyuaNerin/go-krypto/kipher"

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
				c: cipher.NewCTR(kipherimpl.WrapCipher(b), iv),
				k: kipher.NewCTR(kipherimpl.WrapKipher(b), iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).c.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).k.XORKeyStream(dst, src) },
		false,
	)
}
