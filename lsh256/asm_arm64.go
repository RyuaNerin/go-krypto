//go:build arm64 && !purego
// +build arm64,!purego

package lsh256

var (
	hasNEON = true

	simdSetNEON = simdSet{
		init:   __lsh256_neon_init,
		update: __lsh256_neon_update,
		final:  __lsh256_neon_final,
	}
)

//go:noescape
func __lsh256_neon_init(ctx *lsh256ContextAsm, algtype uint64)

//go:noescape
func __lsh256_neon_update(ctx *lsh256ContextAsm, data []byte)

//go:noescape
func __lsh256_neon_final(ctx *lsh256ContextAsm, hashval *byte)
