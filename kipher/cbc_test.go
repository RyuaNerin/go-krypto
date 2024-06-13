package kipher_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"testing"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"

	. "github.com/RyuaNerin/testingutil"
)

func TestCBC(t *testing.T) {
	const blockSize = aes.BlockSize

	key := make([]byte, 16)
	iv := make([]byte, 16)

	src := make([]byte, blocks*blockSize)
	dstEnc := make([]byte, len(src))
	dstDec := make([]byte, len(src))

	b, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	enc := kipher.NewCBCEncrypter(ikipher.WrapKipher(b), iv)
	dec := kipher.NewCBCDecrypter(ikipher.WrapKipher(b), iv)

	rnd.Read(key)
	rnd.Read(iv)

	for i := 0; i < iter; i++ {
		dataSize := (1 + rand.Intn(blocks-1)) * blockSize

		rnd.Read(src[:dataSize])

		enc.CryptBlocks(dstEnc[:dataSize], src[:dataSize])
		dec.CryptBlocks(dstDec[:dataSize], dstEnc[:dataSize])

		if !bytes.Equal(src[:dataSize], dstDec[:dataSize]) {
			t.Errorf("CBC: expected %x, got %x", src[:dataSize], dstDec[:dataSize])
			return
		}
	}
}

func TestCBCDecrypterWitStd(t *testing.T) {
	type ctr struct {
		c, k cipher.BlockMode
	}

	BTTC(
		t,
		128,
		aes.BlockSize, // iv
		aes.BlockSize*16,
		aes.BlockSize,
		func(key, iv []byte) (interface{}, error) {
			bc, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}

			data := &ctr{
				c: cipher.NewCBCDecrypter(ikipher.WrapCipher(bc), iv),
				k: kipher.NewCBCDecrypter(ikipher.WrapKipher(bc), iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).c.CryptBlocks(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).k.CryptBlocks(dst, src) },
		false,
	)
}
