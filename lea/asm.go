//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lea

import (
	"github.com/RyuaNerin/go-krypto/internal/memory"
)

func toAsmFunc(f func(ctx *leaContext, dst, src *byte)) func(ctx *leaContext, dst, src []byte) {
	return func(ctx *leaContext, dst, src []byte) {
		f(ctx, memory.PByte(dst), memory.PByte(src))
	}
}
