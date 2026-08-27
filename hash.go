package krypto

import (
	"fmt"
	"hash"
	"strconv"
)

// KryptoHash identifies a cryptographic hash function that is implemented in another
// package.
type KryptoHash uint

const (
	_ KryptoHash = iota
	// import github.com/RyuaNerin/go-krypto/has160
	//
	// Deprecated: HAS-160 is a legacy algorithm.
	HAS160
	LSH256_224 // import github.com/RyuaNerin/go-krypto/lsh256
	LSH256     // import github.com/RyuaNerin/go-krypto/lsh256
	LSH512_224 // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512_256 // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512_384 // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512     // import github.com/RyuaNerin/go-krypto/lsh512
	maxHash
)

type hashSize struct {
	digest uint8
	block  uint16
	name   string
}

var (
	hashSizes = [maxHash]hashSize{
		HAS160:     {digest: 20, block: 64, name: "HAS-160"},
		LSH256_224: {digest: 28, block: 128, name: "LSH-256-224"},
		LSH256:     {digest: 32, block: 128, name: "LSH-256"},
		LSH512_224: {digest: 28, block: 256, name: "LSH-512-224"},
		LSH512_256: {digest: 32, block: 256, name: "LSH-512-256"},
		LSH512_384: {digest: 48, block: 256, name: "LSH-512-384"},
		LSH512:     {digest: 64, block: 256, name: "LSH-512"},
	}

	hashes [maxHash]func() hash.Hash
)

// New returns a new hash.Hash calculating the given hash function.
//
// If the hash function is not linked into the binary, New panics.
func (h KryptoHash) New() hash.Hash {
	if h > 0 && h < maxHash {
		f := hashes[h]
		if f != nil {
			return f()
		}
	}
	panic(fmt.Sprintf(msgUnavailableFormat, h))
}

// Size returns the digest size of the given hash function.
//
// If the hash function is unknown, Size panics.
func (h KryptoHash) Size() int {
	if h > 0 && h < maxHash {
		return int(hashSizes[h].digest)
	}
	panic(msgSizeOfUnknown)
}

// BlockSize returns the block size of the given hash function.
//
// If the hash function is unknown, BlockSize panics.
func (h KryptoHash) BlockSize() int {
	if h > 0 && h < maxHash {
		return int(hashSizes[h].block)
	}
	panic(msgBlockSizeOfUnknown)
}

// String returns the name of the given hash function.
//
// It returns "unknown hash value <n>" if the hash function is not recognized.
func (h KryptoHash) String() string {
	if h > 0 && h < maxHash {
		return hashSizes[h].name
	}
	return "unknown hash value " + strconv.Itoa(int(h))
}

// Available reports whether the given hash function is linked into the binary.
func (h KryptoHash) Available() bool {
	return h > 0 && h < maxHash && hashes[h] != nil
}

// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.
func RegisterHash(h KryptoHash, f func() hash.Hash) {
	if h == 0 || h >= maxHash {
		panic(msgRegisterHashOfUnknown)
	}
	hashes[h] = f
}
