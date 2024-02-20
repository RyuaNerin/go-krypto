//go:build (amd64 || amd64p32) && !purego
// +build amd64 amd64p32
// +build !purego

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
