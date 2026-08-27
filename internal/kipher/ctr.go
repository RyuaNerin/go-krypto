package kipher

import (
	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
)

func NewCTR(kb krypto.Block8, iv []byte) *CTR {
	ctr := new(CTR)
	ctr.ctr.Init(kb, iv, 0)

	return ctr
}

type CTR struct {
	ctr kipherimpl.CTR
}

var _ krypto.Stream = (*CTR)(nil)

func (x *CTR) Zeroize() {
	x.ctr.Zeroize()
}

func (x *CTR) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}

	x.ctr.Xor(dst, src)
}
