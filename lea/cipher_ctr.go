//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

type ctrAble interface {
	NewCTR(iv []byte) cipher.Stream
}

// Assert that lea_key implements the ctrAble interfaces.
var _ ctrAble = (*leaContext)(nil)

const streamBufferBlockSize = 8 * 4 // * BlockSize = 256 bytes

type ctrContext struct {
	ctx    *leaContext
	ctr    [BlockSize]byte
	out    [BlockSize * streamBufferBlockSize]byte
	outPos int
}

func (leaCtx *leaContext) NewCTR(iv []byte) cipher.Stream {
	ctr := &ctrContext{
		ctx: leaCtx,
	}
	copy(ctr.ctr[:], iv)
	ctr.refill()

	return ctr
}

func (ctr *ctrContext) fillCtr(outIdx int) {
	copy(ctr.out[:][outIdx:], ctr.ctr[:])

	for i := BlockSize - 1; i >= 0; i-- {
		c := ctr.ctr[i]
		c++
		ctr.ctr[i] = c
		if c > 0 {
			return
		}
	}
}

func (ctr *ctrContext) refill() {
	for i := 0; i < streamBufferBlockSize; i++ {
		ctr.fillCtr(BlockSize * i)
	}

	for i := 0; i < streamBufferBlockSize/8; i++ {
		out := ctr.out[0x80*i:]
		leaEnc8(ctr.ctx, out, out)
	}

	ctr.outPos = 0
}

func (ctr *ctrContext) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("krypto/lea: output smaller than input")
	}
	if subtle.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/lea: invalid buffer overlap")
	}

	for len(src) > 0 {
		if len(ctr.out) == ctr.outPos {
			ctr.refill()
		}

		n := xorBytes(dst, src, ctr.out[ctr.outPos:])
		ctr.outPos += n
		dst = dst[n:]
		src = src[n:]
	}
}
