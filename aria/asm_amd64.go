//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package aria

import (
	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"
	"github.com/RyuaNerin/go-krypto/internal/ptr"
)

var hasSSSE3 = cpu.X86.HasSSSE3

func (ctx *ariaContextAsm) initRoundKey(key []byte) {
	ctx.ctx.initRoundKey(key)
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_SSSE3(ptr.PByte(dst), ptr.PByte(src), ptr.PByte(rk), uint64(ctx.ctx.rounds))
}

//go:noescape
func __process_SSSE3(dst, src, rk *byte, rounds uint64)
