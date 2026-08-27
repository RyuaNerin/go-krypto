// Package lsh256 implements the LSH-256, LSH-256-224 hash algorithms as defined in KS X 3262
package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func init() {
	krypto.RegisterHash(krypto.LSH256, New)
	krypto.RegisterHash(krypto.LSH256_224, New224)
}

const (
	// The size of a LSH-256 checksum in bytes.
	Size = 32
	// The size of a LSH-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-256 and LSH-224 in bytes.
	BlockSize = 128
)

// New returns a new hash.Hash computing the LSH-256 checksum.
// The returned value is a [krypto.Hash] and therefore implements Zeroize(): call it
// to erase the hash state when the hasher is no longer needed. After Zeroize() a
// call to Sum panics and Write returns krypto.ErrStateCleared. All New* variants are
// zeroizable. The package-level Sum helpers already auto-zeroize their scratch context.
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

type (
	block         [BlockSize]byte
	lsh256Context struct {
		lsh256ContextData
		z zeroize.ZeroizedFlag
	}
	lsh256ContextData struct {
		// 16 bytes aligned

		algType           uint64
		_                 uint64
		remainDataByteLen uint64
		_                 uint64
		cvL               [8]uint32
		cvR               [8]uint32
		lastBlock         block
	}
)

var _ krypto.Zeroizable = (*lsh256Context)(nil)

func (ctx *lsh256Context) Zeroize() {
	ctx.lsh256ContextData.Zeroize()
	ctx.z.SetInvalid()
}

func (ctx *lsh256ContextData) Zeroize() {
	ctx.algType = 0
	ctx.remainDataByteLen = 0
	clear(ctx.cvL[:])
	clear(ctx.cvR[:])
	clear(ctx.lastBlock[:])
}

func (ctx *lsh256ContextData) Reset() {
	switch ctx.algType {
	case Size224:
		copy(ctx.cvL[:], iv224[:8])
		copy(ctx.cvR[:], iv224[8:])

	case Size:
		copy(ctx.cvL[:], iv256[:8])
		copy(ctx.cvR[:], iv256[8:])

	default:
		panic(msgInvalidSize)
	}
	clear(ctx.lastBlock[:])
	ctx.remainDataByteLen = 0
}
