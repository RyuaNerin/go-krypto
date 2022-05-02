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
	BlockSize = 128
)

type algType int

const (
	lshType256H256 algType = 32 // 256
	lshType256H224         = 24 // 224
)

// New returns a new hash.Hash computing the LSH-256 checksum.
func New() hash.Hash {
	ctx := newContext(lshType256H256)
	ctx.Reset()
	return ctx
}

// New224 returns a new hash.Hash computing the LSH-224 checksum.
func New224() hash.Hash {
	ctx := newContext(lshType256H224)
	ctx.Reset()
	return ctx
}

// Sum256 returns the LSH-256 checksum of the data.
func Sum256(data []byte) (sum [Size]byte) {
	ctx := newContext(lshType256H256)
	ctx.Reset()
	ctx.Write(data)

	hash := ctx.Sum(nil)
	copy(sum[:], hash)

	return
}

// Sum224 returns the LSH-224 checksum of the data.
func Sum224(data []byte) (sum224 [Size224]byte) {
	ctx := newContext(lshType256H224)
	ctx.Reset()
	ctx.Write(data)

	hash := ctx.Sum(nil)
	copy(sum224[:], hash)

	return
}
