package kipher

import (
	"crypto/aes"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_BlockMode_CBC(t *testing.T) { TA(t, as, testCBC, false) }

func testCBC(t *testing.T, inputBlocks int) {
	type ctr struct {
		cbc, cbc2 BlockMode
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
			bk := blockWrap{bc}

			data := &ctr{
				cbc:  (*cbcDecrypter)(newCBC(bc, iv)),
				cbc2: newCBCDecrypter2(bk, iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).cbc.CryptBlocks(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).cbc2.CryptBlocks(dst, src) },
		false,
	)
}
