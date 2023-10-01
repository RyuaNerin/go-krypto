//go:build amd64 && gc && !purego

package lea

import (
	"bytes"
	"crypto/cipher"
	"testing"
)

func Test_BlockMode_CBC_Decrypt(t *testing.T) { testAll(t, testCBC(cipher.NewCBCDecrypter)) }

func testCBC(newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		key := make([]byte, keySize/8)
		iv := make([]byte, BlockSize)
		src := make([]byte, BlockSize*8)
		dstStd := make([]byte, BlockSize*8)
		dstAsm := make([]byte, BlockSize*8)
		rnd.Read(key)
		rnd.Read(iv)
		rnd.Read(src)

		var ctxStd nonCipherContext
		var ctxLib leaContext

		err := ctxStd.ctx.initContext(key)
		if err != nil {
			t.Error(err)
		}
		err = ctxLib.initContext(key)
		if err != nil {
			t.Error(err)
		}

		cbcGo := newBlockMode(&ctxStd, iv)
		cbcAsm := newBlockMode(&ctxLib, iv)

		l := make([]byte, 1)

		remain := testBlocks
		for remain > 0 {
			rnd.Read(l)
			lk := int(BlockSize * (1 + l[0]%7))

			cbcGo.CryptBlocks(dstStd, src[:lk])
			cbcAsm.CryptBlocks(dstAsm, src[:lk])

			if !bytes.Equal(dstStd, dstAsm) {
				t.Fail()
			}

			copy(src, dstAsm)
			remain -= lk
		}
	}
}
