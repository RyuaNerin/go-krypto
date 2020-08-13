package lsh512

import (
	"hash"
)

const (
	Size    = 64
	Size224 = 28
	Size256 = 32
	Size384 = 48

	BLOCKSIZE = 256
)

func New() hash.Hash {
	h := &lsh512{
		outlenbits: 512,
	}
	h.Reset()
	return h
}
func New384() hash.Hash {
	h := &lsh512{
		outlenbits: 384,
	}
	h.Reset()
	return h
}
func New256() hash.Hash {
	h := &lsh512{
		outlenbits: 256,
	}
	h.Reset()
	return h
}
func New224() hash.Hash {
	h := &lsh512{
		outlenbits: 224,
	}
	h.Reset()
	return h
}

func Sum(data []byte) (sum [Size]byte) {
	b := lsh512{
		outlenbits: 512,
	}
	b.Reset()
	b.Write(data)

	return b.checkSum()
}

func Sum384(data []byte) (sum384 [Size384]byte) {
	b := lsh512{
		outlenbits: 384,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum384[:], sum[:Size384])
	return
}

func Sum256(data []byte) (sum256 [Size256]byte) {
	b := lsh512{
		outlenbits: 256,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum256[:], sum[:Size256])
	return
}

func Sum224(data []byte) (sum224 [Size224]byte) {
	b := lsh512{
		outlenbits: 224,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum224[:], sum[:Size224])
	return
}
