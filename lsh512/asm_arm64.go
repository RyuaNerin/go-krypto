//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lsh512

var (
	hasNEON = true

	simdSetNEON = simdSet{
		init:   __lsh512_neon_init,
		update: __lsh512_neon_update,
		final:  __lsh512_neon_final,
	}
)

func init() {
	defaultSimd = simdSetNEON
}

//go:noescape
func __lsh512_neon_init(ctx *lsh512ContextAsm, algtype uint64)

//go:noescape
func __lsh512_neon_update(ctx *lsh512ContextAsm, data []byte)

//go:noescape
func __lsh512_neon_final(ctx *lsh512ContextAsm, hashval *byte)
