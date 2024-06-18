package kipher_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"testing"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"

	. "github.com/RyuaNerin/testingutil"
)

func TestCFB(t *testing.T) {
	test := func(blockSize int) func(t *testing.T) {
		return func(t *testing.T) {
			key := make([]byte, 16)
			iv := make([]byte, 16)

			src := make([]byte, blocks)
			dstEnc := make([]byte, len(src))
			dstDec := make([]byte, len(src))

			rnd.Read(key)
			rnd.Read(iv)

			c, err := aes.NewCipher(key)
			if err != nil {
				t.Fatal(err)
			}

			enc := kipher.NewCFBEncrypterWithBlockSize(c, iv, blockSize)
			dec := kipher.NewCFBDecrypterWithBlockSize(c, iv, blockSize)

			for i := 0; i < iter; i++ {
				dataSize := 1 + rand.Intn(blocks-1)
				rnd.Read(src[:dataSize])

				enc.XORKeyStream(dstEnc[:dataSize], src[:dataSize])
				dec.XORKeyStream(dstDec[:dataSize], dstEnc[:dataSize])

				if !bytes.Equal(src[:dataSize], dstDec[:dataSize]) {
					t.Errorf("CFB: expected %x, got %x", src[:dataSize], dstDec[:dataSize])
					return
				}
			}
		}
	}

	t.Run("CFB8", test(1))
	t.Run("CFB32", test(4))
	t.Run("CFB64", test(8))
	t.Run("CFB128", test(16))
}

func TestCFBDecrypterWithStd(t *testing.T) {
	const bs = lea.BlockSize
	const srcSizeMax = 128 * bs

	key := make([]byte, 16)

	iv := make([]byte, bs)

	src := make([]byte, srcSizeMax)
	dstN := make([]byte, srcSizeMax)
	dstF := make([]byte, srcSizeMax)

	for i := 0; i < 256; i++ {
		rnd.Read(key)
		rnd.Read(iv)

		b, _ := lea.NewCipher(key)

		cfbN := cipher.NewCFBDecrypter(b, iv)
		cfbF := kipher.NewCFBDecrypter(ikipher.WrapKipher(b), iv)

		for j := 0; j < 1024; j++ {
			l := 1 + rand.Intn(srcSizeMax-1)

			rnd.Read(src[:l])

			cfbN.XORKeyStream(dstN[:l], src[:l])
			cfbF.XORKeyStream(dstF[:l], src[:l])

			if !bytes.Equal(dstN[:l], dstF[:l]) {
				t.Errorf("CFB: \nexpect: %x,\nactual: %x", dstN[:l], dstF[:l])
				return
			}
		}
	}
}

func TestCFBWithStd(t *testing.T) {
	type ctr struct {
		c, k cipher.Stream
	}

	test := func(encrypt bool) func(t *testing.T) {
		return func(t *testing.T) {
			BTTC(
				t,
				128,
				aes.BlockSize, // iv
				aes.BlockSize*16,
				1,
				func(key, iv []byte) (interface{}, error) {
					bc, err := aes.NewCipher(key)
					if err != nil {
						return nil, err
					}
					bk := ikipher.WrapKipher(bc)

					data := &ctr{}
					if encrypt {
						data.c = cipher.NewCFBEncrypter(bc, iv)
						data.k = kipher.NewCFBEncrypter(bk, iv)
					} else {
						data.c = cipher.NewCFBDecrypter(bc, iv)
						data.k = kipher.NewCFBDecrypter(bk, iv)
					}

					return data, nil
				},
				func(data interface{}, dst, src []byte) { data.(*ctr).c.XORKeyStream(dst, src) },
				func(data interface{}, dst, src []byte) { data.(*ctr).k.XORKeyStream(dst, src) },
				false,
			)
		}
	}

	t.Run("Encrypt", test(true))
	t.Run("Decrypt", test(false))
}
