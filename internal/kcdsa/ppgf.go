package kcdsa

import (
	"encoding"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
)

type marshalableHash interface {
	hash.Hash
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func ppgf(
	dst []byte,
	nBits int, h hash.Hash, src ...[]byte,
) []byte {
	if hh, ok := h.(marshalableHash); ok {
		return ppgfFast(dst, nBits, hh, src...)
	}

	// p.12
	// from java
	i := internal.BitsToBytes(nBits)
	dst = internal.Grow(dst, i)

	LH := h.Size()

	count := 0
	var iBuf [1]byte
	hbuf := make([]byte, 0, LH)

	for {
		iBuf[0] = byte(count)

		h.Reset()
		for _, v := range src {
			h.Write(v)
		}
		h.Write(iBuf[:])
		hbuf = h.Sum(hbuf[:0])

		if i >= LH {
			i -= LH
			copy(dst[i:], hbuf)
			if i == 0 {
				break
			}
		} else {
			copy(dst, hbuf[len(hbuf)-i:])
			break
		}

		count++
	}

	return internal.RightMost(dst, nBits)
}

func ppgfFast(
	dst []byte,
	nBits int, h marshalableHash, src ...[]byte,
) []byte {
	for _, v := range src {
		h.Write(v)
	}
	marshalled, _ := h.MarshalBinary()

	// p.12
	// from java
	i := internal.BitsToBytes(nBits)
	dst = internal.Grow(dst, i)

	LH := h.Size()

	count := 0
	var iBuf [1]byte
	hbuf := make([]byte, 0, LH)

	for {
		iBuf[0] = byte(count)

		h.UnmarshalBinary(marshalled)
		h.Write(iBuf[:])
		hbuf = h.Sum(hbuf[:0])

		if i >= LH {
			i -= LH
			copy(dst[i:], hbuf)
			if i == 0 {
				break
			}
		} else {
			copy(dst, hbuf[len(hbuf)-i:])
			break
		}

		count++
	}

	return internal.RightMost(dst, nBits)
}
