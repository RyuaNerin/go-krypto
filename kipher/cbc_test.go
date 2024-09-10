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

func TestCBCZero(t *testing.T) {
	key := make([]byte, 16)
	iv := make([]byte, 16)

	b, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	enc := kipher.NewCBCEncrypter(ikipher.WrapKipher(b), iv)
	dec := kipher.NewCBCDecrypter(ikipher.WrapKipher(b), iv)

	rnd.Read(key)
	rnd.Read(iv)

	src := make([]byte, 0)
	dst := make([]byte, 0)

	enc.CryptBlocks(dst, src)
	dec.CryptBlocks(dst, src)
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

func BenchmarkCBCDEcryptor(b *testing.B) {
	bench := func(
		newCipher func([]byte) (cipher.Block, error),
		blocks int,
		useKipher bool,
	) func(b *testing.B) {
		const keySize = 16

		b, _ := newCipher(make([]byte, keySize))
		blockSize := b.BlockSize()

		return func(b *testing.B) {
			BBD(
				b,
				keySize*8,
				blockSize,
				blocks*blockSize,
				func(key, iv []byte) (interface{}, error) {
					block, err := newCipher(key)
					if err != nil {
						return nil, err
					}

					if useKipher {
						return kipher.NewCBCDecrypter(ikipher.WrapKipher(block), iv), nil
					} else {
						return cipher.NewCBCDecrypter(ikipher.WrapCipher(block), iv), nil
					}
				},
				func(c interface{}, dst, src []byte) {
					c.(cipher.BlockMode).CryptBlocks(dst, src)
				},
				false,
			)
		}
	}

	b.Run("AES/1/Std", bench(aes.NewCipher, 1, false))
	b.Run("AES/1/Krypto", bench(aes.NewCipher, 1, true))
	b.Run("AES/4/Std", bench(aes.NewCipher, 4, false))
	b.Run("AES/4/Krypto", bench(aes.NewCipher, 4, true))
	b.Run("AES/8/Std", bench(aes.NewCipher, 8, false))
	b.Run("AES/8/Krypto", bench(aes.NewCipher, 8, true))
	b.Run("AES/64/Std", bench(aes.NewCipher, 64, false))
	b.Run("AES/64/Krypto", bench(aes.NewCipher, 64, true))

	b.Run("LEA/1/Std", bench(lea.NewCipher, 1, false))
	b.Run("LEA/1/Krypto", bench(lea.NewCipher, 1, true))
	b.Run("LEA/4/Std", bench(lea.NewCipher, 4, false))
	b.Run("LEA/4/Krypto", bench(lea.NewCipher, 4, true))
	b.Run("LEA/8/Std", bench(lea.NewCipher, 8, false))
	b.Run("LEA/8/Krypto", bench(lea.NewCipher, 8, true))
	b.Run("LEA/64/Std", bench(lea.NewCipher, 64, false))
	b.Run("LEA/64/Krypto", bench(lea.NewCipher, 64, true))
}
