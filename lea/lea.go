package lea

import (
	"crypto/cipher"
	"fmt"
)

var (
	useASM = false
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/lea: invalid key size %d", int(k))
}

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipher(key []byte) (cipher.Block, error) {
	return leaNew(key)
}

// NewCipherECB creates and returns a new cipher.Block by ECB mode.
// This function can be useful in amd64.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipherECB(key []byte) (cipher.Block, error) {
	return leaNewECB(key)
}
