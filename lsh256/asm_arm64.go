//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lsh256

var (
	hasNEON = true

	simdSetNEON = simdSet{
		init:   __lsh256_neon_init,
		update: __lsh256_neon_update,
		final:  __lsh256_neon_final,
	}
)

func init() {
	defaultSimd = simdSetNEON
}

//go:noescape
func __lsh256_neon_init(ctx *lsh256ContextAsm, algtype uint64)

//go:noescape
func __lsh256_neon_update(ctx *lsh256ContextAsm, data []byte)

//go:noescape
func __lsh256_neon_final(ctx *lsh256ContextAsm, hashval *byte)
