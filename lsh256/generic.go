//go:build !(arm64 || amd64) || purego || (gccgo && !go1.18)
// +build !arm64,!amd64 purego gccgo,!go1.18

package lsh256

import (
	"encoding/binary"
	"hash"
	"math/bits"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

func newContextGo(size int) hash.Hash {
	ctx := &lsh256ContextGo{
		outlenbytes: size,
	}
	ctx.Reset()

	return ctx
}

func sumGo(size int, data []byte) [Size]byte {
	ctx := lsh256ContextGo{
		outlenbytes: size,
	}
	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

const (
	numStep = 26

	alphaEven = 29
	alphaOdd  = 5

	betaEven = 1
	betaOdd  = 17
)

var gamma = [...]int{0, 8, 16, 24, 24, 16, 8, 0}

type lsh256ContextGo struct {
	cv    [16]uint32
	tcv   [16]uint32
	msg   [16 * (numStep + 1)]uint32
	block [BlockSize]byte

	boff        int
	outlenbytes int
}

func (b *lsh256ContextGo) Size() int {
	return b.outlenbytes
}

func (b *lsh256ContextGo) BlockSize() int {
	return BlockSize
}

func (b *lsh256ContextGo) Reset() {
	memory.MemclrU32(b.tcv[:])
	memory.MemclrU32(b.msg[:])
	memory.Memclr(b.block[:])

	b.boff = 0
	switch b.outlenbytes {
	case Size:
		copy(b.cv[:], iv256)
	case Size224:
		copy(b.cv[:], iv224)
	}
}

func (b *lsh256ContextGo) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	plen := len(p)

	gap := BlockSize - b.boff
	if b.boff > 0 && len(p) >= gap {
		copy(b.block[b.boff:], p[:gap])
		b.compress(b.block[:])
		b.boff = 0

		p = p[gap:]
	}

	for len(p) >= BlockSize {
		b.compress(p)
		b.boff = 0
		p = p[BlockSize:]
	}

	if len(p) > 0 {
		copy(b.block[b.boff:], p)
		b.boff += len(p)
	}

	return plen, nil
}

func (b *lsh256ContextGo) Sum(p []byte) []byte {
	b0 := *b
	hash := b0.checkSum()
	return append(p, hash[:b.Size()]...)
}

func (b *lsh256ContextGo) checkSum() [Size]byte {
	b.block[b.boff] = 0x80

	memory.Memclr(b.block[b.boff+1:])
	b.compress(b.block[:])

	var temp [8]uint32
	for i := 0; i < 8; i++ {
		temp[i] = b.cv[i] ^ b.cv[i+8]
	}

	var digest [Size]byte
	for i := 0; i < b.outlenbytes; i++ {
		digest[i] = byte(temp[i>>2] >> ((i << 3) & 0x1f))
	}

	return digest
}

func (b *lsh256ContextGo) compress(data []byte) {
	b.msgExpansion(data)

	for i := 0; i < numStep/2; i++ {
		b.step(2*i+0, alphaEven, betaEven)
		b.step(2*i+1, alphaOdd, betaOdd)
	}

	// b.msg add
	for i := 0; i < 16; i++ {
		b.cv[i] ^= b.msg[16*numStep+i]
	}
}

func (b *lsh256ContextGo) msgExpansion(in []byte) {
	for i := 0; i < 32; i++ {
		b.msg[i] = binary.LittleEndian.Uint32(in[i*4:])
	}

	for i := 2; i <= numStep; i++ {
		idx := 16 * i
		b.msg[idx] = b.msg[idx-16] + b.msg[idx-29]
		b.msg[idx+1] = b.msg[idx-15] + b.msg[idx-30]
		b.msg[idx+2] = b.msg[idx-14] + b.msg[idx-32]
		b.msg[idx+3] = b.msg[idx-13] + b.msg[idx-31]
		b.msg[idx+4] = b.msg[idx-12] + b.msg[idx-25]
		b.msg[idx+5] = b.msg[idx-11] + b.msg[idx-28]
		b.msg[idx+6] = b.msg[idx-10] + b.msg[idx-27]
		b.msg[idx+7] = b.msg[idx-9] + b.msg[idx-26]
		b.msg[idx+8] = b.msg[idx-8] + b.msg[idx-21]
		b.msg[idx+9] = b.msg[idx-7] + b.msg[idx-22]
		b.msg[idx+10] = b.msg[idx-6] + b.msg[idx-24]
		b.msg[idx+11] = b.msg[idx-5] + b.msg[idx-23]
		b.msg[idx+12] = b.msg[idx-4] + b.msg[idx-17]
		b.msg[idx+13] = b.msg[idx-3] + b.msg[idx-20]
		b.msg[idx+14] = b.msg[idx-2] + b.msg[idx-19]
		b.msg[idx+15] = b.msg[idx-1] + b.msg[idx-18]
	}
}

func (b *lsh256ContextGo) step(stepidx, alpha, beta int) {
	var vl, vr uint32

	for colidx := 0; colidx < 8; colidx++ {
		vl = b.cv[colidx] ^ b.msg[16*stepidx+colidx]
		vr = b.cv[colidx+8] ^ b.msg[16*stepidx+colidx+8]
		vl = bits.RotateLeft32(vl+vr, alpha) ^ step[8*stepidx+colidx]
		vr = bits.RotateLeft32(vl+vr, beta)
		b.tcv[colidx] = vr + vl
		b.tcv[colidx+8] = bits.RotateLeft32(vr, gamma[colidx])
	}

	// wordPermutation
	b.cv[0] = b.tcv[6]
	b.cv[1] = b.tcv[4]
	b.cv[2] = b.tcv[5]
	b.cv[3] = b.tcv[7]
	b.cv[4] = b.tcv[12]
	b.cv[5] = b.tcv[15]
	b.cv[6] = b.tcv[14]
	b.cv[7] = b.tcv[13]
	b.cv[8] = b.tcv[2]
	b.cv[9] = b.tcv[0]
	b.cv[10] = b.tcv[1]
	b.cv[11] = b.tcv[3]
	b.cv[12] = b.tcv[8]
	b.cv[13] = b.tcv[11]
	b.cv[14] = b.tcv[10]
	b.cv[15] = b.tcv[9]
}
