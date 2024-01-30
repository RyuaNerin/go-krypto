package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/krypto/kipher/ctr.go

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

type ctr struct {
	b       cipher.Block
	ctr     []byte
	out     []byte
	outUsed int
}

// NewCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The length of iv must be the same as the Block's block size.
func NewCTR(block cipher.Block, iv []byte, bufferBlocks int) Stream {
	if len(iv) != block.BlockSize() {
		panic("cipher.NewCTR: IV length must equal block size")
	}
	if kb, ok := block.(kryptoBlock); ok {
		return newCTR2(kb, iv, bufferBlocks)
	}
	return newCTR(block, iv, bufferBlocks)
}

func newCTR(block cipher.Block, iv []byte, bufferBlocks int) Stream {
	bufSize := block.BlockSize() * bufferBlocks
	return &ctr{
		b:       block,
		ctr:     bytes.Clone(iv),
		out:     make([]byte, 0, bufSize),
		outUsed: 0,
	}
}

func (x *ctr) refill() {
	remain := len(x.out) - x.outUsed
	copy(x.out, x.out[x.outUsed:])
	x.out = x.out[:cap(x.out)]
	bs := x.b.BlockSize()
	for remain <= len(x.out)-bs {
		x.b.Encrypt(x.out[remain:], x.ctr)
		remain += bs

		// Increment counter
		for i := len(x.ctr) - 1; i >= 0; i-- {
			x.ctr[i]++
			if x.ctr[i] != 0 {
				break
			}
		}
	}
	x.out = x.out[:remain]
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
		if x.outUsed >= len(x.out)-x.b.BlockSize() {
			x.refill()
		}
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

func (x *ctr) IV() []byte { return bytes.Clone(x.out) }
