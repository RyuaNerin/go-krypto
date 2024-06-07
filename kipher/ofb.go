package kipher

import "crypto/cipher"

// NewOFB returns a [Stream] that encrypts or decrypts using the block cipher b
// in output feedback mode. The initialization vector iv's length must be equal
// to b's block size.
func NewOFB(b cipher.Block, iv []byte) cipher.Stream {
	blockSize := b.BlockSize()
	if len(iv) != blockSize {
		panic(msgInvalidIVLength)
	}
	return cipher.NewOFB(b, iv)
}
