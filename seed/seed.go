// Package seed implements SEED encryption, as defined in TTAS.KO-12.0004/R1
package seed

import (
	"crypto/cipher"
)

const (
	// The SEED block size in bytes.
	BlockSize = 16
)

// NewCipher creates and returns a new cipher.Block. The key argument should be the SEED key, either 16 or 32 bytes to select SEED-128 or SEED-256.
func NewCipher(key []byte) (cipher.Block, error) {
	if l := len(key); l != 16 {
		return nil, KeySizeError(l)
	}

	return new128(key), nil
}
