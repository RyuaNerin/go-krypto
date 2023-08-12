//go:build amd64

package lea

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func Test_BlockMode_CBC_Encrypt_1Block(t *testing.T) { testAll(t, testCBC(1, cipher.NewCBCEncrypter)) }
func Test_BlockMode_CBC_Decrypt_1Block(t *testing.T) { testAll(t, testCBC(1, cipher.NewCBCDecrypter)) }

func Test_BlockMode_CBC_Encrypt_4Blocks(t *testing.T) { testAll(t, testCBC(4, cipher.NewCBCEncrypter)) }
func Test_BlockMode_CBC_Decrypt_4Blocks(t *testing.T) { testAll(t, testCBC(4, cipher.NewCBCDecrypter)) }

func Test_BlockMode_CBC_Encrypt_8Blocks(t *testing.T) { testAll(t, testCBC(8, cipher.NewCBCEncrypter)) }
func Test_BlockMode_CBC_Decrypt_8Blocks(t *testing.T) { testAll(t, testCBC(8, cipher.NewCBCDecrypter)) }

func testCBC(blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		rnd := bufio.NewReaderSize(rand.Reader, 1<<15)

		var ctxGo leaContext
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

		cbcGo := newBlockMode(&ctxGo, iv)
		cbcAsm := newBlockMode(&ctxAsm, iv)

		src := make([]byte, BlockSize)
		dstGo := make([]byte, BlockSize)
		dstAsm := make([]byte, BlockSize)
		rnd.Read(src)

		for i := 0; i < testBlocks/blocks; i++ {
			cbcGo.CryptBlocks(dstGo, src)
			cbcAsm.CryptBlocks(dstAsm, src)

			if !bytes.Equal(dstGo, dstAsm) {
				t.Fail()
			}

			copy(src, dstAsm)
		}
	}
}
