package kipher

import (
	"crypto/aes"
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_BlockMode_CTR(t *testing.T) { TA(t, as, testCTR, false) }

func testCTR(t *testing.T, bufferBlocks int) {
	type ctr struct {
		std, asm Stream
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
				std: newCTR(bc, iv, bufferBlocks),
				asm: newCTR2(bk, iv, bufferBlocks),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).std.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).asm.XORKeyStream(dst, src) },
		false,
	)
}
