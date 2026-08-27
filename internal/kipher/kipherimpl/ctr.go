package kipherimpl

import (
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
)

const bufferedBlocks = 8

type CTR struct {
	Out         []byte
	cipherBlock []byte

	b         krypto.Block8
	blockSize int

	ctr       []byte
	ctrOffset int
}

func (x *CTR) Zeroize() {
	clear(x.Out)
	x.Out = nil
	clear(x.cipherBlock)
	x.cipherBlock = nil

	x.b = nil
	x.blockSize = 0

	clear(x.ctr)
	x.ctr = nil
	x.ctrOffset = 0
}

func (x *CTR) Init(b krypto.Block8, iv []byte, ctrSize int) {
	x.b = b
	if ctrSize > 0 {
		x.ctrOffset = len(iv) - ctrSize
	}

	bs := b.BlockSize()

	buf := make([]byte, bs+bufferedBlocks*bs)
	x.ctr = buf[:bs]
	x.cipherBlock = buf[bs:]
	copy(x.ctr, iv)

	x.blockSize = bs
}

func (x *CTR) CopyCTR(b []byte) {
	if len(b) != len(x.ctr) {
		panic(msgUnexpectedCTRLength)
	}

	copy(b, x.ctr)
}

func (x *CTR) Refill(blocks int) {
	if len(x.Out) != 0 {
		return
	}

	outIdx := (bufferedBlocks - blocks) * x.blockSize

	for idx := outIdx; idx < len(x.cipherBlock); idx += x.blockSize {
		copy(x.cipherBlock[idx:], x.ctr)
		internal.IncCtr(x.ctr[x.ctrOffset:])
	}

	switch {
	case blocks > 4:
		x.b.Encrypt8(x.cipherBlock, x.cipherBlock)
	case blocks > 1:
		o := x.cipherBlock[4*x.blockSize:]
		x.b.Encrypt4(o, o)
	default:
		o := x.cipherBlock[7*x.blockSize:]
		x.b.Encrypt(o, o)
	}

	x.Out = x.cipherBlock[outIdx:]
}

func (x *CTR) Xor(out, in []byte) {
	for len(in) > 0 {
		if len(x.Out) == 0 {
			remainBlocks := min(internal.CeilDiv(len(in), x.blockSize), 8)

			x.Refill(remainBlocks)
		}

		n := subtle.XORBytes(out, in, x.Out)
		out = out[n:]
		in = in[n:]
		x.Out = x.Out[n:]
	}
}
