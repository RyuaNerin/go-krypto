package kipher

import (
	"errors"

	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

const (
	msgInvalidIVLength = "krypto/kipher: IV length must equal block size"
	msgInvalidCFBSize  = "krypto/kipher: CFB segment size must be between 1 byte and the block size"
	msgIVIsTooLargeFmt = "krypto/kipher: IV is too large, must be less than %d bytes"
)

var (
	ErrNotFullBlocks     = errors.New("krypto/kipher: src not full blocks")
	ErrSmallDst          = errors.New("krypto/kipher: dst smaller than src")
	ErrBufferOverlap     = errors.New("krypto/kipher: invalid buffer overlap")
	ErrInvalidNonceZero  = errors.New("krypto/kipher: the nonce can't have zero length, or the security of the key will be immediately compromised")
	ErrInvalidNonceSize  = errors.New("krypto/kipher: invalid nonce size")
	ErrInvalidTagSizeGCM = errors.New("krypto/kipher: incorrect tag size given")
	ErrInvalidTagSizeCCM = errors.New("krypto/kipher: tagSize must be 4, 6, 8, 10, 12, 14 or 16 in CCM")
	ErrInvalidNonce      = errors.New("krypto/kipher: incorrect nonce length given")
	ErrDataTooLarge      = errors.New("krypto/kipher: message too large")
	ErrRequire128Bits    = errors.New("krypto/kipher: requires 128-bit block cipher")

	// ErrOpenFailed is returned by the Open method of the AEADs in this package
	// when the message fails to authenticate.
	ErrOpenFailed = kipher.ErrOpenFailed
)
