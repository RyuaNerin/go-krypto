// Package aria implements ARIA encryption, as defined in KS X 1213-1
package aria

import (
	"crypto/cipher"
)

const (
	// The ARIA block size in bytes.
	BlockSize = 16
)

const (
	rkSize = 16 * 17
)

type ariaContext struct {
	ek     [rkSize]byte
	dk     [rkSize]byte
	rounds int
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	switch l {
	case 16:
	case 24:
	case 32:
	default:
		return nil, KeySizeError(l)
	}

	return newCipher(key)
}
