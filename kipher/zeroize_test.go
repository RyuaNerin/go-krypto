package kipher_test

import (
	"bytes"
	"crypto/cipher"
	"errors"
	"testing"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

func TestZeroizeRepeatedAndBlockKey(t *testing.T) {
	for name, newMode := range map[string]func(cipher.Block) any{
		"CTR":        func(b cipher.Block) any { return kipher.NewCTR(b, make([]byte, 16)) },
		"OFB":        func(b cipher.Block) any { return kipher.NewOFB(b, make([]byte, 16)) },
		"CBCEncrypt": func(b cipher.Block) any { return kipher.NewCBCEncrypter(b, make([]byte, 16)) },
		"CBCDecrypt": func(b cipher.Block) any { return kipher.NewCBCDecrypter(b, make([]byte, 16)) },
	} {
		t.Run(name, func(t *testing.T) {
			b, _ := lea.NewCipher(bytes.Repeat([]byte{0x42}, 16))
			m := newMode(b)
			switch m := m.(type) {
			case cipher.Stream:
				m.XORKeyStream(make([]byte, 17), make([]byte, 17))
			case cipher.BlockMode:
				m.CryptBlocks(make([]byte, 32), make([]byte, 32))
			}
			z := m.(krypto.Zeroizable)
			z.Zeroize()
			z.Zeroize()

			defer func() {
				e := recover()
				if e == nil {
					return
				}
				if errors.Is(e.(error), zeroize.ErrStateCleared) {
					t.Fatal("operation did not reject zeroized state")
				}
			}()
			b.Encrypt(make([]byte, 16), make([]byte, 16))
		})
	}
}

func TestAEADTemporaryCTRPreservesKey(t *testing.T) {
	for name, newAEAD := range map[string]func(cipher.Block) (cipher.AEAD, error){
		"CCM": func(b cipher.Block) (cipher.AEAD, error) { return kipher.NewCCM(b, 12, 16) },
		"GCM": kipher.NewGCM,
	} {
		t.Run(name, func(t *testing.T) {
			b, _ := lea.NewCipher(bytes.Repeat([]byte{0x42}, 16))
			a, err := newAEAD(b)
			require.NoError(t, err)

			pt := bytes.Repeat([]byte("secret"), 33)
			for i := range 3 {
				nonce := bytes.Repeat([]byte{byte(i)}, a.NonceSize())
				ct := a.Seal(nil, nonce, pt, []byte("aad"))

				got, err := a.Open(nil, nonce, ct, []byte("aad"))
				require.NoErrorf(t, err, "round %d: decryption failed", i)
				require.Equal(t, pt, got, "round %d: decrypted plaintext does not match original", i)

				ct[len(ct)-1] ^= 1

				_, err = a.Open(nil, nonce, ct, []byte("aad"))
				require.Error(t, err, "tampered ciphertext accepted")
			}
			b.Encrypt(make([]byte, 16), make([]byte, 16))
			a.(krypto.Zeroizable).Zeroize()
		})
	}
}
