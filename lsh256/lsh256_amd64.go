//go:build amd64

package lsh256

import (
	"hash"

	"golang.org/x/sys/cpu"
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsmData)
	update func(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)
	final  func(ctx *lsh256ContextAsmData, hashval []byte)
}

var (
	simdSetDefault simdSet

	simdSetSSE2 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSE2,
		final:  lsh256FinalSSE2,
	}
	/**
	simdSetSSSE3 = simdSet{
		init:   lsh256InitSSSE3,
		update: lsh256UpdateSSSE3,
		final:  lsh256FinalSSSE3,
	}
	*/
	simdSetAVX2 = simdSet{
		init:   lsh256InitAVX2,
		update: lsh256UpdateAVX2,
		final:  lsh256FinalAVX2,
	}
)

func init() {
	simdSetDefault = simdSetSSE2

	/**
	if cpu.X86.HasSSSE3 {
		simdSetDefault = simdSetSSSE3
	}
	*/

	if cpu.X86.HasAVX2 {
		simdSetDefault = simdSetAVX2
	}
}

func newContextAsm(algType algType, simd simdSet) hash.Hash {
	ctx := new(lsh256ContextAsm)
	initContextAsm(ctx, algType, simd)
	return ctx
}

type lsh256ContextAsm struct {
	simd simdSet

	data lsh256ContextAsmData
}
type lsh256ContextAsmData struct {
	algtype           uint32
	remain_databitlen uint32
	cv_l              [32 / 4]uint32
	cv_r              [32 / 4]uint32
	last_block        [128]byte
}

func initContextAsm(ctx *lsh256ContextAsm, algtype algType, simd simdSet) {
	ctx.simd = simd
	ctx.data.algtype = uint32(algtype)
	ctx.Reset()
}

func (ctx *lsh256ContextAsm) Size() int {
	return int(ctx.data.algtype)
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Reset() {
	ctx.data.remain_databitlen = 0
	ctx.simd.init(&ctx.data)
}

func (ctx *lsh256ContextAsm) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx.simd.update(&ctx.data, data, uint32(len(data)*8))

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(b []byte) []byte {
	hash := make([]byte, Size)
	ctx.simd.final(&ctx.data, hash)

	return append(b, hash[:ctx.Size()]...)
}
