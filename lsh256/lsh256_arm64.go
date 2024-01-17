//go:build arm64 && !purego
// +build arm64,!purego

package lsh256

var (
	simdSetNEON = simdSet{
		init:   __lsh256_neon_init,
		update: __lsh256_neon_update,
		final:  __lsh256_neon_final,
	}
)

var (
	initContext = simdSetNEON.InitContext
)

//go:noescape
func __lsh256_neon_init(ctx *lsh256Context, algtype uint64)

//go:noescape
func __lsh256_neon_update(ctx *lsh256Context, data []byte)

//go:noescape
func __lsh256_neon_final(ctx *lsh256Context, hashval *byte)
