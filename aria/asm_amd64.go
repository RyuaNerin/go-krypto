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

var processAsm = processGo

func init() {
	if hasSSSE3 {
		processAsm = processSSSE3
	}
}

func initRoundKeyAsm(ctx *ariaContextAsm, key []byte) {
	initRoundKeyGo(&ctx.ariaContext, key)
}

func processSSSE3(dst, src, rk []byte, rounds int) {
	__process_SSSE3(memory.P8(dst), memory.P8(src), memory.P8(rk), uint64(rounds))
}

//go:noescape
func __process_SSSE3(dst, src, rk *byte, rounds uint64)
