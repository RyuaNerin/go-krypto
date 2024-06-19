//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lsh512

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type simdSet struct {
	init   func(ctx *lsh512ContextAsm, algtype uint64)
	update func(ctx *lsh512ContextAsm, data []byte)
	final  func(ctx *lsh512ContextAsm, hashval *byte)
}

var defaultSimd simdSet

func newContextAsm(size int) hash.Hash {
	return defaultSimd.NewContext(size)
}

func sumAsm(size int, data []byte) [Size]byte {
	return defaultSimd.Sum(size, data)
}

func (simd *simdSet) NewContext(size int) hash.Hash {
	ctx := &lsh512ContextAsm{
		simd:    simd,
		algType: size,
	}
	ctx.Reset()
	return ctx
}

func (simd *simdSet) Sum(size int, data []byte) [Size]byte {
	ctx := lsh512ContextAsm{
		simd:    simd,
		algType: size,
	}
	ctx.Reset()
	ctx.Write(data)

	return ctx.checkSum()
}

type lsh512ContextAsm struct {
	// 16 bytes aligned
	algType           int
	_                 uintptr
	remainDataByteLen int
	_                 [1]uintptr
	cvL               [8]uint64
	cvR               [8]uint64
	lastBlock         [BlockSize]byte

	simd *simdSet
}

func (ctx *lsh512ContextAsm) Size() int {
	return ctx.algType
}

func (ctx *lsh512ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh512ContextAsm) Reset() {
	ctx.simd.init(ctx, uint64(ctx.algType))
}

func (ctx *lsh512ContextAsm) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(ctx, data)

	return len(data), nil
}

func (ctx *lsh512ContextAsm) Sum(p []byte) []byte {
	ctx0 := *ctx
	hash := ctx0.checkSum()
	return append(p, hash[:ctx.algType]...)
}

func (ctx *lsh512ContextAsm) checkSum() (hash [Size]byte) {
	ctx.simd.final(ctx, memory.P8(hash[:]))
	return
}
