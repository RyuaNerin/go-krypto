package kipher

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

type ofb struct {
	b       cipher.Block
	cipher  []byte
	out     []byte
	outUsed int
}

// NewOFB returns a Stream that encrypts or decrypts using the block cipher b
// in output feedback mode. The initialization vector iv's length must be equal
// to b's block size.
func NewOFB(b cipher.Block, iv []byte, bufferBlocks int) Stream {
	blockSize := b.BlockSize()
	if len(iv) != blockSize {
		panic("cipher.NewOFB: IV length must equal block size")
	}
	bufSize := blockSize * bufferBlocks
	x := &ofb{
		b:       b,
		cipher:  make([]byte, blockSize),
		out:     make([]byte, 0, bufSize),
		outUsed: 0,
	}

	copy(x.cipher, iv)
	return x
}

func (x *ofb) refill() {
	bs := x.b.BlockSize()
	remain := len(x.out) - x.outUsed
	if remain > x.outUsed {
		return
	}
	copy(x.out, x.out[x.outUsed:])
	x.out = x.out[:cap(x.out)]
	for remain < len(x.out)-bs {
		x.b.Encrypt(x.cipher, x.cipher)
		copy(x.out[remain:], x.cipher)
		remain += bs
	}
	x.out = x.out[:remain]
	x.outUsed = 0
}

func (x *ofb) XORKeyStream(dst, src []byte) {
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

func (x *ofb) IV() []byte { return bytes.Clone(x.out) }
