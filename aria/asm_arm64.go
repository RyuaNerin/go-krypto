//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package aria

import "github.com/RyuaNerin/go-krypto/internal/memory"

var hasNEON = true

func (ctx *ariaContextAsm) initRoundKey(key []byte) {
	__encKeySetup_NEON(memory.PByte(ctx.ctx.ek[:]), memory.PByte(key), uint64(len(key)))

	ctx.ctx.dk = ctx.ctx.ek
	__decKeySetup_NEON(memory.PByte(ctx.ctx.dk[:]), uint64(ctx.ctx.rounds))
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_NEON(memory.PByte(dst), memory.PByte(src), memory.PByte(rk), uint64(ctx.ctx.rounds))
}

//go:noescape
func __encKeySetup_NEON(rk *byte, mk *byte, keyBytes uint64)

//go:noescape
func __decKeySetup_NEON(rk *byte, rounds uint64)

//go:noescape
func __process_NEON(dst, src, rk *byte, round uint64)
