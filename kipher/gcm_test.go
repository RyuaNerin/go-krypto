package kipher_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"math/rand"
	"testing"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
)

func TestGCM(t *testing.T) {
	const iter = 64 * 1024
	const maxLen = 16 * ikipher.GCMBlockSize

	test := func(nonceSize, tagSize int) func(t *testing.T) {
		return func(t *testing.T) {
			key := make([]byte, 16)

			nonce := make([]byte, maxLen)
			input := make([]byte, maxLen)
			additional := make([]byte, maxLen)

			sealed := make([]byte, maxLen+ikipher.GCMBlockSize)
			opened := make([]byte, maxLen)

			for i := 0; i < iter; i++ {
				inputSize := 1 + rand.Intn(maxLen-1)
				additionalSize := rand.Intn(maxLen)

				rnd.Read(key)
				rnd.Read(nonce[:nonceSize])
				rnd.Read(input[:inputSize])
				rnd.Read(additional[:additionalSize])

				b, _ := aes.NewCipher(key)

				gcm, err := kipher.NewGCMWithSize(ikipher.WrapCipher(b), nonceSize, tagSize)
				if err != nil {
					t.Error(err)
					return
				}

				sealed = gcm.Seal(sealed[:0], nonce[:nonceSize], input[:inputSize], additional[:additionalSize])
				opened, err = gcm.Open(opened[:0], nonce[:nonceSize], sealed, additional[:additionalSize])
				if err != nil {
					t.Error(err)
					return
				}
				if !bytes.Equal(opened, input[:inputSize]) {
					t.Errorf("failed to Open\nexpect: %v\nactual: %v", hex.EncodeToString(input[:inputSize]), hex.EncodeToString(opened))
					return
				}
			}
		}
	}

	t.Run("Default", test(ikipher.GCMStandardNonceSize, ikipher.GCMTagSize))
	t.Run("Nonce=14", test(14, ikipher.GCMTagSize))
	t.Run("Tag=14", test(ikipher.GCMStandardNonceSize, 14))
	t.Run("Nonce=14/Tag=14", test(14, 14))
}

func TestGCMWithStd(t *testing.T) {
	const maxLen = blocks * ikipher.GCMBlockSize

	test := func(
		nonceSize int,
		newCipher func(cipher.Block) (cipher.AEAD, error),
		newKipher func(ikipher.Block) (cipher.AEAD, error),
	) func(t *testing.T) {
		return func(t *testing.T) {
			key := make([]byte, keySize)

			nonce := make([]byte, nonceSize)
			input := make([]byte, maxLen)
			additional := make([]byte, maxLen)

			dstCipher := make([]byte, maxLen+ikipher.GCMBlockSize)
			dstKipher := make([]byte, maxLen+ikipher.GCMBlockSize)

			for i := 0; i < iter; i++ {
				inputSize := 1 + rand.Intn(maxLen-1)
				additionalSize := rand.Intn(maxLen)

				rnd.Read(key)
				rnd.Read(nonce)
				rnd.Read(input[:inputSize])
				rnd.Read(additional[:additionalSize])

				b, _ := aes.NewCipher(key)

				gcmCipher, err := newCipher(ikipher.WrapCipher(b))
				if err != nil {
					t.Error(err)
					return
				}
				gcmKipher, err := newKipher(ikipher.WrapKipher(b))
				if err != nil {
					t.Error(err)
					return
				}

				dstCipher = gcmCipher.Seal(dstCipher[:0], nonce, input[:inputSize], additional[:additionalSize])
				dstKipher = gcmKipher.Seal(dstKipher[:0], nonce, input[:inputSize], additional[:additionalSize])

				if !bytes.Equal(dstKipher, dstCipher) {
					t.Errorf("failed to Seal\nexpect: %s\nactual: %s", hex.EncodeToString(dstCipher), hex.EncodeToString(dstKipher))
					return
				}
			}
		}
	}

	t.Run("Default", test(
		ikipher.GCMStandardNonceSize,
		cipher.NewGCM,
		func(b ikipher.Block) (cipher.AEAD, error) { return kipher.NewGCM(b) }),
	)
	t.Run("Nonce=14", test(
		14,
		func(b cipher.Block) (cipher.AEAD, error) { return cipher.NewGCMWithNonceSize(b, 14) },
		func(b ikipher.Block) (cipher.AEAD, error) { return kipher.NewGCMWithNonceSize(b, 14) },
	))
	t.Run("Tag=14", test(
		ikipher.GCMStandardNonceSize,
		func(b cipher.Block) (cipher.AEAD, error) { return cipher.NewGCMWithTagSize(b, 14) },
		func(b ikipher.Block) (cipher.AEAD, error) { return kipher.NewGCMWithTagSize(b, 14) },
	))
}
