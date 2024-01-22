// Package seed implements SEED encryption, as defined in TTAS.KO-12.0004/R1
package seed

import (
	"crypto/cipher"
	"fmt"
)

const (
	// The SEED block size in bytes.
	BlockSize = 16
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/seed: invalid key size %d", int(k))
}

// NewCipher creates and returns a new cipher.Block. The key argument should be the SEED key, either 16 or 32 bytes to select SEED-128 or SEED-256.
func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	switch l {
	case 16:
		return new128(key), nil
	}

	return nil, KeySizeError(l)
}
