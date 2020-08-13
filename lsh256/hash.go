package lsh256

import (
	"hash"
)

const (
	Size    = 32
	Size224 = 28

	BLOCKSIZE = 128
)

func New() hash.Hash {
	h := &lsh256{
		outlenbits: 256,
	}
	h.Reset()
	return h
}
func New224() hash.Hash {
	h := &lsh256{
		outlenbits: 224,
	}
	h.Reset()
	return h
}

func Sum256(data []byte) (sum [Size]byte) {
	b := lsh256{
		outlenbits: 256,
	}
	b.Reset()
	b.Write(data)

	return b.checkSum()
}

func Sum224(data []byte) (sum224 [Size224]byte) {
	b := lsh256{
		outlenbits: 224,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum224[:], sum[:Size224])

	return
}
