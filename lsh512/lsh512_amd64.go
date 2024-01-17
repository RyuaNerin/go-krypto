//go:build amd64 && !purego
// +build amd64,!purego

package lsh512

import (
	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasSSSE3 && cpu.X86.HasAVX2 && cpu.X86.HasAVX
)

var (
	simdSetSSE2 = simdSet{
		init:   __lsh512_sse2_init,
		update: __lsh512_sse2_update,
		final:  __lsh512_sse2_final,
	}
	simdSetSSSE3 = simdSet{
		init:   __lsh512_sse2_init,
		update: __lsh512_ssse3_update,
		final:  __lsh512_ssse3_final,
	}
	simdSetAVX2 = simdSet{
		init:   __lsh512_avx2_init,
		update: __lsh512_avx2_update,
		final:  __lsh512_avx2_final,
	}
)

var (
	initContext = simdSetSSE2.InitContext
)

func init() {
	if hasAVX2 {
		initContext = simdSetAVX2.InitContext
	} else if hasSSSE3 {
		initContext = simdSetSSSE3.InitContext
	}
}

//go:noescape
func __lsh512_sse2_init(ctx *lsh512Context, algtype uint64)

//go:noescape
func __lsh512_sse2_update(ctx *lsh512Context, data []byte)

//go:noescape
func __lsh512_sse2_final(ctx *lsh512Context, hashval *byte)

//go:noescape
//func __lsh512_ssse3_init(ctx *lsh512Context, algtype uint64)

//go:noescape
func __lsh512_ssse3_update(ctx *lsh512Context, data []byte)

//go:noescape
func __lsh512_ssse3_final(ctx *lsh512Context, hashval *byte)

//go:noescape
func __lsh512_avx2_init(ctx *lsh512Context, algtype uint64)

//go:noescape
func __lsh512_avx2_update(ctx *lsh512Context, data []byte)

//go:noescape
func __lsh512_avx2_final(ctx *lsh512Context, hashval *byte)
