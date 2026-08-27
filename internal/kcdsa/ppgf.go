package kcdsa

import (
	"fmt"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func ppgf(
	dst []byte,
	nBits int, h hash.Hash, src ...[]byte,
) []byte {
	return newPPGF(h).Generate(dst, nBits, src...)
}

type ppgfCtx struct {
	h           hash.Hash
	initVectors []byte // either initVectors or initSources

	mh          internal.Hash
	initSources [][]byte // either initVectors or initSources

	lh int
}

var _ krypto.Zeroizable = (*ppgfCtx)(nil)

func (p *ppgfCtx) Zeroize() {
	zeroize.Zeroize(p.h)
	zeroize.Zeroize(p.initVectors)
	zeroize.Zeroize(p.mh)
	// ignore initSources
	p.initSources = nil
	zeroize.Zeroize(p.lh)
}

func newPPGF(h hash.Hash, src ...[]byte) (ppgf ppgfCtx) {
	mh, ok := h.(internal.Hash)
	if !ok {
		panic(msgHashMustImplementBinaryMarshaling)
	}
	var err error

	h.Reset()
	for _, v := range src {
		_, _ = h.Write(v)
	}
	ppgf.mh = mh
	ppgf.initVectors, err = mh.MarshalBinary()
	if err != nil {
		panic(fmt.Sprintf("kcdsa: failed to marshal hash: %v", err))
	}

	ppgf.lh = h.Size()
	return
}

func (p ppgfCtx) Generate(dst []byte, nBits int, src ...[]byte) []byte {
	// p.12
	// from java
	i := internal.BitsToBytes(nBits)
	dst = internal.Grow(dst, i)

	var iBuf [1]byte
	hbuf := make([]byte, 0, p.lh)
	defer zeroize.Zeroize(iBuf[:], hbuf)

	var saved []byte
	if len(src) > 0 {
		err := p.mh.UnmarshalBinary(p.initVectors)
		if err != nil {
			panic(fmt.Sprintf("kcdsa: failed to unmarshal hash: %v", err))
		}
		for _, v := range src {
			_, _ = p.mh.Write(v)
		}
		saved, err = p.mh.MarshalBinary()
		if err != nil {
			panic(fmt.Sprintf("kcdsa: failed to marshal hash: %v", err))
		}
	} else {
		saved = p.initVectors
	}
	defer zeroize.Zeroize(saved)

	for {
		if err := p.mh.UnmarshalBinary(saved); err != nil {
			panic(fmt.Sprintf("kcdsa: failed to unmarshal hash: %v", err))
		}
		_, _ = p.mh.Write(iBuf[:])
		hbuf = p.mh.Sum(hbuf[:0])

		if i >= p.lh {
			i -= p.lh
			copy(dst[i:], hbuf)
			if i == 0 {
				break
			}
		} else {
			copy(dst, hbuf[len(hbuf)-i:])
			break
		}

		iBuf[0]++
	}

	return internal.RightMost(dst, nBits)
}
