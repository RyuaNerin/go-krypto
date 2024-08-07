//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

import "github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"

var (
	hasSSE2  = true
	hasSSSE3 = cpu.X86.HasSSSE3
	hasAVX2  = cpu.X86.HasSSSE3 && cpu.X86.HasAVX && cpu.X86.HasAVX2

	simdSetSSE2 = simdSet{
		init:   __lsh256_sse2_init,
		update: __lsh256_sse2_update,
		final:  __lsh256_sse2_final,
	}
	simdSetSSSE3 = simdSet{
		init:   __lsh256_sse2_init,
		update: __lsh256_ssse3_update,
		final:  __lsh256_ssse3_final,
	}
	simdSetAVX2 = simdSet{
		init:   __lsh256_avx2_init,
		update: __lsh256_avx2_update,
		final:  __lsh256_avx2_final,
	}
)

func init() {
	switch {
	// avx2 is slower than ssse3
	// case hasAVX2:
	//	defaultSimd = simdSetAVX2

	case hasSSSE3:
		defaultSimd = simdSetSSSE3

	default:
		defaultSimd = simdSetSSE2
	}
}

//go:noescape
func __lsh256_sse2_init(ctx *lsh256ContextAsm, algtype uint64)

//go:noescape
func __lsh256_sse2_update(ctx *lsh256ContextAsm, data []byte)

//go:noescape
func __lsh256_sse2_final(ctx *lsh256ContextAsm, hashval *byte)

//go:noescape
//func __lsh256_ssse3_init(ctx *lsh256Context, algtype uint64)

//go:noescape
func __lsh256_ssse3_update(ctx *lsh256ContextAsm, data []byte)

//go:noescape
func __lsh256_ssse3_final(ctx *lsh256ContextAsm, hashval *byte)

//go:noescape
func __lsh256_avx2_init(ctx *lsh256ContextAsm, algtype uint64)

//go:noescape
func __lsh256_avx2_update(ctx *lsh256ContextAsm, data []byte)

//go:noescape
func __lsh256_avx2_final(ctx *lsh256ContextAsm, hashval *byte)
