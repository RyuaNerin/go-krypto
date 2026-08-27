package internal_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	"github.com/stretchr/testify/require"
)

func TestWrapAndUnwrap(t *testing.T) {
	t.Run("Hash", func(t *testing.T) {
		h := sha256.New()
		require.Equal(t, h, internal.UnwrapHash(internal.WrapHash(h)))
	})
	t.Run("BlockMode", func(t *testing.T) {
		c, _ := aes.NewCipher(make([]byte, 16))
		bm := cipher.NewCBCEncrypter(c, make([]byte, 16))
		require.Same(t, bm, internal.UnwrapBlockMode(internal.WrapBlockMode(bm)))
	})
	t.Run("Stream", func(t *testing.T) {
		c, _ := aes.NewCipher(make([]byte, 16))
		s := cipher.NewCTR(c, make([]byte, 16))
		require.Equal(t, s, internal.UnwrapStream(internal.WrapStream(s)))
	})
	t.Run("AEAD", func(t *testing.T) {
		c, _ := aes.NewCipher(make([]byte, 16))
		aead, _ := cipher.NewGCM(c)
		require.Same(t, aead, internal.UnwrapAEAD(internal.WrapAEAD(aead)))
	})
	t.Run("Block", func(t *testing.T) {
		c, _ := aes.NewCipher(make([]byte, 16))
		require.Equal(t, c, internal.UnwrapBlock(internal.JustBlock(c)))
		require.Equal(t, c, internal.UnwrapBlock(internal.JustBlock4(c)))
		require.Equal(t, c, internal.UnwrapBlock(internal.ToBlock8(c)))
	})
}
