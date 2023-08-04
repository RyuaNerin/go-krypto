package eckcdsa

import "hash"

func kisa_sngf(seed []byte, h hash.Hash, nBits int) []byte {
	LH := h.BlockSize()

	// p.12
	// from java
	i := ((nBits + 7) & 0xFFFFFFF8) / 8
	iBuf := make([]byte, 1)

	count := 0

	U := make([]byte, i)

	var hbuf []byte
	for {
		iBuf[0] = byte(count)

		h.Reset()
		h.Write(seed)
		h.Write(iBuf)
		hbuf = h.Sum(hbuf[:0])

		if i >= LH {
			i -= LH
			for j := 0; j < LH; j++ {
				U[j+i] = hbuf[j]
			}
			if i == 0 {
				break
			}
		} else {
			for j := 0; j < i; j++ {
				U[j] = hbuf[j+LH-i]
			}
			break
		}

		count++
	}

	i = nBits & 0x07
	if i != 0 {
		U[0] &= byte((1 << i) - 1)
	}

	return U
}
