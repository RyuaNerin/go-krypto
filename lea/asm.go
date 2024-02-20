//go:build (arm64 || amd64 || amd64p32) && !purego
// +build arm64 amd64 amd64p32
// +build !purego

package lea

import "github.com/RyuaNerin/go-krypto/internal/ptr"

func toAsmFunc(f func(ctx *leaContext, dst, src *byte)) func(ctx *leaContext, dst, src []byte) {
	return func(ctx *leaContext, dst, src []byte) {
		f(ctx, ptr.PByte(dst), ptr.PByte(src))
	}
}
