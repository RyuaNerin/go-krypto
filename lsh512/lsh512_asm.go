//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lsh512

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/ptr"
)

type simdSet struct {
	init   func(ctx *lsh512Context, algtype uint64)
	update func(ctx *lsh512Context, data []byte)
	final  func(ctx *lsh512Context, hashval *byte)
}

func (simd *simdSet) NewContext(size int) hash.Hash {
	ctx := new(lsh512Context)
	simd.InitContext(ctx, size)
	return ctx
}

func (simd *simdSet) InitContext(ctx *lsh512Context, size int) {
	ctx.simd = simd
	ctx.size = size
	ctx.Reset()
}

type lsh512Context struct {
	data [16 + 16 + 8*8 + 8*8 + 256]byte // 최상단으로 배치하여 aligned 문제 수정...

	simd *simdSet
	size int
}

func initContextAsm(ctx *lsh512Context, simd *simdSet, size int) {
	ctx.simd = simd
	ctx.size = size
	ctx.Reset()
}

func (ctx *lsh512Context) Size() int {
	return ctx.size
}

func (ctx *lsh512Context) BlockSize() int {
	return BlockSize
}

func (ctx *lsh512Context) Reset() {
	ctx.simd.init(ctx, uint64(ctx.size))
}

func (ctx *lsh512Context) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(ctx, data)

	return len(data), nil
}

func (ctx *lsh512Context) Sum(p []byte) []byte {
	ctx0 := *ctx
	hash := ctx0.checkSum()
	return append(p, hash[:ctx.size]...)
}

func (ctx *lsh512Context) checkSum() (hash [Size]byte) {
	ctx.simd.final(ctx, ptr.BytePtr(hash[:]))
	return
}
