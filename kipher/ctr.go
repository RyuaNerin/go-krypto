package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

// NewCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
//
// Zeroize() on the returned value also zeroizes the block cipher block you
// passed in (in place), so do not reuse block afterward. Further use of the mode
// panics with krypto.ErrStateCleared.
func NewCTR(block cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}

	return wrapStream(block, kipher.NewCTR(internal.ToBlock8(block), iv))
}
