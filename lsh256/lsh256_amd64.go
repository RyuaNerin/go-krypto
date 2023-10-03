//go:build amd64 && gc && !purego

package lsh256

import (
	"hash"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasSSSE3 && cpu.X86.HasAVX2 && cpu.X86.HasAVX

	useAVX2 = false
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsmData)
	update func(ctx *lsh256ContextAsmData, data []byte)
	final  func(ctx *lsh256ContextAsmData, hashval []byte)
}

var (
	simdSetDefault simdSet

	simdSetSSE2 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSE2,
		final:  lsh256FinalSSE2,
	}
	simdSetSSSE3 = simdSet{
		init:   lsh256InitSSE2, // lsh256InitSSSE3,
		update: lsh256UpdateSSSE3,
		final:  lsh256FinalSSSE3,
	}
	simdSetAVX2 = simdSet{
		init:   lsh256InitAVX2,
		update: lsh256UpdateAVX2,
		final:  lsh256FinalAVX2,
	}
)

func init() {
	simdSetDefault = simdSetSSE2

	if hasSSSE3 {
		simdSetDefault = simdSetSSSE3
	}

	if hasAVX2 && useAVX2 {
		simdSetDefault = simdSetAVX2
	}
}

func init() {
	newContext = func(size int) hash.Hash {
		return newContextAsm(size, simdSetDefault)
	}
}

func newContextAsm(size int, simd simdSet) hash.Hash {
	ctx := new(lsh256ContextAsm)
	initContextAsm(ctx, size, simd)
	return ctx
}

func sumAsm(size int, data []byte) [Size]byte {
	var b lsh256ContextAsm
	initContextAsm(&b, size, simdSetDefault)
	b.Reset()
	b.Write(data)

	return b.checkSum()
}

type lsh256ContextAsm struct {
	simd simdSet

	data lsh256ContextAsmData
}
type lsh256ContextAsmData struct {
	// 16 aligned
	algtype            uint32
	_                  [4]byte
	remain_databytelen uint64

	cv_l       [32]byte
	cv_r       [32]byte
	last_block [128]byte
}

func initContextAsm(ctx *lsh256ContextAsm, size int, simd simdSet) {
	ctx.simd = simd
	ctx.data.algtype = uint32(size)
	ctx.Reset()
}

func (ctx *lsh256ContextAsm) Size() int {
	return int(ctx.data.algtype)
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Reset() {
	ctx.data.remain_databytelen = 0
	ctx.simd.init(&ctx.data)
}

func (ctx *lsh256ContextAsm) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(&ctx.data, data)

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(p []byte) []byte {
	ctx0 := *ctx
	hash := ctx0.checkSum()
	return append(p, hash[:ctx.Size()]...)
}

func (ctx *lsh256ContextAsm) checkSum() (hash [Size]byte) {
	ctx.simd.final(&ctx.data, hash[:])
	return
}
