package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
)

// NewCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block cipher.Block, iv []byte) cipher.Stream {
	kb, ok := block.(ikipher.Block)
	if !ok {
		if ctr, ok := block.(ikipher.CTRAble); ok {
			return ctr.NewCTR(iv)
		}
	}
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}

	if ok {
		return newCTR(kb, iv)
	}
	return cipher.NewCTR(block, iv)
}

func newCTR(kb ikipher.Block, iv []byte) cipher.Stream {
	ctr := new(ctr)
	ctr.ctr.Init(kb, iv, 0)

	return ctr
}

type ctr struct {
	ctr ikipher.CTR
}

func (x *ctr) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgBufferOverlap)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgSmallDst)
	}

	x.ctr.Xor(dst, src)
}
