package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

const CCMBlockSize = kipher.CCMBlockSize

// NewCCM returns the given 128-bit, block cipher wrapped in Counter with CBC-MAC (CCM) mode.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewCCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize <= 0 {
		return nil, ErrInvalidNonceZero
	}
	if nonceSize < 7 || nonceSize > 13 {
		return nil, ErrInvalidNonceSize
	}
	if tagSize < 4 || 16 < tagSize || tagSize%2 != 0 {
		return nil, ErrInvalidTagSizeCCM
	}

	if b.BlockSize() != kipher.CCMBlockSize {
		return nil, ErrRequire128Bits
	}

	return wrapAEAD(b, kipher.NewCCM(internal.ToBlock8(b), nonceSize, tagSize)), nil
}
