//go:build amd64 && !purego
// +build amd64,!purego

package aria

import (
	"kryptosimd/internal/ptr"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
)

//go:noescape
func __process_SSSE3(dst, src, rk *byte, rounds uint64)

var newCipher = newCipherGo

func init() {
	if hasSSSE3 {
		newCipher = newCipherAsm
	}
}

func (ctx *ariaContextAsm) initRoundKey(key []byte) {
	ctx.ctx.initRoundKey(key)
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_SSSE3(ptr.BytePtr(dst), ptr.BytePtr(src), ptr.BytePtr(rk), uint64(ctx.ctx.rounds))
}
