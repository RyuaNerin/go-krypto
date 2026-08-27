// Package kipher implements standard block cipher modes that can be wrapped around low-level block cipher implementations.
//
// See NIST SP 800-38A, NIST SP 800-38C, NIST SP 800-38D
package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

type (
	blockMode struct{ wrap[krypto.BlockMode] }
	stream    struct{ wrap[krypto.Stream] }
	aead      struct{ wrap[krypto.AEAD] }

	wrap[T zeroize.Zeroizier] struct {
		z        zeroize.ZeroizedFlag
		cipher   cipher.Block
		fallback T
	}
)

var (
	_ krypto.BlockMode = (*blockMode)(nil)
	_ krypto.Stream    = (*stream)(nil)
	_ krypto.AEAD      = (*aead)(nil)
)

func wrapBlockMode(cipher cipher.Block, fallback krypto.BlockMode) krypto.BlockMode {
	return &blockMode{
		wrap[krypto.BlockMode]{
			cipher:   cipher,
			fallback: fallback,
		},
	}
}

func wrapStream(cipher cipher.Block, fallback krypto.Stream) krypto.Stream {
	return &stream{
		wrap[krypto.Stream]{
			cipher:   cipher,
			fallback: fallback,
		},
	}
}

func wrapAEAD(cipher cipher.Block, fallback krypto.AEAD) krypto.AEAD {
	return &aead{
		wrap[krypto.AEAD]{
			cipher:   cipher,
			fallback: fallback,
		},
	}
}

func (b *wrap[T]) Zeroize() {
	var zero T

	zeroize.Zeroize(b.cipher, b.fallback)
	b.cipher = nil
	b.fallback = zero
	b.z.SetInvalid()
}

func (b *blockMode) BlockSize() int { return b.fallback.BlockSize() }

func (b *blockMode) CryptBlocks(dst, src []byte) {
	b.z.ThrowIfZeroized()
	b.fallback.CryptBlocks(dst, src)
}

func (s *stream) XORKeyStream(dst, src []byte) {
	s.z.ThrowIfZeroized()
	s.fallback.XORKeyStream(dst, src)
}

func (a *aead) NonceSize() int { return a.fallback.NonceSize() }
func (a *aead) Overhead() int  { return a.fallback.Overhead() }

func (a *aead) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	a.z.ThrowIfZeroized()
	return a.fallback.Seal(dst, nonce, plaintext, additionalData)
}

func (a *aead) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	a.z.ThrowIfZeroized()
	return a.fallback.Open(dst, nonce, ciphertext, additionalData)
}
