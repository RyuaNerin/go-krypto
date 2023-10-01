//go:build amd64 && gc && !purego

package lea

import (
	"bytes"
	"crypto/cipher"
	"testing"
)

func Test_BlockMode_CTR(t *testing.T) { testAll(t, testCTR) }

func testCTR(t *testing.T, keySize int) {
	key := make([]byte, keySize/8)
	iv := make([]byte, BlockSize)
	src := make([]byte, BlockSize*8)
	dstStd := make([]byte, BlockSize*8)
	dstAsm := make([]byte, BlockSize*8)
	rnd.Read(key)
	rnd.Read(iv)
	rnd.Read(src)

	var ctxStd nonCipherContext
	var ctxAsm leaContext

	err := ctxStd.ctx.initContext(key)
	if err != nil {
		t.Error(err)
	}

	err = ctxAsm.initContext(key)
	if err != nil {
		t.Error(err)
	}

	ctrStd := cipher.NewCTR(&ctxStd, iv)
	ctrAsm := cipher.NewCTR(&ctxAsm, iv)

	l := make([]byte, 1)

	remain := testBlocks * 8
	for remain > 0 {
		rnd.Read(l)
		lk := int(1 + l[0]%(BlockSize*8-1))

		ctrStd.XORKeyStream(dstStd, src[:lk])
		ctrAsm.XORKeyStream(dstAsm, src[:lk])

		if !bytes.Equal(dstStd, dstAsm) {
			t.Fail()
		}

		copy(src, dstAsm[:lk])
		remain -= lk
	}
}
