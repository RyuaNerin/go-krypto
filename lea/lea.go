// Package lea implements LEA encryption, as defined in TTAK.KO-12.0223
package lea

import (
	"crypto/cipher"
)

const (
	// The LEA block size in bytes.
	BlockSize = 16
)

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
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
