//go:build amd64

package lea

import (
	"bytes"
	"crypto/cipher"
	"testing"
)

func Test_LEA128_CTR_Asm(t *testing.T) { testCTR(t, 128) }
func Test_LEA196_CTR_Asm(t *testing.T) { testCTR(t, 196) }
func Test_LEA256_CTR_Asm(t *testing.T) { testCTR(t, 256) }

func testCTR(t *testing.T, keySize int) {
	var ctxGo leaContextGo
	var ctxAsm leaContextAsm

	key := make([]byte, keySize/8)

	err := ctxGo.initContext(key)
	if err != nil {
		t.Error(err)
	}

	err = ctxAsm.g.initContext(key)
	if err != nil {
		t.Error(err)
	}

	iv := make([]byte, BlockSize)
	src := make([]byte, BlockSize)
	dstGo := make([]byte, BlockSize)
	dstAsm := make([]byte, BlockSize)

	ctrGo := cipher.NewCTR(&ctxGo, iv)
	ctrAsm := cipher.NewCTR(&ctxGo, iv)

	for i := 0; i < testBlocks; i++ {
		ctrGo.XORKeyStream(dstGo, src)
		ctrAsm.XORKeyStream(dstAsm, src)

		if !bytes.Equal(dstGo, dstAsm) {
			t.Fail()
		}

		copy(src, dstAsm)
	}
}
