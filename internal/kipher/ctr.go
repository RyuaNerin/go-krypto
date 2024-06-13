package kipher

import (
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

const bufferedBlocks = 8

type CTR struct {
	Out []byte
	out []byte

	b         Block
	blockSize int

	ctr       []byte
	ctrOffset int
}

func (x *CTR) Init(b Block, iv []byte, ctrSize int) {
	x.b = b
	if ctrSize > 0 {
		x.ctrOffset = len(iv) - ctrSize
	}

	bs := b.BlockSize()

	buf := make([]byte, bs+bufferedBlocks*bs)
	x.ctr = buf[:bs]
	x.out = buf[bs:]
	copy(x.ctr, iv)

	x.blockSize = bs
}

func (x *CTR) CopyCTR(b []byte) {
	if len(b) != len(x.ctr) {
		panic("internal: unexpected ctr length")
	}

	copy(b, x.ctr)
}

func (x *CTR) Refill(blocks int) {
	if len(x.Out) != 0 {
		return
	}

	outIdx := (bufferedBlocks - blocks) * x.blockSize

	for idx := outIdx; idx < len(x.out); idx += x.blockSize {
		copy(x.out[idx:], x.ctr)
		internal.IncCtr(x.ctr[x.ctrOffset:])
	}

	switch {
	case blocks > 4:
		x.b.Encrypt8(x.out, x.out)
	case blocks > 1:
		o := x.out[4*x.blockSize:]
		x.b.Encrypt4(o, o)
	default:
		o := x.out[7*x.blockSize:]
		x.b.Encrypt(o, o)
	}

	x.Out = x.out[outIdx:]
}

func (x *CTR) Xor(out, in []byte) {
	for len(in) > 0 {
		if len(x.Out) == 0 {
			remainBlocks := internal.CeilDiv(len(in), x.blockSize)
			if remainBlocks > 8 {
				remainBlocks = 8
			}

			x.Refill(remainBlocks)
		}

		n := subtle.XORBytes(out, in, x.Out)
		out = out[n:]
		in = in[n:]
		x.Out = x.Out[n:]
	}
}
