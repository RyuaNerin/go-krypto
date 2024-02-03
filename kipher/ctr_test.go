package kipher

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_BlockMode_CTR(t *testing.T) { TA(t, as, testCTR, false) }

func testCTR(t *testing.T, bufferBlocks int) {
	type ctr struct {
		c, k cipher.Stream
	}

	BTTC(
		t,
		128,
		aes.BlockSize, // iv
		aes.BlockSize*bufferBlocks,
		1,
		func(key, iv []byte) (interface{}, error) {
			bc, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}
			bk := blockWrap{bc}

			data := &ctr{
				c: cipher.NewCTR(bc, iv),
				k: NewCTR(bk, iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).c.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).k.XORKeyStream(dst, src) },
		false,
	)
}
