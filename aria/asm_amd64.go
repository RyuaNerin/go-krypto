//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package aria

import (
	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"
	"github.com/RyuaNerin/go-krypto/internal/memory"
)

var hasSSSE3 = cpu.X86.HasSSSE3

func (ctx *ariaContextAsm) initRoundKey(key []byte) {
	ctx.ctx.initRoundKey(key)
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_SSSE3(memory.PByte(dst), memory.PByte(src), memory.PByte(rk), uint64(ctx.ctx.rounds))
}

//go:noescape
func __process_SSSE3(dst, src, rk *byte, rounds uint64)
