package kipher

import (
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

const (
	MaxOFBBlockSize = 16
	ofbBatch        = 8
)

func NewOFB(b krypto.Block, iv []byte) *OFB {
	x := &OFB{
		cipher:  b,
		outUsed: 0,
	}

	copy(x.cipherBlock[:], iv)
	return x
}

type OFB struct {
	cipher      krypto.Block
	cipherBlock [MaxOFBBlockSize]byte
	outBlock    [MaxOFBBlockSize]byte
	outLen      int
	outUsed     int
}

var _ krypto.Stream = (*OFB)(nil)

func (x *OFB) Zeroize() {
	x.cipher.Zeroize()
	x.cipher = nil
	clear(x.cipherBlock[:])
	clear(x.outBlock[:])
	x.outLen = 0
	x.outUsed = 0
}

func (x *OFB) refill() {
	bs := x.cipher.BlockSize()
	x.cipher.Encrypt(x.cipherBlock[:bs], x.cipherBlock[:bs])
	copy(x.outBlock[:bs], x.cipherBlock[:bs])
	x.outLen = bs
	x.outUsed = 0
}

func (x *OFB) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	for len(src) > 0 {
		if x.outUsed == x.outLen {
			x.refill()
		}
		n := subtle.XORBytes(dst, src, x.outBlock[x.outUsed:x.outLen])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}
