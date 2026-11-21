// Package has160 implements the HAS-160 hash, as defined in TTAS.KO-12.0011/R2.
//
// Deprecated: HAS-160 is a legacy algorithm.
package has160

import (
	"encoding/binary"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func init() {
	krypto.RegisterHash(krypto.HAS160, New)
}

const (
	// The size of a HAS-160 checksum in bytes.
	Size = 20

	// The block size of HAS-160 in bytes.
	BlockSize = 64
)

// New returns a new hash.Hash computing the HAS-160 checksum.
// The returned value is a [krypto.Hash] and therefore implements Zeroize(): call it
// to erase the hash state when the hasher is no longer needed. After Zeroize() a
// call to Sum panics and Write returns krypto.ErrStateCleared.
func New() hash.Hash {
	ctx := new(has160Context)
	ctx.z.SetValid()
	ctx.Reset()
	return ctx
}

// Sum returns the HAS-160 checksum of the data.
func Sum(data []byte) (sum [Size]byte) {
	var ctx has160Context
	ctx.z.SetValid()
	defer ctx.Zeroize()

	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

type has160Context struct {
	z zeroize.ZeroizedFlag
	has160ContextData
}
type has160ContextData struct {
	h      [5]uint32
	block  [BlockSize]byte
	boff   int
	length uint64
}

var _ krypto.Hash = (*has160Context)(nil)

func (ctx *has160Context) Zeroize() {
	clear(ctx.h[:])
	clear(ctx.block[:])
	ctx.boff = 0
	ctx.length = 0

	ctx.z.Zeroize()
}

func (ctx *has160Context) Size() int {
	return Size
}

func (ctx *has160Context) BlockSize() int {
	return BlockSize
}

func (ctx *has160Context) Reset() {
	ctx.h = initH
	clear(ctx.block[:])
	ctx.boff = 0
	ctx.length = 0
}

func (ctx *has160Context) Write(p []byte) (n int, err error) {
	if err := ctx.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(p) == 0 {
		return 0, nil
	}

	lenP := len(p)
	boff := ctx.boff
	if boff < 0 || boff >= BlockSize {
		panic(msgInvalidInternalState)
	}

	ctx.length += uint64(len(p))

	gap := BlockSize - boff
	if 0 < boff && gap <= len(p) {
		copy(ctx.block[boff:], p[:gap])
		ctx.stepBlock(&ctx.block)
		ctx.boff = 0
		boff = 0
		p = p[gap:]
	}

	for len(p) >= BlockSize {
		ctx.stepBlock((*[BlockSize]byte)(p))
		ctx.boff = 0
		p = p[BlockSize:]
	}

	if len(p) > 0 {
		copy(ctx.block[boff:], p)
		ctx.boff = boff + len(p)
	}

	return lenP, nil
}

func (ctx *has160Context) Sum(p []byte) []byte {
	ctx.z.ThrowIfZeroized()

	ctx0 := ctx.has160ContextData
	defer zeroize.Zeroize(&ctx0)
	hash := ctx0.checkSum()
	return append(p, hash[:]...)
}

func (ctx *has160ContextData) checkSum() (hash [20]byte) {
	boff := ctx.boff
	if boff < 0 || boff >= BlockSize {
		panic(msgInvalidInternalState)
	}
	ctx.block[boff] = 0x80
	boff++

	if BlockSize-8 < boff {
		clear(ctx.block[boff:])
		ctx.stepBlock(&ctx.block)
		boff = 0
	}

	clear(ctx.block[boff:])
	binary.LittleEndian.PutUint64(ctx.block[BlockSize-8:], ctx.length*8)
	ctx.stepBlock(&ctx.block)

	binary.LittleEndian.PutUint32(hash[0:], ctx.h[0])
	binary.LittleEndian.PutUint32(hash[4:], ctx.h[1])
	binary.LittleEndian.PutUint32(hash[8:], ctx.h[2])
	binary.LittleEndian.PutUint32(hash[12:], ctx.h[3])
	binary.LittleEndian.PutUint32(hash[16:], ctx.h[4])

	return
}
