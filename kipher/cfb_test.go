package kipher

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	. "github.com/RyuaNerin/testingutil"
)

func Test_CFB(t *testing.T) {
	rnd := bufio.NewReaderSize(rand.Reader, 1<<10)

	const iter = 1024
	const maxSize = 40960

	test := func(blockSize int) func(t *testing.T) {
		return func(t *testing.T) {
			key := make([]byte, 16)
			iv := make([]byte, 16)
			src := make([]byte, maxSize)

			dstEnc := make([]byte, len(src))
			dstDec := make([]byte, len(src))

			for i := 0; i < iter; i++ {
				rnd.Read(src[:4])
				dataSize := 1 + binary.BigEndian.Uint32(src[:4])%(maxSize-1)

				rnd.Read(key)
				rnd.Read(iv)
				rnd.Read(src[:dataSize])

				c, err := aes.NewCipher(key)
				if err != nil {
					t.Fatal(err)
				}
				enc := newCFB(c, iv, false, blockSize)
				dec := newCFB(c, iv, true, blockSize)

				enc.XORKeyStream(dstEnc[:dataSize], src[:dataSize])
				dec.XORKeyStream(dstDec[:dataSize], dstEnc[:dataSize])

				if !bytes.Equal(src[:dataSize], dstDec[:dataSize]) {
					t.Errorf("CFB: expected %x, got %x", src[:dataSize], dstDec[:dataSize])
					return
				}
			}
		}
	}

	t.Run("AES-CFB8", test(1))
	t.Run("AES-CFB32", test(4))
	t.Run("AES-CFB64", test(8))
	t.Run("AES-CFB128", test(16))
}

func Test_CFB_FullBlock_Encrypt(t *testing.T) { TA(t, as, testCFB(false), false) }
func Test_CFB_FullBlock_Decrypt(t *testing.T) { TA(t, as, testCFB(true), false) }

func testCFB(decrypt bool) func(t *testing.T, inputBlocks int) {
	fnCipher := cipher.NewCFBEncrypter
	if decrypt {
		fnCipher = cipher.NewCFBDecrypter
	}

	return func(t *testing.T, inputBlocks int) {
		type cfb struct {
			c, k cipher.Stream
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
				bk := internal.WrapBlock(bc)

				data := &cfb{
					c: fnCipher(bc, iv),
					k: newCFB(bk, iv, decrypt, bc.BlockSize()),
				}
				return data, nil
			},
			func(data interface{}, dst, src []byte) { data.(*cfb).c.XORKeyStream(dst, src) },
			func(data interface{}, dst, src []byte) { data.(*cfb).k.XORKeyStream(dst, src) },
			false,
		)
	}
}
