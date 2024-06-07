package kcdsa

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
)

type ppgfCtx struct {
	h           hash.Hash
	initVectors []byte // either initVectors or initSources

	mh          internal.Hash
	initSources [][]byte // either initVectors or initSources

	lh int
}

func ppgf(
	dst []byte,
	nBits int, h hash.Hash, src ...[]byte,
) []byte {
	return newPPGF(h).Read(dst, nBits, src...)
}

func newPPGF(h hash.Hash, src ...[]byte) (ppgf ppgfCtx) {
	mh, ok := h.(internal.Hash)
	if ok {
		h.Reset()
		for _, v := range src {
			h.Write(v)
		}
		ppgf.mh = mh
		ppgf.initVectors, _ = mh.MarshalBinary()
	} else {
		ppgf.h = h
		ppgf.initSources = src
	}

	ppgf.lh = h.Size()
	return
}

func (p ppgfCtx) Read(dst []byte, nBits int, src ...[]byte) []byte {
	// p.12
	// from java
	i := internal.BitsToBytes(nBits)
	dst = internal.Grow(dst, i)

	var iBuf [1]byte
	hbuf := make([]byte, 0, p.lh)

	if p.mh != nil {
		var saved []byte
		if len(src) > 0 {
			p.mh.UnmarshalBinary(p.initVectors)
			for _, v := range src {
				p.mh.Write(v)
			}
			saved, _ = p.mh.MarshalBinary()
		} else {
			saved = p.initVectors
		}

		for {
			p.mh.UnmarshalBinary(saved)
			p.mh.Write(iBuf[:])
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
	} else {
		for {
			p.h.Reset()
			for _, v := range p.initSources {
				p.h.Write(v)
			}
			for _, v := range src {
				p.h.Write(v)
			}
			p.h.Write(iBuf[:])
			hbuf = p.h.Sum(hbuf[:0])

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
	}

	return internal.RightMost(dst, nBits)
}
