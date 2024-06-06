package kipher

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

// NewCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block cipher.Block, iv []byte) cipher.Stream {
	kb, ok := block.(internal.Block)
	if !ok {
		if ctr, ok := block.(internal.CTRAble); ok {
			return ctr.NewCTR(iv)
		}
	}
	if len(iv) != block.BlockSize() {
		panic("krypto/kipher.NewCTR: IV length must equal block size")
	}

	if ok {
		ctr := &ctr{
			b:   kb,
			ctr: bytes.Clone(iv),
			out: make([]byte, 8*kb.BlockSize()),
		}
		ctr.refill()

		return ctr
	}
	return cipher.NewCTR(block, iv)
}

type ctr struct {
	b       internal.Block
	ctr     []byte
	out     []byte
	outUsed int
}

func (x *ctr) refill() {
	blockSize := x.b.BlockSize()

	for i := 0; i < 8; i++ {
		copy(x.out[blockSize*i:], x.ctr)
		internal.IncCtr(x.ctr)
	}

	x.b.Encrypt8(x.out, x.out)
	x.outUsed = 0
}

func (x *ctr) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}

	for len(src) > 0 {
		if len(x.out) == x.outUsed {
			x.refill()
		}
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}
