// Package lsh512 implements the LSH-512, LSH-384, LSH-512-256, LSH-512-224 hash algorithms as defined in KS X 3262
package lsh512

import (
	"errors"
	"hash"
)

var ErrInvalidDataBitLen = errors.New("krypto/lsh512: bit level update is not allowed")

const (
	// The size of a LSH-512 checksum in bytes.
	Size = 64
	// The size of a LSH-384 checksum in bytes.
	Size384 = 48
	// The size of a LSH-512-256 checksum in bytes.
	Size256 = 32
	// The size of a LSH-512-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-512, LSH-384, LSH-512-256 and LSH-512-224 in bytes.
	BlockSize = 256
)

func newHash(size int) hash.Hash {
	h := &lsh512{
		outlenbits: size * 8,
	}
	h.Reset()
	return h
}

// New returns a new hash.Hash computing the LSH-512 checksum.
func New() hash.Hash { return newHash(Size) }

// New384 returns a new hash.Hash computing the LSH-384 checksum.
func New384() hash.Hash { return newHash(Size384) }

// New256 returns a new hash.Hash computing the LSH-512-256 checksum.
func New256() hash.Hash { return newHash(Size256) }

// New224 returns a new hash.Hash computing the LSH-512-224 checksum.
func New224() hash.Hash { return newHash(Size224) }

// Sum512 returns the LSH-512 checksum of the data.
func Sum512(data []byte) (sum [Size]byte) {
	b := lsh512{
		outlenbits: Size * 8,
	}
	b.Reset()
	b.Write(data)

	return b.checkSum()
}

// Sum384 returns the LSH-384 checksum of the data.
func Sum384(data []byte) (sum384 [Size384]byte) {
	b := lsh512{
		outlenbits: Size384 * 8,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum384[:], sum[:Size384])
	return
}

// Sum256 returns the LSH-512-256 checksum of the data.
func Sum256(data []byte) (sum256 [Size256]byte) {
	b := lsh512{
		outlenbits: Size256 * 8,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum256[:], sum[:Size256])
	return
}

// Sum224 returns the LSH-512-224 checksum of the data.
func Sum224(data []byte) (sum224 [Size224]byte) {
	b := lsh512{
		outlenbits: Size224 * 8,
	}
	b.Reset()
	b.Write(data)

	sum := b.checkSum()
	copy(sum224[:], sum[:Size224])
	return
}
