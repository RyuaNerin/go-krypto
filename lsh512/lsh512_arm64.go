//go:build arm64 && !purego
// +build arm64,!purego

package lsh512

var (
	simdSetNEON = simdSet{
		init:   __lsh512_neon_init,
		update: __lsh512_neon_update,
		final:  __lsh512_neon_final,
	}
)

var (
	initContext = simdSetNEON.InitContext
)

//go:noescape
func __lsh512_neon_init(ctx *lsh512Context, algtype uint64)

//go:noescape
func __lsh512_neon_update(ctx *lsh512Context, data []byte)

//go:noescape
func __lsh512_neon_final(ctx *lsh512Context, hashval *byte)
