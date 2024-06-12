//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsm, algtype uint64)
	update func(ctx *lsh256ContextAsm, data []byte)
	final  func(ctx *lsh256ContextAsm, hashval *byte)
}

func (simd *simdSet) NewContext(size int) hash.Hash {
	ctx := &lsh256ContextAsm{
		simd:    simd,
		algType: size,
	}
	ctx.Reset()
	return ctx
}

func (simd *simdSet) Sum(size int, data []byte) [Size]byte {
	ctx := lsh256ContextAsm{
		simd:    simd,
		algType: size,
	}
	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

type lsh256ContextAsm struct {
	// 16 bytes aligned
	algType           int
	_                 uintptr
	remainDataByteLen int
	_                 uintptr
	cvL               [4]uint64
	cvR               [4]uint64
	lastBlock         [BlockSize]byte

	simd *simdSet
}

func (ctx *lsh256ContextAsm) Size() int {
	return ctx.algType
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Reset() {
	ctx.simd.init(ctx, uint64(ctx.algType))
}

func (ctx *lsh256ContextAsm) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(ctx, data)

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(p []byte) []byte {
	ctx0 := *ctx
	hash := ctx0.checkSum()
	return append(p, hash[:ctx.algType]...)
}

func (ctx *lsh256ContextAsm) checkSum() (hash [Size]byte) {
	ctx.simd.final(ctx, memory.PByte(hash[:]))
	return
}
