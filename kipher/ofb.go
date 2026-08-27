package kipher

import (
	"crypto/cipher"
	"fmt"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

// NewOFB returns a [Stream] that encrypts or decrypts using the block cipher b
// in output feedback mode. The initialization vector iv's length must be equal
// to b's block size.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewOFB(b cipher.Block, iv []byte) cipher.Stream {
	blockSize := b.BlockSize()
	if blockSize > kipher.MaxOFBBlockSize {
		panic(fmt.Sprintf(msgIVIsTooLargeFmt, kipher.MaxOFBBlockSize))
	}
	if len(iv) != blockSize {
		panic(msgInvalidIVLength)
	}

	return wrapStream(b, kipher.NewOFB(internal.ToBlock8(b), iv))
}
