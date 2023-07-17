//go:build amd64

package lea

import (
	"bytes"
	"crypto/cipher"
	"testing"
)

var (
	cbcEnc = cipher.NewCBCEncrypter
	cbcDec = cipher.NewCBCDecrypter
)

func Test_LEA128_CBC_Encrypt_1Block(t *testing.T) { testCBC(t, 128, 1, cbcEnc) }
func Test_LEA196_CBC_Encrypt_1Block(t *testing.T) { testCBC(t, 196, 1, cbcEnc) }
func Test_LEA256_CBC_Encrypt_1Block(t *testing.T) { testCBC(t, 256, 1, cbcEnc) }

func Test_LEA128_CBC_Encrypt_4Blocks(t *testing.T) { testCBC(t, 128, 4, cbcEnc) }
func Test_LEA196_CBC_Encrypt_4Blocks(t *testing.T) { testCBC(t, 196, 4, cbcEnc) }
func Test_LEA256_CBC_Encrypt_4Blocks(t *testing.T) { testCBC(t, 256, 4, cbcEnc) }

func Test_LEA128_CBC_Encrypt_8Blocks(t *testing.T) { testCBC(t, 128, 8, cbcEnc) }
func Test_LEA196_CBC_Encrypt_8Blocks(t *testing.T) { testCBC(t, 196, 8, cbcEnc) }
func Test_LEA256_CBC_Encrypt_8Blocks(t *testing.T) { testCBC(t, 256, 8, cbcEnc) }

func Test_LEA128_CBC_Decrypt_1Block(t *testing.T) { testCBC(t, 128, 1, cbcDec) }
func Test_LEA196_CBC_Decrypt_1Block(t *testing.T) { testCBC(t, 196, 1, cbcDec) }
func Test_LEA256_CBC_Decrypt_1Block(t *testing.T) { testCBC(t, 256, 1, cbcDec) }

func Test_LEA128_CBC_Decrypt_4Blocks(t *testing.T) { testCBC(t, 128, 4, cbcDec) }
func Test_LEA196_CBC_Decrypt_4Blocks(t *testing.T) { testCBC(t, 196, 4, cbcDec) }
func Test_LEA256_CBC_Decrypt_4Blocks(t *testing.T) { testCBC(t, 256, 4, cbcDec) }

func Test_LEA128_CBC_Decrypt_8Blocks(t *testing.T) { testCBC(t, 128, 8, cbcDec) }
func Test_LEA196_CBC_Decrypt_8Blocks(t *testing.T) { testCBC(t, 196, 8, cbcDec) }
func Test_LEA256_CBC_Decrypt_8Blocks(t *testing.T) { testCBC(t, 256, 8, cbcDec) }

func testCBC(t *testing.T, keySize, blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) {
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

	cbcGo := newBlockMode(&ctxGo, iv)
	cbcAsm := newBlockMode(&ctxAsm, iv)

	for i := 0; i < testBlocks/blocks; i++ {
		cbcGo.CryptBlocks(dstGo, src)
		cbcAsm.CryptBlocks(dstAsm, src)

		if !bytes.Equal(dstGo, dstAsm) {
			t.Fail()
		}

		copy(src, dstAsm)
	}

}
