//go:build amd64

package lsh256

import (
	"hash"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasAVX2
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsmData)
	update func(ctx *lsh256ContextAsmData, data []byte)
	final  func(ctx *lsh256ContextAsmData, hashval []byte)
}

var (
	simdSetDefault simdSet

	SimdSetSSE2 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSE2,
		final:  lsh256FinalSSE2,
	}
	SimdSetSSE2_v2 = simdSet{
		init: lsh256_sse2_init,
	}
	SimdSetSSSE3 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSSE3,
		final:  lsh256FinalSSSE3,
	}
	SimdSetAVX2 = simdSet{
		init:   lsh256InitAVX2,
		update: lsh256UpdateAVX2,
		final:  lsh256FinalAVX2,
	}
)

func init() {
	simdSetDefault = SimdSetSSE2

	if hasSSSE3 {
		simdSetDefault = SimdSetSSSE3
	}

	if hasAVX2 {
		simdSetDefault = SimdSetAVX2
	}
}

func NewContextAsm(algType int, simd simdSet) hash.Hash {
	ctx := new(lsh256ContextAsm)
	initContextAsm(ctx, algType, simd)
	return ctx
}

type lsh256ContextAsm struct {
	simd simdSet

	data lsh256ContextAsmData
}
type lsh256ContextAsmData struct {
	// 16 aligned
	algtype uint32
	_pad0   [16 - 4]byte
	// 16 aligned
	remain_databitlen uint32
	_pad1             [16 - 4]byte

	cv_l       [32]byte
	cv_r       [32]byte
	last_block [128]byte
}

func initContextAsm(ctx *lsh256ContextAsm, algtype int, simd simdSet) {
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
	ctx.simd.update(&ctx.data, data)

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(b []byte) []byte {
	hash := make([]byte, Size)
	ctx.simd.final(&ctx.data, hash)

	return append(b, hash[:ctx.Size()]...)
}
