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

			enc := kipher.NewCFBEncrypter(c, iv, blockSize)
			dec := kipher.NewCFBDecrypter(c, iv, blockSize)

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
						data.k = kipher.NewCFBEncrypter(bk, iv, bk.BlockSize())
					} else {
						data.c = cipher.NewCFBDecrypter(bc, iv)
						data.k = kipher.NewCFBDecrypter(bk, iv, bk.BlockSize())
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

func BenchmarkCFB128Encrypter(b *testing.B) {
	b.Run("AES/std", benchCFB(aes.NewCipher, aes.BlockSize, true, false))
	b.Run("AES/krypto", benchCFB(aes.NewCipher, aes.BlockSize, true, true))
	b.Run("LEA/std", benchCFB(lea.NewCipher, lea.BlockSize, true, false))
	b.Run("LEA/krypto", benchCFB(lea.NewCipher, lea.BlockSize, true, true))
}

func BenchmarkCFB128Decrypter(b *testing.B) {
	b.Run("AES/std", benchCFB(aes.NewCipher, aes.BlockSize, false, false))
	b.Run("AES/krypto", benchCFB(aes.NewCipher, aes.BlockSize, false, true))
	b.Run("LEA/std", benchCFB(lea.NewCipher, lea.BlockSize, false, false))
	b.Run("LEA/krypto", benchCFB(lea.NewCipher, lea.BlockSize, false, true))
}

func BenchmarkCFB64(b *testing.B) {
	b.Run("LEA/Encrypt", benchCFB(lea.NewCipher, 8, true, true))
	b.Run("LEA/Decrypt", benchCFB(lea.NewCipher, 8, false, true))
}

func BenchmarkCFB32(b *testing.B) {
	b.Run("LEA/Encrypt", benchCFB(lea.NewCipher, 4, true, true))
	b.Run("LEA/Decrypt", benchCFB(lea.NewCipher, 4, false, true))
}

func BenchmarkCFB8(b *testing.B) {
	b.Run("LEA/Encrypt", benchCFB(lea.NewCipher, 1, true, true))
	b.Run("LEA/Decrypt", benchCFB(lea.NewCipher, 1, false, true))
}

func benchCFB(
	newCipher func([]byte) (cipher.Block, error),
	cfbBlockSize int,
	encryptMode bool,
	useKipher bool,
) func(b *testing.B) {
	const keySize = 16
	const blocks = 8 + 4 + 1

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
					if encryptMode {
						return kipher.NewCFBEncrypter(ikipher.WrapKipher(block), iv, cfbBlockSize), nil
					} else {
						return kipher.NewCFBDecrypter(ikipher.WrapKipher(block), iv, cfbBlockSize), nil
					}
				} else {
					if encryptMode {
						return cipher.NewCFBEncrypter(ikipher.WrapCipher(block), iv), nil
					} else {
						return cipher.NewCFBDecrypter(ikipher.WrapCipher(block), iv), nil
					}
				}
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.Stream).XORKeyStream(dst, src)
			},
			false,
		)
	}
}
