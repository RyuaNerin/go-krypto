//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/ptr"
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsm, algtype uint64)
	update func(ctx *lsh256ContextAsm, data []byte)
	final  func(ctx *lsh256ContextAsm, hashval *byte)
}

func (simd *simdSet) NewContext(size int) hash.Hash {
	ctx := new(lsh256ContextAsm)
	ctx.simd = simd
	ctx.size = size
	ctx.Reset()
	return ctx
}

func (simd *simdSet) Sum(size int, data []byte) [Size]byte {
	var ctx lsh256ContextAsm
	ctx.simd = simd
	ctx.size = size
	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

const contextDataSize = 16 + 16 + 4*8 + 4*8 + 128

type lsh256ContextAsm struct {
	//nolint:unused
	data [contextDataSize]byte // 최상단으로 배치하여 aligned 문제 수정...

	simd *simdSet
	size int
}

func (ctx *lsh256ContextAsm) Size() int {
	return ctx.size
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Reset() {
	ctx.simd.init(ctx, uint64(ctx.size))
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
	return append(p, hash[:ctx.size]...)
}

func (ctx *lsh256ContextAsm) checkSum() (hash [Size]byte) {
	ctx.simd.final(ctx, ptr.PByte(hash[:]))
	return
}
