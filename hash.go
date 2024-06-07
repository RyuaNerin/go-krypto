package krypto

import (
	"fmt"
	"hash"
	"strconv"
)

// Hash identifies a cryptographic hash function that is implemented in another
// package.
type KryptoHash uint

func (h KryptoHash) String() string {
	switch h {
	case HAS160:
		return "HAS-160"
	case LSH256_224:
		return "LSH-256-224"
	case LSH256:
		return "LSH-256"
	case LSH512_224:
		return "LSH-512-224"
	case LSH512_256:
		return "LSH-512-256"
	case LSH512_384:
		return "LSH-512-384"
	case LSH512:
		return "LSH-512"
	default:
		return "unknown hash value " + strconv.Itoa(int(h))
	}
}

const (
	_          KryptoHash = iota
	HAS160                // import github.com/RyuaNerin/go-krypto/has160
	LSH256_224            // import github.com/RyuaNerin/go-krypto/lsh256
	LSH256                // import github.com/RyuaNerin/go-krypto/lsh256
	LSH512_224            // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512_256            // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512_384            // import github.com/RyuaNerin/go-krypto/lsh512
	LSH512                // import github.com/RyuaNerin/go-krypto/lsh512
	maxHash
)

var digestSizes = []uint8{
	HAS160:     20,
	LSH256_224: 28,
	LSH256:     32,
	LSH512_224: 28,
	LSH512_256: 32,
	LSH512_384: 48,
	LSH512:     64,
}

func (h KryptoHash) Size() int {
	if h > 0 && h < maxHash {
		return int(digestSizes[h])
	}
	panic(msgSizeOfUnknwon)
}

var hashes = make([]func() hash.Hash, maxHash)

// New returns a new hash.Hash calculating the given hash function. New panics
// if the hash function is not linked into the binary.
func (h KryptoHash) New() hash.Hash {
	if h > 0 && h < maxHash {
		f := hashes[h]
		if f != nil {
			return f()
		}
	}
	panic(fmt.Sprintf(msgUnavailableFormat, h))
}

// Available reports whether the given hash function is linked into the binary.
func (h KryptoHash) Available() bool {
	return h < maxHash && hashes[h] != nil
}

// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.
func RegisterHash(h KryptoHash, f func() hash.Hash) {
	if h >= maxHash {
		panic(msgRegisterHashOfUnknown)
	}
	hashes[h] = f
}
