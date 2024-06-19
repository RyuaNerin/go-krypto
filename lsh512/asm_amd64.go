//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lsh512

import "github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"

var (
	hasSSE2  = true
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasSSSE3 && cpu.X86.HasAVX && cpu.X86.HasAVX2

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

func init() {
	switch {
	case hasAVX2:
		defaultSimd = simdSetAVX2

	case hasSSSE3:
		defaultSimd = simdSetSSSE3

	default:
		defaultSimd = simdSetSSE2
	}
}

//go:noescape
func __lsh512_sse2_init(ctx *lsh512ContextAsm, algtype uint64)

//go:noescape
func __lsh512_sse2_update(ctx *lsh512ContextAsm, data []byte)

//go:noescape
func __lsh512_sse2_final(ctx *lsh512ContextAsm, hashval *byte)

//go:noescape
//func __lsh512_ssse3_init(ctx *lsh512Context, algtype uint64)

//go:noescape
func __lsh512_ssse3_update(ctx *lsh512ContextAsm, data []byte)

//go:noescape
func __lsh512_ssse3_final(ctx *lsh512ContextAsm, hashval *byte)

//go:noescape
func __lsh512_avx2_init(ctx *lsh512ContextAsm, algtype uint64)

//go:noescape
func __lsh512_avx2_update(ctx *lsh512ContextAsm, data []byte)

//go:noescape
func __lsh512_avx2_final(ctx *lsh512ContextAsm, hashval *byte)
