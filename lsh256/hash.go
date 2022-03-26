// Package lsh256 implements the LSH-256 hash algorithms as defined in TTAK.KO-12.0276
package lsh256

import (
	"errors"
	"hash"
)

var ErrInvalidDataBitLen = errors.New("krypto/lsh256: bit level update is not allowed")

const (
	// The size of a LSH-256 checksum in bytes.
	Size = 32
	// The size of a LSH-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-256 and LSH-224 in bytes.
	BLOCKSIZE = 128
)

// New returns a new hash.Hash computing the LSH-256 checksum.
func New() hash.Hash {
	h := &lsh256{
		outlenbits: 256,
	}
	h.Reset()
	return h
}

// New224 returns a new hash.Hash computing the LSH-224 checksum.
func New224() hash.Hash {
	h := &lsh256{
		outlenbits: 224,
	}
	h.Reset()
	return h
}

// Sum256 returns the LSH-256 checksum of the data.
func Sum256(data []byte) (sum [Size]byte) {
	b := lsh256{
		outlenbits: 256,
	}
	b.Reset()
	b.Write(data)

	return b.checkSum()
}

// Sum224 returns the LSH-224 checksum of the data.
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
