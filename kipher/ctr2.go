package kipher

import (
	"bytes"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

const streamBufferBlockSize = 8 * 4 // * BlockSize = 256 bytes

type ctr2 struct {
	b            kryptoBlock
	ctr          []byte
	out          []byte
	outUsed      int
	bufferBlocks int
}

func newCTR2(block kryptoBlock, iv []byte, bufferBlocks int) Stream {
	blockSize := block.BlockSize()
	ctr := &ctr2{
		b:            block,
		ctr:          bytes.Clone(iv),
		out:          make([]byte, blockSize*bufferBlocks),
		bufferBlocks: bufferBlocks,
	}
	ctr.refill()

	return ctr
}

func (x *ctr2) fillCtr(outIdx int) {
	copy(x.out[outIdx:], x.ctr)

	for i := x.b.BlockSize() - 1; i >= 0; i-- {
		c := x.ctr[i]
		c++
		x.ctr[i] = c
		if c > 0 {
			return
		}
	}
}

func (x *ctr2) refill() {
	blockSize := x.b.BlockSize()

	for i := 0; i < x.bufferBlocks; i++ {
		x.fillCtr(blockSize * i)
	}

	var (
		bs8 = blockSize * 8
		bs4 = blockSize * 4
		bs1 = blockSize * 1
	)

	out := x.out
	for len(out) >= bs8 {
		x.b.Encrypt8(out[:bs8], out[:bs8])
		out = out[bs8:]
	}
	for len(out) >= bs4 {
		x.b.Encrypt4(out[:bs4], out[:bs4])
		out = out[bs4:]
	}
	for len(out) > 0 {
		x.b.Encrypt(out[:bs1], out[:bs1])
		out = out[bs1:]
	}

	x.outUsed = 0
}

func (x *ctr2) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("krypto/lea: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/lea: invalid buffer overlap")
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

func (x *ctr2) IV() []byte { return bytes.Clone(x.out) }
