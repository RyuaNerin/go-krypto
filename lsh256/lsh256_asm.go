//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lsh256

import (
	"hash"
)

type simdSet struct {
	init   func(ctx *lsh256Context, algtype uint64)
	update func(ctx *lsh256Context, data []byte)
	final  func(ctx *lsh256Context, hashval []byte)
}

func (simd *simdSet) NewContext(size int) hash.Hash {
	ctx := new(lsh256Context)
	simd.InitContext(ctx, size)
	return ctx
}

func (simd *simdSet) InitContext(ctx *lsh256Context, size int) {
	ctx.simd = simd
	ctx.size = size
	ctx.Reset()
}

type lsh256Context struct {
	data [16 + 16 + 4*8 + 4*8 + 128]byte // 최상단으로 배치하여 aligned 문제 수정...

	simd *simdSet
	size int
}

func (ctx *lsh256Context) Size() int {
	return ctx.size
}

func (ctx *lsh256Context) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256Context) Reset() {
	ctx.simd.init(ctx, uint64(ctx.size))
}

func (ctx *lsh256Context) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(ctx, data)

	return len(data), nil
}

func (ctx *lsh256Context) Sum(p []byte) []byte {
	ctx0 := *ctx
	hash := ctx0.checkSum()
	return append(p, hash[:ctx.size]...)
}

func (ctx *lsh256Context) checkSum() (hash [Size]byte) {
	ctx.simd.final(ctx, hash[:])
	return
}
