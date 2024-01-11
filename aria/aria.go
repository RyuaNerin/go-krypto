// Package aria implements ARIA encryption, as defined in KS X 1213-1
package aria

import (
	"crypto/cipher"
	"fmt"
)

const (
	// The ARIA block size in bytes.
	BlockSize = 16
)

var (
	newCipher = newCipherGo
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/aria: invalid key size %d", int(k))
}

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
