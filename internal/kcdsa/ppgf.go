package kcdsa

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
)

func ppgf(
	in []byte,
	nBits int, h hash.Hash, seed ...[]byte,
) []byte {
	// p.12
	// from java
	i := internal.Bytes(nBits)
	in = internal.Expand(in, i)

	LH := h.Size()

	count := 0
	var iBuf [1]byte
	hbuf := make([]byte, 0, LH)

	for {
		iBuf[0] = byte(count)

		h.Reset()
		for _, v := range seed {
			h.Write(v)
		}
		h.Write(iBuf[:])
		hbuf = h.Sum(hbuf[:0])

		if i >= LH {
			i -= LH
			copy(in[i:], hbuf)
			if i == 0 {
				break
			}
		} else {
			copy(in, hbuf[len(hbuf)-i:])
			break
		}

		count++
	}

	return internal.TruncateLeft(in, nBits)
}
