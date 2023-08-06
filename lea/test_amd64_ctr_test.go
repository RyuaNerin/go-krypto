//go:build amd64

package lea

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func Test_BlockMode_CTR_Asm(t *testing.T) { testAll(t, testCTR) }

func testCTR(t *testing.T, keySize int) {
	rnd := bufio.NewReaderSize(rand.Reader, 1<<15)

	var ctxGo leaContextGo
	var ctxAsm leaContextAsm

	key := make([]byte, keySize/8)
	rnd.Read(key)

	err := ctxGo.initContext(key)
	if err != nil {
		t.Error(err)
	}

	err = ctxAsm.g.initContext(key)
	if err != nil {
		t.Error(err)
	}

	iv := make([]byte, BlockSize)
	rnd.Read(iv)

	ctrGo := cipher.NewCTR(&ctxGo, iv)
	ctrAsm := cipher.NewCTR(&ctxGo, iv)

	src := make([]byte, BlockSize)
	dstGo := make([]byte, BlockSize)
	dstAsm := make([]byte, BlockSize)

	rnd.Read(src)
	for i := 0; i < testBlocks; i++ {
		ctrGo.XORKeyStream(dstGo, src)
		ctrAsm.XORKeyStream(dstAsm, src)

		if !bytes.Equal(dstGo, dstAsm) {
			t.Fail()
		}

		copy(src, dstAsm)
	}
}
