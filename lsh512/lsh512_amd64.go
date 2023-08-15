//go:build amd64 && gc && !purego

package lsh512

import (
	"hash"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasSSSE3 && cpu.X86.HasAVX && cpu.X86.HasAVX2
)

type simdSet struct {
	init   func(ctx *lsh512ContextAsmData)
	update func(ctx *lsh512ContextAsmData, data []byte)
	final  func(ctx *lsh512ContextAsmData, hashval []byte)
}

var (
	simdSetDefault simdSet

	simdSetSSE2 = simdSet{
		init:   lsh512InitSSE2,
		update: lsh512UpdateSSE2,
		final:  lsh512FinalSSE2,
	}
	simdSetSSSE3 = simdSet{
		init:   lsh512InitSSE2,
		update: lsh512UpdateSSSE3,
		final:  lsh512FinalSSSE3,
	}
	simdSetAVX2 = simdSet{
		init:   lsh512InitAVX2,
		update: lsh512UpdateAVX2,
		final:  lsh512FinalAVX2,
	}
)

func init() {
	simdSetDefault = simdSetSSE2

	if hasSSSE3 {
		simdSetDefault = simdSetSSSE3
	}

	if hasAVX2 {
		simdSetDefault = simdSetAVX2
	}
}

func newContextAsm(algType int, simd simdSet) hash.Hash {
	ctx := new(lsh512ContextAsm)
	initContextAsm(ctx, algType, simd)
	return ctx
}

type lsh512ContextAsm struct {
	simd simdSet

	data lsh512ContextAsmData
}
type lsh512ContextAsmData struct {
	// 16 aligned
	algtype            uint32
	_                  [4]byte
	remain_databytelen uint64

	cv_l         [8]uint64
	cv_r         [8]uint64
	i_last_block [256]byte
}

func initContextAsm(ctx *lsh512ContextAsm, algtype int, simd simdSet) {
	ctx.simd = simd
	ctx.data.algtype = uint32(algtype)
	ctx.Reset()
}

func (ctx *lsh512ContextAsm) Size() int {
	return int(ctx.data.algtype)
}

func (ctx *lsh512ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh512ContextAsm) Reset() {
	ctx.data.remain_databytelen = 0
	ctx.simd.init(&ctx.data)
}

func (ctx *lsh512ContextAsm) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(&ctx.data, data)

	return len(data), nil
}

func (ctx *lsh512ContextAsm) Sum(b []byte) []byte {
	hash := make([]byte, Size)
	ctx.simd.final(&ctx.data, hash)

	return append(b, hash[:ctx.Size()]...)
}
