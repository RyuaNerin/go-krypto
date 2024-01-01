//go:build amd64 && gc && !purego

package lea

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2 && cpu.X86.HasAVX
)

func init() {
	leaEnc4 = leaEnc4SSE2
	leaDec4 = leaDec4SSE2

	leaEnc8 = leaEnc8SSE2
	leaDec8 = leaDec8SSE2

	if hasAVX2 {
		leaEnc8 = leaEnc8AVX2
		leaDec8 = leaDec8AVX2
	}
}

func leaEnc4SSE2(ctx *leaContext, dst, src []byte) {
	__lea_encrypt_4block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}
func leaDec4SSE2(ctx *leaContext, dst, src []byte) {
	__lea_decrypt_4block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}

func leaEnc8AVX2(ctx *leaContext, dst, src []byte) {
	__lea_encrypt_8block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}
func leaDec8AVX2(ctx *leaContext, dst, src []byte) {
	__lea_decrypt_8block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}

func leaEnc8SSE2(ctx *leaContext, dst, src []byte) {
	leaEnc4SSE2(ctx, dst[0x00:], src[0x00:])
	leaEnc4SSE2(ctx, dst[0x40:], src[0x40:])
}
func leaDec8SSE2(ctx *leaContext, dst, src []byte) {
	leaDec4SSE2(ctx, dst[0x00:], src[0x00:])
	leaDec4SSE2(ctx, dst[0x40:], src[0x40:])
}
