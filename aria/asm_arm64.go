//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package aria

import "github.com/RyuaNerin/go-krypto/internal/memory"

var hasNEON = true

func initRoundKeyAsm(ctx *ariaContextAsm, key []byte) {
	ctx.rounds = (len(key) + 32) / 4

	__encKeySetup_NEON(memory.P8(ctx.ek[:]), memory.P8(key), uint64(len(key)))

	ctx.dk = ctx.ek
	__decKeySetup_NEON(memory.P8(ctx.dk[:]), uint64(ctx.rounds))
}

func processAsm(dst, src, rk []byte, rounds int) {
	__process_NEON(memory.P8(dst), memory.P8(src), memory.P8(rk), uint64(rounds))
}

//go:noescape
func __encKeySetup_NEON(rk *byte, mk *byte, keyBytes uint64)

//go:noescape
func __decKeySetup_NEON(rk *byte, rounds uint64)

//go:noescape
func __process_NEON(dst, src, rk *byte, round uint64)
