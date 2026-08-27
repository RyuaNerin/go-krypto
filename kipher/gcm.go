package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
)

const (
	GCMBlockSize         = 16
	GCMTagSize           = 16
	GCMMinimumTagSize    = 2
	GCMStandardNonceSize = kipherimpl.GCMStandardNonceSize
)

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
// with the standard nonce length.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared. This applies to all NewGCM* variants.
func NewGCM(b cipher.Block) (cipher.AEAD, error) {
	return NewGCMWithSize(b, GCMStandardNonceSize, GCMTagSize)
}

// NewGCMWithNonceSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode, which accepts nonces of the given length. The length must not
// be zero.
//
// Only use this function if you require compatibility with an existing
// cryptosystem that uses non-standard nonce lengths. All other users should use
// NewGCM, which is faster and more resistant to misuse.
//
// if nonceSize = 0, nonceSize = [GCMStandardNonceSize]
func NewGCMWithNonceSize(b cipher.Block, size int) (cipher.AEAD, error) {
	return NewGCMWithSize(b, size, GCMTagSize)
}

// NewGCMWithTagSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode, which generates tags with the given length.
//
// Tag lengths must be between 2 and 16 bytes, inclusive. Tags shorter than 12
// bytes weaken authentication; NIST SP 800-38D does not allow tags shorter than
// 4 bytes, and allows 4- and 8-byte tags only with restrictions on their use.
//
// Only use this function if you require compatibility with an existing
// cryptosystem that uses non-standard tag lengths. All other users should use
// NewGCM, which is more resistant to misuse.
//
// if tagSize = 0, tagSize = [GCMTagSize]
func NewGCMWithTagSize(b cipher.Block, tagSize int) (cipher.AEAD, error) {
	return NewGCMWithSize(b, GCMStandardNonceSize, tagSize)
}

// NewGCMWithSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode with the given nonce and tag lengths. See NewGCMWithNonceSize and
// NewGCMWithTagSize for the limits.
//
// if nonceSize = 0, nonceSize = [GCMStandardNonceSize]
// if tagSize = 0, tagSize = [GCMTagSize]
func NewGCMWithSize(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if b.BlockSize() != GCMBlockSize {
		return nil, ErrRequire128Bits
	}

	if nonceSize < 0 {
		return nil, ErrInvalidNonceZero
	}
	if nonceSize == 0 {
		nonceSize = GCMStandardNonceSize
	}

	if tagSize == 0 {
		tagSize = GCMTagSize
	}
	if tagSize < GCMMinimumTagSize || tagSize > GCMBlockSize {
		return nil, ErrInvalidTagSizeGCM
	}

	return wrapAEAD(b, kipher.NewGCM(internal.ToBlock8(b), nonceSize, tagSize)), nil
}
