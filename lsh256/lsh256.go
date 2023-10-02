// Package lsh256 implements the LSH-256, LSH-256-224 hash algorithms as defined in KS X 3262
package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto"
)

func init() {
	krypto.RegisterHash(krypto.LSH256, New)
	krypto.RegisterHash(krypto.LSH256_224, New224)
}

var (
	newContext func(size int) hash.Hash               = newContextGo
	sum        func(size int, data []byte) [Size]byte = sumGo
)

const (
	// The size of a LSH-256 checksum in bytes.
	Size = 32
	// The size of a LSH-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-256 and LSH-224 in bytes.
	BlockSize = 128
)

// New returns a new hash.Hash computing the LSH-256 checksum.
func New() hash.Hash { return newContext(Size) }

// New224 returns a new hash.Hash computing the LSH-224 checksum.
func New224() hash.Hash { return newContext(Size224) }

// Sum256 returns the LSH-256 checksum of the data.
func Sum256(data []byte) (sum256 [Size]byte) {
	sum := sum(Size, data)
	copy(sum256[:], sum[:Size])
	return
}

// Sum224 returns the LSH-224 checksum of the data.
func Sum224(data []byte) (sum224 [Size224]byte) {
	sum := sum(Size224, data)
	copy(sum224[:], sum[:Size224])
	return
}
