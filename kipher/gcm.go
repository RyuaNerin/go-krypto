package kipher

import (
	"crypto/cipher"
)

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
func NewGCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	return newGCMWithNonceAndTagSize(b, nonceSize, tagSize)
}
